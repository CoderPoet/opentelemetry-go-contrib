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

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama/example"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

var (
	brokers = flag.String("brokers", os.Getenv("KAFKA_PEERS"), "The Kafka brokers to connect to, as a comma separated list")

	resourceAttrs = []attribute.KeyValue{
		semconv.ServiceNameKey.String("kafka-producer"),
		semconv.ServiceNamespaceKey.String("kafka"),
		semconv.DeploymentEnvironmentKey.String("dev"),
	}
)

func newOSSignalContext() (context.Context, func()) {
	// trap Ctrl+C and call cancel on the context
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

	return ctx, func() {
		signal.Stop(c)
		cancel()
	}
}

func main() {
	shutdown := example.InitProvider(resourceAttrs)
	defer shutdown()

	ctx, cancel := newOSSignalContext()
	defer cancel()

	flag.Parse()

	if *brokers == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	brokerList := strings.Split(*brokers, ",")
	log.Printf("Kafka brokers: %s", strings.Join(brokerList, ", "))

	producer := newAccessLogProducer(brokerList)

	for {
		select {
		case <-ctx.Done():
			err := producer.Close()
			if err != nil {
				log.Fatalln("Failed to close producer:", err)
			}
		default:
			produceMsg(producer)
			<-time.After(1 * time.Second)
		}
	}
}

func produceMsg(producer sarama.AsyncProducer) {
	rand.Seed(time.Now().Unix())

	// Create root span
	tr := otel.Tracer("producer")
	ctx, span := tr.Start(context.Background(), "produce message")
	defer span.End()

	// Inject tracing info into message
	msg := sarama.ProducerMessage{
		Topic: example.KafkaTopic,
		Key:   sarama.StringEncoder("random_number"),
		Value: sarama.StringEncoder(fmt.Sprintf("%d", rand.Intn(1000))),
	}
	otel.GetTextMapPropagator().Inject(ctx, otelsarama.NewProducerMessageCarrier(&msg))

	producer.Input() <- &msg
	successMsg := <-producer.Successes()
	log.Println("Successful to write message, offset:", successMsg.Offset)
}

func newAccessLogProducer(brokerList []string) sarama.AsyncProducer {
	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0
	// So we can know the partition and offset of messages.
	config.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	// Wrap instrumentation
	producer = otelsarama.WrapAsyncProducer(
		config,
		producer,
		otelsarama.WithResource(resourceAttrs...),
		otelsarama.WithAddress(brokerList),
		otelsarama.WithMessagingCloudProvider(otelsarama.MessagingCloudProvider{
			Provider:         "volcengine",
			AccountID:        "1400000035",
			Region:           "cn-guilin-boe",
			AvailabilityZone: "cn-guilin-boe-a",
			Platform:         "volcengine_kafka",
			InstanceID:       "mock-instance-id",
		}),
	)

	// We will log to STDOUT if we're not able to produce messages.
	go func() {
		for err := range producer.Errors() {
			log.Println("Failed to write message:", err)
		}
	}()

	return producer
}
