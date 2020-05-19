package main

import (
	"fmt"
	scCache "github.com/hello2mao/go-common/second_chance_cache"
	"os"
)

func main() {
	cache := scCache.NewSecondChanceCache(2)
	cache.Add("a", 123)
	value, exist := cache.Get("a")
	if !exist {
		os.Exit(-1)
	}

	fmt.Printf("a: %v\n", value.(int))
}

