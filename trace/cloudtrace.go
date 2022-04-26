package trace

import (
	"context"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel/trace"

	"github.com/trivelaapp/go-kit/errors"
)

// CloudTraceClientParams encapsulates the necessary parameters to initialize a CloudTraceClient.
type CloudTraceClientParams struct {
	ApplicationName    string
	ApplicationVersion string
	ProjectID          string

	// TraceRatio indicates how often the system should collect traces.
	// Use it with caution: It may overload the system and also be too expensive to mantain its value too high in a high throwput system
	// Values vary between 0 and 1, with 0 meaning No Sampling and 1 meaning Always Sampling.
	// Values lower than 0 are treated as 0 and values greater than 1 are treated as 1.
	TraceRatio float64
}

// CloudTraceClient creates CloudTrace tracers.
type CloudTraceClient struct {
	applicationName    string
	applicationVersion string
	projectID          string
	traceRatio         float64
}

// NewCloudTraceClient create a new instance of a CloudTraceClient.
func NewCloudTraceClient(params CloudTraceClientParams) (TraceProvider, error) {
	if params.ApplicationName == "" {
		return nil, errors.NewMissingRequiredDependency("ApplicationName")
	}

	if params.ApplicationVersion == "" {
		params.ApplicationVersion = "Unknown"
	}

	if params.ProjectID == "" {
		return nil, errors.NewMissingRequiredDependency("ProjectID")
	}

	return &CloudTraceClient{
		applicationName:    params.ApplicationName,
		applicationVersion: params.ApplicationVersion,
		projectID:          params.ProjectID,
		traceRatio:         params.TraceRatio,
	}, nil
}

// MustNewCloudTraceClient create a new instance of a CloudTraceClient.
// It panics if any error is found.
func MustNewCloudTraceClient(params CloudTraceClientParams) TraceProvider {
	client, err := NewCloudTraceClient(params)
	if err != nil {
		panic(err)
	}

	return client
}

// Tracer produces a new Trace tracer and a Flush function.
// The flush function is designed to flush all pending tracer into provider. Usually used during application's shutdown.
func (c CloudTraceClient) Tracer(ctx context.Context) (trace.Tracer, func(context.Context) error, error) {
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
