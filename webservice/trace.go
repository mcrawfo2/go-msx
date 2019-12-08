package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"github.com/emicklei/go-restful"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"net/http"
)

func tracingFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	ctx := req.Request.Context()
	operationName := RouteOperationFromContext(ctx)

	var opts []opentracing.StartSpanOption

	// Grab the incoming trace
	wireContext, err := opentracing.GlobalTracer().Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Request.Header))
	if err == nil {
		opts = append(opts, ext.RPCServerOption(wireContext))
	}

	ctx, span := trace.NewSpan(ctx, operationName, opts...)
	defer span.Finish()
	req.Request = req.Request.WithContext(ctx)

	span.SetTag(trace.FieldOperation, operationName)
	span.SetTag(trace.FieldHttpMethod, req.Request.Method)
	span.SetTag(trace.FieldHttpUrl, req.Request.URL.Path)

	chain.ProcessFilter(req, resp)

	logContext := log.LogContext{
		"operation": operationName,
		"method":    req.Request.Method,
		"path":      req.Request.URL.Path,
		"code":      resp.StatusCode(),
	}

	span.LogFields(trace.Int(trace.FieldHttpCode, resp.StatusCode()))
	if resp.Error() != nil {
		span.LogFields(trace.Error(resp.Error()))
		logger.WithLogContext(logContext).WithError(resp.Error()).Errorf("Incoming request failed: %s", http.StatusText(resp.StatusCode()))
	} else if resp.StatusCode() < 399 {
		logger.WithLogContext(logContext).Infof("Incoming request succeeded: %s", http.StatusText(resp.StatusCode()))
	} else {
		logger.WithLogContext(logContext).Errorf("Incoming request failed: %s", http.StatusText(resp.StatusCode()))
	}
}
