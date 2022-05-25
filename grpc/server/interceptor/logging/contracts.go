package logging

import "context"

// Logger defines how the application logs data into the system
type Logger interface {
	// Debug logs debug data.
	Debug(ctx context.Context, msg string, args ...interface{})

	// Info logs info data.
	Info(ctx context.Context, msg string, args ...interface{})

	// Warning logs warning data.
	Warning(ctx context.Context, msg string, args ...interface{})

	// Error logs error data.
	Error(ctx context.Context, err error)

	// Critical logs critical data.
	Critical(ctx context.Context, err error)

	// Fatal logs critical data and exists current program execution.
	Fatal(ctx context.Context, err error)
}
