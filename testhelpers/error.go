// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package testhelpers

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func ReportErrors(t *testing.T, name string, errs []error) {
	for _, err := range errs {
		assert.Fail(t, err.Error(), "Failed %s validator", name)
	}
}

func CheckErr(got []interface{}, hasErr, wantErr bool) error {
	if wantErr {
		if len(got) == 0 {
			return errors.Errorf("Wanted Error, No values returned")
		} else if !hasErr {
			return errors.Errorf("Wanted Error, Method not flagged as returning error")
		} else if nil == got[len(got)-1] {
			return errors.Errorf("Wanted Error, No error returned")
		} else if _, ok := got[len(got)-1].(error); !ok {
			return errors.Errorf("Wanted Error, Non-error value returned")
		}
	} else if hasErr {
		if len(got) == 0 {
			return errors.Errorf("Method flagged as returning error, No values returned")
		} else if err, ok := got[len(got)-1].(error); ok && err != nil {
			return errors.Errorf("Unwanted Error, Error returned:\n%s", Dump(err))
		} else if nil == got[len(got)-1] {
			// No error returned
		} else if _, ok := got[len(got)-1].(error); !ok {
			return errors.Errorf("Non-error value returned")
		}
	}

	return nil
}
