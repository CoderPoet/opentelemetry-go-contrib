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

package otelgrpc // import "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

import (
	"context"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/contrib/internal/semantic"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	grpc_codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc/internal"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

type gRPCContextKey struct{}

type gRPCContext struct {
	messagesReceived int64
	messagesSent     int64
	metricAttrs      []attribute.KeyValue
}

type serverHandler struct {
	*config
}

// NewServerHandler creates a stats.Handler for gRPC server.
func NewServerHandler(opts ...Option) stats.Handler {
	h := &serverHandler{
		config: newConfig(opts, "server"),
	}

	return h
}

// TagConn can attach some information to the given context.
func (h *serverHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	span := trace.SpanFromContext(ctx)
	attrs := peerAttr(peerFromCtx(ctx))
	span.SetAttributes(attrs...)
	return ctx
}

// HandleConn processes the Conn stats.
func (h *serverHandler) HandleConn(ctx context.Context, info stats.ConnStats) {
}

// TagRPC can attach some information to the given context.
func (h *serverHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	ctx = extract(ctx, h.config.Propagators)

	name, attrs := internal.ParseFullMethod(info.FullMethodName)
	attrs = append(attrs, RPCSystemGRPC)
	ctx, span := h.tracer.Start(
		trace.ContextWithRemoteSpanContext(ctx, trace.SpanContextFromContext(ctx)),
		name,
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(attrs...),
	)

	// carrier from metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	carrier := &metadataSupplier{metadata: &md}
	// source service resource attrs
	sourceServiceAttrs := semantic.ExtractPeerMetadataAttributes(carrier)
	// set source service resource attrs into span
	span.SetAttributes(sourceServiceAttrs...)
	// set into attrs
	attrs = append(attrs, sourceServiceAttrs...)

	// dest service attrs for client
	if canonicalService := semantic.CanonicalServiceFromSpan(span); canonicalService != "" {
		attrs = append(attrs,
			semantic.CanonicalServiceKey.String(canonicalService),
			semantic.DestinationCanonicalServiceKey.String(canonicalService),
		)
		span.SetAttributes(attrs...)
	}

	gctx := gRPCContext{
		metricAttrs: attrs,
	}
	return context.WithValue(ctx, gRPCContextKey{}, &gctx)
}

// HandleRPC processes the RPC stats.
func (h *serverHandler) HandleRPC(ctx context.Context, rs stats.RPCStats) {
	isServer := true
	h.handleRPC(ctx, rs, isServer)
}

type clientHandler struct {
	*config
}

// NewClientHandler creates a stats.Handler for gRPC client.
func NewClientHandler(opts ...Option) stats.Handler {
	h := &clientHandler{
		config: newConfig(opts, "client"),
	}

	return h
}

// TagRPC can attach some information to the given context.
func (h *clientHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	name, attrs := internal.ParseFullMethod(info.FullMethodName)
	attrs = append(attrs, RPCSystemGRPC)
	ctx, span := h.tracer.Start(
		ctx,
		name,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attrs...),
	)

	// carrier from metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	carrier := &metadataSupplier{metadata: &md}
	// read only span
	readOnlySpan := span.(sdktrace.ReadOnlySpan)
	// inject resource into metadata
	semantic.InjectPeerMetadata(carrier, readOnlySpan.Resource().Attributes())
	// inject metadata into ctx
	ctx = metadata.NewOutgoingContext(ctx, md)

	// dest service
	attrs = append(attrs, semantic.DestinationCanonicalServiceKey.String(h.config.PeerService))

	// source service attrs for client
	if canonicalService := semantic.CanonicalServiceFromSpan(span); canonicalService != "" {
		attrs = append(attrs,
			semantic.CanonicalServiceKey.String(canonicalService),
			semantic.SourceCanonicalServiceKey.String(canonicalService),
		)
		span.SetAttributes(attrs...)
	}

	gctx := gRPCContext{
		metricAttrs: attrs,
	}

	return inject(context.WithValue(ctx, gRPCContextKey{}, &gctx), h.config.Propagators)
}

// HandleRPC processes the RPC stats.
func (h *clientHandler) HandleRPC(ctx context.Context, rs stats.RPCStats) {
	isServer := false
	h.handleRPC(ctx, rs, isServer)
}

