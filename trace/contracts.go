package trace

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// TracerProvider defines how providers of Traces should behavior.
type TracerProvider interface {
	// Tracer produces a new Trace tracer and a Flush function.
	// The flush function is designed to flush all pending tracer into provider. Usually used during application's shutdown.
	Tracer(context.Context) (trace.Tracer, func(context.Context) error, error)
}
