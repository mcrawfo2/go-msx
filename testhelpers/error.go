package testhelpers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func ReportErrors(t *testing.T, name string, errs []error) {
	for _, err := range errs {
		t.Errorf("%s: Validator failed: %s", name, err.Error())
	}
	assert.Len(t, errs, 0)
}
