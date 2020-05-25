package mocks

import (
	"errors"
	"hash"
	"reflect"

	"github.com/hello2mao/go-common/incubator/csp"
)

type Encryptor struct {
	KeyArg       csp.Key
	PlaintextArg []byte
	OptsArg      csp.EncrypterOpts

	EncValue []byte
	EncErr   error
}

func (e *Encryptor) Encrypt(k csp.Key, plaintext []byte, opts csp.EncrypterOpts) (ciphertext []byte, err error) {
	if !reflect.DeepEqual(e.KeyArg, k) {
		return nil, errors.New("invalid key")
	}
	if !reflect.DeepEqual(e.PlaintextArg, plaintext) {
		return nil, errors.New("invalid plaintext")
	}
	if !reflect.DeepEqual(e.OptsArg, opts) {
		return nil, errors.New("invalid opts")
	}

	return e.EncValue, e.EncErr
}

type Decryptor struct {
}

func (*Decryptor) Decrypt(k csp.Key, ciphertext []byte, opts csp.DecrypterOpts) (plaintext []byte, err error) {
	panic("implement me")
}

type Signer struct {
	KeyArg    csp.Key
	DigestArg []byte
	OptsArg   csp.SignerOpts

	Value []byte
	Err   error
}

func (s *Signer) Sign(k csp.Key, digest []byte, opts csp.SignerOpts) (signature []byte, err error) {
	if !reflect.DeepEqual(s.KeyArg, k) {
		return nil, errors.New("invalid key")
	}
	if !reflect.DeepEqual(s.DigestArg, digest) {
		return nil, errors.New("invalid digest")
	}
	if !reflect.DeepEqual(s.OptsArg, opts) {
		return nil, errors.New("invalid opts")
	}

	return s.Value, s.Err
}

type Verifier struct {
	KeyArg       csp.Key
	SignatureArg []byte
	DigestArg    []byte
	OptsArg      csp.SignerOpts

	Value bool
	Err   error
}

func (s *Verifier) Verify(k csp.Key, signature, digest []byte, opts csp.SignerOpts) (valid bool, err error) {
	if !reflect.DeepEqual(s.KeyArg, k) {
		return false, errors.New("invalid key")
	}
	if !reflect.DeepEqual(s.SignatureArg, signature) {
		return false, errors.New("invalid signature")
	}
	if !reflect.DeepEqual(s.DigestArg, digest) {
		return false, errors.New("invalid digest")
	}
	if !reflect.DeepEqual(s.OptsArg, opts) {
		return false, errors.New("invalid opts")
	}

	return s.Value, s.Err
}

type Hasher struct {
	MsgArg  []byte
	OptsArg csp.HashOpts

	Value     []byte
	ValueHash hash.Hash
	Err       error
}

func (h *Hasher) Hash(msg []byte, opts csp.HashOpts) (hash []byte, err error) {
	if !reflect.DeepEqual(h.MsgArg, msg) {
		return nil, errors.New("invalid message")
	}
	if !reflect.DeepEqual(h.OptsArg, opts) {
		return nil, errors.New("invalid opts")
	}

	return h.Value, h.Err
}

func (h *Hasher) GetHash(opts csp.HashOpts) (hash.Hash, error) {
	if !reflect.DeepEqual(h.OptsArg, opts) {
		return nil, errors.New("invalid opts")
	}

	return h.ValueHash, h.Err
}

type KeyGenerator struct {
	OptsArg csp.KeyGenOpts

	Value csp.Key
	Err   error
}

func (kg *KeyGenerator) KeyGen(opts csp.KeyGenOpts) (k csp.Key, err error) {
	if !reflect.DeepEqual(kg.OptsArg, opts) {
		return nil, errors.New("invalid opts")
	}

	return kg.Value, kg.Err
}

type KeyDeriver struct {
	KeyArg  csp.Key
	OptsArg csp.KeyDerivOpts

	Value csp.Key
	Err   error
}

func (kd *KeyDeriver) KeyDeriv(k csp.Key, opts csp.KeyDerivOpts) (dk csp.Key, err error) {
	if !reflect.DeepEqual(kd.KeyArg, k) {
		return nil, errors.New("invalid key")
	}
	if !reflect.DeepEqual(kd.OptsArg, opts) {
		return nil, errors.New("invalid opts")
	}

	return kd.Value, kd.Err
}

type KeyImporter struct {
	RawArg  []byte
	OptsArg csp.KeyImportOpts

	Value csp.Key
	Err   error
}

func (ki *KeyImporter) KeyImport(raw interface{}, opts csp.KeyImportOpts) (k csp.Key, err error) {
	if !reflect.DeepEqual(ki.RawArg, raw) {
		return nil, errors.New("invalid raw")
	}
	if !reflect.DeepEqual(ki.OptsArg, opts) {
		return nil, errors.New("invalid opts")
	}

	return ki.Value, ki.Err
}
