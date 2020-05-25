package signer

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"testing"

	"github.com/hello2mao/go-common/pkg/csp/mocks"
	"github.com/hello2mao/go-common/pkg/csp/utils"
	"github.com/stretchr/testify/assert"
)

func TestInitFailures(t *testing.T) {
	_, err := New(nil, &mocks.MockKey{})
	assert.Error(t, err)

	_, err = New(&mocks.MockCSP{}, nil)
	assert.Error(t, err)

	_, err = New(&mocks.MockCSP{}, &mocks.MockKey{Symm: true})
	assert.Error(t, err)

	_, err = New(&mocks.MockCSP{}, &mocks.MockKey{PKErr: errors.New("No PK")})
	assert.Error(t, err)
	assert.Equal(t, "failed getting public key: No PK", err.Error())

	_, err = New(&mocks.MockCSP{}, &mocks.MockKey{PK: &mocks.MockKey{BytesErr: errors.New("No bytes")}})
	assert.Error(t, err)
	assert.Equal(t, "failed marshalling public key: No bytes", err.Error())

	_, err = New(&mocks.MockCSP{}, &mocks.MockKey{PK: &mocks.MockKey{BytesValue: []byte{0, 1, 2, 3}}})
	assert.Error(t, err)
}

func TestInit(t *testing.T) {
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)
	pkRaw, err := utils.PublicKeyToDER(&k.PublicKey)
	assert.NoError(t, err)

	signer, err := New(&mocks.MockCSP{}, &mocks.MockKey{PK: &mocks.MockKey{BytesValue: pkRaw}})
	assert.NoError(t, err)
	assert.NotNil(t, signer)

	// Test public key
	R, S, err := ecdsa.Sign(rand.Reader, k, []byte{0, 1, 2, 3})
	assert.NoError(t, err)

	assert.True(t, ecdsa.Verify(signer.Public().(*ecdsa.PublicKey), []byte{0, 1, 2, 3}, R, S))
}

func TestPublic(t *testing.T) {
	pk := &mocks.MockKey{}
	signer := &cspCryptoSigner{pk: pk}

	pk2 := signer.Public()
	assert.NotNil(t, pk, pk2)
}

func TestSign(t *testing.T) {
	expectedSig := []byte{0, 1, 2, 3, 4}
	expectedKey := &mocks.MockKey{}
	expectedDigest := []byte{0, 1, 2, 3, 4, 5}
	expectedOpts := &mocks.SignerOpts{}

	signer := &cspCryptoSigner{
		key: expectedKey,
		csp: &mocks.MockCSP{
			SignArgKey: expectedKey, SignDigestArg: expectedDigest, SignOptsArg: expectedOpts,
			SignValue: expectedSig}}
	signature, err := signer.Sign(nil, expectedDigest, expectedOpts)
	assert.NoError(t, err)
	assert.Equal(t, expectedSig, signature)

	signer = &cspCryptoSigner{
		key: expectedKey,
		csp: &mocks.MockCSP{
			SignArgKey: expectedKey, SignDigestArg: expectedDigest, SignOptsArg: expectedOpts,
			SignErr: errors.New("no signature")}}
	_, err = signer.Sign(nil, expectedDigest, expectedOpts)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "no signature")

	signer = &cspCryptoSigner{
		key: nil,
		csp: &mocks.MockCSP{SignArgKey: expectedKey, SignDigestArg: expectedDigest, SignOptsArg: expectedOpts}}
	_, err = signer.Sign(nil, expectedDigest, expectedOpts)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "invalid key")

	signer = &cspCryptoSigner{
		key: expectedKey,
		csp: &mocks.MockCSP{SignArgKey: expectedKey, SignDigestArg: expectedDigest, SignOptsArg: expectedOpts}}
	_, err = signer.Sign(nil, nil, expectedOpts)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "invalid digest")

	signer = &cspCryptoSigner{
		key: expectedKey,
		csp: &mocks.MockCSP{SignArgKey: expectedKey, SignDigestArg: expectedDigest, SignOptsArg: expectedOpts}}
	_, err = signer.Sign(nil, expectedDigest, nil)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "invalid opts")
}
