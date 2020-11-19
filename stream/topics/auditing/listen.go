package auditing

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
)

type MessageHandler func(ctx context.Context, message Message) error

func NewMessageListener(fn MessageHandler, filters []MessageFilter) stream.ListenerAction {
	return func(msg *message.Message) error {
		var m Message
		err := json.Unmarshal(msg.Payload, &m)
		if err != nil {
			return errors.Wrap(err, "Failed to unmarshal message payload to Message")
		}

		if !FilterMessage(msg.Context(), m, filters) {
			return nil
		}

		return fn(msg.Context(), m)
	}
}

type MessageFilter func(ctx context.Context, message Message) bool

func FilterMessage(ctx context.Context, msg Message, filters []MessageFilter) bool {
	for _, filter := range filters {
		if !filter(ctx, msg) {
			return false
		}
	}
	return true
}

func FilterByService(serviceType string) MessageFilter {
	return func(ctx context.Context, message Message) bool {
		return message.Service == serviceType
	}
}

func FilterByAction(action string) MessageFilter {
	return func(ctx context.Context, message Message) bool {
		return message.Action == action
	}
}
