package stream

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
)

type Publisher interface {
	Publish(message *message.Message) error
	Close() error
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
	if publisher, err = NewPublisher(cfg, topic); err != nil && err != ErrBinderNotEnabled {
		return errors.Wrap(err, "Failed to create stream publisher")
	} else if err == ErrBinderNotEnabled {
		return err
	}
	defer publisher.Close()

	if err = publisher.Publish(msg); err != nil {
		return errors.Wrap(err, "Failed to publish message")
	}

	return nil
}

func PublishObject(ctx context.Context, topic string, payload interface{}, metadata map[string]string) (err error) {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return Publish(ctx, topic, bytes, metadata)

}

type TopicPublisher struct {
	cfg       *BindingConfiguration
	publisher message.Publisher
}

func (p *TopicPublisher) Publish(message *message.Message) error {
	return p.publisher.Publish(p.cfg.Destination, message)
}

func (p *TopicPublisher) Close() error {
	return p.publisher.Close()
}

func NewTopicPublisher(publisher message.Publisher, cfg *BindingConfiguration) Publisher {
	return NewTracePublisher(
		&TopicPublisher{
			publisher: publisher,
			cfg:       cfg,
		},
		cfg)
}

type IntransientPublisher struct {
	publisher Publisher
}

func (n *IntransientPublisher) Publish(msg *message.Message) error {
	return n.publisher.Publish(msg)
}

func (n *IntransientPublisher) Close() error {
	return nil
}

func NewIntransientPublisher(publisher Publisher) Publisher {
	return &IntransientPublisher{
		publisher: publisher,
	}
}
