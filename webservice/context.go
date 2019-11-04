package webservice

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
	"net/http"
)

const (
	configRootWebServer = "server"
)

var (
	ErrDisabled = errors.New("Web server disabled")
	service     *WebServer
)

func Start(ctx context.Context) error {
	return service.Serve(ctx)
}

func Stop(ctx context.Context) error {
	return service.StopServing(ctx)
}

func ConfigureWebServer(cfg *config.Config, ctx context.Context) (err error) {
	service, err = NewWebServerFromConfig(cfg, ctx)
	return
}

func NewWebServerFromConfig(cfg *config.Config, ctx context.Context) (*WebServer, error) {
	var webServerConfig WebServerConfig
	if err := cfg.Populate(&webServerConfig, configRootWebServer); err != nil {
		return nil, err
	}

	if !webServerConfig.Enabled {
		return nil, ErrDisabled
	}

	return NewWebServer(&webServerConfig, ctx), nil
}

type webServerContextKey int

const (
	contextKeyWebServer webServerContextKey = iota
	contextKeyContainer
	contextKeyRouter
	contextKeyService
	contextKeyRoute
	contextKeyOperation
	contextKeyPathParameters
	contextKeyHttpRequest
	contextKeySecurityProvider
)

func ContextWithWebServer(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKeyWebServer, service)
}

func WebServerFromContext(ctx context.Context) *WebServer {
	return ctx.Value(contextKeyWebServer).(*WebServer)
}

func ContextWithContainer(ctx context.Context, container *restful.Container) context.Context {
	return context.WithValue(ctx, contextKeyContainer, container)
}

func ContainerFromContext(ctx context.Context) *restful.Container {
	return ctx.Value(contextKeyContainer).(*restful.Container)
}

func ContextWithRouter(ctx context.Context, router restful.RouteSelector) context.Context {
	return context.WithValue(ctx, contextKeyRouter, router)
}

func RouterFromContext(ctx context.Context) restful.RouteSelector {
	return ctx.Value(contextKeyRouter).(restful.RouteSelector)
}

func ContextWithService(ctx context.Context, service *restful.WebService) context.Context {
	return context.WithValue(ctx, contextKeyService, service)
}

func ServiceFromContext(ctx context.Context) *restful.WebService {
	return ctx.Value(contextKeyService).(*restful.WebService)
}

func ContextWithRoute(ctx context.Context, route *restful.Route) context.Context {
	return context.WithValue(ctx, contextKeyRoute, route)
}

func RouteFromContext(ctx context.Context) *restful.Route {
	return ctx.Value(contextKeyRoute).(*restful.Route)
}

func ContextWithPathParameters(ctx context.Context, pathParameters map[string]string) context.Context {
	return context.WithValue(ctx, contextKeyPathParameters, pathParameters)
}

func PathParametersFromContext(ctx context.Context) map[string]string {
	return ctx.Value(contextKeyPathParameters).(map[string]string)
}

func ContextWithRouteOperation(ctx context.Context, operation string) context.Context {
	return context.WithValue(ctx, contextKeyOperation, operation)
}

func RouteOperationFromContext(ctx context.Context) string {
	return ctx.Value(contextKeyOperation).(string)
}

func ContextWithHttpRequest(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, contextKeyHttpRequest, req)
}

func HttpRequestFromContext(ctx context.Context) *http.Request {
	return ctx.Value(contextKeyHttpRequest).(*http.Request)
}

func ContextWithSecurityProvider(ctx context.Context, provider SecurityProvider) context.Context {
	return context.WithValue(ctx, contextKeySecurityProvider, provider)
}

func SecurityProviderFromContext(ctx context.Context) SecurityProvider {
	providerInterface := ctx.Value(contextKeySecurityProvider)
	if providerInterface == nil {
		return nil
	} else {
		return providerInterface.(SecurityProvider)
	}
}
