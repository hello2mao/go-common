package core

import (
	"time"
)

type Queue interface {
	Publish(payload []byte) bool
	StartConsuming(prefetchLimit int, pollDuration time.Duration) bool
	StopConsuming() <-chan struct{}
	AddConsumer(tag string, consumer Consumer) string
	AddConsumerFunc(tag string, consumerFunc ConsumerFunc) string
	Close() bool
}
