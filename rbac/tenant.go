package rbac

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
)

var logger = log.NewLogger("msx.rbac")
var ErrUserDoesNotHaveTenantAccess = errors.New("User does not have access to the tenant")

func HasTenant(ctx context.Context, tenantId types.UUID) error {
	logger.WithContext(ctx).Debugf("Verifying tenant access for tenantId %q", tenantId.String())

	if allTenants, err := HasAccessAllTenants(ctx); err != nil {
		return err
	} else if allTenants {
		return nil
	}

	userContext := security.UserContextFromContext(ctx)
	if userContext.TenantId.Equals(tenantId) {
		return nil
	}

	userContextDetails, err := security.NewUserContextDetails(ctx)
	if err != nil {
		return err
	}

	if !userContextDetails.HasTenantId(tenantId) {
		// TODO: Check remotely

		logger.WithContext(ctx).
			WithError(ErrUserDoesNotHaveTenantAccess).
			Errorf("Tenant access check failed for tenantId %q", tenantId.String())
		return ErrUserDoesNotHaveTenantAccess
	}

	return nil
}
