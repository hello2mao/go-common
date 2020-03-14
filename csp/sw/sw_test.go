package sw

import (
	"errors"
	"reflect"
	"testing"

	"github.com/hello2mao/go-common/csp"
	"github.com/hello2mao/go-common/csp/mocks"
	mocks2 "github.com/hello2mao/go-common/csp/sw/mocks"
	"github.com/stretchr/testify/assert"
)

func TestKeyGenInvalidInputs(t *testing.T) {
	// Init a CSP instance with a key store that returns an error on store
	cspInstance, err := NewWithParams(256, "SHA2", &mocks.KeyStore{StoreKeyErr: errors.New("cannot store key")})
	assert.NoError(t, err)

	_, err = cspInstance.KeyGen(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid Opts parameter. It must not be nil.")

	_, err = cspInstance.KeyGen(&mocks.KeyGenOpts{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Unsupported 'KeyGenOpts' provided [")

	_, err = cspInstance.KeyGen(&csp.ECDSAP256KeyGenOpts{})
	assert.Error(t, err, "Generation of a non-ephemeral key must fail. KeyStore is programmed to fail.")
	assert.Contains(t, err.Error(), "cannot store key")
}

func TestKeyDerivInvalidInputs(t *testing.T) {
	cspInstance, err := NewWithParams(256, "SHA2", &mocks.KeyStore{StoreKeyErr: errors.New("cannot store key")})
	assert.NoError(t, err)

	_, err = cspInstance.KeyDeriv(nil, &csp.ECDSAReRandKeyOpts{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid Key. It must not be nil.")

	_, err = cspInstance.KeyDeriv(&mocks.MockKey{}, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid opts. It must not be nil.")

	_, err = cspInstance.KeyDeriv(&mocks.MockKey{}, &csp.ECDSAReRandKeyOpts{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Unsupported 'Key' provided [")

	keyDerivers := make(map[reflect.Type]KeyDeriver)
	keyDerivers[reflect.TypeOf(&mocks.MockKey{})] = &mocks2.KeyDeriver{
		KeyArg:  &mocks.MockKey{},
		OptsArg: &mocks.KeyDerivOpts{EphemeralValue: false},
		Value:   nil,
		Err:     nil,
	}
	cspInstance.(*CSP).KeyDerivers = keyDerivers
	_, err = cspInstance.KeyDeriv(&mocks.MockKey{}, &mocks.KeyDerivOpts{EphemeralValue: false})
	assert.Error(t, err, "KeyDerivation of a non-ephemeral key must fail. KeyStore is programmed to fail.")
	assert.Contains(t, err.Error(), "cannot store key")
}

func TestKeyImportInvalidInputs(t *testing.T) {
	cspInstance, err := NewWithParams(256, "SHA2", &mocks.KeyStore{})
	assert.NoError(t, err)

	_, err = cspInstance.KeyImport(nil, &csp.AES256ImportKeyOpts{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid raw. It must not be nil.")

	_, err = cspInstance.KeyImport([]byte{0, 1, 2, 3, 4}, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid opts. It must not be nil.")

	_, err = cspInstance.KeyImport([]byte{0, 1, 2, 3, 4}, &mocks.KeyImportOpts{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Unsupported 'KeyImportOpts' provided [")
}

func TestGetKeyInvalidInputs(t *testing.T) {
	// Init a CSP instance with a key store that returns an error on get
	cspInstance, err := NewWithParams(256, "SHA2", &mocks.KeyStore{GetKeyErr: errors.New("cannot get key")})
	assert.NoError(t, err)

	_, err = cspInstance.GetKey(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot get key")

	// Init a CSP instance with a key store that returns a given key
	k := &mocks.MockKey{}
	cspInstance, err = NewWithParams(256, "SHA2", &mocks.KeyStore{GetKeyValue: k})
	assert.NoError(t, err)
	// No SKI is needed here
	k2, err := cspInstance.GetKey(nil)
	assert.NoError(t, err)
	assert.Equal(t, k, k2, "Keys must be the same.")
}

func TestSignInvalidInputs(t *testing.T) {
	cspInstance, err := NewWithParams(256, "SHA2", &mocks.KeyStore{})
	assert.NoError(t, err)

	_, err = cspInstance.Sign(nil, []byte{1, 2, 3, 5}, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid Key. It must not be nil.")

	_, err = cspInstance.Sign(&mocks.MockKey{}, nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid digest. Cannot be empty.")

	_, err = cspInstance.Sign(&mocks.MockKey{}, []byte{1, 2, 3, 5}, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Unsupported 'SignKey' provided [")
}

func TestVerifyInvalidInputs(t *testing.T) {
	cspInstance, err := NewWithParams(256, "SHA2", &mocks.KeyStore{})
	assert.NoError(t, err)

	_, err = cspInstance.Verify(nil, []byte{1, 2, 3, 5}, []byte{1, 2, 3, 5}, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid Key. It must not be nil.")

	_, err = cspInstance.Verify(&mocks.MockKey{}, nil, []byte{1, 2, 3, 5}, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid signature. Cannot be empty.")

	_, err = cspInstance.Verify(&mocks.MockKey{}, []byte{1, 2, 3, 5}, nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid digest. Cannot be empty.")

	_, err = cspInstance.Verify(&mocks.MockKey{}, []byte{1, 2, 3, 5}, []byte{1, 2, 3, 5}, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Unsupported 'VerifyKey' provided [")
}

func TestEncryptInvalidInputs(t *testing.T) {
	cspInstance, err := NewWithParams(256, "SHA2", &mocks.KeyStore{})
	assert.NoError(t, err)

	_, err = cspInstance.Encrypt(nil, []byte{1, 2, 3, 4}, &csp.AESCBCPKCS7ModeOpts{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid Key. It must not be nil.")

	_, err = cspInstance.Encrypt(&mocks.MockKey{}, []byte{1, 2, 3, 4}, &csp.AESCBCPKCS7ModeOpts{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Unsupported 'EncryptKey' provided [")
}

func TestDecryptInvalidInputs(t *testing.T) {
	cspInstance, err := NewWithParams(256, "SHA2", &mocks.KeyStore{})
	assert.NoError(t, err)

	_, err = cspInstance.Decrypt(nil, []byte{1, 2, 3, 4}, &csp.AESCBCPKCS7ModeOpts{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid Key. It must not be nil.")

	_, err = cspInstance.Decrypt(&mocks.MockKey{}, []byte{1, 2, 3, 4}, &csp.AESCBCPKCS7ModeOpts{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Unsupported 'DecryptKey' provided [")
}

func TestHashInvalidInputs(t *testing.T) {
	cspInstance, err := NewWithParams(256, "SHA2", &mocks.KeyStore{})
	assert.NoError(t, err)

	_, err = cspInstance.Hash(nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid opts. It must not be nil.")

	_, err = cspInstance.Hash(nil, &mocks.HashOpts{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Unsupported 'HashOpt' provided [")
}

func TestGetHashInvalidInputs(t *testing.T) {
	cspInstance, err := NewWithParams(256, "SHA2", &mocks.KeyStore{})
	assert.NoError(t, err)

	_, err = cspInstance.GetHash(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid opts. It must not be nil.")

	_, err = cspInstance.GetHash(&mocks.HashOpts{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Unsupported 'HashOpt' provided [")
}
