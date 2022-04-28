package format

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
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

// LogAttributeSet is a set of LogAttributes.
type LogAttributeSet map[LogAttribute]bool

const (
	// LogAttributeRootError defines the name of the RooError attribute attached into logs.
	LogAttributeRootError LogAttribute = "root_error"

	// LogAttributeErrorKind defines the name of the ErrorKind attribute attached into logs.
	LogAttributeErrorKind LogAttribute = "err_kind"

	// LogAttributeErrorCode defines the name of the ErrorCode attribute attached into logs.
	LogAttributeErrorCode LogAttribute = "err_code"
)

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
