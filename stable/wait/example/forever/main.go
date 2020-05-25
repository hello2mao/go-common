package main

import (
	"fmt"
	"time"

	"github.com/hello2mao/go-common/stable/wait"
)

func main() {
	wait.Forever(func() {
		fmt.Printf("now: %s\n", time.Now())
	}, 2*time.Second)
}
