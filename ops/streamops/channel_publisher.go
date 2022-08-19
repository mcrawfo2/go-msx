package streamops

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"
	"github.com/pkg/errors"
)

// ChannelPublisher maps to asyncapi.Operation
type ChannelPublisher struct {
	channel          *Channel
	name             string
	publisherService stream.PublisherService
	documentors      ops.Documentors[ChannelPublisher]
}

func (p *ChannelPublisher) AddDocumentor(d ...ops.Documentor[ChannelPublisher]) {
	p.documentors = p.documentors.WithDocumentor(d...)
}

func (p ChannelPublisher) Documentor(predicate ops.DocumentorPredicate[ChannelPublisher]) ops.Documentor[ChannelPublisher] {
	return p.documentors.Documentor(predicate)
}

func (p ChannelPublisher) Channel() *Channel {
	return p.channel
}

func (p ChannelPublisher) Name() string {
	return p.name
}

func (p *ChannelPublisher) Publish(ctx context.Context, payload []byte, metadata map[string]string) error {
	return p.publisherService.Publish(ctx, p.channel.Name(), payload, metadata)
}

func NewChannelPublisher(ctx context.Context, channel *Channel, name string) (*ChannelPublisher, error) {
	if channel == nil {
		return nil, errors.Errorf("Nil channel passed to publisher %q", name)
	}

	publisherService := stream.NewPublisherService(ctx)

	result := &ChannelPublisher{
		channel:          channel,
		name:             name,
		publisherService: publisherService,
	}

	RegisterChannelPublisher(result)

	return result, nil
}

var registeredChannelPublishers = make(map[string]*ChannelPublisher)

func RegisterChannelPublisher(p *ChannelPublisher) {
	registeredChannelPublishers[p.channel.Name()] = p
}

func RegisteredChannelPublisher(channel string) *ChannelPublisher {
	return registeredChannelPublishers[channel]
}
