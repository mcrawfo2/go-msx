package tracetest

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
)

type CheckError struct {
	Span      opentracing.Span
	Validator SpanPredicate
}

func (c CheckError) Error() string {
	return fmt.Sprintf("Failed validator: %s - %+v", c.Validator.Description, c.Span)
}

type Check []SpanPredicate

func (c Check) Check(span opentracing.Span) []error {
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
