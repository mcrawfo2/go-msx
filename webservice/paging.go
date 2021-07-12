package webservice

import "github.com/emicklei/go-restful/v3"

var QueryParamPageNumber = restful.QueryParameter("page", "Page number (0-based)").
	Required(true).
	DefaultValue("0")

var QueryParamPageSize = restful.QueryParameter("pageSize", "Page size").
	Required(true).
	DefaultValue("100")

func Paginated(b *restful.RouteBuilder) {
	b.Param(QueryParamPageNumber)
	b.Param(QueryParamPageSize)
}
