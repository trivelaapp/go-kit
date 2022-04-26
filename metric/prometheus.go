package metric

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"

	"github.com/trivelaapp/go-kit/errors"
)

// PrometheusMeterProviderParams encapsulates the necessary parameters to initialize a PrometheusMeterProvider.
type PrometheusMeterProviderParams struct {
	ApplicationName     string
	ApplicationVersion  string
	MetricsServerPort   int
	HistogramBoundaries []float64
}

// PrometheusMeterProvider creates Metric meters.
type PrometheusMeterProvider struct {
	applicationName     string
	applicationVersion  string
	metricsServerPort   int
	histogramBoundaries []float64
}

// NewPrometheusMeterProvider create a new instance of a PrometheusMeterProvider.
func NewPrometheusMeterProvider(params PrometheusMeterProviderParams) (MeterProvider, error) {
	if params.ApplicationName == "" {
		return nil, errors.NewMissingRequiredDependency("ApplicationName")
	}

	if params.ApplicationVersion == "" {
		params.ApplicationVersion = "unknown_version"
	}

	if params.MetricsServerPort == 0 {
		return nil, errors.NewMissingRequiredDependency("MetricsServerPort")
	}

	if len(params.HistogramBoundaries) == 0 {
		return nil, errors.NewMissingRequiredDependency("HistogramBoundaries")
	}

	return &PrometheusMeterProvider{
		applicationName:     params.ApplicationName,
		applicationVersion:  params.ApplicationVersion,
		metricsServerPort:   params.MetricsServerPort,
		histogramBoundaries: params.HistogramBoundaries,
	}, nil
}

// MustNewPrometheusMeterProvider create a new instance of a PrometheusMeterProvider.
// It panics if any error is found.
func MustNewPrometheusMeterProvider(params PrometheusMeterProviderParams) MeterProvider {
	client, err := NewPrometheusMeterProvider(params)
	if err != nil {
		panic(err)
	}

	return client
}

// Meter produces a new Prometheus meter.
// Its flush function doesn't do anything, since it works in a Pull model.
func (c PrometheusMeterProvider) Meter(ctx context.Context) (metric.Meter, func(context.Context) error, error) {
	config := prometheus.Config{
		DefaultHistogramBoundaries: c.histogramBoundaries,
	}

	controller := controller.New(
		processor.NewFactory(
			selector.NewWithHistogramDistribution(
				histogram.WithExplicitBoundaries(config.DefaultHistogramBoundaries),
			),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
		controller.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(c.applicationName),
				semconv.ServiceVersionKey.String(c.applicationVersion),
			),
		),
	)

	exporter, err := prometheus.New(config, controller)
	if err != nil {
		return nil, nil, errors.New("failed to initialize prometheus exporter").WithRootError(err)
	}

	meterProvider := exporter.MeterProvider()
	global.SetMeterProvider(meterProvider)

	go func() {
		http.HandleFunc("/metrics", exporter.ServeHTTP)
		http.ListenAndServe(fmt.Sprintf(":%d", c.metricsServerPort), nil)
	}()

	return meterProvider.Meter(c.applicationName), func(context.Context) error { return nil }, nil
}
