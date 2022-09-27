package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub"
)

// TopicWrapper envelopes a pubsub topic type.
type TopicWrapper struct {
	topic *pubsub.Topic
}

// Publish envelopes a pubsub topic publish method.
func (tw TopicWrapper) Publish(ctx context.Context, msg *pubsub.Message) error {
	_, err := tw.topic.Publish(ctx, msg).Get(ctx)
	return err
}
