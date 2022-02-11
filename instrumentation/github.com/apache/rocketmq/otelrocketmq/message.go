package otelrocketmq

import (
	"sort"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.opentelemetry.io/otel/propagation"
)

var _ propagation.TextMapCarrier = (*ProducerMessageCarrier)(nil)
var _ propagation.TextMapCarrier = (*ConsumerMessageCarrier)(nil)

// ProducerMessageCarrier injects and extracts traces from a primitive.Message.
type ProducerMessageCarrier struct {
	msg *primitive.Message
}

// NewProducerMessageCarrier creates a new ProducerMessageCarrier.
func NewProducerMessageCarrier(msg *primitive.Message) ProducerMessageCarrier {
	return ProducerMessageCarrier{msg}
}

// Get retrieves a single value for a given key.
func (p ProducerMessageCarrier) Get(key string) string {
	return p.msg.GetProperty(key)
}

// Set sets a property.
func (p ProducerMessageCarrier) Set(key string, value string) {
	p.msg.WithProperty(key, value)
}

// Keys returns a slice of all key identifiers in the carrier.
func (p ProducerMessageCarrier) Keys() []string {
	out := make([]string, 0, len(p.msg.GetProperties()))
	for k, _ := range p.msg.GetProperties() {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// ConsumerMessageCarrier injects and extracts traces from a primitive.Message.
type ConsumerMessageCarrier struct {
	msg *primitive.MessageExt
}

// NewConsumerMessageCarrier creates a new ConsumerMessageCarrier.
func NewConsumerMessageCarrier(msg *primitive.MessageExt) ConsumerMessageCarrier {
	return ConsumerMessageCarrier{msg}
}

// Get retrieves a single value for a given key.
func (c ConsumerMessageCarrier) Get(key string) string {
	return c.msg.GetProperty(key)
}

// Set sets a property.
func (c ConsumerMessageCarrier) Set(key string, value string) {
	c.msg.WithProperty(key, value)
}

// Keys returns a slice of all key identifiers in the carrier.
func (c ConsumerMessageCarrier) Keys() []string {
	out := make([]string, 0, len(c.msg.GetProperties()))
	for k, _ := range c.msg.GetProperties() {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
