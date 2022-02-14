// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package otelsarama

import (
	"context"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama/internal"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

type consumerMessagesDispatcher interface {
	Messages() <-chan *sarama.ConsumerMessage
}

type consumerMessagesDispatcherWrapper struct {
	d        consumerMessagesDispatcher
	messages chan *sarama.ConsumerMessage

	cfg config
}

func newConsumerMessagesDispatcherWrapper(d consumerMessagesDispatcher, cfg config) *consumerMessagesDispatcherWrapper {
	return &consumerMessagesDispatcherWrapper{
		d:        d,
		messages: make(chan *sarama.ConsumerMessage),
		cfg:      cfg,
	}
}

// Messages returns the read channel for the messages that are returned by
// the broker.
func (w *consumerMessagesDispatcherWrapper) Messages() <-chan *sarama.ConsumerMessage {
	return w.messages
}

func (w *consumerMessagesDispatcherWrapper) Run() {
	msgs := w.d.Messages()

	for msg := range msgs {
		requestStartTime := time.Now()

		// Extract a span context from message to link.
		carrier := NewConsumerMessageCarrier(msg)
		parentSpanContext := w.cfg.Propagators.Extract(context.Background(), carrier)

		// span and metric common attrs
		attrs := []attribute.KeyValue{
			semconv.MessagingSystemKey.String("kafka"),
			semconv.PeerServiceKey.String(w.cfg.peerService),
			semconv.MessagingDestinationKindTopic,
			semconv.MessagingDestinationKey.String(msg.Topic),
			semconv.MessagingOperationReceive,
			semconv.MessageTypeReceived,
			internal.KafkaPartitionKey.Int64(int64(msg.Partition)),
		}

		attrs = AppendConsumerAttributes(w.cfg, attrs)
		attrs = AppendCommonAttributes(w.cfg, attrs)

		// get SourceCanonicalServiceKey from carrier
		sourceCanonicalService := SourceCanonicalServiceFromCarrier(carrier)
		if len(sourceCanonicalService) > 0 {
			attrs = append(attrs, internal.SourceCanonicalServiceKey.String(sourceCanonicalService))
		}

		// Create a span.
		opts := []trace.SpanStartOption{
			trace.WithAttributes(attrs...),
			// Only set on span to attrs to prevent high cardinality in metrics
			trace.WithAttributes(semconv.MessagingMessageIDKey.String(strconv.FormatInt(msg.Offset, 10))),
			trace.WithSpanKind(trace.SpanKindConsumer),
		}
		newCtx, span := w.cfg.Tracer.Start(parentSpanContext, "kafka.consume", opts...)

		// Inject current span context, so consumers can use it to propagate span.
		w.cfg.Propagators.Inject(newCtx, carrier)

		// Send messages back to user.
		w.messages <- msg

		// Use floating point division here for higher precision (instead of Millisecond method).
		elapsedTime := float64(time.Since(requestStartTime)) / float64(time.Millisecond)

		// span kind
		attrs = append(attrs, internal.SpanKindKey.String(trace.SpanKindConsumer.String()))

		// record request counter
		w.cfg.instruments.requestsCounter.Add(parentSpanContext, 1, attrs...)

		// record request latency
		w.cfg.instruments.requestLatencyHistogram.Record(parentSpanContext, elapsedTime, attrs...)

		span.End()
	}
	close(w.messages)
}
