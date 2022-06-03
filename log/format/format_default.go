package format

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/trivelaapp/go-kit/errors"
)

type defaultLogFormatter struct{}

// NewDefault creates a new default LogFormatter.
func NewDefault() *defaultLogFormatter {
	return &defaultLogFormatter{}
}

// Format formats the log payload that will be rendered.
func (b defaultLogFormatter) Format(ctx context.Context, in LogInput) any {
	payload := map[string]any{
		"level":     in.Level,
		"timestamp": in.Timestamp.Format(time.RFC3339),
		"message":   in.Message,
	}

	if in.Payload != nil {
		payload["payload"] = in.Payload
	}

	attrs := extractLogAttributesFromContext(ctx, in.Attributes)
	if in.Err != nil {
		attrs[LogAttributeRootError] = errors.RootError(in.Err)
		attrs[LogAttributeErrorKind] = string(errors.Kind(in.Err))
		attrs[LogAttributeErrorCode] = string(errors.Code(in.Err))
	}
	if len(attrs) > 0 {
		payload["attributes"] = attrs
	}

	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().TraceID().IsValid() {
		return payload
	}

	span.AddEvent("log", trace.WithAttributes(buildOtelAttributes(attrs, "log")...))
	if in.Err != nil {
		span.RecordError(in.Err, trace.WithAttributes(buildOtelAttributes(attrs, "exception")...))
		span.SetStatus(codes.Error, in.Err.Error())
	}

	payload["trace_id"] = span.SpanContext().TraceID().String()
	payload["span_id"] = span.SpanContext().SpanID().String()

	return payload
}
