package sw

import (
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/sha512"
	"reflect"

	"github.com/hello2mao/go-common/csp"
	"github.com/pkg/errors"
	"golang.org/x/crypto/sha3"
)

// NewDefaultSecurityLevel returns a new instance of the software-based CSP
// at security level 256, hash family SHA2 and using FolderBasedKeyStore as KeyStore.
func NewDefaultSecurityLevel(keyStorePath string) (csp.CSP, error) {
	ks := &fileBasedKeyStore{}
	if err := ks.Init(nil, keyStorePath, false); err != nil {
		return nil, errors.Wrapf(err, "Failed initializing key store at [%v]", keyStorePath)
	}

	return NewWithParams(256, "SHA2", ks)
}

// NewDefaultSecurityLevel returns a new instance of the software-based CSP
// at security level 256, hash family SHA2 and using the passed KeyStore.
func NewDefaultSecurityLevelWithKeystore(keyStore csp.KeyStore) (csp.CSP, error) {
	return NewWithParams(256, "SHA2", keyStore)
}

// NewWithParams returns a new instance of the software-based CSP
// set at the passed security level, hash family and KeyStore.
func NewWithParams(securityLevel int, hashFamily string, keyStore csp.KeyStore) (csp.CSP, error) {
	// Init config
	conf := &config{}
	err := conf.setSecurityLevel(securityLevel, hashFamily)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed initializing configuration at [%v,%v]", securityLevel, hashFamily)
	}

	swcsp, err := New(keyStore)
	if err != nil {
		return nil, err
	}

	// Notice that errors are ignored here because some test will fail if one
	// of the following call fails.

	// Set the Encryptors
	swcsp.AddWrapper(reflect.TypeOf(&aesPrivateKey{}), &aescbcpkcs7Encryptor{})

	// Set the Decryptors
	swcsp.AddWrapper(reflect.TypeOf(&aesPrivateKey{}), &aescbcpkcs7Decryptor{})

	// Set the Signers
	swcsp.AddWrapper(reflect.TypeOf(&ecdsaPrivateKey{}), &ecdsaSigner{})

	// Set the Verifiers
	swcsp.AddWrapper(reflect.TypeOf(&ecdsaPrivateKey{}), &ecdsaPrivateKeyVerifier{})
	swcsp.AddWrapper(reflect.TypeOf(&ecdsaPublicKey{}), &ecdsaPublicKeyKeyVerifier{})

	// Set the Hashers
	swcsp.AddWrapper(reflect.TypeOf(&csp.SHAOpts{}), &hasher{hash: conf.hashFunction})
	swcsp.AddWrapper(reflect.TypeOf(&csp.SHA256Opts{}), &hasher{hash: sha256.New})
	swcsp.AddWrapper(reflect.TypeOf(&csp.SHA384Opts{}), &hasher{hash: sha512.New384})
	swcsp.AddWrapper(reflect.TypeOf(&csp.SHA3_256Opts{}), &hasher{hash: sha3.New256})
	swcsp.AddWrapper(reflect.TypeOf(&csp.SHA3_384Opts{}), &hasher{hash: sha3.New384})

	// Set the key generators
	swcsp.AddWrapper(reflect.TypeOf(&csp.ECDSAKeyGenOpts{}), &ecdsaKeyGenerator{curve: conf.ellipticCurve})
	swcsp.AddWrapper(reflect.TypeOf(&csp.ECDSAP256KeyGenOpts{}), &ecdsaKeyGenerator{curve: elliptic.P256()})
	swcsp.AddWrapper(reflect.TypeOf(&csp.ECDSAP384KeyGenOpts{}), &ecdsaKeyGenerator{curve: elliptic.P384()})
	swcsp.AddWrapper(reflect.TypeOf(&csp.AESKeyGenOpts{}), &aesKeyGenerator{length: conf.aesBitLength})
	swcsp.AddWrapper(reflect.TypeOf(&csp.AES256KeyGenOpts{}), &aesKeyGenerator{length: 32})
	swcsp.AddWrapper(reflect.TypeOf(&csp.AES192KeyGenOpts{}), &aesKeyGenerator{length: 24})
	swcsp.AddWrapper(reflect.TypeOf(&csp.AES128KeyGenOpts{}), &aesKeyGenerator{length: 16})

	// Set the key deriver
	swcsp.AddWrapper(reflect.TypeOf(&ecdsaPrivateKey{}), &ecdsaPrivateKeyKeyDeriver{})
	swcsp.AddWrapper(reflect.TypeOf(&ecdsaPublicKey{}), &ecdsaPublicKeyKeyDeriver{})
	swcsp.AddWrapper(reflect.TypeOf(&aesPrivateKey{}), &aesPrivateKeyKeyDeriver{conf: conf})

	// Set the key importers
	swcsp.AddWrapper(reflect.TypeOf(&csp.AES256ImportKeyOpts{}), &aes256ImportKeyOptsKeyImporter{})
	swcsp.AddWrapper(reflect.TypeOf(&csp.HMACImportKeyOpts{}), &hmacImportKeyOptsKeyImporter{})
	swcsp.AddWrapper(reflect.TypeOf(&csp.ECDSAPKIXPublicKeyImportOpts{}), &ecdsaPKIXPublicKeyImportOptsKeyImporter{})
	swcsp.AddWrapper(reflect.TypeOf(&csp.ECDSAPrivateKeyImportOpts{}), &ecdsaPrivateKeyImportOptsKeyImporter{})
	swcsp.AddWrapper(reflect.TypeOf(&csp.ECDSAGoPublicKeyImportOpts{}), &ecdsaGoPublicKeyImportOptsKeyImporter{})
	swcsp.AddWrapper(reflect.TypeOf(&csp.X509PublicKeyImportOpts{}), &x509PublicKeyImportOptsKeyImporter{csp: swcsp})

	return swcsp, nil
}
