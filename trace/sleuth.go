package trace

import (
	"bytes"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/zipkin"
	"strings"
)

var (
	ErrNoTrace      = errors.New("No trace context found")
	ErrInvalidTrace = errors.New("Invalid trace context found")
)

type SleuthTextMapPropagator struct {
	zipkin zipkin.Propagator
}

func (s SleuthTextMapPropagator) Inject(spanContext jaeger.SpanContext, abstractCarrier interface{}) error {
	textMapWriter, ok := abstractCarrier.(opentracing.TextMapWriter)
	if !ok {
		return opentracing.ErrInvalidCarrier
	}

	byteBuffer := bytes.NewBuffer(nil)
	if spanContext.TraceID().High != 0 {
		byteBuffer.WriteString(fmt.Sprintf("%08x", spanContext.TraceID().High))
	}
	byteBuffer.WriteString(fmt.Sprintf("%08x", spanContext.TraceID().Low))
	byteBuffer.WriteRune('-')
	byteBuffer.WriteString(fmt.Sprintf("%016x", uint64(spanContext.SpanID())))
	byteBuffer.WriteRune('-')
	switch {
	case spanContext.IsDebug():
		byteBuffer.WriteRune('d')
	case spanContext.IsSampled():
		byteBuffer.WriteRune('1')
	default:
		byteBuffer.WriteRune('0')
	}

	if spanContext.ParentID() != 0 {
		byteBuffer.WriteRune('-')
		byteBuffer.WriteString(fmt.Sprintf("%016x", uint64(spanContext.ParentID())))
	}

	textMapWriter.Set("b3", byteBuffer.String())

	return nil
}

func (s SleuthTextMapPropagator) Extract(abstractCarrier interface{}) (jaeger.SpanContext, error) {
	textMapReader, ok := abstractCarrier.(opentracing.TextMapReader)
	if !ok {
		return jaeger.SpanContext{}, opentracing.ErrInvalidCarrier
	}

	httpHeaders := opentracing.HTTPHeadersCarrier{}

	var headerValue string
	err := textMapReader.ForeachKey(func(key, val string) error {
		if key == "b3" {
			headerValue = val
		} else {
			headerName := strings.ToLower(key)
			switch headerName {
			case "x-b3-traceid", "x-b3-spanid", "x-b3-sampled", "x-b3-parentspanid":
				httpHeaders.Set(headerName, strings.Trim(val, "\""))
			}
		}
		return nil
	})

	if err != nil {
		return jaeger.SpanContext{}, err
	} else if headerValue != "" {
		headerValueParts := strings.SplitN(headerValue, "-", 4)
		if len(headerValueParts) < 3 {
			return jaeger.SpanContext{}, ErrInvalidTrace
		}

		httpHeaders.Set("x-b3-traceid", headerValueParts[0])
		httpHeaders.Set("x-b3-spanid", headerValueParts[1])

		if headerValueParts[2] != "" {
			httpHeaders.Set("x-b3-sampled", headerValueParts[2])
		}

		if len(headerValueParts) == 4 && headerValueParts[3] != "" {
			httpHeaders.Set("x-b3-parentspanid", headerValueParts[3])
		}
	}

	if len(httpHeaders) == 0 {
		return jaeger.SpanContext{}, ErrNoTrace
	} else {
		return s.zipkin.Extract(httpHeaders)
	}
}

func NewSleuthTextMapPropagator(zipkinPropagator zipkin.Propagator) SleuthTextMapPropagator {
	return SleuthTextMapPropagator{
		zipkin: zipkinPropagator,
	}
}
