package testhelpers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func ReportErrors(t *testing.T, name string, errs []error) {
	for _, err := range errs {
		assert.Fail(t, err.Error(), "Failed %s validator", name)
	}
}
