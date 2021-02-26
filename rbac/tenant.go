package rbac

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration/usermanagement"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
	"strings"
	"sync/atomic"
)

var logger = log.NewLogger("msx.rbac")
var ErrUserDoesNotHaveTenantAccess = errors.New("User does not have access to the tenant")
var ErrTenantDoesNotExist = errors.New("Tenant does not exist.")
var rootTenantId atomic.Value

// HasTenant Validates if the user is associated with the tenant
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

// HasAccessToTenant validates that the user has access to the tenant.
// Error returned if user does not have access to tenant (or child/descendent of tenant), nil otherwise
func HasAccessToTenant(ctx context.Context, tenantId types.UUID) error {

	//check if this is an PermissionAccessAllTenants account
	allTenants, _ := HasAccessAllTenants(ctx)

	if allTenants {
		return ValidateTenant(ctx, tenantId)
	}

	return HasTenant(ctx, tenantId)
}

// ValidateTenant Checks if the tenant id is a valid tenant. See securitytest.MockedTenantValidation for test mocking
// Must NOT be used for checking access control.
func ValidateTenant(ctx context.Context, tenantId types.UUID) error {
	userManagementApi, err := usermanagement.NewIntegration(ctx)
	if err != nil {
		return errors.Wrap(err, "Error initializing usermanagement.")
	}

	//all tenants belong to the root tenantid, or have a parent
	rootTenant, err := GetRootTenant(ctx)
	if err != nil {
		return errors.Wrap(err, "Error getting Root TenantId.")
	}

	if rootTenant.Equals(tenantId) {
		return nil
	}

	//check if the token has a parent
	parentTenant, err := userManagementApi.GetTenantHierarchyParent(tenantId)
	if err != nil {
		return errors.Wrap(err, "Error getting Parent TenantId.")
	}
	if len(strings.TrimSpace(parentTenant.BodyString)) != 0 {
		return nil
	}

	return ErrTenantDoesNotExist
}

// GetRootTenant fetches and locally caches the root tenant id
// Services should generally avoid using this method and not depend on the value of the root tenant ID
func GetRootTenant(ctx context.Context) (types.UUID, error) {
	//fetch the root token; it's system wide so we only need to do this once
	if rootTenantId.Load() == nil {
		userManagementApi, err := usermanagement.NewIntegration(ctx)

		if err != nil {
			return nil, err
		}

		resp, err := userManagementApi.GetTenantHierarchyRoot()
		if err == nil {
			rootTenantId.Store(types.MustParseUUID(resp.BodyString))
		} else {
			return nil, err
		}
	}

	return rootTenantId.Load().(types.UUID), nil
}
