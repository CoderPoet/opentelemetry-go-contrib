package request

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const (
	metricNameRequestsTotal         = "aos_requests_total"
	MetricsNameRequestDuration      = "aos_request_duration_milliseconds"
	metricNameGraphRequestsTotal    = "traces_service_graph_request_total"
	MetricsNameGraphRequestDuration = "traces_service_graph_request_client_seconds"
)

// Measure for request
type Measure interface {
	// Record request increase and duration
	Record(ctx context.Context, requestIncrease int64, requestDuration time.Duration, attrs ...attribute.KeyValue)
}

type provider struct {
	config *config
	meter  metric.Meter

	counters          map[string]metric.Int64Counter
	histogramRecorder map[string]metric.Float64Histogram
}

// New create a request metrics provider
func New(opts ...Option) Measure {
	cfg := newConfig(opts)

	p := &provider{config: cfg}

	p.meter = cfg.meter

	p.counters = make(map[string]metric.Int64Counter)
	p.histogramRecorder = make(map[string]metric.Float64Histogram)

	requestsCounter, err := p.meter.Int64Counter(
		metricNameRequestsTotal,
		metric.WithDescription("This is a COUNTER incremented for every request handled"),
		metric.WithUnit("requests"),
	)
	handleErr(err)

	graphRequestsCounter, err := p.meter.Int64Counter(
		metricNameGraphRequestsTotal,
		metric.WithDescription("This is a COUNTER incremented for every request handled"),
		metric.WithUnit("requests"),
	)
	handleErr(err)

	durationRecorder, err := p.meter.Float64Histogram(
		MetricsNameRequestDuration,
		metric.WithDescription("This is a DISTRIBUTION which measures the duration of requests"),
		metric.WithUnit("ms"),
	)
	handleErr(err)

	graphDurationRecorder, err := p.meter.Float64Histogram(
		MetricsNameGraphRequestDuration,
		metric.WithDescription("This is a DISTRIBUTION which measures the duration of requests"),
		//metric.WithUnit("ms"),
	)
	handleErr(err)

	p.counters[metricNameRequestsTotal] = requestsCounter
	p.counters[metricNameGraphRequestsTotal] = graphRequestsCounter
	p.histogramRecorder[MetricsNameRequestDuration] = durationRecorder
	p.histogramRecorder[MetricsNameGraphRequestDuration] = graphDurationRecorder

	return p
}

// Record request metric
func (m *provider) Record(ctx context.Context, requestIncrease int64, requestDuration time.Duration, attrs ...attribute.KeyValue) {
	m.counters[metricNameRequestsTotal].Add(ctx, requestIncrease, metric.WithAttributes(attrs...))
	m.counters[metricNameGraphRequestsTotal].Add(ctx, requestIncrease, metric.WithAttributes(attrs...))
	m.histogramRecorder[MetricsNameRequestDuration].Record(ctx, float64(requestDuration/time.Millisecond), metric.WithAttributes(attrs...))
	m.histogramRecorder[MetricsNameGraphRequestDuration].Record(ctx, float64(requestDuration/time.Millisecond), metric.WithAttributes(attrs...))
}

func handleErr(err error) {
	if err != nil {
		otel.Handle(err)
	}
}
