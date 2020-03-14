package sw

import (
	"hash"

	"github.com/hello2mao/go-common/csp"
)

// KeyGenerator is a CSP-like interface that provides key generation algorithms
type KeyGenerator interface {

	// KeyGen generates a key using opts.
	KeyGen(opts csp.KeyGenOpts) (k csp.Key, err error)
}

// KeyDeriver is a CSP-like interface that provides key derivation algorithms
type KeyDeriver interface {

	// KeyDeriv derives a key from k using opts.
	// The opts argument should be appropriate for the primitive used.
	KeyDeriv(k csp.Key, opts csp.KeyDerivOpts) (dk csp.Key, err error)
}

// KeyImporter is a CSP-like interface that provides key import algorithms
type KeyImporter interface {

	// KeyImport imports a key from its raw representation using opts.
	// The opts argument should be appropriate for the primitive used.
	KeyImport(raw interface{}, opts csp.KeyImportOpts) (k csp.Key, err error)
}

// Encryptor is a CSP-like interface that provides encryption algorithms
type Encryptor interface {

	// Encrypt encrypts plaintext using key k.
	// The opts argument should be appropriate for the algorithm used.
	Encrypt(k csp.Key, plaintext []byte, opts csp.EncrypterOpts) (ciphertext []byte, err error)
}

// Decryptor is a CSP-like interface that provides decryption algorithms
type Decryptor interface {

	// Decrypt decrypts ciphertext using key k.
	// The opts argument should be appropriate for the algorithm used.
	Decrypt(k csp.Key, ciphertext []byte, opts csp.DecrypterOpts) (plaintext []byte, err error)
}

// Signer is a CSP-like interface that provides signing algorithms
type Signer interface {

	// Sign signs digest using key k.
	// The opts argument should be appropriate for the algorithm used.
	//
	// Note that when a signature of a hash of a larger message is needed,
	// the caller is responsible for hashing the larger message and passing
	// the hash (as digest).
	Sign(k csp.Key, digest []byte, opts csp.SignerOpts) (signature []byte, err error)
}

// Verifier is a CSP-like interface that provides verifying algorithms
type Verifier interface {

	// Verify verifies signature against key k and digest
	// The opts argument should be appropriate for the algorithm used.
	Verify(k csp.Key, signature, digest []byte, opts csp.SignerOpts) (valid bool, err error)
}

// Hasher is a CSP-like interface that provides hash algorithms
type Hasher interface {

	// Hash hashes messages msg using options opts.
	// If opts is nil, the default hash function will be used.
	Hash(msg []byte, opts csp.HashOpts) (hash []byte, err error)

	// GetHash returns and instance of hash.Hash using options opts.
	// If opts is nil, the default hash function will be returned.
	GetHash(opts csp.HashOpts) (h hash.Hash, err error)
}
