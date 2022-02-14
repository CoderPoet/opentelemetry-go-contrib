package otelsarama

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/unit"
)

const (
	// MessagingCounterMetricName counter metric
	MessagingCounterMetricName = "messaging.io.counter"
	// MessagingDurationHistogramMetricName duration histogram metric
	MessagingDurationHistogramMetricName = "messaging.io.duration"
	// MessagingSizeBytesHistogramMetricName size histogram metric
	MessagingSizeBytesHistogramMetricName = "messaging.io.bytes"
)

type instruments struct {
	// requestsCounter is the number of queries executed.
	requestsCounter metric.Int64Counter

	// requestLatencyHistogram is the sum of attempt latencies.
	requestLatencyHistogram metric.Float64Histogram
}

func newInstruments(p metric.MeterProvider) *instruments {
	meter := p.Meter(defaultInstrumentName, metric.WithInstrumentationVersion(SemVersion()))
	instruments := &instruments{}
	var err error

	instruments.requestsCounter, err = meter.NewInt64Counter(
		MessagingCounterMetricName,
		metric.WithDescription("Cumulative number of requests"),
	)
	handleErr(err)

	instruments.requestLatencyHistogram, err = meter.NewFloat64Histogram(
		MessagingDurationHistogramMetricName,
		metric.WithUnit(unit.Milliseconds),
		metric.WithDescription("Time-consuming distribution of statistical requests"),
	)
	handleErr(err)

	return instruments
}

func handleErr(err error) {
	if err != nil {
		otel.Handle(err)
	}
}
