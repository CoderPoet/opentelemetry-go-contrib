package semantic

import (
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

// PeerMetadataCarrier peer metadata carrier, such as http.header, kitex metainfo, grpc metadata
type PeerMetadataCarrier interface {
	// Get returns the value associated with the passed key.
	Get(key string) string
	// Set stores the key-value pair.
	Set(key string, value string)
	// Keys lists the keys stored in this carrier.
	Keys() []string
}

// ExtractPeerMetadataAttributes extract peer metadata from carrier
// request metadata:
// 1、service-name
// 2、service-namespace
// 3、deployment-environment
// -> source_canonical_service
func ExtractPeerMetadataAttributes(carrier PeerMetadataCarrier) []attribute.KeyValue {
	var (
		attrs            []attribute.KeyValue
		canonicalService string
	)

	serviceName, serviceNamespace, deploymentEnv := carrier.Get(semconvAttributeKeyToHTTPHeader(string(semconv.ServiceNameKey))),
		carrier.Get(semconvAttributeKeyToHTTPHeader(string(semconv.ServiceNamespaceKey))),
		carrier.Get(semconvAttributeKeyToHTTPHeader(string(semconv.DeploymentEnvironmentKey)))

	if serviceName != "" {
		canonicalService = serviceName
		//attrs = append(attrs, semconv.PeerServiceKey.String(serviceName))
	}

	if serviceNamespace != "" {
		canonicalService = canonicalService + "." + serviceNamespace
		//attrs = append(attrs, PeerServiceNamespaceKey.String(serviceNamespace))
	}

	if deploymentEnv != "" {
		canonicalService = canonicalService + "." + deploymentEnv
		//attrs = append(attrs, PeerDeploymentEnvironmentKey.String(deploymentEnv))
	}

	if canonicalService != "" {
		attrs = append(attrs, SourceCanonicalServiceKey.String(canonicalService))
	}

	return attrs
}

// InjectPeerMetadata inject peer metadata
// request metadata:
// 1、service-name
// 2、service-namespace
// 3、deployment-environment
func InjectPeerMetadata(carrier PeerMetadataCarrier, rAttrs []attribute.KeyValue) {
	serviceName, serviceNamespace, deploymentEnv := getServiceFromResourceAttributes(rAttrs)
	if serviceName != "" {
		carrier.Set(semconvAttributeKeyToHTTPHeader(string(semconv.ServiceNameKey)), serviceName)
	}

	if serviceNamespace != "" {
		carrier.Set(semconvAttributeKeyToHTTPHeader(string(semconv.ServiceNamespaceKey)), serviceNamespace)
	}

	if deploymentEnv != "" {
		carrier.Set(semconvAttributeKeyToHTTPHeader(string(semconv.DeploymentEnvironmentKey)), deploymentEnv)
	}
}
