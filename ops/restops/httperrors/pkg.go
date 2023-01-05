// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package httperrors

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/ops/restops"
	"cto-github.cisco.com/NFV-BU/go-msx/rbac"
	"cto-github.cisco.com/NFV-BU/go-msx/repository"
	"cto-github.cisco.com/NFV-BU/go-msx/schema/js"
	"net/http"
)

func init() {
	restops.SetMappedErrorStatusCode(js.ErrValidationFailed, http.StatusBadRequest)

	restops.SetMappedErrorStatusCode(ops.ErrMissingRequiredValue, http.StatusBadRequest)

	restops.SetMappedErrorStatusCode(rbac.ErrTenantDoesNotExist, http.StatusBadRequest)
	restops.SetMappedErrorStatusCode(rbac.ErrUserDoesNotHaveTenantAccess, http.StatusBadRequest)

	restops.SetMappedErrorStatusCode(repository.ErrAlreadyExists, http.StatusConflict)
	restops.SetMappedErrorStatusCode(repository.ErrNotFound, http.StatusNotFound)
}
