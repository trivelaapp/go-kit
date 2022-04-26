package trace

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// TraceProvider defines how providers of Traces should behavior.
type TraceProvider interface {
	// Tracer produces a new Trace tracer and a Flush function.
	// The flush function is designed to flush all pending tracer into provider. Usually used during application's shutdown.
	Tracer(context.Context) (trace.Tracer, func(context.Context) error, error)
}
