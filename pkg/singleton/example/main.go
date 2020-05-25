package main

import (
	"fmt"

	"github.com/hello2mao/go-common/pkg/singleton"
)

// ------------------------------------------

type TestConfig struct {
	Name string
}

func initConfig() (interface{}, error) {
	c := TestConfig{
		Name: "Tom",
	}
	return c, nil
}

var instance = singleton.NewSingleton(initConfig)

func GetInstance() TestConfig {
	s, err := instance.Get()
	if err != nil {
		panic("get instance err:" + err.Error())
	}
	return s.(TestConfig)
}

// ------------------------------------------

func main() {
	fmt.Printf("Name: %s\n", GetInstance().Name)
}
