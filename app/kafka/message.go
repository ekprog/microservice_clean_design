package kafka

import (
	"github.com/Shopify/sarama"
)

type Committable func()

type Message[T any] struct {
	Value   T
	Details *sarama.ConsumerMessage
}
