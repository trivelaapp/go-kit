package log

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/trivelaapp/go-kit/errors"
)

type payloadBuilder interface {
	build(context.Context, logInput) interface{}
}

type defaultLogBuilder struct{}

func (b defaultLogBuilder) build(ctx context.Context, in logInput) interface{} {
	payload := map[string]interface{}{
		"level":     in.level.String(),
		"timestamp": in.timestamp.Format(time.RFC3339),
		"message":   in.message,
	}

	attrs := extractLogAttributesFromContext(ctx, in.attributes)
	if in.err != nil {
		attrs[LogAttributeRootError] = errors.RootError(in.err)
		attrs[LogAttributeErrorKind] = string(errors.Kind(in.err))
		attrs[LogAttributeErrorCode] = string(errors.Code(in.err))
	}
	if len(attrs) > 0 {
		payload["attributes"] = attrs
	}

	return payload
}

type defaultWithTraceLogBuilder struct{}

func (b defaultWithTraceLogBuilder) build(ctx context.Context, in logInput) interface{} {
	payload := map[string]interface{}{
		"level":     in.level.String(),
		"timestamp": in.timestamp.Format(time.RFC3339),
		"message":   in.message,
	}

	attrs := extractLogAttributesFromContext(ctx, in.attributes)
	if in.err != nil {
		attrs[LogAttributeRootError] = errors.RootError(in.err)
		attrs[LogAttributeErrorKind] = string(errors.Kind(in.err))
		attrs[LogAttributeErrorCode] = string(errors.Code(in.err))
	}
	if len(attrs) > 0 {
		payload["attributes"] = attrs
	}

	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().TraceID().IsValid() {
		return payload
	}

	span.AddEvent("log", trace.WithAttributes(buildOtelAttributes(attrs, "log")...))
	if in.err != nil {
		span.RecordError(in.err, trace.WithAttributes(buildOtelAttributes(attrs, "exception")...))
		span.SetStatus(codes.Error, in.err.Error())
	}

	payload["trace_id"] = span.SpanContext().TraceID().String()
	payload["span_id"] = span.SpanContext().SpanID().String()

	return payload
}

type gcpCloudLoggingWithTraceLogBuilder struct {
	projectID string
}

func (b gcpCloudLoggingWithTraceLogBuilder) build(ctx context.Context, in logInput) interface{} {
	payload := map[string]interface{}{
		"severity": in.level.String(),
		"time":     in.timestamp.Format(time.RFC3339),
		"message":  in.message,
	}

	attrs := extractLogAttributesFromContext(ctx, in.attributes)
	if in.err != nil {
		attrs[LogAttributeRootError] = errors.RootError(in.err)
		attrs[LogAttributeErrorKind] = string(errors.Kind(in.err))
		attrs[LogAttributeErrorCode] = string(errors.Code(in.err))

		// Necessary to link error to Cloud Error Reporting.
		// More details in: https://cloud.google.com/error-reporting/docs/formatting-error-messages#log-entry-examples
		attrs["@type"] = "type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent"
	}
	if len(attrs) > 0 {
		payload["attributes"] = attrs
	}

	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().TraceID().IsValid() {
		return payload
	}

	span.AddEvent("log", trace.WithAttributes(buildOtelAttributes(attrs, "log")...))
	if in.err != nil {
		span.RecordError(in.err, trace.WithAttributes(buildOtelAttributes(attrs, "exception")...))
		span.SetStatus(codes.Error, in.err.Error())
	}

	// Necessary to link with Cloud Trace
	// More details in: https://cloud.google.com/logging/docs/structured-logging
	payload["logging.googleapis.com/trace"] = fmt.Sprintf("projects/%s/traces/%s", b.projectID, span.SpanContext().TraceID().String())
	payload["logging.googleapis.com/spanId"] = span.SpanContext().SpanID().String()
	payload["logging.googleapis.com/trace_sampled"] = span.SpanContext().IsSampled()

	return payload
}

func extractLogAttributesFromContext(ctx context.Context, attrSet LogAttributeSet) map[LogAttribute]interface{} {
	attributes := map[LogAttribute]interface{}{}

	for attr := range attrSet {
		if value := ctx.Value(string(attr)); value != nil {
			attributes[attr] = value
		}
	}

	return attributes
}

func buildOtelAttributes(attrs map[LogAttribute]interface{}, prefix string) []attribute.KeyValue {
	eAttrs := []attribute.KeyValue{}
	for k, v := range attrs {
		eAttrs = append(eAttrs, attribute.String(fmt.Sprintf("%s.%s", prefix, k), v.(string)))
	}

	return eAttrs
}
