// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package testhelpers

import (
	"context"
	"testing"
)

type Testable interface {
	Test(t *testing.T)
}

type TestFunc func(t *testing.T)

func (f TestFunc) Test(t *testing.T) {
	f(t)
}

type ServiceTestFunc func(t *testing.T, ctx context.Context)

type ServiceTestable interface {
	Test(t *testing.T, fn ServiceTestFunc)
}
