package stream

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
)

type Publisher interface {
	Publish(messages ...*message.Message) error
	Close() error
}

type TopicPublisher struct {
	cfg *BindingConfiguration
	publisher message.Publisher
}

func (p *TopicPublisher) Publish (messages ...*message.Message) error {
	return p.publisher.Publish(p.cfg.Destination, messages...)
}

func (p *TopicPublisher) Close() error {
	return p.publisher.Close()
}

func NewTopicPublisher(publisher message.Publisher, cfg *BindingConfiguration) *TopicPublisher {
	return &TopicPublisher{
		publisher: publisher,
		cfg:       cfg,
	}
}

func Publish(ctx context.Context, topic string, payload []byte, metadata map[string]string) (err error) {
	msg := message.NewMessage(watermill.NewUUID(), payload)
	for k, v := range metadata {
		msg.Metadata.Set(k, v)
	}
	msg.SetContext(ctx)

	var cfg *config.Config
	if cfg = config.FromContext(ctx); cfg == nil {
		return errors.New("Failed to obtain application config")
	}

	var publisher Publisher
	if publisher, err = NewPublisher(cfg, topic); err != nil {
		return errors.Wrap(err, "Failed to create stream publisher")
	}
	defer publisher.Close()

	if err = publisher.Publish(msg); err != nil {
		return errors.Wrap(err, "Failed to publish message")
	}

	return nil
}
