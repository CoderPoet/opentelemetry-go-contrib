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

package config // import "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc/example/config"

import (
	"context"
	"os"
	"sync"

	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/otel"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

var resource *sdkresource.Resource
var initResourcesOnce sync.Once

// Init configures an OpenTelemetry exporter and trace provider.
func Init(serviceName string) (*sdktrace.TracerProvider, *sdkmetric.MeterProvider, error) {
	tp, err := initTracerProvider(serviceName)
	if err != nil {
		return nil, nil, err
	}
	mp := initMeterProvider(serviceName)
	return tp, mp, nil
}

func initResource(serviceName string) *sdkresource.Resource {
	initResourcesOnce.Do(func() {
		extraResources, _ := sdkresource.New(
			context.Background(),
			sdkresource.WithOS(),
			sdkresource.WithProcess(),
			sdkresource.WithContainer(),
			sdkresource.WithHost(),
			sdkresource.WithAttributes(semconv.ServiceName(serviceName)),
		)
		resource, _ = sdkresource.Merge(
			sdkresource.Default(),
			extraResources,
		)
	})
	return resource
}

func initTracerProvider(serviceName string) (*sdktrace.TracerProvider, error) {
	exporter, err := stdout.New(stdout.WithPrettyPrint())
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(initResource(serviceName)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp, nil
}

func initMeterProvider(serviceName string) *sdkmetric.MeterProvider {
	ctx := context.Background()

	os.Setenv("OTEL_METRICS_EXPORTER", "prometheus")

	r, err := autoexport.NewMetricReader(ctx)
	if err != nil {
		return nil
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(r),
		sdkmetric.WithResource(initResource(serviceName)),
		sdkmetric.WithView(sdkmetric.NewView(
			sdkmetric.Instrument{Scope: instrumentation.Scope{Name: "go.opentelemetry.io/contrib/google.golang.org/grpc/otelgrpc"}},
			sdkmetric.Stream{Aggregation: sdkmetric.AggregationDrop{}},
		)),
	)
	otel.SetMeterProvider(mp)
	return mp
}
