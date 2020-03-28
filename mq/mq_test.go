package mq

import (
	"testing"

	"github.com/hello2mao/go-common/mq/core"
)

func TestOpenDefaultQueue(t *testing.T) {
	queue := OpenDefaultQueue("")
	queue.Publish([]byte(""))

	queue.AddConsumer("", NewTestConsumer())
}

type testConsumer struct {
}

func NewTestConsumer() *testConsumer {
	return &testConsumer{}
}

func (t *testConsumer) Consume(delivery core.Delivery) {
	delivery.Ack()
}
