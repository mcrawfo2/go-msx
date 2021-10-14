package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/stats"
	"github.com/emicklei/go-restful"
	"net/http"
	"strconv"
	"time"
)

const (
	statsSubsystemHttpServer = "httpserver"
	statsHistogramCallTime   = "call_time"   // Distribution of response times
	statsGaugeCalls          = "calls"       // Active call count
	statsCounterCallErrors   = "call_errors" // Responses returning an error
	statsCounterCallCodes    = "call_codes"  // Response codes by API
)

var (
	histogramVecCallTime = stats.NewHistogramVec(statsSubsystemHttpServer, statsHistogramCallTime, nil, "operation", "path")
	gaugeVecCalls        = stats.NewGaugeVec(statsSubsystemHttpServer, statsGaugeCalls, "operation", "path")
	counterVecCallErrors = stats.NewCounterVec(statsSubsystemHttpServer, statsCounterCallErrors, "operation", "path")
	counterVecCallCodes  = stats.NewCounterVec(statsSubsystemHttpServer, statsCounterCallCodes, "operation", "path", "code")
)

func statsFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	operation := "unknown"

	route := RouteFromContext(req.Request.Context())
	if route != nil {
		operation = route.Operation
	} else if req.Request.Method == http.MethodOptions {
		operation = "options"
	}

	path := req.Request.URL.Path

	start := time.Now()
	gaugeVecCalls.WithLabelValues(operation, path).Inc()

	var err error
	var code int

	defer func() {
		gaugeVecCalls.WithLabelValues(operation, path).Dec()
		histogramVecCallTime.WithLabelValues(operation, path).Observe(float64(time.Since(start)) / float64(time.Millisecond))
		if err != nil {
			counterVecCallErrors.WithLabelValues(operation, path).Inc()
		}
		counterVecCallCodes.WithLabelValues(operation, path, strconv.Itoa(code)).Inc()
	}()

	chain.ProcessFilter(req, resp)

	err = resp.Error()
	if err == nil {
		errInterface := req.Attribute(AttributeError)
		if errInterface != nil {
			err = errInterface.(error)
		}
	}

	code = resp.StatusCode()
}
