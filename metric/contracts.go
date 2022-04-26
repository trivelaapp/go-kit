package metric

import (
	"context"

	"go.opentelemetry.io/otel/metric"
)

type logger interface {
	Error(ctx context.Context, err error)
}

// MeterProvider defines how providers of Meters should behavior.
type MeterProvider interface {
	// Meter produces a new Meter.
	// Meters that works in a Push model, exposes a flush function that flushes all pending metrics into provider.
	// Usually used during application's shutdown.
	Meter(context.Context) (metric.Meter, func(context.Context) error, error)
}
