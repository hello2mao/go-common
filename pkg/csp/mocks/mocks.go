package mocks

import (
	"bytes"
	"crypto"
	"errors"
	"hash"
	"reflect"

	"github.com/hello2mao/go-common/pkg/csp"
)

type MockCSP struct {
	SignArgKey    csp.Key
	SignDigestArg []byte
	SignOptsArg   csp.SignerOpts

	SignValue []byte
	SignErr   error

	VerifyValue bool
	VerifyErr   error

	ExpectedSig []byte

	KeyImportValue csp.Key
	KeyImportErr   error

	EncryptError error
	DecryptError error

	HashVal []byte
	HashErr error
}

func (*MockCSP) KeyGen(opts csp.KeyGenOpts) (csp.Key, error) {
	panic("Not yet implemented")
}

func (*MockCSP) KeyDeriv(k csp.Key, opts csp.KeyDerivOpts) (csp.Key, error) {
	panic("Not yet implemented")
}

func (m *MockCSP) KeyImport(raw interface{}, opts csp.KeyImportOpts) (csp.Key, error) {
	return m.KeyImportValue, m.KeyImportErr
}

func (*MockCSP) GetKey(ski []byte) (csp.Key, error) {
	panic("Not yet implemented")
}

func (m *MockCSP) Hash(msg []byte, opts csp.HashOpts) ([]byte, error) {
	return m.HashVal, m.HashErr
}

func (*MockCSP) GetHash(opts csp.HashOpts) (hash.Hash, error) {
	panic("Not yet implemented")
}

func (b *MockCSP) Sign(k csp.Key, digest []byte, opts csp.SignerOpts) ([]byte, error) {
	if !reflect.DeepEqual(b.SignArgKey, k) {
		return nil, errors.New("invalid key")
	}
	if !reflect.DeepEqual(b.SignDigestArg, digest) {
		return nil, errors.New("invalid digest")
	}
	if !reflect.DeepEqual(b.SignOptsArg, opts) {
		return nil, errors.New("invalid opts")
	}

	return b.SignValue, b.SignErr
}

func (b *MockCSP) Verify(k csp.Key, signature, digest []byte, opts csp.SignerOpts) (bool, error) {
	// we want to mock a success
	if b.VerifyValue {
		return b.VerifyValue, nil
	}

	// we want to mock a failure because of an error
	if b.VerifyErr != nil {
		return b.VerifyValue, b.VerifyErr
	}

	// in neither case, compare the signature with the expected one
	return bytes.Equal(b.ExpectedSig, signature), nil
}

func (m *MockCSP) Encrypt(k csp.Key, plaintext []byte, opts csp.EncrypterOpts) ([]byte, error) {
	if m.EncryptError == nil {
		return plaintext, nil
	} else {
		return nil, m.EncryptError
	}
}

func (m *MockCSP) Decrypt(k csp.Key, ciphertext []byte, opts csp.DecrypterOpts) ([]byte, error) {
	if m.DecryptError == nil {
		return ciphertext, nil
	} else {
		return nil, m.DecryptError
	}
}

type MockKey struct {
	BytesValue []byte
	BytesErr   error
	Symm       bool
	PK         csp.Key
	PKErr      error
	Pvt        bool
}

func (m *MockKey) Bytes() ([]byte, error) {
	return m.BytesValue, m.BytesErr
}

func (*MockKey) SKI() []byte {
	panic("Not yet implemented")
}

func (m *MockKey) Symmetric() bool {
	return m.Symm
}

func (m *MockKey) Private() bool {
	return m.Pvt
}

func (m *MockKey) PublicKey() (csp.Key, error) {
	return m.PK, m.PKErr
}

type SignerOpts struct {
	HashFuncValue crypto.Hash
}

func (o *SignerOpts) HashFunc() crypto.Hash {
	return o.HashFuncValue
}

type KeyGenOpts struct {
	EphemeralValue bool
}

func (*KeyGenOpts) Algorithm() string {
	return "Mock KeyGenOpts"
}

func (o *KeyGenOpts) Ephemeral() bool {
	return o.EphemeralValue
}

type KeyStore struct {
	GetKeyValue csp.Key
	GetKeyErr   error
	StoreKeyErr error
}

func (*KeyStore) ReadOnly() bool {
	panic("Not yet implemented")
}

func (ks *KeyStore) GetKey(ski []byte) (csp.Key, error) {
	return ks.GetKeyValue, ks.GetKeyErr
}

func (ks *KeyStore) StoreKey(k csp.Key) error {
	return ks.StoreKeyErr
}

type KeyImportOpts struct{}

func (*KeyImportOpts) Algorithm() string {
	return "Mock KeyImportOpts"
}

func (*KeyImportOpts) Ephemeral() bool {
	panic("Not yet implemented")
}

type EncrypterOpts struct{}
type DecrypterOpts struct{}

type HashOpts struct{}

func (HashOpts) Algorithm() string {
	return "Mock HashOpts"
}

type KeyDerivOpts struct {
	EphemeralValue bool
}

func (*KeyDerivOpts) Algorithm() string {
	return "Mock KeyDerivOpts"
}

func (o *KeyDerivOpts) Ephemeral() bool {
	return o.EphemeralValue
}
