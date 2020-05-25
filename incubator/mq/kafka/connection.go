package kafka

import (
	"github.com/hello2mao/go-common/incubator/mq/core"
)

type Connection struct {
}

func (connection *Connection) OpenQueue(name string) core.Queue {
	return &Queue{}
}
