package metric

import (
	"context"
	"time"

	mexporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric"
	"go.opentelemetry.io/otel/metric"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"

	"github.com/trivelaapp/go-kit/errors"
)

// CloudMetricsMeterProviderParams encapsulates the necessary parameters to initialize a CloudMetricsMeterProvider.
type CloudMetricsMeterProviderParams struct {
	ApplicationName    string
	ApplicationVersion string
	ProjectID          string
	Logger             logger
}

// CloudMetricsMeterProvider creates Metric meters.
type CloudMetricsMeterProvider struct {
	applicationName    string
	applicationVersion string
	projectID          string
	logger             logger
}

// NewCloudMetricsMeterProvider create a new instance of a CloudMetricsMeterProvider.
func NewCloudMetricsMeterProvider(params CloudMetricsMeterProviderParams) (MeterProvider, error) {
	if params.ApplicationName == "" {
		return nil, errors.NewMissingRequiredDependency("ApplicationName")
	}

	if params.ApplicationVersion == "" {
		params.ApplicationVersion = "unknown_version"
	}

	if params.ProjectID == "" {
		return nil, errors.NewMissingRequiredDependency("ProjectID")
	}

	if params.Logger == nil {
		return nil, errors.NewMissingRequiredDependency("Logger")
	}

	return &CloudMetricsMeterProvider{
		applicationName:    params.ApplicationName,
		applicationVersion: params.ApplicationVersion,
		projectID:          params.ProjectID,
		logger:             params.Logger,
	}, nil
}

// MustNewCloudMetricsMeterProvider create a new instance of a CloudMetricsMeterProvider.
// It panics if any error is found.
func MustNewCloudMetricsMeterProvider(params CloudMetricsMeterProviderParams) MeterProvider {
	client, err := NewCloudMetricsMeterProvider(params)
	if err != nil {
		panic(err)
	}

	return client
}

// Meter produces a new CloudMetrics meter.
// The produced Meter pushes metrics to Cloud Metrics every 30 seconds.
func (c CloudMetricsMeterProvider) Meter(ctx context.Context) (metric.Meter, func(context.Context) error, error) {
	opts := []mexporter.Option{
		mexporter.WithProjectID(c.projectID),
		mexporter.WithOnError(func(err error) {
			c.logger.Error(ctx, errors.New("could not push metrics to Metrics Provider").WithRootError(err))
		}),
		mexporter.WithInterval(time.Second * 30),
	}

	copts := []controller.Option{
		controller.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(c.applicationName),
				semconv.ServiceVersionKey.String(c.applicationVersion),
			),
		),
	}

	exporter, err := mexporter.InstallNewPipeline(opts, copts...)
	if err != nil {
		return nil, nil, errors.New("could not install CloudTrace exporter").WithRootError(err)
	}

	return exporter.Meter(c.applicationName), exporter.Stop, nil
}
