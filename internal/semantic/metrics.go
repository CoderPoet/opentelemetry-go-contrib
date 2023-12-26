package semantic

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// MetricsAttributesFromSpanAttributes metrics attrs from span attributes
// Extract the following attributes:
// 1、http attributes
// 2、peer attributes
// 3、request attributes
// 4、metric resource attributes
// 5、span status
func MetricsAttributesFromSpanAttributes(sAttrs []attribute.KeyValue, code codes.Code) []attribute.KeyValue {
	var attrs []attribute.KeyValue

	// span attributes
	for _, attr := range sAttrs {
		if matchAttributeKey(attr.Key, HTTPMetricsAttributes) {
			attrs = append(attrs, attr)
		}
		if matchAttributeKey(attr.Key, PeerMetricsAttributes) {
			attrs = append(attrs, attr)
		}
		if matchAttributeKey(attr.Key, RequestAttributes) {
			attrs = append(attrs, attr)
		}
		if matchAttributeKey(attr.Key, DBMetricsAttributes) {
			attrs = append(attrs, attr)
		}
	}

	// span resource attributes
	for _, attr := range sAttrs {
		if matchAttributeKey(attr.Key, MetricResourceAttributes) {
			attrs = append(attrs, attr)
		}
	}

	// status code
	if code == codes.Error {
		attrs = append(attrs, StatusCodeError)
	} else {
		attrs = append(attrs, StatusCodeOk)
	}

	return attrs
}

func matchAttributeKey(key attribute.Key, toMatchKeys []attribute.Key) bool {
	for _, attrKey := range toMatchKeys {
		if attrKey == key {
			return true
		}
	}
	return false
}
