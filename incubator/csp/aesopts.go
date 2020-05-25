package csp

import "io"

// AES128KeyGenOpts contains options for AES key generation at 128 security level
type AES128KeyGenOpts struct {
	Temporary bool
}

// Algorithm returns the key generation algorithm identifier (to be used).
func (opts *AES128KeyGenOpts) Algorithm() string {
	return AES128
}

// Ephemeral returns true if the key to generate has to be ephemeral,
// false otherwise.
func (opts *AES128KeyGenOpts) Ephemeral() bool {
	return opts.Temporary
}

// AES192KeyGenOpts contains options for AES key generation at 192  security level
type AES192KeyGenOpts struct {
	Temporary bool
}

// Algorithm returns the key generation algorithm identifier (to be used).
func (opts *AES192KeyGenOpts) Algorithm() string {
	return AES192
}

// Ephemeral returns true if the key to generate has to be ephemeral,
// false otherwise.
func (opts *AES192KeyGenOpts) Ephemeral() bool {
	return opts.Temporary
}

// AES256KeyGenOpts contains options for AES key generation at 256 security level
type AES256KeyGenOpts struct {
	Temporary bool
}

// Algorithm returns the key generation algorithm identifier (to be used).
func (opts *AES256KeyGenOpts) Algorithm() string {
	return AES256
}

// Ephemeral returns true if the key to generate has to be ephemeral,
// false otherwise.
func (opts *AES256KeyGenOpts) Ephemeral() bool {
	return opts.Temporary
}

// AESCBCPKCS7ModeOpts contains options for AES encryption in CBC mode
// with PKCS7 padding.
// Notice that both IV and PRNG can be nil. In that case, the CSP implementation
// is supposed to sample the IV using a cryptographic secure PRNG.
// Notice also that either IV or PRNG can be different from nil.
type AESCBCPKCS7ModeOpts struct {
	// IV is the initialization vector to be used by the underlying cipher.
	// The length of IV must be the same as the Block's block size.
	// It is used only if different from nil.
	IV []byte
	// PRNG is an instance of a PRNG to be used by the underlying cipher.
	// It is used only if different from nil.
	PRNG io.Reader
}
