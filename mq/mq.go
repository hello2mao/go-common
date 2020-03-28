package mq

import (
	"github.com/hello2mao/go-common/mq/core"
	"github.com/hello2mao/go-common/mq/kafka"
)

var defaultQueue core.Queue

func OpenDefaultQueue(name string) core.Queue {
	connection := kafka.Connection{}
	return connection.OpenQueue(name)
}

func GetDefaultQueue() core.Queue {
	return defaultQueue
}
