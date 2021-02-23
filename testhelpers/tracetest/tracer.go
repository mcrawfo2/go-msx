package tracetest

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
)

func RecordTracing() *mocktracer.MockTracer {
	var result = mocktracer.New()

	opentracing.SetGlobalTracer(result)

	return result
}
