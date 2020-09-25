package notification

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
)

const TopicName = "NOTIFICATION_TOPIC"

func Publish(ctx context.Context, message Message) error {
	return stream.PublishObject(ctx, TopicName, message, nil)
}

func PublishFromProducer(ctx context.Context, producer MessageProducer) error {
	message, err := producer.Message(ctx)
	if err != nil {
		return err
	}
	return Publish(ctx, message)
}
