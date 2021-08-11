package scheduled

import "context"

type contextKey int

const (
	contextKeyService contextKey = iota
)

func SchedulerServiceFromContext(ctx context.Context) SchedulerServiceApi {
	service, _ := ctx.Value(contextKeyService).(SchedulerServiceApi)
	return service
}

func ContextWithSchedulerService(ctx context.Context, service SchedulerServiceApi) context.Context {
	return context.WithValue(ctx, contextKeyService, service)
}
