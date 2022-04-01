// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package statsinterceptor

import (
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient"
	"cto-github.cisco.com/NFV-BU/go-msx/stats"
	"net/http"
	"strconv"
	"time"
)

const (
	statsSubsystemHttpClient = "httpclient"
	statsHistogramCallTime   = "call_time"   // Distribution of response times
	statsGaugeCalls          = "calls"       // Active call count
	statsCounterCallErrors   = "call_errors" // Responses returning an error
	statsCounterCallCodes    = "call_codes"  // Response codes by API
)

var (
	histogramVecCallTime = stats.NewHistogramVec(statsSubsystemHttpClient, statsHistogramCallTime, nil, "operation", "host", "path")
	gaugeVecCalls        = stats.NewGaugeVec(statsSubsystemHttpClient, statsGaugeCalls, "operation", "host", "path")
	counterVecCallErrors = stats.NewCounterVec(statsSubsystemHttpClient, statsCounterCallErrors, "operation", "host", "path")
	counterVecCallCodes  = stats.NewCounterVec(statsSubsystemHttpClient, statsCounterCallCodes, "operation", "host", "path", "code")
)

func NewInterceptor(fn httpclient.DoFunc) httpclient.DoFunc {
	return func(req *http.Request) (response *http.Response, err error) {
		ctx := req.Context()
		operationName := httpclient.OperationNameFromContext(ctx)
		if operationName == "" {
			operationName = "custom"
		}
		host := req.URL.Host
		path := req.URL.Path

		gaugeVecCalls.WithLabelValues(operationName, host, path).Inc()
		startTime := time.Now()

		defer func() {
			callTime := time.Now().Sub(startTime).Milliseconds()
			histogramVecCallTime.WithLabelValues(operationName, host, path).Observe(float64(callTime))

			gaugeVecCalls.WithLabelValues(operationName, host, path).Dec()

			if err != nil {
				counterVecCallErrors.WithLabelValues(operationName, host, path).Inc()
			} else {
				counterVecCallCodes.WithLabelValues(operationName, host, path, strconv.Itoa(response.StatusCode)).Inc()
			}
		}()

		response, err = fn(req)
		return
	}
}
