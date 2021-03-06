package trace

import (
	"context"

	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/trace"

	"github.com/trivelaapp/go-kit/errors"
)

// JaegerTracerProviderParams encapsulates the necessary parameters to initialize a JaegerTracerProvider.
type JaegerTracerProviderParams struct {
	ApplicationName    string
	ApplicationVersion string
	Endpoint           string

	// TraceRatio indicates how often the system should collect traces.
	// Use it with caution: It may overload the system and also be too expensive to mantain its value too high in a high throwput system
	// Values vary between 0 and 1, with 0 meaning No Sampling and 1 meaning Always Sampling.
	// Values lower than 0 are treated as 0 and values greater than 1 are treated as 1.
	TraceRatio float64
}

// JaegerTracerProvider creates Jaeger tracers.
type JaegerTracerProvider struct {
	applicationName    string
	applicationVersion string
	endpoint           string
	traceRatio         float64
}

// NewJaegerTracerProvider create a new instance of a JaegerTracerProvider.
func NewJaegerTracerProvider(params JaegerTracerProviderParams) (TracerProvider, error) {
	if params.ApplicationName == "" {
		return nil, errors.NewMissingRequiredDependency("ApplicationName")
	}

	if params.ApplicationVersion == "" {
		params.ApplicationVersion = "Unknown"
	}

	if params.Endpoint == "" {
		return nil, errors.NewMissingRequiredDependency("Endpoint")
	}

	return &JaegerTracerProvider{
		applicationName:    params.ApplicationName,
		applicationVersion: params.ApplicationVersion,
		endpoint:           params.Endpoint,
		traceRatio:         params.TraceRatio,
	}, nil
}

// MustNewJaegerTracerProvider create a new instance of a JaegerTracerProvider.
// It panics if any error is found.
func MustNewJaegerTracerProvider(params JaegerTracerProviderParams) TracerProvider {
	client, err := NewJaegerTracerProvider(params)
	if err != nil {
		panic(err)
	}

	return client
}

// Tracer produces a new Jaeger tracer and a Flush function.
// The flush function is designed to flush all pending tracer into provider. Usually used during application's shutdown.
func (c JaegerTracerProvider) Tracer(ctx context.Context) (trace.Tracer, func(context.Context) error, error) {
	jaegerExporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(c.endpoint)))
	if err != nil {
		return nil, nil, errors.New("can't initialize Jaeger exporter").WithRootError(err)
	}

	trace, err := NewTraceClient(TraceClientParams{
		ApplicationName:    c.applicationName,
		ApplicationVersion: c.applicationVersion,
		Exporter:           jaegerExporter,
		TraceRatio:         c.traceRatio,
	})
	if err != nil {
		return nil, nil, err
	}

	tracer, flush := trace.Tracer(ctx)
	return tracer, flush, nil
}
