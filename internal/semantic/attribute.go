package semantic

import (
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	// OperationKey operation attribute key
	// Type: String
	// Required: No
	// Stability: stable
	OperationKey = attribute.Key("operation")

	// DestinationCanonicalServiceKey destination service key
	// Type: String
	// Required: Yes
	// Stability: stable
	DestinationCanonicalServiceKey = attribute.Key("destination_canonical_service")

	// SourceCanonicalServiceKey source service key
	//
	// Type: String
	// Required: Yes
	// Stability: stable
	SourceCanonicalServiceKey = attribute.Key("source_canonical_service")

	// CanonicalServiceKey service key
	// Type: String
	// Required: Yes
	// Stability: stable
	CanonicalServiceKey = attribute.Key("canonical_service")
)

const (
	// SourceOperationKey source operation
	//
	// Type: string
	// Required: Optional
	// Examples: '/operation1'
	SourceOperationKey = attribute.Key("source_operation")
)

const (
	// StatusKey status key
	//
	// Type: Enum
	// Required: No
	// Stability: stable
	// Examples: 'STATUS_CODE_OK'
	StatusKey = attribute.Key("status_code")
)

const (
	// RPCSystemKitexRecvSize recv_size
	RPCSystemKitexRecvSize = attribute.Key("kitex_recv_size")
	// RPCSystemKitexSendSize send_size
	RPCSystemKitexSendSize = attribute.Key("kitex_send_size")
)

const (
	// RequestProtocolKey protocol of the request.
	//
	// Type: string
	// Required: Always
	// Examples:
	// http: 'http'
	// rpc: 'grpc', 'java_rmi', 'wcf', 'kitex'
	// db: mysql, postgresql
	// mq: 'rabbitmq', 'activemq', 'AmazonSQS'
	RequestProtocolKey = attribute.Key("request_protocol")

	// ProtocolKey protocol of the request.
	ProtocolKey = attribute.Key("protocol")
)

var (
	// RequestProtocolKeyHTTP http protocol
	RequestProtocolKeyHTTP = RequestProtocolKey.String("http")
	// RequestProtocolKeyHertzHTTP hertz http protocol
	RequestProtocolKeyHertzHTTP = RequestProtocolKey.String("hertz-http")
	// RequestProtocolKeyKitexThrift kitex thrift protocol
	RequestProtocolKeyKitexThrift = RequestProtocolKey.String("kitex-thrift")
	// RequestProtocolKeyKitexGRPC kitex grpc protocol
	RequestProtocolKeyKitexGRPC = RequestProtocolKey.String("kitex-grpc")
	// RequestProtocolKeyNoSQL nosql protocol
	RequestProtocolKeyNoSQL = RequestProtocolKey.String("nosql")
	// RequestProtocolKeySQL sql protocol
	RequestProtocolKeySQL = RequestProtocolKey.String("sql")
)

const (
	// SpanKindKey key, Equivalent to SpanKind
	//
	// Type: Enum
	// Required: Yes
	// Stability: stable
	SpanKindKey = attribute.Key("span_kind")
)

var (
	// SpanKindKeyClient client
	SpanKindKeyClient = SpanKindKey.String(trace.SpanKindClient.String())
	// SpanKindKeyServer server
	SpanKindKeyServer = SpanKindKey.String(trace.SpanKindServer.String())
	// SpanKindKeyConsumer consumer
	SpanKindKeyConsumer = SpanKindKey.String(trace.SpanKindConsumer.String())
	// SpanKindKeyProducer producer
	SpanKindKeyProducer = SpanKindKey.String(trace.SpanKindProducer.String())
	// SpanKindKeyInternal internal
	SpanKindKeyInternal = SpanKindKey.String(trace.SpanKindInternal.String())
)

var (
	// StatusCodeOk ok status
	StatusCodeOk = StatusKey.String("STATUS_CODE_OK")
	// StatusCodeError error status
	StatusCodeError = StatusKey.String("STATUS_CODE_ERROR")
)

var (
	// RPCSystemKitex Semantic convention for kitex as the remoting system.
	RPCSystemKitex = semconv.RPCSystemKey.String("kitex")
)
