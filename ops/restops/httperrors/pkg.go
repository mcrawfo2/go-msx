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
