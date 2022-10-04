// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package tokeninterceptor

import (
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient"
	"cto-github.cisco.com/NFV-BU/go-msx/security/httprequest"
	"net/http"
)

// NewInterceptor returns an HTTP transport interceptor to inject
// the authorization header based on the current user context.
func NewInterceptor(fn httpclient.DoFunc) httpclient.DoFunc {
	return func(req *http.Request) (resp *http.Response, err error) {
		httprequest.InjectToken(req)

		return fn(req)
	}
}
