package main

import (
	"fmt"
	"os"

	"github.com/hello2mao/go-common/csp"
	"github.com/hello2mao/go-common/csp/factory"
)

func main() {

	// hash
	hash, err := factory.GetDefault().Hash([]byte("data"), &csp.SHA256Opts{})
	if err != nil {
		fmt.Printf("Hash err: %s\n", err)
		os.Exit(-1)
	}
	fmt.Printf("hash: %x\n", hash)

	// aes
	key, err := factory.GetDefault().KeyGen(&csp.AES256KeyGenOpts{Temporary: true})
	if err != nil {
		fmt.Printf("KeyGen err: %s\n", err)
		os.Exit(-1)
	}
	ciphertext, err := factory.GetDefault().Encrypt(key, []byte("data"), &csp.AESCBCPKCS7ModeOpts{})
	if err != nil {
		fmt.Printf("Encrypt err: %s\n", err)
		os.Exit(-1)
	}
	fmt.Printf("ciphertext: %x\n", ciphertext)
	plaintext, err := factory.GetDefault().Decrypt(key, ciphertext, &csp.AESCBCPKCS7ModeOpts{})
	if err != nil {
		fmt.Printf("Decrypt err: %s\n", err)
		os.Exit(-1)
	}
	fmt.Printf("plaintext: %s\n", plaintext)

	// ecdsa
	key, err = factory.GetDefault().KeyGen(&csp.ECDSAKeyGenOpts{Temporary: true})
	if err != nil {
		fmt.Printf("KeyGen err: %s\n", err)
		os.Exit(-1)
	}
	signature, err := factory.GetDefault().Sign(key, hash, nil)
	if err != nil {
		fmt.Printf("Sign err: %s\n", err)
		os.Exit(-1)
	}
	fmt.Printf("signature: %x\n", signature)
	valid, err := factory.GetDefault().Verify(key, signature, hash, nil)
	if err != nil {
		fmt.Printf("Verify err: %s\n", err)
		os.Exit(-1)
	}
	fmt.Printf("valid: %v\n", valid)
}
