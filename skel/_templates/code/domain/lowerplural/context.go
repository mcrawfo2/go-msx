package lowerplural

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
)

type contextKey int

const (
	contextKeyUpperCamelSingularController contextKey = iota
	contextKeyUpperCamelSingularService
	contextKeyUpperCamelSingularRepository
)

func lowerCamelSingularControllerFromContext(ctx context.Context) webservice.RestController {
	value, _ := ctx.Value(contextKeyUpperCamelSingularController).(webservice.RestController)
	return value
}

func contextWithUpperCamelSingularController(ctx context.Context, controller webservice.RestController) context.Context {
	return context.WithValue(ctx, contextKeyUpperCamelSingularController, controller)
}

func lowerCamelSingularServiceFromContext(ctx context.Context) lowerCamelSingularServiceApi {
	value, _ := ctx.Value(contextKeyUpperCamelSingularService).(lowerCamelSingularServiceApi)
	return value
}

func contextWithUpperCamelSingularService(ctx context.Context, service lowerCamelSingularServiceApi) context.Context {
	return context.WithValue(ctx, contextKeyUpperCamelSingularService, service)
}

func lowerCamelSingularRepositoryFromContext(ctx context.Context) lowerCamelSingularRepositoryApi {
	value, _ := ctx.Value(contextKeyUpperCamelSingularRepository).(lowerCamelSingularRepositoryApi)
	return value
}

func contextWithUpperCamelSingularRepository(ctx context.Context, repository lowerCamelSingularRepositoryApi) context.Context {
	return context.WithValue(ctx, contextKeyUpperCamelSingularRepository, repository)
}
