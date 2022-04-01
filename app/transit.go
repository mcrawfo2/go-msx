// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package app

import (
	"cto-github.cisco.com/NFV-BU/go-msx/transit"
	"cto-github.cisco.com/NFV-BU/go-msx/transit/vaultprovider"
)

func init() {
	OnEvent(EventConfigure, PhaseAfter, transit.ConfigureEncrypterFactory)
	OnEvent(EventConfigure, PhaseAfter, vaultprovider.RegisterVaultTransitProvider)
}
