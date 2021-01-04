package webservice

import (
	"context"
	"github.com/emicklei/go-restful"
)

const (
	configRootWebServer = "server"
)

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

func ContextWithWebServerValue(ctx context.Context, server *WebServer) context.Context {
	return context.WithValue(ctx, contextKeyWebServer, server)
}

func WebServerFromContext(ctx context.Context) *WebServer {
	if v := ctx.Value(contextKeyWebServer); v != nil {
		if r, ok := v.(*WebServer); ok {
			return r
		}
		logger.WithContext(ctx).Errorf("Invalid web server in context: %+v", v)
	}
	return nil
}

func ContextWithContainer(ctx context.Context, container *restful.Container) context.Context {
	return context.WithValue(ctx, contextKeyContainer, container)
}

func ContainerFromContext(ctx context.Context) *restful.Container {
	if v := ctx.Value(contextKeyContainer); v != nil {
		if r, ok := v.(*restful.Container); ok {
			return r
		}
		logger.WithContext(ctx).Errorf("Invalid container in context: %+v", v)
	}
	return nil
}

func ContextWithRouter(ctx context.Context, router restful.RouteSelector) context.Context {
	return context.WithValue(ctx, contextKeyRouter, router)
}

func RouterFromContext(ctx context.Context) restful.RouteSelector {
	if v := ctx.Value(contextKeyRouter); v != nil {
		if r, ok := v.(restful.RouteSelector); ok {
			return r
		}
		logger.WithContext(ctx).Errorf("Invalid router in context: %+v", v)
	}
	return nil
}

func ContextWithService(ctx context.Context, service *restful.WebService) context.Context {
	return context.WithValue(ctx, contextKeyService, service)
}

func ServiceFromContext(ctx context.Context) *restful.WebService {
	if v := ctx.Value(contextKeyService); v != nil {
		if r, ok := v.(*restful.WebService); ok {
			return r
		}
		logger.WithContext(ctx).Errorf("Invalid web service in context: %+v", v)
	}
	return nil
}

func ContextWithRoute(ctx context.Context, route *restful.Route) context.Context {
	return context.WithValue(ctx, contextKeyRoute, route)
}

func RouteFromContext(ctx context.Context) *restful.Route {
	if v := ctx.Value(contextKeyRoute); v != nil {
		if r, ok := v.(*restful.Route); ok {
			return r
		}
		logger.WithContext(ctx).Errorf("Invalid route in context: %+v", v)
	}
	return nil
}

func ContextWithRouteOperation(ctx context.Context, operation string) context.Context {
	return context.WithValue(ctx, contextKeyOperation, operation)
}

func RouteOperationFromContext(ctx context.Context) string {
	if v := ctx.Value(contextKeyOperation); v != nil {
		if r, ok := v.(string); ok {
			return r
		}
		logger.WithContext(ctx).Errorf("Invalid route operation in context: %+v", v)
	}
	return ""
}

func ContextWithSecurityProvider(ctx context.Context, provider AuthenticationProvider) context.Context {
	return context.WithValue(ctx, contextKeySecurityProvider, provider)
}

func AuthenticationProviderFromContext(ctx context.Context) AuthenticationProvider {
	if v := ctx.Value(contextKeySecurityProvider); v != nil {
		if r, ok := v.(AuthenticationProvider); ok {
			return r
		}
		logger.WithContext(ctx).Errorf("Invalid security provider in context: %+v", v)
	}
	return nil
}
