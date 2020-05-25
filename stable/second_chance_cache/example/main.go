package main

import (
	"fmt"
	scCache "github.com/hello2mao/go-common/stable/second_chance_cache"
	"os"
)

func main() {
	// init cache
	cache := scCache.NewSecondChanceCache(2)
	// add item
	cache.Add("a", 123)

	// get from cache
	value, exist := cache.Get("a")
	if !exist {
		os.Exit(-1)
	}

	fmt.Printf("a: %v\n", value.(int))
}

