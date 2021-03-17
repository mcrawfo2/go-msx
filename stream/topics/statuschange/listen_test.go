package statuschange

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/streamtest"
	"github.com/pkg/errors"
	"testing"
)

func TestNewMessageListener(t *testing.T) {
	type args struct {
		fn      MessageHandler
		filters []MessageFilter
	}
	tests := []struct {
		name string
		args args
		test *streamtest.TopicReceiveTest
	}{
		{
			name: "NoFilter",
			args: args{
				fn: func(ctx context.Context, message Message) error {
					streamtest.TopicReceiveTestFromContext(ctx).Received()
					return nil
				},
				filters: []MessageFilter{},
			},
			test: streamtest.NewTopicReceiveTest().
				WithWantReceive(true),
		},
		{
			name: "MatchingFilter",
			args: args{
				fn: func(ctx context.Context, message Message) error {
					streamtest.TopicReceiveTestFromContext(ctx).Received()
					return nil
				},
				filters: []MessageFilter{
					FilterMessageByEntityType("entity-type"),
				},
			},
			test: streamtest.NewTopicReceiveTest().
				WithPayload([]byte(`{"entityType":"entity-type"}`)).
				WithWantReceive(true),
		},
		{
			name: "NonMatchingFilter",
			args: args{
				fn: func(ctx context.Context, message Message) error {
					streamtest.TopicReceiveTestFromContext(ctx).Received()
					return nil
				},
				filters: []MessageFilter{
					FilterMessageByEntityType("entity-type"),
					FilterMessageByStatus("some-status"),
				},
			},
			test: streamtest.NewTopicReceiveTest().
				WithPayload([]byte(`{"entityType":"entity-type","status":"some-other-status"}`)).
				WithWantReceive(false),
		},
		{
			name: "PayloadError",
			args: args{
				fn: func(ctx context.Context, message Message) error {
					streamtest.TopicReceiveTestFromContext(ctx).Received()
					return nil
				},
			},
			test: streamtest.NewTopicReceiveTest().
				WithPayload([]byte("[")).
				WithWantReceive(false).
				WithWantError(true),
		},
		{
			name: "ListenerError",
			args: args{
				fn: func(ctx context.Context, message Message) error {
					streamtest.TopicReceiveTestFromContext(ctx).Received()
					return errors.New("listener-error")
				},
			},
			test: streamtest.NewTopicReceiveTest().
				WithWantReceive(true).
				WithWantError(true),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test.
			WithTopic(TopicName).
			WithAction(NewMessageListener(tt.args.fn, tt.args.filters)).
			Test)
	}
}
