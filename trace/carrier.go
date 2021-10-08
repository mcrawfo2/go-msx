package trace

import (
	"bytes"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

const (
	HeaderTraceId      = "x-b3-traceid"
	HeaderSpanId       = "x-b3-spanid"
	HeaderParentSpanId = "x-b3-parentspanid"
	HeaderSampled      = "x-b3-sampled"
)

var (
	ErrNoTrace      = errors.New("No trace context found")
	ErrInvalidTrace = errors.New("Invalid trace context found")
	ErrNoTracer     = errors.New("No tracer activated")
)

type carrier interface {
	Inject(spanContext SpanContext) error
	Extract() (SpanContext, error)
}

var _ carrier = make(TextMapCarrier)
var _ carrier = make(HttpHeadersCarrier)

type TextMapCarrier map[string]string

func (t TextMapCarrier) expand() (TextMapCarrier, error) {
	expanded := make(TextMapCarrier)

	if headerValue, ok := t["b3"]; ok {
		headerValueParts := strings.SplitN(headerValue, "-", 4)
		if len(headerValueParts) < 3 {
			return nil, errors.Wrap(ErrInvalidTrace, "Invalid b3 header")
		}

		expanded.Set(HeaderTraceId, headerValueParts[0])
		expanded.Set(HeaderSpanId, headerValueParts[1])

		if headerValueParts[2] != "" {
			expanded.Set(HeaderSampled, headerValueParts[2])
		}

		if len(headerValueParts) == 4 && headerValueParts[3] != "" {
			expanded.Set(HeaderParentSpanId, headerValueParts[3])
		}
	} else {
		err := t.ForeachKey(func(key, val string) error {
			headerName := strings.ToLower(key)
			switch headerName {
			case HeaderTraceId, HeaderSpanId, HeaderSampled, HeaderParentSpanId:
				expanded.Set(headerName, strings.Trim(val, "\""))
			}
			return nil
		})
		if err != nil {
			return nil, errors.Wrap(ErrInvalidTrace, err.Error())
		}
	}

	if len(expanded) == 0 {
		return nil, errors.Wrap(ErrInvalidTrace, "No fields found")
	}

	return expanded, nil
}

func (t TextMapCarrier) Inject(spanContext SpanContext) error {
	err := tracer.Inject(spanContext, t)
	if err != nil {
		return err
	}

	expanded, err := t.expand()
	if err != nil {
		return err
	}

	byteBuffer := bytes.NewBuffer(nil)
	byteBuffer.WriteString(expanded[HeaderTraceId])
	byteBuffer.WriteRune('-')
	byteBuffer.WriteString(expanded[HeaderSpanId])

	if sampled, ok := expanded[HeaderSampled]; ok {
		byteBuffer.WriteRune('-')
		byteBuffer.WriteString(sampled)
	}

	if parentId, ok := expanded[HeaderParentSpanId]; ok {
		byteBuffer.WriteRune('-')
		byteBuffer.WriteString(parentId)
	}

	t.Set("b3", byteBuffer.String())
	return nil
}

func (t TextMapCarrier) Set(key, value string) {
	t[key] = value
}

func (t TextMapCarrier) ForeachKey(fn func(key, value string) error) error {
	for k, v := range t {
		err := fn(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t TextMapCarrier) Extract() (SpanContext, error) {
	expanded, err := t.expand()
	if err != nil {
		return nil, err
	}

	return tracer.Extract(expanded)
}

type HttpHeadersCarrier map[string][]string

func (t HttpHeadersCarrier) toTextMapCarrier() TextMapCarrier {
	var textMapCarrier = make(TextMapCarrier)
	_ = t.ForeachKey(func(key, value string) error {
		textMapCarrier[key] = value
		return nil
	})
	return textMapCarrier
}

func (t HttpHeadersCarrier) Inject(sc SpanContext) error {
	textMapCarrier := make(TextMapCarrier)
	err := tracer.Inject(sc, textMapCarrier)
	if err != nil {
		return err
	}

	for k, v := range textMapCarrier {
		t.Set(k, v)
	}

	return nil
}

func (t HttpHeadersCarrier) Extract() (SpanContext, error) {
	textMapCarrier, err := t.toTextMapCarrier().expand()
	if err != nil {
		return nil, err
	}

	return tracer.Extract(textMapCarrier)
}

func (t HttpHeadersCarrier) Set(key, value string) {
	http.Header(t).Set(key, value)
}

func (t HttpHeadersCarrier) ForeachKey(fn func(key, value string) error) error {
	for k, v := range t {
		if len(v) == 0 {
			continue
		}

		err := fn(k, v[0])
		if err != nil {
			return err
		}
	}
	return nil
}
