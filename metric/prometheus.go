package metric

import (
	"context"
	"log"

	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"

	"github.com/trivelaapp/go-kit/errors"
)

// PrometheusClientParams encapsulates the necessary parameters to initialize a PrometheusClient.
type PrometheusClientParams struct {
	ApplicationName     string
	ApplicationVersion  string
	MetricsServerPort   int
	HistogramBoundaries []float64
}

// PrometheusClient creates Metric meters.
type PrometheusClient struct {
	applicationName     string
	applicationVersion  string
	metricsServerPort   int
	histogramBoundaries []float64
}

// NewPrometheusClient create a new instance of a PrometheusClient.
func NewPrometheusClient(params PrometheusClientParams) (*PrometheusClient, error) {
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

	return &PrometheusClient{
		applicationName:     params.ApplicationName,
		applicationVersion:  params.ApplicationVersion,
		metricsServerPort:   params.MetricsServerPort,
		histogramBoundaries: params.HistogramBoundaries,
	}, nil
}

// MustNewPrometheusClient create a new instance of a PrometheusClient.
// It panics if any error is found.
func MustNewPrometheusClient(params PrometheusClientParams) *PrometheusClient {
	client, err := NewPrometheusClient(params)
	if err != nil {
		panic(err)
	}

	return client
}

// Meter produces a new Prometheus meter.
func (c PrometheusClient) Meter(ctx context.Context) (metric.Meter, error) {
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
		log.Panicf("failed to initialize prometheus exporter %v", err)
	}

	metric, err := NewMetricClient(MetricClientParams{
		ApplicationName:      c.applicationName,
		MetricsServerPort:    c.metricsServerPort,
		MetricsServerHandler: exporter,
		MeterProvider:        exporter.MeterProvider(),
	})
	if err != nil {
		errors.New("can't initialize Prometheus exporter").WithRootError(err)
	}

	return metric.Meter(ctx), nil
}
