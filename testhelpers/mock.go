// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package testhelpers

import (
	"github.com/stretchr/testify/mock"
)

func MockExpectedCallByMethod(m *mock.Mock, methodName string) (call *mock.Call) {
	for _, c := range m.ExpectedCalls {
		if c.Method == methodName {
			call = c
		}
	}

	return
}

func UnsetMockExpectedCallsByMethod(m *mock.Mock, methodName string) (count int) {
	expectedCalls := append([]*mock.Call{}, m.ExpectedCalls...)
	for _, c := range expectedCalls {
		if c.Method == methodName {
			c.Unset()
			count++
		}
	}

	return
}
