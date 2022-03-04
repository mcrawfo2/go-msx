//go:generate mockery --inpackage --name=UpperCamelSingularPublisherApi --structname=MockUpperCamelSingularPublisher --filename mock_publisher_lowersingular.go
//go:generate mockery --inpackage --name=UpperCamelSingularMessageProducer --structname=MockUpperCamelSingularMessageProducer --filename mock_producer_lowersingular.go
package lowerplural

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/stream"

	"cto-github.cisco.com/NFV-BU/go-msx/skel/_templates/code/topic/api"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

const (
	contextKeyUpperCamelSingularPublisher = contextKey("UpperCamelSingularPublisher")
)

type UpperCamelSingularPublisherApi interface {
	PublishUpperCamelSingular(ctx context.Context, id types.UUID, data string) error
}

type lowerCamelSingularPublisher struct {
	publisherService stream.PublisherService
}

type UpperCamelSingularMessageProducer interface {
	Produce() api.UpperCamelSingularMessage
}

func (p *lowerCamelSingularPublisher) PublishUpperCamelSingular(ctx context.Context, id types.UUID, data string) error {
	return p.PublishUpperCamelSingularMessage(ctx, api.UpperCamelSingularMessage{
		Id:   id,
		Data: data,
	})
}

func (p *lowerCamelSingularPublisher) PublishUpperCamelSingularMessage(ctx context.Context, msg api.UpperCamelSingularMessage) error {
	logger.WithContext(ctx).Debugf("Publishing message for lowerplural %q", msg.Id.String())
	return p.publisherService.PublishObject(ctx, topicUpperCamelSingular, msg, nil)
}

func (p *lowerCamelSingularPublisher) PublishUpperCamelSingularFromProducer(ctx context.Context, producer UpperCamelSingularMessageProducer) error {
	return p.PublishUpperCamelSingularMessage(ctx, producer.Produce())
}

func newUpperCamelSingularPublisher(ctx context.Context) UpperCamelSingularPublisherApi {
	publisher := UpperCamelSingularPublisherFromContext(ctx)
	if publisher == nil {
		return &lowerCamelSingularPublisher{
			publisherService: stream.PublisherServiceFromContext(ctx),
		}
	}
	return publisher
}

func UpperCamelSingularPublisherFromContext(ctx context.Context) UpperCamelSingularPublisherApi {
	value, _ := ctx.Value(contextKeyUpperCamelSingularPublisher).(UpperCamelSingularPublisherApi)
	return value
}

func ContextWithUpperCamelSingularPublisher(ctx context.Context, publisher UpperCamelSingularPublisherApi) context.Context {
	return context.WithValue(ctx, contextKeyUpperCamelSingularPublisher, publisher)
}
