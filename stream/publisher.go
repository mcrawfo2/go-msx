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

func Publish(ctx context.Context, binding string, payload []byte, metadata map[string]string) (err error) {
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
	publisher, err = NewPublisher(cfg, binding)
	if err == ErrBinderNotEnabled {
		return err
	} else if err != nil {
		return errors.Wrap(err, "Failed to create stream publisher")
	}
	defer publisher.Close()

	if err = publisher.Publish(msg); err != nil {
		return errors.Wrap(err, "Failed to publish message")
	}

	return nil
}

func PublishObject(ctx context.Context, binding string, payload interface{}, metadata map[string]string) (err error) {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return Publish(ctx, binding, bytes, metadata)
}

// TopicPublisher adapts a message.Publisher to publish to a pre-determined topic
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

// NewTopicPublisher creates a new TopicPublisher instance
func NewTopicPublisher(publisher message.Publisher, cfg *BindingConfiguration) Publisher {
	return NewTracePublisher(
		&TopicPublisher{
			publisher: publisher,
			cfg:       cfg,
		},
		cfg)
}

// IntransientPublisher adapts a Publisher to ignore the Close signal
type IntransientPublisher struct {
	publisher Publisher
}

func (n *IntransientPublisher) Publish(msg *message.Message) error {
	return n.publisher.Publish(msg)
}

func (n *IntransientPublisher) Close() error {
	return nil
}

// NewIntransientPublisher creates a new IntransientPublisher instance from the supplied Publisher
func NewIntransientPublisher(publisher Publisher) Publisher {
	return &IntransientPublisher{
		publisher: publisher,
	}
}

// MessagePublisher is the low-level watermill Publisher interface
type MessagePublisher interface {
	message.Publisher
}
