package eventlog

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
)

const TopicName = "EVENT_LOG_TOPIC"

func Publish(ctx context.Context, message Message) error {
	return stream.PublishObject(ctx, TopicName, message, nil)
}

func PublishFromProducer(ctx context.Context, producer MessageProducer) error {
	return Publish(ctx, producer.Message(ctx))
}
