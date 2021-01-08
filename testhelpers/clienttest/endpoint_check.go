package clienttest

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"fmt"
	"testing"
)

type ServiceEndpointVerifier func(t *testing.T, endpoint integration.MsxServiceEndpoint)

type ServiceEndpointCheck struct {
	Validators []ServiceEndpointPredicate
}

func (r ServiceEndpointCheck) Check(endpoint integration.MsxServiceEndpoint) []error {
	var results []error

	for _, predicate := range r.Validators {
		if !predicate.Matches(endpoint) {
			results = append(results, ServiceEndpointCheckError{
				Validator: predicate,
			})
		}
	}

	return results

}

type ServiceEndpointCheckError struct {
	Validator ServiceEndpointPredicate
}

func (c ServiceEndpointCheckError) Error() string {
	return fmt.Sprintf("Failed Request validator: %s", c.Validator.Description)
}

type ServiceEndpointPredicate struct {
	Description string
	Matches     func(integration.MsxServiceEndpoint) bool
}

func ServiceEndpointHasMethod(method string) ServiceEndpointPredicate {
	return ServiceEndpointPredicate{
		Description: fmt.Sprintf("request.Method == %q", method),
		Matches: func(endpoint integration.MsxServiceEndpoint) bool {
			return endpoint.Method == method
		},
	}
}

func ServiceEndpointHasPath(path string) ServiceEndpointPredicate {
	return ServiceEndpointPredicate{
		Description: fmt.Sprintf("request.Path == %q", path),
		Matches: func(endpoint integration.MsxServiceEndpoint) bool {
			return endpoint.Path == path
		},
	}
}
