package trace

import (
	"context"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel/trace"

	"github.com/trivelaapp/go-kit/errors"
)

// CloudTraceTracerProviderParams encapsulates the necessary parameters to initialize a CloudTraceTracerProvider.
type CloudTraceTracerProviderParams struct {
	ApplicationName    string
	ApplicationVersion string
	ProjectID          string

	// TraceRatio indicates how often the system should collect traces.
	// Use it with caution: It may overload the system and also be too expensive to mantain its value too high in a high throwput system
	// Values vary between 0 and 1, with 0 meaning No Sampling and 1 meaning Always Sampling.
	// Values lower than 0 are treated as 0 and values greater than 1 are treated as 1.
	TraceRatio float64
}

// CloudTraceTracerProvider creates CloudTrace tracers.
type CloudTraceTracerProvider struct {
	applicationName    string
	applicationVersion string
	projectID          string
	traceRatio         float64
}

// NewCloudTraceTracerProvider create a new instance of a CloudTraceTracerProvider.
func NewCloudTraceTracerProvider(params CloudTraceTracerProviderParams) (TracerProvider, error) {
	if params.ApplicationName == "" {
		return nil, errors.NewMissingRequiredDependency("ApplicationName")
	}

	if params.ApplicationVersion == "" {
		params.ApplicationVersion = "Unknown"
	}

	if params.ProjectID == "" {
		return nil, errors.NewMissingRequiredDependency("ProjectID")
	}

	return &CloudTraceTracerProvider{
		applicationName:    params.ApplicationName,
		applicationVersion: params.ApplicationVersion,
		projectID:          params.ProjectID,
		traceRatio:         params.TraceRatio,
	}, nil
}

// MustNewCloudTraceTracerProvider create a new instance of a CloudTraceTracerProvider.
// It panics if any error is found.
func MustNewCloudTraceTracerProvider(params CloudTraceTracerProviderParams) TracerProvider {
	client, err := NewCloudTraceTracerProvider(params)
	if err != nil {
		panic(err)
	}

	return client
}

// Tracer produces a new CloudTrace tracer and a Flush function.
// The flush function is designed to flush all pending tracer into provider. Usually used during application's shutdown.
func (c CloudTraceTracerProvider) Tracer(ctx context.Context) (trace.Tracer, func(context.Context) error, error) {
	exporter, err := texporter.New(texporter.WithProjectID(c.projectID))
	if err != nil {
		return nil, nil, errors.New("can't initialize CloudTrace exporter").WithRootError(err)
	}

	trace, err := NewTraceClient(TraceClientParams{
		ApplicationName:    c.applicationName,
		ApplicationVersion: c.applicationVersion,
		Exporter:           exporter,
		TraceRatio:         c.traceRatio,
	})
	if err != nil {
		return nil, nil, err
	}

	tracer, flush := trace.Tracer(ctx)
	return tracer, flush, nil
}
