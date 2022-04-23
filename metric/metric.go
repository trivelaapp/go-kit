package metric

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"

	"github.com/trivelaapp/go-kit/errors"
)

// MetricClientParams encapsulates the necessary parameters to initialize a MetricClient.
type MetricClientParams struct {
	ApplicationName   string
	MeterProvider     metric.MeterProvider
	MetricsServerPort int

	// Handler that will serve metrics.
	MetricsServerHandler http.Handler
}

// MetricClient creates Metric meters.
type MetricClient struct {
	applicationName      string
	meterProvider        metric.MeterProvider
	metricsServerPort    int
	metricsServerHandler http.Handler
}

// NewMetricClient create a new instance of a MetricClient.
func NewMetricClient(params MetricClientParams) (*MetricClient, error) {
	if params.ApplicationName == "" {
		return nil, errors.NewMissingRequiredDependency("ApplicationName")
	}

	if params.MeterProvider == nil {
		return nil, errors.NewMissingRequiredDependency("MeterProvider")
	}

	if params.MetricsServerPort == 0 {
		return nil, errors.NewMissingRequiredDependency("MetricsServerPort")
	}

	if params.MetricsServerHandler == nil {
		return nil, errors.NewMissingRequiredDependency("MetricsServerHandler")
	}

	return &MetricClient{
		applicationName:      params.ApplicationName,
		meterProvider:        params.MeterProvider,
		metricsServerPort:    params.MetricsServerPort,
		metricsServerHandler: params.MetricsServerHandler,
	}, nil
}

// MustNewMetricClient create a new instance of a MetricClient.
// It panics if any error is found.
func MustNewMetricClient(params MetricClientParams) *MetricClient {
	client, err := NewMetricClient(params)
	if err != nil {
		panic(err)
	}

	return client
}

// Meter produces a new Metric meter.
func (c MetricClient) Meter(ctx context.Context) metric.Meter {
	global.SetMeterProvider(c.meterProvider)

	go func() {
		http.HandleFunc("/metrics", c.metricsServerHandler.ServeHTTP)
		http.ListenAndServe(fmt.Sprintf(":%d", c.metricsServerPort), c.metricsServerHandler)
	}()

	return c.meterProvider.Meter(c.applicationName)
}
