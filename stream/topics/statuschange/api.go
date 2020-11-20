package statuschange

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
)

type MessageProducer interface {
	Message(context.Context) (Message, error)
}

func PublishFromProducer(ctx context.Context, producer MessageProducer) (err error) {
	msg, err := producer.Message(ctx)
	if err != nil {
		return err
	}

	return Publish(ctx, msg)
}

func Publish(ctx context.Context, message Message) error {
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return stream.Publish(ctx, TopicStatusChange, bytes, nil)
}

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


func AddListener(fn MessageHandler, filters []MessageFilter) error {
	listener := NewMessageListener(fn, filters)
	return stream.AddListener(TopicStatusChange, listener)
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

func FilterMessageByEntityType(entityType string) MessageFilter {
	return func(ctx context.Context, message Message) bool {
		return message.EntityType == entityType
	}
}

func FilterMessageByStatus(status ...string) MessageFilter {
	return func(ctx context.Context, message Message) bool {
		for _, s := range status {
			if message.Status == s {
				return true
			}
		}
		return false
	}
}
