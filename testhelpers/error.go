// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

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
