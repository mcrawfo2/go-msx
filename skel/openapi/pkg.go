// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package openapi

import (
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/skel"
)

var logger = log.NewPackageLogger()

func init() {
	skel.AddTarget("generate-webservices", "Create web services from swagger manifest", GenerateDomainOpenApi)
	skel.AddTarget("generate-domain-openapi", "Create domains from OpenAPI 3.0 manifest", GenerateDomainOpenApi)
}
