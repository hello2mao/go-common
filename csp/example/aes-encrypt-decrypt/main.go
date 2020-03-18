package main

import (
	"fmt"
	"os"

	"github.com/hello2mao/go-common/csp"
	"github.com/hello2mao/go-common/csp/factory"
)

func main() {

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
}
