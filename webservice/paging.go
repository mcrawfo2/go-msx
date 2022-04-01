// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package webservice

import "github.com/emicklei/go-restful"

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
