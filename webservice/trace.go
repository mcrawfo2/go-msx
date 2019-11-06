package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"github.com/emicklei/go-restful"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func tracingFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	ctx := req.Request.Context()
	operationName := RouteOperationFromContext(ctx)

	// Grab the incoming trace
	wireContext, err := opentracing.GlobalTracer().Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Request.Header))
	if err != nil {
		logger.WithError(err).Info("No parent tracing.  Starting new trace")
	}

	ctx, span := trace.NewSpan(ctx, operationName, ext.RPCServerOption(wireContext))
	defer span.Finish()
	req.Request = req.Request.WithContext(ctx)

	chain.ProcessFilter(req, resp)
}