// TagConn can attach some information to the given context.
func (h *clientHandler) TagConn(ctx context.Context, cti *stats.ConnTagInfo) context.Context {
	span := trace.SpanFromContext(ctx)
	attrs := peerAttr(cti.RemoteAddr.String())
	span.SetAttributes(attrs...)
	return ctx
}

// HandleConn processes the Conn stats.
func (h *clientHandler) HandleConn(context.Context, stats.ConnStats) {
	// no-op
}

func (c *config) handleRPC(ctx context.Context, rs stats.RPCStats, isServer bool) { // nolint: revive  // isServer is not a control flag.
	span := trace.SpanFromContext(ctx)
	gctx, _ := ctx.Value(gRPCContextKey{}).(*gRPCContext)
	var messageId int64
	metricAttrs := make([]attribute.KeyValue, 0, len(gctx.metricAttrs)+1)
	metricAttrs = append(metricAttrs, gctx.metricAttrs...)
	wctx := withoutCancel(ctx)

	switch rs := rs.(type) {
	case *stats.Begin:
	case *stats.InPayload:
		if gctx != nil {
			messageId = atomic.AddInt64(&gctx.messagesReceived, 1)
			//c.rpcRequestSize.Record(wctx, int64(rs.Length), metric.WithAttributes(metricAttrs...))
		}

		if c.ReceivedEvent {
			span.AddEvent("message",
				trace.WithAttributes(
					semconv.MessageTypeReceived,
					semconv.MessageIDKey.Int64(messageId),
					semconv.MessageCompressedSizeKey.Int(rs.CompressedLength),
					semconv.MessageUncompressedSizeKey.Int(rs.Length),
				),
			)
		}
	case *stats.OutPayload:
		if gctx != nil {
			messageId = atomic.AddInt64(&gctx.messagesSent, 1)
			//c.rpcResponseSize.Record(wctx, int64(rs.Length), metric.WithAttributes(metricAttrs...))
		}

		if c.SentEvent {
			span.AddEvent("message",
				trace.WithAttributes(
					semconv.MessageTypeSent,
					semconv.MessageIDKey.Int64(messageId),
					semconv.MessageCompressedSizeKey.Int(rs.CompressedLength),
					semconv.MessageUncompressedSizeKey.Int(rs.Length),
				),
			)
		}
	case *stats.OutTrailer:
	case *stats.End:
		var rpcStatusAttr attribute.KeyValue

		if rs.Error != nil {
			s, _ := status.FromError(rs.Error)
			if isServer {
				statusCode, msg := serverStatus(s)
				span.SetStatus(statusCode, msg)
			} else {
				span.SetStatus(codes.Error, s.Message())
			}
			rpcStatusAttr = semconv.RPCGRPCStatusCodeKey.Int(int(s.Code()))
		} else {
			rpcStatusAttr = semconv.RPCGRPCStatusCodeKey.Int(int(grpc_codes.OK))
		}
		span.SetAttributes(rpcStatusAttr)
		span.End()

		metricAttrs = append(metricAttrs, rpcStatusAttr)

		// Record R.E.D topology metrics
		if readOnlySpan, ok := span.(sdktrace.ReadOnlySpan); ok {
			metricAttrs = append(metricAttrs,
				semantic.MetricsAttributesFromSpanAttributes(readOnlySpan.Attributes(), readOnlySpan.Status().Code)...)
		}
		if isServer {
			metricAttrs = append(metricAttrs, semantic.SpanKindKeyServer)
		} else {
			metricAttrs = append(metricAttrs, semantic.SpanKindKeyClient)
		}
		c.measure.Record(wctx, 1, time.Since(rs.BeginTime), metricAttrs...)
	default:
		return
	}
}

func withoutCancel(parent context.Context) context.Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	return withoutCancelCtx{parent}
}

type withoutCancelCtx struct {
	c context.Context
}

func (withoutCancelCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (withoutCancelCtx) Done() <-chan struct{} {
	return nil
}

func (withoutCancelCtx) Err() error {
	return nil
}

func (w withoutCancelCtx) Value(key any) any {
	return w.c.Value(key)
}

func (w withoutCancelCtx) String() string {
	return "withoutCancel"
}
