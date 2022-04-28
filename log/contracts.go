package log

import (
	"context"

	"github.com/trivelaapp/go-kit/log/format"
)

// LogFormatter defines how structures that formats logs should behavior.
type LogFormatter interface {
	// Format formats the log payload that will be rendered.
	Format(context.Context, format.LogInput) any
}
