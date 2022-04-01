// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package app

import (
	"cto-github.cisco.com/NFV-BU/go-msx/certificate/cacheprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/certificate/fileprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/certificate/vaultprovider"
)

func init() {
	OnEvent(EventConfigure, PhaseAfter, fileprovider.RegisterFactory)
	OnEvent(EventConfigure, PhaseAfter, vaultprovider.RegisterFactory)
	OnEvent(EventConfigure, PhaseAfter, cacheprovider.RegisterFactory)
}
