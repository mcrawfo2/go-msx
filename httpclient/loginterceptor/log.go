// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package loginterceptor

import (
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"encoding/json"
	"net/http"
)

var (
	logger = log.NewLogger("msx.httpclient.loginterceptor")
)

func NewInterceptor(fn httpclient.DoFunc) httpclient.DoFunc {
	return func(req *http.Request) (response *http.Response, err error) {
		ctx := req.Context()
		response, err = fn(req)
		if response == nil {
			logger.WithContext(ctx).WithError(err).Errorf("000 : %s %s", req.Method, req.URL.String())
		} else if response.StatusCode > 399 {
			// Fully log the response
			logger.WithContext(ctx).Errorf("%d %s : %s %s", response.StatusCode, response.Status, req.Method, req.URL.String())
			var responseBytes []byte
			responseBytes, _ = json.Marshal(response)
			logger.WithContext(ctx).Error(string(responseBytes))
		} else {
			logger.WithContext(ctx).Infof("%d %s : %s %s", response.StatusCode, response.Status, req.Method, req.URL.String())
		}

		return response, err
	}
}
