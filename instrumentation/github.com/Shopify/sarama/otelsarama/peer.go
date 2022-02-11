package otelsarama

import (
	"fmt"

	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama/internal"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

const (
	defaultServiceName    = "unknown"
	defaultAttributeValue = "default"
)

// SourceCanonicalServiceFromCarrier get source canonical service frm carrier
func SourceCanonicalServiceFromCarrier(carrier propagation.TextMapCarrier) string {
	return carrier.Get(string(internal.SourceCanonicalServiceKey))
}

// InjectSourceCanonicalService inject source canonical service
func InjectSourceCanonicalService(carrier propagation.TextMapCarrier, sourceCanonicalService string) {
	carrier.Set(string(internal.SourceCanonicalServiceKey), sourceCanonicalService)
}

// SourceCanonicalServiceFromResource source canonical service
func SourceCanonicalServiceFromResource(resourceAttrs []attribute.KeyValue) string {
	serviceName, serviceNamespace, deploymentEnvironment := defaultServiceName, defaultAttributeValue, defaultAttributeValue
	for _, resourceAttribute := range resourceAttrs {
		switch resourceAttribute.Key {
		case semconv.ServiceNameKey:
			serviceName = resourceAttribute.Value.AsString()
		case semconv.ServiceNamespaceKey:
			serviceNamespace = resourceAttribute.Value.AsString()
		case semconv.DeploymentEnvironmentKey:
			deploymentEnvironment = resourceAttribute.Value.AsString()
		}
	}

	return fmt.Sprintf("%s.%s.%s", serviceName, serviceNamespace, deploymentEnvironment)
}
