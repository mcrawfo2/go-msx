package trace

import "context"

type contextKey string

const contextKeySpan = contextKey("Span")

func SpanFromContext(ctx context.Context) Span {
	span, _ := ctx.Value(contextKeySpan).(Span)
	return span
}

func ContextWithSpan(ctx context.Context, span Span) context.Context {
	return context.WithValue(ctx, contextKeySpan, span)
}
