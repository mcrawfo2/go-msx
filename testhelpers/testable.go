package testhelpers

import "testing"

type Testable interface {
	Test(t *testing.T)
}
