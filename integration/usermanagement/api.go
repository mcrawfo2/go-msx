// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package usermanagement

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration/auth"
	"cto-github.cisco.com/NFV-BU/go-msx/integration/idm"
	"cto-github.cisco.com/NFV-BU/go-msx/integration/secrets"
)

//go:generate mockery --inpackage --name=Api --structname=MockUserManagement
// Deprecated: use secrets.Api, idm.Api or auth.Api directly
type Api interface {
	auth.Api
	idm.Api
	secrets.Api
}
