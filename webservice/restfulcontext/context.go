// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restfulcontext

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/emicklei/go-restful"
)

type contextKey string

const (
	contextKeyContainer       = contextKey("RestfulContainer")
	contextKeyRouteSelector   = contextKey("RestfulRouteSelector")
	contextKeyWebService      = contextKey("RestfulWebService")
	contextKeyRoute           = contextKey("RestfulRoute")
	contextKeyFilterFunctions = contextKey("RestfulFilterFunctions")
	contextKeyOperation       = contextKey("Operation")
)

func ContextContainer() types.ContextKeyAccessor[*restful.Container] {
	return types.NewContextKeyAccessor[*restful.Container](contextKeyContainer)
}

func ContextRouteSelector() types.ContextKeyAccessor[restful.RouteSelector] {
	return types.NewContextKeyAccessor[restful.RouteSelector](contextKeyRouteSelector)
}

func ContextWebService() types.ContextKeyAccessor[*restful.WebService] {
	return types.NewContextKeyAccessor[*restful.WebService](contextKeyWebService)
}

func ContextRoute() types.ContextKeyAccessor[*restful.Route] {
	return types.NewContextKeyAccessor[*restful.Route](contextKeyRoute)
}

func ContextFilterFunctions() types.ContextKeyAccessor[[]restful.FilterFunction] {
	return types.NewContextKeyAccessor[[]restful.FilterFunction](contextKeyFilterFunctions)
}

func ContextOperationName() types.ContextKeyAccessor[string] {
	return types.NewContextKeyAccessor[string](contextKeyOperation)
}
