package pubsub

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
	"go.opentelemetry.io/otel/trace"

	"github.com/trivelaapp/go-kit/errors"
)

// TraceIDContextKey defines the trace id key in a context.
const TraceIDContextKey string = "trace_id"

// PubSubClientParams encapsulates the necessary params to build a PubSubClient.
type PubSubClientParams[T any] struct {
	Topic Publisher
}

// PubSubClient is a client of Google Pubsub topic with schema of type T.
type PubSubClient[T any] struct {
	topic Publisher
}

// NewPubSubClient creates a new PubSubClient instance.
func NewPubSubClient[T any](params *PubSubClientParams[T]) (*PubSubClient[T], error) {
	if params.Topic == nil {
		return nil, errors.NewMissingRequiredDependency("Topic")
	}

	return &PubSubClient[T]{params.Topic}, nil
}

// MustNewPubSubClient creates a new PubSubClient instance.
// It panics if any error is found.
func MustNewPubSubClient[T any](params *PubSubClientParams[T]) *PubSubClient[T] {
	s, err := NewPubSubClient(params)
	if err != nil {
		panic(err)
	}

	return s
}

// PublishInput is the input for publishing data into a topic with schema of type T.
type PublishInput[T any] struct {
	Data       T
	Attributes map[string]string
}

// Publish publishes messages in a pubsub topic with schema of type T.
func (c PubSubClient[T]) Publish(ctx context.Context, in ...PublishInput[T]) []error {
	var errs []error

	traceID := getTraceID(trace.SpanFromContext(ctx))

	for _, message := range in {
		data, err := json.Marshal(message.Data)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		message.Attributes[TraceIDContextKey] = traceID
		msg := &pubsub.Message{
			Data:       data,
			Attributes: message.Attributes,
		}

		if err := c.topic.Publish(ctx, msg); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func getTraceID(span trace.Span) string {
	if !span.SpanContext().HasTraceID() {
		return ""
	}

	return span.SpanContext().TraceID().String()
}
