package kafka

import (
	"time"

	"github.com/hello2mao/go-common/pkg/mq/core"
)

type Queue struct {
}

func (q *Queue) Publish(payload []byte) bool {
	return true
}

func (q *Queue) StartConsuming(prefetchLimit int, pollDuration time.Duration) bool {
	return true
}

func (q *Queue) StopConsuming() <-chan struct{} {
	return nil
}

func (q *Queue) AddConsumer(tag string, consumer core.Consumer) string {
	return ""
}

func (q *Queue) AddConsumerFunc(tag string, consumerFunc core.ConsumerFunc) string {
	return ""
}

func (q *Queue) Close() bool {
	return true
}
