package trace

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/trivelaapp/go-kit/errors"
)

// TraceClientParams encapsulates the necessary parameters to initialize a TraceClient.
type TraceClientParams struct {
	ApplicationName    string
	ApplicationVersion string
	Exporter           sdktrace.SpanExporter

	// TraceRatio indicates how often the system should collect traces.
	// Use it with caution: It may overload the system and also be too expensive to mantain its value too high in a high throwput system
	// Values vary between 0 and 1, with 0 meaning No Sampling and 1 meaning Always Sampling.
	// Values lower than 0 are treated as 0 and values greater than 1 are treated as 1.
	TraceRatio float64
}

// TraceClient creates Trace tracers.
type TraceClient struct {
	applicationName    string
	applicationVersion string
	exporter           sdktrace.SpanExporter
	traceRatio         float64
}

// NewTraceClient create a new instance of a TraceClient.
func NewTraceClient(params TraceClientParams) (*TraceClient, error) {
	if params.ApplicationName == "" {
		return nil, errors.NewMissingRequiredDependency("ApplicationName")
	}

	if params.ApplicationVersion == "" {
		params.ApplicationVersion = "unknown_version"
	}

	if params.Exporter == nil {
		return nil, errors.NewMissingRequiredDependency("Exporter")
	}

	return &TraceClient{
		applicationName:    params.ApplicationName,
		applicationVersion: params.ApplicationVersion,
		exporter:           params.Exporter,
		traceRatio:         params.TraceRatio,
	}, nil
}

// MustNewTraceClient create a new instance of a TraceClient.
// It panics if any error is found.
func MustNewTraceClient(params TraceClientParams) *TraceClient {
	client, err := NewTraceClient(params)
	if err != nil {
		panic(err)
	}

	return client
}

// Tracer produces a new Trace tracer and a Flush function.
// The flush function is designed to flush all pending tracer into provider. Usually used during application's shutdown.
func (c TraceClient) Tracer(ctx context.Context) (trace.Tracer, func(context.Context) error) {
	tOpts := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(c.applicationName),
			semconv.ServiceVersionKey.String(c.applicationVersion),
		)),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(c.traceRatio)),
		sdktrace.WithBatcher(c.exporter),
	}
	tp := sdktrace.NewTracerProvider(tOpts...)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp.Tracer(c.applicationName), tp.Shutdown
}
