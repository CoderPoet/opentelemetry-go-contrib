package otelsarama

import (
	"encoding/binary"

	"github.com/Shopify/sarama"
)

const (
	isTransactionalMask   = 0x10
	controlMask           = 0x20
	maximumRecordOverhead = 5*binary.MaxVarintLen32 + binary.MaxVarintLen64 + 1
)

const producerMessageOverhead = 26 // the metadata overhead of CRC, flags, etc.

// ByteSize get msg byte size
func ByteSize(version sarama.KafkaVersion, m *sarama.ProducerMessage) int {
	intVersion := 1
	if version.IsAtLeast(sarama.V0_11_0_0) {
		intVersion = 2
	}

	var size int
	if intVersion >= 2 {
		size = maximumRecordOverhead
		for _, h := range m.Headers {
			size += len(h.Key) + len(h.Value) + 2*binary.MaxVarintLen32
		}
	} else {
		size = producerMessageOverhead
	}
	if m.Key != nil {
		size += m.Key.Length()
	}
	if m.Value != nil {
		size += m.Value.Length()
	}
	return size
}
