package otelrocketmq

import (
	"context"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

func startProducerSpan(cfg config, pCtx *primitive.ProducerCtx, msg *primitive.Message) trace.Span {
	carrier := NewProducerMessageCarrier(msg)
	ctx := cfg.Propagators.Extract(context.Background(), carrier)

	// Create a span.
	attrs := []attribute.KeyValue{
		RocketMessagingSystem,
		semconv.MessagingDestinationKindTopic,
		semconv.MessagingDestinationKey.String(msg.Topic),
		RocketmqMessageKeysKey.String(msg.GetKeys()),
		RocketmqMessageTagKey.String(msg.GetTags()),
		RocketmqBrokerAddressKey.String(pCtx.BrokerAddr),
		RocketmqMessageTypeKey.String(string(ConvertMessageType(pCtx.MsgType))),
	}

	opts := []trace.SpanStartOption{
		trace.WithAttributes(attrs...),
		trace.WithSpanKind(trace.SpanKindProducer),
	}

	ctx, span := cfg.Tracer.Start(ctx, "rocketmq.produce", opts...)

	cfg.Propagators.Inject(ctx, carrier)

	return span
}

func finishProducerSpan(span trace.Span, result *primitive.SendResult) {
	span.SetAttributes(
		semconv.MessagingMessageIDKey.String(result.MsgID),
		RocketmqQueueIDKey.Int(result.MessageQueue.QueueId),
		RocketmqRegionIDKey.String(result.RegionID),
		RocketmqOffsetMsgIDKey.String(result.OffsetMsgID),
		RocketmqTransactionIDKey.String(result.TransactionID),
		RocketmqQueueOffsetKey.Int64(result.QueueOffset),
		RocketmqBrokerNameKey.String(result.MessageQueue.BrokerName),
		RocketmqSendResultKey.String(result.String()),
	)

	if result.Status != primitive.SendOK {
		err := fmt.Errorf("failed to produce message: %s", result.String())
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}

	span.End()
}
