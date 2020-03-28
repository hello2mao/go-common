package fabric_client

import (
	"encoding/hex"
	"hash"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/cryptosuite"
	"github.com/pkg/errors"
	"github.com/tjfoc/gmsm/sm3"
)

type transactor struct {
	creatorFunc func() ([]byte, error)
	hashFunc    func() (hash.Hash, error)
	nonceFunc   func() ([]byte, error)
}

func (t *transactor) createTxid() (string, []byte, error) {
	creator, err := t.creatorFunc()
	if err != nil {
		return "", nil, errors.WithMessage(err, "failed to run creatorFunc")
	}
	nonce, err := t.nonceFunc()
	if err != nil {
		return "", nil, errors.WithMessage(err, "failed to run nonceFunc")
	}
	hashFunc, err := t.hashFunc()
	if err != nil {
		return "", nil, errors.WithMessage(err, "failed to run hashFunc")
	}
	computeTxid, err := computeTxnID(nonce, creator, hashFunc)
	if err != nil {
		return "", nil, errors.WithMessage(err, "failed to computeTxnID")
	}
	return computeTxid, nonce, nil
}

func (t *transactor) getCreator() ([]byte, error) {
	return t.creatorFunc()
}

func newTransactor(channelContext context.Channel, supportGM bool) *transactor {
	return &transactor{
		creatorFunc: channelContext.Serialize,
		nonceFunc:   GetRandomNonce,
		hashFunc: func() (i hash.Hash, e error) {
			if supportGM {
				return sm3.New(), nil
			} else {
				ho := cryptosuite.GetSHA256Opts()
				return channelContext.CryptoSuite().GetHash(ho)
			}
		},
	}
}

func computeTxnID(nonce, creator []byte, h hash.Hash) (string, error) {
	h.Reset()
	b := append(nonce, creator...)
	_, err := h.Write(b)
	if err != nil {
		return "", err
	}
	digest := h.Sum(nil)
	id := hex.EncodeToString(digest)
	return id, nil
}
