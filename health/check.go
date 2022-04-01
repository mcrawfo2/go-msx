// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package health

import "context"

type Check func(context.Context) CheckResult

var healthChecks = make(map[string]Check)

func RegisterCheck(name string, check Check) {
	healthChecks[name] = check
}
