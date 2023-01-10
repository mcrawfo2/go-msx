// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import "net/http"

type Middleware func(next http.Handler) http.Handler

type Middlewares []Middleware

func (m Middlewares) Compose(final http.Handler) http.Handler {
	here := final
	for _, mw := range m {
		here = mw(here)
	}
	return here
}
