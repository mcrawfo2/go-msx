package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"github.com/emicklei/go-restful"
)

const (
	HeaderNameAuthorization = "Authorization"
)

var (
	HeaderAuthorization *restful.Parameter
	logger              = log.NewLogger("webservice")
)

func init() {
	HeaderAuthorization = restful.
		HeaderParameter(HeaderNameAuthorization, "Authentication token in form 'Bearer {token}'").
		Required(false)
}

func StandardRoute(b *restful.RouteBuilder) {
	StandardUserContext(b)
	StandardReturns(b)
}

func StandardUserContext(b *restful.RouteBuilder) {
	b.Filter(RequireAuthorizedFilter).
		Notes("This endpoint is secured with JWT").
		Param(HeaderAuthorization)
}

func StandardReturns(b *restful.RouteBuilder) {
	b.Do(Returns(200, 400, 401, 403))
}
