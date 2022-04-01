// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package traceinterceptor

import (
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"net/http"
)

var (
	logger = log.NewLogger("msx.httpclient.traceinterceptor")
)

func NewInterceptor(fn httpclient.DoFunc) httpclient.DoFunc {
	return func(req *http.Request) (*http.Response, error) {
		ctx := req.Context()

		operationName := httpclient.OperationNameFromContext(ctx)
		ctx, span := trace.NewSpan(ctx, operationName,
			trace.StartWithTag(trace.FieldOperation, operationName),
			trace.StartWithTag(trace.FieldHttpMethod, req.Method),
			trace.StartWithTag(trace.FieldHttpUrl, req.URL.String()),
			trace.StartWithTag(trace.FieldSpanType, "web"))
		defer span.Finish()
		req = req.WithContext(ctx)

		err := trace.HttpHeadersCarrier(req.Header).Inject(span.Context())
		if err != nil {
			logger.WithContext(ctx).WithError(err).Error("Failed to inject tracing into request")
		}

		response, err := fn(req)
		if response != nil {
			span.LogFields(trace.HttpCode(response.StatusCode))
		}
		if err != nil {
			span.LogFields(trace.Error(err))
		}

		return response, err
	}
}
