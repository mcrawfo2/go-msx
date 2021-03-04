package notification

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/streamtest"
	"github.com/pkg/errors"
	"testing"
)

type MessageProducerFunc func(context.Context) (Message, error)

func (f MessageProducerFunc) Message(ctx context.Context) (Message, error) {
	return f(ctx)
}

func TestPublish(t *testing.T) {
	tests := []struct {
		name string
		test *streamtest.TopicPublishTest
	}{
		{
			name: "Success",
			test: streamtest.NewTopicPublishTest().
				WithTopic(TopicName),
		},
		{
			name: "PublishError",
			test: streamtest.NewTopicPublishTest().
				WithTopic(TopicName).
				WithPublishError(errors.New("publish error")),
		},
	}
	for _, tt := range tests {
		tt.test.WithCall(func(t *testing.T, ctx context.Context) error {
			return Publish(ctx, Message{})
		})

		t.Run(tt.name, tt.test.Test)
	}
}

func TestPublishFromProducer(t *testing.T) {
	callError := errors.New("call error")

	tests := []struct {
		name     string
		test     *streamtest.TopicPublishTest
		producer MessageProducerFunc
		wantErr  bool
	}{
		{
			name: "Success",
			test: streamtest.NewTopicPublishTest().WithTopic(TopicName),
			producer: func(ctx context.Context) (Message, error) {
				return Message{}, nil
			},
		},
		{
			name: "PublishError",
			test: streamtest.NewTopicPublishTest().WithTopic(TopicName).WithCallError(callError),
			producer: func(ctx context.Context) (Message, error) {
				return Message{}, callError
			},
		},
	}
	for _, tt := range tests {
		tt.test.WithCall(func(t *testing.T, ctx context.Context) error {
			return PublishFromProducer(ctx, tt.producer)
		})
		t.Run(tt.name, tt.test.Test)
	}
}
