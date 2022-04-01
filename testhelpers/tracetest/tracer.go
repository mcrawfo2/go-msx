// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package tracetest

import (
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	mocktracer "cto-github.cisco.com/NFV-BU/go-msx/trace/mock"
)

func RecordTracing() *mocktracer.MockTracer {
	t := mocktracer.NewMockTracer()
	trace.SetTracer(t)
	return t
}
