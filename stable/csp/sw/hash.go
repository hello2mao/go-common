package sw

import (
	"hash"

	"github.com/hello2mao/go-common/stable/csp"
)

type hasher struct {
	hash func() hash.Hash
}

func (c *hasher) Hash(msg []byte, opts csp.HashOpts) ([]byte, error) {
	h := c.hash()
	h.Write(msg)
	return h.Sum(nil), nil
}

func (c *hasher) GetHash(opts csp.HashOpts) (hash.Hash, error) {
	return c.hash(), nil
}
