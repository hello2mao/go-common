package main

import (
	"fmt"
	"os"

	"github.com/hello2mao/go-common/pkg/csp"
	"github.com/hello2mao/go-common/pkg/csp/factory"
)

func main() {

	// hash
	hash, err := factory.GetDefault().Hash([]byte("data"), &csp.SHA256Opts{})
	if err != nil {
		fmt.Printf("Hash err: %s\n", err)
		os.Exit(-1)
	}
	fmt.Printf("hash: %x\n", hash)
}
