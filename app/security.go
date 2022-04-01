// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package app

import (
	"cto-github.cisco.com/NFV-BU/go-msx/security/certprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/security/idmdetailsprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/security/jwttokenprovider"
)

func init() {
	OnEvent(EventConfigure, PhaseAfter, jwttokenprovider.RegisterTokenProvider)
	OnEvent(EventConfigure, PhaseAfter, idmdetailsprovider.RegisterTokenDetailsProvider)
	OnEvent(EventConfigure, PhaseAfter, certprovider.RegisterCertificateProvider)
}
