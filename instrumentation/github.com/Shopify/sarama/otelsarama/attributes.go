package otelsarama

import (
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama/internal"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

// AppendCommonAttributes append common attributes
func AppendCommonAttributes(cfg config, attrs []attribute.KeyValue) []attribute.KeyValue {
	if len(cfg.messagingCloudProvider.InstanceID) > 0 {
		attrs = append(attrs, internal.MessagingCloudInstanceIDKey.String(cfg.messagingCloudProvider.InstanceID))
	}
	if len(cfg.messagingCloudProvider.Provider) > 0 {
		attrs = append(attrs, internal.MessagingCloudProviderKey.String(cfg.messagingCloudProvider.Provider))
	}
	if len(cfg.messagingCloudProvider.Platform) > 0 {
		attrs = append(attrs, internal.MessagingCloudPlatformKey.String(cfg.messagingCloudProvider.Platform))
	}
	if len(cfg.messagingCloudProvider.Region) > 0 {
		attrs = append(attrs, internal.MessagingCloudRegionKey.String(cfg.messagingCloudProvider.Region))
	}
	if len(cfg.messagingCloudProvider.AvailabilityZone) > 0 {
		attrs = append(attrs, internal.MessagingCloudAvailabilityZoneKey.String(cfg.messagingCloudProvider.AvailabilityZone))
	}
	if len(cfg.addresses) > 0 {
		attrs = append(attrs, semconv.MessagingURLKey.String(strings.Join(cfg.addresses, ",")))
	}

	return attrs
}

// AppendConsumerAttributes append consumer attributes
func AppendConsumerAttributes(cfg config, attrs []attribute.KeyValue) []attribute.KeyValue {
	if len(cfg.consumerGroupID) > 0 {
		attrs = append(attrs,
			semconv.MessagingKafkaConsumerGroupKey.String(cfg.consumerGroupID),
			internal.MessagingConsumerGroup.String(cfg.consumerGroupID),
		)

		if len(cfg.consumerClientID) > 0 {
			attrs = append(attrs, semconv.MessagingConsumerIDKey.String(cfg.consumerGroupID+"-"+cfg.consumerClientID))
		}
	}

	if len(cfg.consumerClientID) > 0 {
		attrs = append(attrs, semconv.MessagingKafkaClientIDKey.String(cfg.consumerClientID))
	}

	return attrs
}
