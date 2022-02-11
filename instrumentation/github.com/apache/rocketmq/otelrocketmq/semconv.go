package otelrocketmq

import (
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

var (
	RocketMessagingSystem = semconv.MessagingSystemKey.String("rocketmq")
)

const (
	RocketmqQueueIDKey                = attribute.Key("messaging.rocketmq.queue_id")
	RocketmqBrokerNameKey             = attribute.Key("messaging.rocketmq.broker_name")
	RocketmqRegionIDKey               = attribute.Key("messaging.rocketmq.region_id")
	RocketmqTransactionIDKey          = attribute.Key("messaging.rocketmq.transaction_id")
	RocketmqOffsetMsgIDKey            = attribute.Key("messaging.rocketmq.offset_msg_id")
	RocketmqQueueOffsetKey            = attribute.Key("messaging.rocketmq.queue_offset")
	RocketmqMessageKeysKey            = attribute.Key("messaging.rocketmq.message_keys")
	RocketmqMessageTagKey             = attribute.Key("messaging.rocketmq.message_tag")
	RocketmqMessageTypeKey            = attribute.Key("messaging.rocketmq.message_type")
	RocketmqClientIDKey               = attribute.Key("messaging.rocketmq.client_id")
	RocketmqClientGroupKey            = attribute.Key("messaging.rocketmq.client_group")
	RocketmqNamespaceKey              = attribute.Key("messaging.rocketmq.namespace")
	RocketmqConsumptionModelKey       = attribute.Key("messaging.rocketmq.consumption_model")
	RocketmqSendResultKey             = attribute.Key("messaging.rocketmq.send_result")
	RocketmqBrokerAddressKey          = attribute.Key("messaging.rocketmq.broker_address")
	MessagingRocketmqConsumerGroupKey = attribute.Key("messaging.rocketmq.consumer_group")
)

type MessageTypeStr string

const (
	NormalMsg      MessageTypeStr = "normal"
	TransMsgHalf                  = "transaction_half"
	TransMsgCommit                = "transaction_commit"
	DelayMsg                      = "delay"
)

func ConvertMessageType(msgType primitive.MessageType) MessageTypeStr {
	switch msgType {
	case primitive.NormalMsg:
		return NormalMsg
	case primitive.TransMsgHalf:
		return TransMsgHalf
	case primitive.TransMsgCommit:
		return TransMsgCommit
	case primitive.DelayMsg:
		return DelayMsg
	default:
		return NormalMsg
	}
}
