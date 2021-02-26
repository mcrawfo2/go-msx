package tracetest

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
)

type SpanPredicate struct {
	Description string
	Matches     func(opentracing.Span) bool
}

func HasBaggage(name string, val string) SpanPredicate {
	return SpanPredicate{
		Description: fmt.Sprintf("span.Baggage[%q] == %q", name, val),
		Matches: func(span opentracing.Span) bool {
			spanContext := span.Context()
			matches := false
			spanContext.ForeachBaggageItem(func(k, v string) bool {
				if k == name {
					matches = true
					return false
				}
				return true
			})
			return matches
		},
	}
}

func HasTag(name string, val interface{}) SpanPredicate {
	return SpanPredicate{
		Description: fmt.Sprintf("span.Tags[%q] == %q", name, val),
		Matches: func(span opentracing.Span) bool {
			var tags = span.(*mocktracer.MockSpan).Tags()
			return tags[name] == val
		},
	}
}

func HasLogWithField(key string, value string) SpanPredicate {
	return SpanPredicate{
		Description: fmt.Sprintf("span.Logs.*.Fields[%q] == %s", key, value),
		Matches: func(span opentracing.Span) bool {
			var logs = span.(*mocktracer.MockSpan).Logs()
			for _, log := range logs {
				for _, logField := range log.Fields {
					if logField.Key == key && logField.ValueString == value {
						return true
					}
				}
			}
			return false
		},
	}
}