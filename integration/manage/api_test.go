// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package manage

import "testing"

func Test_Implementations(t *testing.T) {
	// Ensure MockManage is up to date
	var _ Api = new(MockManage)
}
