package otelrocketmq

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func NewProducerTraceInterceptor(opts ...Option) primitive.Interceptor {
	cfg := newConfig(opts...)
	return func(ctx context.Context, req, reply interface{}, next primitive.Invoker) error {
		var producerCtx *primitive.ProducerCtx
		if ctx != nil {
			producerCtx = primitive.GetProducerCtx(ctx)
		}

		msg := req.(*primitive.Message)

		span := startProducerSpan(cfg, producerCtx, msg)

		err := next(ctx, msg, reply)

		// SendOneway && SendAsync has no reply.
		// TODO async span
		if reply == nil {
			span.End()
			return err
		}

		result := reply.(*primitive.SendResult)
		if result.RegionID == "" || !result.TraceOn {
			span.End()
			return err
		}

		finishProducerSpan(span, result)

		return err
	}
}

func NewConsumerTraceInterceptor(opts ...Option) primitive.Interceptor {
	cfg := newConfig(opts...)
	return func(ctx context.Context, req, reply interface{}, next primitive.Invoker) error {
		consumerCtx, exist := primitive.GetConsumerCtx(ctx)
		if !exist || len(consumerCtx.Msgs) == 0 {
			return next(ctx, req, reply)
		}

		for _, msg := range consumerCtx.Msgs {
			if msg == nil {
				continue
			}

			span := startConsumerSpan(cfg, msg)

			span.SetAttributes(
				MessagingRocketmqConsumerGroupKey.String(consumerCtx.ConsumerGroup),
			)

			span.End()
		}

		err := next(ctx, req, reply)

		return err
	}
}
