package traceinterceptor

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/httpclient"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/tracetest"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"github.com/pkg/errors"
	"net/http"
	"testing"
)

func TestNewInterceptor(t *testing.T) {
	var requestOperation = "getDevices"
	var requestMethod = "GET"
	var requestUrl = "http://manageddevice/api/v1/devices"
	var err = errors.New("error")

	tests := []struct {
		name     string
		response *http.Response
		err      error
		wantSpan tracetest.Check
	}{
		{
			name: "Success",
			response: &http.Response{
				Status:     "OK",
				StatusCode: 200,
			},
			wantSpan: tracetest.Check{
				tracetest.HasTag(trace.FieldOperation, requestOperation),
				tracetest.HasTag(trace.FieldHttpMethod, requestMethod),
				tracetest.HasTag(trace.FieldHttpUrl, requestUrl),
				tracetest.HasLogWithField(trace.FieldHttpCode, "200"),
			},
		},
		{
			name:     "Error",
			response: nil,
			err:      err,
			wantSpan: tracetest.Check{
				tracetest.HasTag(trace.FieldOperation, requestOperation),
				tracetest.HasTag(trace.FieldHttpMethod, requestMethod),
				tracetest.HasTag(trace.FieldHttpUrl, requestUrl),
				tracetest.HasLogWithField("error", err.Error()),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var span trace.Span
			tracetest.RecordTracing()

			ctx := context.Background()
			ctx = httpclient.ContextWithOperationName(ctx, requestOperation)

			request, _ := http.NewRequest(requestMethod, requestUrl, http.NoBody)
			request = request.WithContext(ctx)

			got := NewInterceptor(func(req *http.Request) (*http.Response, error) {
				span = trace.SpanFromContext(req.Context())
				if tt.response != nil {
					tt.response.Request = req
				}
				return tt.response, tt.err
			})

			_, _ = got(request)

			testhelpers.ReportErrors(t, "Span", tt.wantSpan.Check(span))
		})
	}
}
