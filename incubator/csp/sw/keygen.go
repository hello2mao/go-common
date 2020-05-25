package sw

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"

	"github.com/hello2mao/go-common/incubator/csp"
)

type ecdsaKeyGenerator struct {
	curve elliptic.Curve
}

func (kg *ecdsaKeyGenerator) KeyGen(opts csp.KeyGenOpts) (csp.Key, error) {
	privKey, err := ecdsa.GenerateKey(kg.curve, rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("Failed generating ECDSA key for [%v]: [%s]", kg.curve, err)
	}

	return &ecdsaPrivateKey{privKey}, nil
}

type aesKeyGenerator struct {
	length int
}

func (kg *aesKeyGenerator) KeyGen(opts csp.KeyGenOpts) (csp.Key, error) {
	lowLevelKey, err := GetRandomBytes(int(kg.length))
	if err != nil {
		return nil, fmt.Errorf("Failed generating AES %d key [%s]", kg.length, err)
	}

	return &aesPrivateKey{lowLevelKey, false}, nil
}
