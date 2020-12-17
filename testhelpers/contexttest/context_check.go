package contexttest

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"fmt"
	"reflect"
)

type ContextCheck struct {
	Validators []ContextPredicate
}

func (c ContextCheck) Check(ctx context.Context) []error {
	var results []error

	for _, predicate := range c.Validators {
		if !predicate.Matches(ctx) {
			results = append(results, ContextCheckError{
				Validator: predicate,
			})
		}
	}

	return results

}

type ContextPredicate struct {
	Description string
	Matches func(ctx context.Context) bool
}

type ContextCheckError struct {
	Validator ContextPredicate
}

func (c ContextCheckError) Error() string {
	return fmt.Sprintf("Failed response validator: %s", c.Validator.Description)
}

func ContextGetterHasAnyValue(fn interface{}) ContextPredicate {
	return ContextPredicate{
		Description: "",
		Matches: func(ctx context.Context) bool {
			fnValue := reflect.ValueOf(fn)
			arguments := []reflect.Value{reflect.ValueOf(ctx)}
			fnResults := fnValue.Call(arguments)

			switch fnResults[0].Type().Kind() {
			case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map:
				return !fnResults[0].IsNil()
			default:
				return !fnResults[0].IsZero()
			}
		},
	}
}

func ContextHasNamedUserContext(name string) ContextPredicate {
	return ContextPredicate{
		Description: fmt.Sprintf("Context has UserContext with name %q", name),
		Matches: func(ctx context.Context) bool {
			userContext := security.UserContextFromContext(ctx)
			return userContext != nil && userContext.UserName == name
		},
	}
}
