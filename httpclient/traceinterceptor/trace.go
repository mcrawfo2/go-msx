package traceinterceptor

import (
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"github.com/opentracing/opentracing-go"
	"net/http"
)

var (
	logger = log.NewLogger("msx.httpclient.traceinterceptor")
)

func NewInterceptor(fn httpclient.DoFunc) httpclient.DoFunc {
	return func(req *http.Request) (*http.Response, error) {
		ctx := req.Context()
		operationName := httpclient.OperationNameFromContext(ctx)
		ctx, span := trace.NewSpan(ctx, operationName)
		defer span.Finish()
		req = req.WithContext(ctx)

		span.SetTag(trace.FieldOperation, operationName)
		span.SetTag(trace.FieldHttpMethod, req.Method)
		span.SetTag(trace.FieldHttpUrl, req.URL.String())

		err := opentracing.GlobalTracer().Inject(
			span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(req.Header))
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
