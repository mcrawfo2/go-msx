// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package app

import (
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/sanitize"
)

var logger = log.NewLogger("msx.app")

func init() {
	OnEvent(EventConfigure, PhaseAfter, sanitize.ConfigureSecretSanitizer)
}
