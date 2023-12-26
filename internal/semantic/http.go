package semantic

import (
	"context"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

type httpCarrier struct {
	headers http.Header
}

func (h *httpCarrier) Get(key string) string {
	return h.headers.Get(key)
}

func (h *httpCarrier) Set(key string, value string) {
	h.headers.Set(key, value)
}

func (h *httpCarrier) Keys() []string {
	return []string{
		semconvAttributeKeyToHTTPHeader(string(semconv.ServiceNameKey)),
		semconvAttributeKeyToHTTPHeader(string(semconv.ServiceNamespaceKey)),
		semconvAttributeKeyToHTTPHeader(string(semconv.DeploymentEnvironmentKey)),
	}
}

// PeerAttributesFromServerHTTPHeader extract peer attributes from server http header
// Extract the following attributes from the http header on the server side:
// 1. service.name
// 2. service.namespace
// 3. deployment.environment
func PeerAttributesFromServerHTTPHeader(headers http.Header) []attribute.KeyValue {
	return ExtractPeerMetadataAttributes(&httpCarrier{headers: headers})
}

// InjectPeerToMetadata inject peer attributes into metadata
// request metadata:
// 1、service-name
// 2、service-namespace
// 3、deployment-environment
func InjectPeerToMetadata(headers http.Header, rAttrs []attribute.KeyValue) {
	InjectPeerMetadata(&httpCarrier{headers: headers}, rAttrs)
}

// PeerAttributesFromClientHTTPHeader extract peer attributes from client http header
// First, get canonical service from resource attributes
// Then, get destination_canonical_service
// the value logic of the destination_canonical_service is as follows:
// 1. header x-destination-canonical-service
// 2. attributes peer.service or http.host
// 3. header host
func PeerAttributesFromClientHTTPHeader(_ context.Context, headers http.Header, host string, spanAttrs []attribute.KeyValue, rAttrs []attribute.KeyValue) []attribute.KeyValue {
	var attrs []attribute.KeyValue

	canonicalService := canonicalServiceFromResourceAttributes(rAttrs)

	attrs = append(
		attrs,
		CanonicalServiceKey.String(canonicalService),
		SourceCanonicalServiceKey.String(canonicalService),
	)

	// destination canonical service priority: header -> peer service -> host
	if destinationCanonicalService := headers.Get(DestinationCanonicalServiceMetadataKey); destinationCanonicalService != "" {
		// from header: x-destination-canonical-service
		attrs = append(attrs, DestinationCanonicalServiceKey.String(destinationCanonicalService))
	} else if peerService := PeerServiceFromSpanAttributes(spanAttrs); peerService != "" {
		// fallback to span attributes
		attrs = append(attrs, DestinationCanonicalServiceKey.String(peerService))
	} else {
		// fall back to req host
		attrs = append(attrs, DestinationCanonicalServiceKey.String(host))
	}

	return attrs
}

// PeerServiceFromSpanAttributes get peer service from span attributes
// 1. attributes peer.service
// 2. attributes http.host
func PeerServiceFromSpanAttributes(attrs []attribute.KeyValue) string {
	var peerService, host string
	for _, attr := range attrs {
		switch attr.Key {
		// fall back to peer.service
		case semconv.PeerServiceKey:
			peerService = attr.Value.AsString()
		case semconv.HTTPHostKey:
			host = attr.Value.AsString()
		}
	}

	if peerService != "" {
		return peerService
	}

	return host
}

func semconvAttributeKeyToHTTPHeader(key string) string {
	return strings.ReplaceAll(key, ".", "-")
}

func getServiceFromResourceAttributes(attrs []attribute.KeyValue) (serviceName, serviceNamespace, deploymentEnv string) {
	for _, attr := range attrs {
		switch attr.Key {
		case semconv.ServiceNameKey:
			serviceName = attr.Value.AsString()
		case semconv.ServiceNamespaceKey:
			serviceNamespace = attr.Value.AsString()
		case semconv.DeploymentEnvironmentKey:
			deploymentEnv = attr.Value.AsString()
		}
	}
	return
}
