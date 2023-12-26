package semantic

import (
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// CanonicalServiceFromSpan get canonical service from span
// 1. get span resource attributes
// 2. get {service.name}.{service.namespace}.{deployment.environment} from resource attributes as canonical_service
func CanonicalServiceFromSpan(span oteltrace.Span) string {
	if readonlySpan, ok := span.(sdktrace.ReadOnlySpan); ok {
		return canonicalServiceFromResourceAttributes(readonlySpan.Resource().Attributes())
	}

	return ""
}

// canonicalServiceFromResourceAttributes get canonical service from resource attrs
func canonicalServiceFromResourceAttributes(rAttrs []attribute.KeyValue) string {
	var canonicalService, serviceName, serviceNamespace, deploymentEnv string

	for _, attr := range rAttrs {
		switch attr.Key {
		case semconv.ServiceNameKey:
			serviceName = attr.Value.AsString()
		case semconv.ServiceNamespaceKey:
			serviceNamespace = attr.Value.AsString()
		case semconv.DeploymentEnvironmentKey:
			deploymentEnv = attr.Value.AsString()
		}
	}

	if serviceName != "" {
		canonicalService = serviceName
	}

	if serviceNamespace != "" {
		canonicalService = canonicalService + "." + serviceNamespace
	}

	if deploymentEnv != "" {
		canonicalService = canonicalService + "." + deploymentEnv
	}

	return canonicalService
}
