// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package webservice

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/restfulcontext"
	"github.com/emicklei/go-restful"
)

const (
	configRootWebServer = "server"
)

type webServerContextKey int

const (
	contextKeyWebServer webServerContextKey = iota
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

// Deprecated.
func ContextWithContainer(ctx context.Context, container *restful.Container) context.Context {
	return restfulcontext.ContextContainer().Set(ctx, container)
}

// Deprecated.
func ContainerFromContext(ctx context.Context) *restful.Container {
	return restfulcontext.ContextContainer().Get(ctx)
}

// Deprecated.
func ContextWithRouter(ctx context.Context, router restful.RouteSelector) context.Context {
	return restfulcontext.ContextRouteSelector().Set(ctx, router)
}

// Deprecated.
func RouterFromContext(ctx context.Context) restful.RouteSelector {
	return restfulcontext.ContextRouteSelector().Get(ctx)
}

// Deprecated.
func ContextWithService(ctx context.Context, service *restful.WebService) context.Context {
	return restfulcontext.ContextWebService().Set(ctx, service)
}

// Deprecated.
func ServiceFromContext(ctx context.Context) *restful.WebService {
	return restfulcontext.ContextWebService().Get(ctx)
}

// Deprecated.
func ContextWithRoute(ctx context.Context, route *restful.Route) context.Context {
	return restfulcontext.ContextRoute().Set(ctx, route)
}

// Deprecated.
func RouteFromContext(ctx context.Context) *restful.Route {
	return restfulcontext.ContextRoute().Get(ctx)
}

// Deprecated.
func ContextWithRouteOperation(ctx context.Context, operation string) context.Context {
	return restfulcontext.ContextOperationName().Set(ctx, operation)
}

// Deprecated.
func RouteOperationFromContext(ctx context.Context) string {
	return restfulcontext.ContextOperationName().Get(ctx)
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

// Deprecated.
func ContextWithFilters(ctx context.Context, filters ...restful.FilterFunction) context.Context {
	return restfulcontext.ContextFilterFunctions().Set(ctx, filters)
}

// Deprecated.
func FiltersFromContext(ctx context.Context) []restful.FilterFunction {
	return restfulcontext.ContextFilterFunctions().Get(ctx)
}
