// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

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
