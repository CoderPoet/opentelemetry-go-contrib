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

package internal

import "go.opentelemetry.io/otel/attribute"

const (
	MessagingCloudProviderKey         = attribute.Key("messaging.cloud.provider")
	MessagingCloudAccountIDKey        = attribute.Key("messaging.cloud.account.id")
	MessagingCloudRegionKey           = attribute.Key("messaging.cloud.region")
	MessagingCloudAvailabilityZoneKey = attribute.Key("messaging.cloud.availability_zone")
	MessagingCloudPlatformKey         = attribute.Key("messaging.cloud.platform")
	MessagingCloudInstanceIDKey       = attribute.Key("messaging.cloud.instance_id")
)

const (
	KafkaPartitionKey     = attribute.Key("messaging.kafka.partition")
	MessagingPartitionKey = attribute.Key("messaging.partition")
)

const (
	MessagingConsumerGroup = attribute.Key("messaging.consumer_group")
)

const (
	SpanKindKey               = attribute.Key("span.kind")
	StatusCodeKey             = attribute.Key("status.code")
	SourceCanonicalServiceKey = attribute.Key("source_canonical_service")
	SourceKafkaPartitionKey   = attribute.Key("messaging.source.kafka.partition")
)
