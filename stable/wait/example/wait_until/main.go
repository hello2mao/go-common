package main

import (
	"fmt"
	"time"

	"github.com/hello2mao/go-common/stable/wait"
)

func main() {
	startTime :=  time.Now().Unix()

	// loops until stop channel is closed, running f every period.
	interrupt := make(chan struct{})
	wait.Until(func() {
		now := time.Now().Unix()
		fmt.Printf("now: %v\n", now)
		if now - startTime > 10 {
			close(interrupt)
		}
	}, 2*time.Second, interrupt)
}
