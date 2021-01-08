package clienttest

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"fmt"
	"testing"
)

type ResponseVerifier func(t *testing.T, req *integration.MsxResponse)

type ResponseCheck struct {
	Validators []ResponsePredicate
}

func (r ResponseCheck) Check(req *integration.MsxResponse) []error {
	var results []error

	for _, predicate := range r.Validators {
		if !predicate.Matches(req) {
			results = append(results, ResponseCheckError{
				Validator: predicate,
			})
		}
	}

	return results

}

type ResponseCheckError struct {
	Validator ResponsePredicate
}

func (c ResponseCheckError) Error() string {
	return fmt.Sprintf("Failed response validator: %s", c.Validator.Description)
}

type ResponsePredicate struct {
	Description string
	Matches     func(*integration.MsxResponse) bool
}
