package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/webservicetest"
	"github.com/emicklei/go-restful"
	"testing"
)

func TestPaginated(t *testing.T) {
	new(RouteBuilderTest).
		WithRouteBuilderDo(Paginated).
		WithRoutePredicate(webservicetest.RouteHasParameter(restful.QueryParameterKind, "page")).
		WithRoutePredicate(webservicetest.RouteHasParameter(restful.QueryParameterKind, "pageSize"))
}
