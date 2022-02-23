package lowerplural

import (
	"context"
)

type contextKey string

const (
	contextKeyUpperCamelSingularPublisher = contextKey("UpperCamelSingularPublisher")
)

func lowerCamelSingularPublisherFromContext(ctx context.Context) lowerCamelSingularPublisherApi {
	value, _ := ctx.Value(contextKeyUpperCamelSingularPublisher).(lowerCamelSingularPublisherApi)
	return value
}

func contextWithUpperCamelSingularPublisher(ctx context.Context, publisher lowerCamelSingularPublisherApi) context.Context {
	return context.WithValue(ctx, contextKeyUpperCamelSingularPublisher, publisher)
}
