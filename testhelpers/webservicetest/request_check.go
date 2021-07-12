package webservicetest

import (
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"testing"
)

type RequestVerifier func(t *testing.T, req *restful.Request)

type RequestCheck struct {
	Validators []RequestPredicate
}

func (r RequestCheck) Check(req *restful.Request) []error {
	var results []error

	for _, predicate := range r.Validators {
		if !predicate.Matches(req) {
			results = append(results, RequestCheckError{
				Validator: predicate,
			})
		}
	}

	return results

}

type RequestCheckError struct {
	Validator RequestPredicate
}

func (c RequestCheckError) Error() string {
	return fmt.Sprintf("Failed Request validator: %s", c.Validator.Description)
}

type RequestPredicate struct {
	Description string
	Matches     func(*restful.Request) bool
}

func RequestHasAttribute(name string, value interface{}) RequestPredicate {
	return RequestPredicate{
		Description: fmt.Sprintf("request.Attribute(%q) == %v", name, value),
		Matches: func(request *restful.Request) bool {
			return request.Attribute(name) == value
		},
	}
}

func RequestHasPathParameter(name string, value interface{}) RequestPredicate {
	return RequestPredicate{
		Description: fmt.Sprintf("request.Attribute(%q) == %v", name, value),
		Matches: func(request *restful.Request) bool {
			return request.PathParameter(name) == value
		},
	}
}
