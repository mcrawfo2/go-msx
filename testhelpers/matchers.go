package testhelpers

import "context"

// AnyContext returns true if the specified value implements context.Context
func AnyContext(value any) bool {
	return Implements[context.Context](value)
}

// Implements returns true if the specified value can be cast to the type parameter
func Implements[A interface{}](value any) bool {
	_, ok := value.(A)
	return ok
}
