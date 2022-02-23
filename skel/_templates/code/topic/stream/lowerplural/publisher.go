package lowerplural

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"

	"cto-github.cisco.com/NFV-BU/go-msx/skel/_templates/code/topic/api"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

const (
	topicUpperCamelSingular = "SCREAMING_SNAKE_SINGULAR_TOPIC"
)

type lowerCamelSingularPublisherApi interface {
	PublishUpperCamelSingular(ctx context.Context, id types.UUID, data string) error
}

type lowerCamelSingularPublisher struct {
	publisherService stream.PublisherService
}

func (p *lowerCamelSingularPublisher) PublishUpperCamelSingular(ctx context.Context, id types.UUID, data string) error {
	return p.PublishUpperCamelSingularMessage(ctx, api.UpperCamelSingularMessage{
		Id:   id,
		Data: data,
	})
}

func (p *lowerCamelSingularPublisher) PublishUpperCamelSingularMessage(ctx context.Context, msg api.UpperCamelSingularMessage) error {
	return p.publisherService.PublishObject(ctx, topicUpperCamelSingular, msg, nil)
}

func (p *lowerCamelSingularPublisher) PublishUpperCamelSingularFromProducer(ctx context.Context, producer UpperCamelSingularMessageProducer) error {
	return p.PublishUpperCamelSingularMessage(ctx, producer.Produce())
}

func newUpperCamelSingularPublisher(ctx context.Context) lowerCamelSingularPublisherApi {
	publisher := lowerCamelSingularPublisherFromContext(ctx)
	if publisher == nil {
		return &lowerCamelSingularPublisher{
			publisherService: stream.PublisherServiceFromContext(ctx),
		}
	}
	return publisher
}
