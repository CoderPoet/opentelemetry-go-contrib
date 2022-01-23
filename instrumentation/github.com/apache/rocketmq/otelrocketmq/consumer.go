package otelrocketmq

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

func startConsumerSpan(cfg config, msg *primitive.MessageExt) trace.Span {
	// Extract a span context from message to link.
	carrier := NewConsumerMessageCarrier(msg)
	parentSpanContext := cfg.Propagators.Extract(context.Background(), carrier)

	// Create a span.
	attrs := []attribute.KeyValue{
		RocketMessagingSystem,
		semconv.MessagingDestinationKindTopic,
		semconv.MessagingDestinationKey.String(msg.Topic),
		semconv.MessagingOperationReceive,
		semconv.MessagingMessageIDKey.String(msg.MsgId),
		RocketmqQueueIDKey.Int(msg.Queue.QueueId),
		RocketmqMessageKeysKey.String(msg.GetKeys()),
		RocketmqMessageTagKey.String(msg.GetTags()),
	}

	if msg.Compress {
		attrs = append(attrs, semconv.MessageCompressedSizeKey.Int64(int64(msg.StoreSize)))
	} else {
		attrs = append(attrs, semconv.MessagingMessagePayloadSizeBytesKey.Int64(int64(msg.StoreSize)))
	}

	opts := []trace.SpanStartOption{
		trace.WithAttributes(attrs...),
		trace.WithSpanKind(trace.SpanKindConsumer),
	}
	newCtx, span := cfg.Tracer.Start(parentSpanContext, "rocketmq.consume", opts...)

	// Inject current span context, so consumers can use it to propagate span.
	cfg.Propagators.Inject(newCtx, carrier)

	return span
}
