// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

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
		WithRoutePredicate(webservicetest.RouteHasParameter(restful.QueryParameterKind, "pageSize")).
		Test(t)
}
