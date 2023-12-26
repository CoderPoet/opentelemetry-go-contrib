package request

import (
	"go.opentelemetry.io/otel/metric"
)

// Option request option
type Option interface {
	apply(cfg *config)
}

type option func(cfg *config)

func (fn option) apply(cfg *config) {
	fn(cfg)
}

type config struct {
	//MeterProvider metric.MeterProvider
	MetricsPrefix string
	meter         metric.Meter
}

func newConfig(opts []Option) *config {
	cfg := &config{
		//MeterProvider: otel.GetMeterProvider(),
		MetricsPrefix: "aos_",
	}

	for _, opt := range opts {
		opt.apply(cfg)
	}

	return cfg
}

// WithMeter configures metric meter provider
func WithMeter(provider metric.Meter) Option {
	return option(func(cfg *config) {
		cfg.meter = provider
	})
}

//// WithMeterProvider configures metric meter provider
//func WithMeterProvider(provider metric.MeterProvider) Option {
//	return option(func(cfg *config) {
//		cfg.MeterProvider = provider
//	})
//}

// WithMetricsPrefix with metrics prefix
func WithMetricsPrefix(prefix string) Option {
	return option(func(cfg *config) {
		cfg.MetricsPrefix = prefix
	})
}
