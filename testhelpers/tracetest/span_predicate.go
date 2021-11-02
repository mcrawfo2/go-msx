package tracetest

import (
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"cto-github.cisco.com/NFV-BU/go-msx/trace/mock"
	"fmt"
)

type SpanPredicate struct {
	Description string
	Matches     func(trace.Span) bool
}

func HasBaggage(name string, val string) SpanPredicate {
	return SpanPredicate{
		Description: fmt.Sprintf("span.Baggage[%q] == %q", name, val),
		Matches: func(span trace.Span) bool {
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
		Matches: func(span trace.Span) bool {
			var tags = span.(*mock.MockSpan).Tags()
			return tags[name] == val
		},
	}
}

func HasLogWithField(key string, value string) SpanPredicate {
	return SpanPredicate{
		Description: fmt.Sprintf("span.Logs.*.Fields[%q] == %s", key, value),
		Matches: func(span trace.Span) bool {
			var logs = span.(*mock.MockSpan).Logs()
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
