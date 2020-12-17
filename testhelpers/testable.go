package testhelpers

import "testing"

type Testable interface {
	Test(t *testing.T)
}

type TestFunc func(t *testing.T)

func (f TestFunc) Test(t *testing.T) {
	f(t)
}