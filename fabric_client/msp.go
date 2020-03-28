package fabric_client

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/tjfoc/gmsm/sm2"
)

// RegisterAndEnrollUser register and enroll a user
// This user is new
func (c *FabricClient) RegisterAndEnrollUser(newUserName, sdkPath, keyStorePath string) error {
	if err := c.createMspClientIfNotExist(); err != nil {
		return err
	}
	// check user
	if FileOrDirExist(path.Join(sdkPath, c.FabricClientConfig.OrgName, "users", newUserName+"@"+c.FabricClientConfig.OrgName)) {
		return fmt.Errorf("target user already exist")
	}
	userRegisterRequest := &msp.RegistrationRequest{
		Name:        newUserName,
		Type:        "user",
		Affiliation: c.FabricClientConfig.OrgName,
	}
	// Register
	secret, err := c.MspClient.Register(userRegisterRequest)
	if err != nil {
		return errors.WithMessage(err, "failed to register")
	}
	// Enroll
	err = c.MspClient.Enroll(userRegisterRequest.Name, msp.WithSecret(secret))
	if err != nil {
		return errors.WithMessage(err, "failed to enroll")
	}
	identity, err := c.MspClient.GetSigningIdentity(userRegisterRequest.Name)
	if err != nil {
		return errors.WithMessage(err, "failed to get signing identity")
	}
	cert := identity.EnrollmentCertificate()

	// Start to gen file
	// ├── msp
	// │   ├── admincerts
	// │   │   └── {user}@{orgName}-cert.pem
	// │   ├── cacerts
	// │   │   └── ca.{orgName}-cert.pem
	// │   ├── keystore
	// │   │   └──{privateKeyName}_sk
	// │   ├── signcerts
	// │   │   └── {user}@{orgName}-cert.pem
	// │   └── tlscacerts
	// │       └── tlsca.{orgName}-cert.pem
	// └── tls
	//    ├── ca.crt
	//    ├── fabric-client.crt
	//    └── fabric-client.key
	log.Infof("[RegisterAndEnrollUser]start to gen file")

	// msp admincerts
	admincertsPath := path.Join(
		sdkPath,
		c.FabricClientConfig.OrgName,
		"users",
		newUserName+"@"+c.FabricClientConfig.OrgName,
		"msp/admincerts",
	)
	admincertsFile := newUserName + "@" + c.FabricClientConfig.OrgName + "-cert.pem"
	if err = os.MkdirAll(admincertsPath, os.ModePerm); err != nil {
		return fmt.Errorf("MkdirAll msp admincerts failed: %v", err)
	}
	err = ioutil.WriteFile(
		path.Join(
			admincertsPath,
			admincertsFile,
		),
		cert,
		os.ModePerm,
	)
	if err != nil {
		return fmt.Errorf("write msp admincerts file failed: %v", err)
	}

	// msp cacerts copied from admin
	cacertsPath := path.Join(
		sdkPath,
		c.FabricClientConfig.OrgName,
		"users",
		newUserName+"@"+c.FabricClientConfig.OrgName,
		"msp/cacerts",
	)
	cacertsFile := "ca." + c.FabricClientConfig.OrgName + "-cert.pem"
	if err = os.MkdirAll(cacertsPath, os.ModePerm); err != nil {
		return fmt.Errorf("MkdirAll msp cacerts failed: %v", err)
	}
	err = CopyFile(
		path.Join(
			sdkPath,
			c.FabricClientConfig.OrgName,
			"users",
			"Admin@"+c.FabricClientConfig.OrgName,
			"msp/cacerts",
			"ca."+c.FabricClientConfig.OrgName+"-cert.pem",
		),
		path.Join(
			cacertsPath,
			cacertsFile,
		),
	)
	if err != nil {
		return fmt.Errorf("write msp cacerts file failed: %v", err)
	}

	// msp keystore
	files, _ := ioutil.ReadDir(keyStorePath)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if f.Size() > (1 << 16) { //64k, somewhat arbitrary limit, considering even large RSA keys
			continue
		}
		raw, err := ioutil.ReadFile(path.Join(keyStorePath, f.Name()))
		if err != nil {
			continue
		}
		key, err := PEMtoPrivateKey(raw, nil)
		if err != nil {
			continue
		}

		if c.FabricClientConfig.EnableGM {
			k := &sm2PrivateKey{key.(*sm2.PrivateKey)}
			if !bytes.Equal(k.SKI(), identity.PrivateKey().SKI()) {
				continue
			}
		} else {
			k := &ecdsaPrivateKey{key.(*ecdsa.PrivateKey)}
			if !bytes.Equal(k.SKI(), identity.PrivateKey().SKI()) {
				continue
			}
		}

		// write file
		keystorePath := path.Join(
			sdkPath,
			c.FabricClientConfig.OrgName,
			"users",
			newUserName+"@"+c.FabricClientConfig.OrgName,
			"msp/keystore",
		)
		keystoreFile := f.Name()
		if err = os.MkdirAll(keystorePath, os.ModePerm); err != nil {
			return fmt.Errorf("MkdirAll msp keystore failed: %v", err)
		}
		err = ioutil.WriteFile(
			path.Join(
				keystorePath,
				keystoreFile,
			),
			raw,
			os.ModePerm,
		)
		if err != nil {
			return fmt.Errorf("write msp keystore file failed: %v", err)
		}
		break
	}

	// msp signcerts
	signcertsPath := path.Join(
		sdkPath,
		c.FabricClientConfig.OrgName,
		"users",
		newUserName+"@"+c.FabricClientConfig.OrgName,
		"msp/signcerts",
	)
	signcertsFile := newUserName + "@" + c.FabricClientConfig.OrgName + "-cert.pem"
	if err = os.MkdirAll(signcertsPath, os.ModePerm); err != nil {
		return fmt.Errorf("MkdirAll msp admincerts failed: %v", err)
	}
	err = ioutil.WriteFile(
		path.Join(
			signcertsPath,
			signcertsFile,
		),
		cert,
		os.ModePerm,
	)
	if err != nil {
		return fmt.Errorf("write msp admincerts file failed: %v", err)
	}

	// msp tlscacerts copied from admin
	tlscacertsPath := path.Join(
		sdkPath,
		c.FabricClientConfig.OrgName,
		"users",
		newUserName+"@"+c.FabricClientConfig.OrgName,
		"msp/tlscacerts",
	)
	tlscacertsFile := "tlsca." + c.FabricClientConfig.OrgName + "-cert.pem"
	if err = os.MkdirAll(tlscacertsPath, os.ModePerm); err != nil {
		return fmt.Errorf("MkdirAll msp cacerts failed: %v", err)
	}
	err = CopyFile(
		path.Join(
			sdkPath,
			c.FabricClientConfig.OrgName,
			"users",
			"Admin@"+c.FabricClientConfig.OrgName,
			"msp/tlscacerts",
			"tlsca."+c.FabricClientConfig.OrgName+"-cert.pem",
		),
		path.Join(
			tlscacertsPath,
			tlscacertsFile,
		),
	)
	if err != nil {
		return fmt.Errorf("write msp tlscacerts file failed: %v", err)
	}

	// tls copied from admin
	tlsPath := path.Join(
		sdkPath,
		c.FabricClientConfig.OrgName,
		"users",
		newUserName+"@"+c.FabricClientConfig.OrgName,
		"tls",
	)
	caCertFile := "ca.crt"
	clientCertFile := "client.crt"
	clientKeyFile := "client.key"
	if err = os.MkdirAll(tlsPath, os.ModePerm); err != nil {
		return fmt.Errorf("MkdirAll tls failed: %v", err)
	}
	err = CopyFile(
		path.Join(
			sdkPath,
			c.FabricClientConfig.OrgName,
			"users",
			"Admin@"+c.FabricClientConfig.OrgName,
			"tls",
			caCertFile,
		),
		path.Join(
			tlsPath,
			caCertFile,
		),
	)
	if err != nil {
		return fmt.Errorf("write caCertFile file failed: %v", err)
	}
	err = CopyFile(
		path.Join(
			sdkPath,
			c.FabricClientConfig.OrgName,
			"users",
			"Admin@"+c.FabricClientConfig.OrgName,
			"tls",
			clientCertFile,
		),
		path.Join(
			tlsPath,
			clientCertFile,
		),
	)
	if err != nil {
		return fmt.Errorf("write clientCertFile file failed: %v", err)
	}
	err = CopyFile(
		path.Join(
			sdkPath,
			c.FabricClientConfig.OrgName,
			"users",
			"Admin@"+c.FabricClientConfig.OrgName,
			"tls",
			clientKeyFile,
		),
		path.Join(
			tlsPath,
			clientKeyFile,
		),
	)
	if err != nil {
		return fmt.Errorf("write clientKeyFile file failed: %v", err)
	}

	log.Infof("[RegisterAndEnrollUser]end of gen file")

	return nil
}

