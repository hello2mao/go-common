package signer

import (
	"crypto"
	"io"

	"github.com/hello2mao/go-common/csp"
	"github.com/hello2mao/go-common/csp/utils"
	"github.com/pkg/errors"
)

// cspCryptoSigner is the csp-based implementation of a crypto.Signer
type cspCryptoSigner struct {
	csp csp.CSP
	key csp.Key
	pk  interface{}
}

// New returns a new csp-based crypto.Signer
// for the given csp instance and key.
func New(csp csp.CSP, key csp.Key) (crypto.Signer, error) {
	// Validate arguments
	if csp == nil {
		return nil, errors.New("csp instance must be different from nil.")
	}
	if key == nil {
		return nil, errors.New("key must be different from nil.")
	}
	if key.Symmetric() {
		return nil, errors.New("key must be asymmetric.")
	}

	// Marshall the csp public key as a crypto.PublicKey
	pub, err := key.PublicKey()
	if err != nil {
		return nil, errors.Wrap(err, "failed getting public key")
	}

	raw, err := pub.Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "failed marshalling public key")
	}

	pk, err := utils.DERToPublicKey(raw)
	if err != nil {
		return nil, errors.Wrap(err, "failed marshalling der to public key")
	}

	return &cspCryptoSigner{csp, key, pk}, nil
}

// Public returns the public key corresponding to the opaque,
// private key.
func (s *cspCryptoSigner) Public() crypto.PublicKey {
	return s.pk
}

// Sign signs digest with the private key, possibly using entropy from rand.
// For an (EC)DSA key, it should be a DER-serialised, ASN.1 signature
// structure.
//
// Hash implements the SignerOpts interface and, in most cases, one can
// simply pass in the hash function used as opts. Sign may also attempt
// to type assert opts to other types in order to obtain algorithm
// specific values. See the documentation in each package for details.
//
// Note that when a signature of a hash of a larger message is needed,
// the caller is responsible for hashing the larger message and passing
// the hash (as digest) and the hash function (as opts) to Sign.
func (s *cspCryptoSigner) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
	return s.csp.Sign(s.key, digest, opts)
}
