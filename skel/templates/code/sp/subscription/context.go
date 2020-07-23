package subscription

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
)

var logger = log.NewLogger("${app.name}.subscription")

type contextKey int

const (
	contextKeyController contextKey = iota
	contextKeyService
	contextKeyRepository
)

func controllerFromContext(ctx context.Context) webservice.RestController {
	value, _ := ctx.Value(contextKeyController).(webservice.RestController)
	return value
}

func contextWithController(ctx context.Context, controller webservice.RestController) context.Context {
	return context.WithValue(ctx, contextKeyController, controller)
}

func serviceFromContext(ctx context.Context) subscriptionServiceApi {
	value, _ := ctx.Value(contextKeyService).(subscriptionServiceApi)
	return value
}

func contextWithService(ctx context.Context, service subscriptionServiceApi) context.Context {
	return context.WithValue(ctx, contextKeyService, service)
}
