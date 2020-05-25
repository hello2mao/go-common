package main

import (
	"context"
	"fmt"
	"time"

	"github.com/hello2mao/go-common/stable/wait"
)

func main() {
	startTime :=  time.Now().Unix()

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 12 * time.Second)
	defer cancel()
	wait.UntilWithContext(timeoutCtx, func(ctx context.Context) {
		now := time.Now().Unix()
		fmt.Printf("now: %v\n", now)
		if now - startTime > 10 {
			ctx.Done()
		}
	}, 2*time.Second)
}
