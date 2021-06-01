package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
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

	traceContext, _ := log.LogContextFromContext(req.Request.Context())
	for k, v := range traceContext {
		logContext[k] = v
	}
	if traceContext == nil {
		traceContext = make(log.LogContext)
	}

	span.LogFields(trace.Int(trace.FieldHttpCode, resp.StatusCode()))

	err = resp.Error()
	if err == nil {
		errInterface := req.Attribute(AttributeError)
		if errInterface != nil {
			err = errInterface.(error)
		}
	}

	if err != nil {
		span.LogFields(trace.Error(err))

		bt := types.BackTraceFromError(err)
		logger.
			WithLogContext(logContext).
			WithError(err).
			WithField(log.FieldStack, bt.Stanza()).
			Errorf("Incoming request failed: %s: %s", http.StatusText(resp.StatusCode()), err.Error())
		log.Stack(logger, ctx, bt)
	} else if resp.StatusCode() < 399 {
		var silenced = false
		silencedAttribute := req.Attribute(AttributeSilenceLog)
		if silencedAttributeValue, ok := silencedAttribute.(bool); ok {
			silenced = silencedAttributeValue
		}
		if !silenced {
			logger.WithLogContext(logContext).Infof("Incoming request succeeded: %s", http.StatusText(resp.StatusCode()))
		}
	} else {
		logger.WithLogContext(logContext).Errorf("Incoming request failed: %s", http.StatusText(resp.StatusCode()))
	}
}
