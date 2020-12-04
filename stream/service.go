package stream

import "context"

type PublisherService interface {
	Publish(ctx context.Context, topic string, payload []byte, metadata map[string]string) (err error)
}

type publisherFunc func(ctx context.Context, topic string, payload []byte, metadata map[string]string) (err error)

func (p publisherFunc) Publish(ctx context.Context, topic string, payload []byte, metadata map[string]string) (err error) {
	return p(ctx, topic, payload, metadata)
}

var ProductionPublisherService PublisherService = publisherFunc(Publish)

func NewPublisherService(ctx context.Context) PublisherService {
	service := PublisherServiceFromContext(ctx)
	if service == nil {
		service = ProductionPublisherService
	}
	return service
}

// Ensure MockPublisherService is up-to-date
var _ PublisherService = new(MockPublisherService)
