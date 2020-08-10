package oss

import "context"

type contextKey int

const (
	contextKeyIntegration contextKey = iota
)

func IntegrationFromContext(ctx context.Context) Api {
	value, _ := ctx.Value(contextKeyIntegration).(Api)
	return value
}

func ContextWithIntegration(ctx context.Context, api Api) context.Context {
	return context.WithValue(ctx, contextKeyIntegration, api)
}
