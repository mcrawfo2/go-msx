package webservice

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
)

const (
	configRootWebServer = "server"
)

var (
	ErrDisabled = errors.New("Web server disabled")
	service     *WebServer
)

func Start(ctx context.Context) error {
	if service == nil {
		return nil
	}
	return service.Serve(ctx)
}

func Stop(ctx context.Context) error {
	if service == nil {
		return nil
	}
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

	actuatorConfig, err := NewManagementSecurityConfig(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read management security config")
	}

	return NewWebServer(&webServerConfig, actuatorConfig, ctx)
}

type webServerContextKey int

const (
	contextKeyWebServer webServerContextKey = iota
	contextKeyContainer
	contextKeyRouter
	contextKeyService
	contextKeyRoute
	contextKeyOperation
	contextKeySecurityProvider
)

func ContextWithWebServer(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKeyWebServer, service)
}

func WebServerFromContext(ctx context.Context) *WebServer {
	webServerInterface := ctx.Value(contextKeyWebServer)
	if webServerInterface != nil {
		return webServerInterface.(*WebServer)
	}
	return nil
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

func ContextWithRouteOperation(ctx context.Context, operation string) context.Context {
	return context.WithValue(ctx, contextKeyOperation, operation)
}

func RouteOperationFromContext(ctx context.Context) string {
	return ctx.Value(contextKeyOperation).(string)
}

func ContextWithSecurityProvider(ctx context.Context, provider AuthenticationProvider) context.Context {
	return context.WithValue(ctx, contextKeySecurityProvider, provider)
}

func AuthenticationProviderFromContext(ctx context.Context) AuthenticationProvider {
	providerInterface := ctx.Value(contextKeySecurityProvider)
	if providerInterface != nil {
		return providerInterface.(AuthenticationProvider)
	}
	return nil
}
