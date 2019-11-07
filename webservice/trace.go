package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"github.com/emicklei/go-restful"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
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

	span.SetTag(trace.FieldOperation, RouteOperationFromContext(ctx))
	span.SetTag(trace.FieldHttpMethod, req.Request.Method)
	span.SetTag(trace.FieldHttpUrl, req.Request.URL.Path)

	chain.ProcessFilter(req, resp)

	span.LogFields(log.Int(trace.FieldHttpCode, resp.StatusCode()))
	if resp.Error() != nil {
		span.LogFields(log.String(trace.FieldError, resp.Error().Error()))
	}
}
