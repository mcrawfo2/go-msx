package tracetest

import (
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"fmt"
)

type CheckError struct {
	Span      trace.Span
	Validator SpanPredicate
}

func (c CheckError) Error() string {
	return fmt.Sprintf("Failed validator: %s - %+v", c.Validator.Description, c.Span)
}

type Check []SpanPredicate

func (c Check) Check(span trace.Span) []error {
	var results []error

	for _, predicate := range c {
		if !predicate.Matches(span) {
			results = append(results, CheckError{
				Span:      span,
				Validator: predicate,
			})
		}
	}

	return results
}
