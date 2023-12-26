package semantic

import (
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

const (
	// MetricsNameRequestDuration request latency histogram metric
	MetricsNameRequestDuration = "request_duration_milliseconds"
)

var (
	// HTTPMetricsAttributes http attributes
	HTTPMetricsAttributes = []attribute.Key{
		semconv.HTTPHostKey,
		semconv.HTTPRouteKey,
		semconv.HTTPMethodKey,
	}

	// RequestAttributes requests attributes
	RequestAttributes = []attribute.Key{
		SpanKindKey,
		StatusKey,
		RequestProtocolKey,
		ProtocolKey,
		SourceOperationKey,
		OperationKey,
		CanonicalServiceKey,
		SourceCanonicalServiceKey,
		DestinationCanonicalServiceKey,
	}

	// PeerMetricsAttributes peer attributes
	PeerMetricsAttributes = []attribute.Key{
		semconv.PeerServiceKey,
		//PeerServiceNamespaceKey,
		//PeerDeploymentEnvironmentKey,
	}

	// DBMetricsAttributes database attributes
	DBMetricsAttributes = []attribute.Key{
		semconv.DBSystemKey,
		semconv.DBStatementKey,
		semconv.DBNameKey,
		semconv.DBSQLTableKey,
	}

	// MetricResourceAttributes resource attributes
	MetricResourceAttributes = []attribute.Key{
		semconv.ServiceNameKey,
		semconv.ServiceNamespaceKey,
		semconv.DeploymentEnvironmentKey,
		semconv.ServiceInstanceIDKey,
		semconv.ServiceVersionKey,
		semconv.TelemetrySDKLanguageKey,
		semconv.TelemetrySDKVersionKey,
		semconv.ProcessPIDKey,
		semconv.HostNameKey,
		semconv.HostIDKey,
	}
)
