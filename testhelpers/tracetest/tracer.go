package tracetest

import (
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	mocktracer "cto-github.cisco.com/NFV-BU/go-msx/trace/mock"
)

func RecordTracing() *mocktracer.MockTracer {
	t := mocktracer.NewMockTracer()
	trace.SetTracer(t)
	return t
}