// DeleteUser delete target user
func (c *FabricClient) DeleteUser(userName, sdkPath string) {
	log.Debugf("start to delete user: %s", userName)

	// delete local file
	userFilePath := path.Join(sdkPath, c.FabricClientConfig.OrgName, "users", userName+"@"+c.FabricClientConfig.OrgName)
	if FileOrDirExist(userFilePath) {
		err := os.RemoveAll(userFilePath)
		if err != nil {
			log.Errorf("DeleteUser RemoveAll err: %v", err)
			// continue
		}
	}

	if err := c.createMspClientIfNotExist(); err != nil {
		log.Fatalf("createMspClientIfNotExist err: %v", err)
		return
	}

	// remove identity
	_, err := c.MspClient.RemoveIdentity(&msp.RemoveIdentityRequest{
		ID:    userName,
		Force: true,
	})
	if err != nil {
		log.Debugf("RemoveIdentity err: %+v", err)
		return
	}
	log.Debugf("DeleteUser %s success.", userName)
}

func (c *FabricClient) createMspClientIfNotExist() error {
	if c.MspClient != nil {
		return nil
	}
	if c.FabricClientConfig.OrgName == "" {
		return fmt.Errorf("OrgName must not empty")
	}
	mspClient, err := msp.New(c.Sdk.Context(), msp.WithOrg(c.FabricClientConfig.OrgName))
	if err != nil {
		return errors.WithMessage(err, "failed to create MSP fabric-client")
	}
	c.MspClient = mspClient
	return nil
}

type ecdsaPrivateKey struct {
	privKey *ecdsa.PrivateKey
}

func (k *ecdsaPrivateKey) SKI() (ski []byte) {
	if k.privKey == nil {
		return nil
	}

	// Marshall the public key
	raw := elliptic.Marshal(k.privKey.Curve, k.privKey.PublicKey.X, k.privKey.PublicKey.Y)

	// Hash it
	hash := sha256.New()
	hash.Write(raw)
	return hash.Sum(nil)
}

type sm2PrivateKey struct {
	privKey *sm2.PrivateKey
}

func (k *sm2PrivateKey) SKI() (ski []byte) {
	if k.privKey == nil {
		return nil
	}

	//Marshall the public key
	raw := elliptic.Marshal(k.privKey.Curve, k.privKey.PublicKey.X, k.privKey.PublicKey.Y)

	// Hash it
	hash := sha256.New()
	hash.Write(raw)
	return hash.Sum(nil)
}
