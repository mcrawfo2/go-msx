package stream

import (
	"context"
)

type contextKey int

const (
	contextKeyService contextKey = iota
)

func PublisherServiceFromContext(ctx context.Context) PublisherService {
	value, _ := ctx.Value(contextKeyService).(PublisherService)
	return value
}

func ContextWithPublisherService(ctx context.Context, service PublisherService) context.Context {
	return context.WithValue(ctx, contextKeyService, service)
}
