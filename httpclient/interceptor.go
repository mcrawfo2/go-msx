// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package httpclient

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"net/http"
)

type NewInterceptor func(fn DoFunc) DoFunc

func ClientInterceptor(fn NewInterceptor) ClientConfigurationFunc {
	return func(c *http.Client) {
		ApplyInterceptor(c, fn)
	}
}

func ApplyInterceptor(c *http.Client, fn NewInterceptor) {
	c.Transport = fn(c.Transport.RoundTrip)
}

func ApplyRecoveryErrorInterceptor(c *http.Client) {
	c.Transport = func(fn DoFunc) DoFunc {
		return func(req *http.Request) (resp *http.Response, err error) {
			err = types.RecoverErrorDecorator(func(ctx context.Context) error {
				resp, err = fn(req)
				return err
			})(req.Context())
			return
		}
	}(c.Transport.RoundTrip)
}
