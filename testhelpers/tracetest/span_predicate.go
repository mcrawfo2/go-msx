package tracetest

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	tracelog "github.com/opentracing/opentracing-go/log"
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
			var tags = span.(*Span).Tags
			return tags[name] == val
		},
	}
}

func HasLogWithField(field tracelog.Field) SpanPredicate {
	return SpanPredicate{
		Description: fmt.Sprintf("span.Logs.*.Fields[%q] == %+v", field.Key(), field.Value()),
		Matches: func(span opentracing.Span) bool {
			var logs = span.(*Span).Logs
			for _, log := range logs {
				for _, logField := range log.Fields {
					if logField.Key() == field.Key() && logField.Value() == field.Value() {
						return true
					}
				}
			}
			return false
		},
	}
}