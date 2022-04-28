package format

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"

	"github.com/trivelaapp/go-kit/http/server/middleware"
)

// LogInput is the input given to a LogFormatter that is used to produce log payload.
type LogInput struct {
	Level      string
	Message    string
	Err        error
	Attributes LogAttributeSet
	Timestamp  time.Time
}

// LogAttribute represents an information to be extracted from the context and included into the log.
type LogAttribute string

const (
	// LogAttributeRootError defines the name of the RooError attribute attached into logs.
	LogAttributeRootError LogAttribute = "root_error"

	// LogAttributeErrorKind defines the name of the ErrorKind attribute attached into logs.
	LogAttributeErrorKind LogAttribute = "err_kind"

	// LogAttributeErrorCode defines the name of the ErrorCode attribute attached into logs.
	LogAttributeErrorCode LogAttribute = "err_code"
)

// LogAttributeSet is a set of LogAttributes.
type LogAttributeSet map[LogAttribute]bool

// Add creates a new LogAttributeSet with the given LogAttribute attached.
func (s LogAttributeSet) Add(attr LogAttribute) LogAttributeSet {
	s[attr] = true
	return s
}

// Merge merges the given LogAttributeSet into the existing one and returns a copy of the result.
func (s LogAttributeSet) Merge(set LogAttributeSet) LogAttributeSet {
	for attr := range set {
		s[attr] = true
	}
	return s
}

// DefaultHTTPServerAttributeSet defines some useful LogAttributes usually used in HTTP Servers context.
var DefaultHTTPServerAttributeSet LogAttributeSet = LogAttributeSet{
	LogAttribute(semconv.ServiceNameKey):               true,
	LogAttribute(semconv.ServiceVersionKey):            true,
	LogAttribute(semconv.NetPeerIPKey):                 true,
	LogAttribute(semconv.HTTPMethodKey):                true,
	LogAttribute(semconv.HTTPRouteKey):                 true,
	LogAttribute(semconv.HTTPStatusCodeKey):            true,
	LogAttribute(semconv.HTTPResponseContentLengthKey): true,
	LogAttribute(middleware.HTTPResponseLatencyKey):    true,
}

func extractLogAttributesFromContext(ctx context.Context, attrSet LogAttributeSet) map[LogAttribute]any {
	attributes := map[LogAttribute]any{}

	for attr := range attrSet {
		if value := ctx.Value(string(attr)); value != nil {
			attributes[attr] = value
		}
	}

	return attributes
}

func buildOtelAttributes(attrs map[LogAttribute]any, prefix string) []attribute.KeyValue {
	eAttrs := []attribute.KeyValue{}
	for k, v := range attrs {
		eAttrs = append(eAttrs, attribute.String(fmt.Sprintf("%s.%s", prefix, k), v.(string)))
	}

	return eAttrs
}
