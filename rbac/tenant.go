// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

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

// HasTenant Validates if the user is associated with the tenant (i.e. the tenant is assigned to the
// user, or a descendent of an assigned tenant).
//
// Returns nil for a user with ACCESS_ALL_TENANTS, even for a non-existent tenant. (Use HasAccessToTenant to
// ensure a tenant exists)
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

	if userContextDetails.HasTenantId(tenantId) {
		return nil
	}

	// Fall back to querying the tenant hierarchy to see if an assigned/designated tenant is an
	// ancestor of the given tenantId.  This will grant access to descendent tenants that are created
	// during the user's session.
	tenantHierarchy, err := GetTenantHierarchyApi(ctx)
	if err != nil {
		return err
	}
	ancestors, err := tenantHierarchy.Ancestors(ctx, tenantId)
	if err != nil {
		return err
	}
	for _, id := range ancestors {
		for _, assignedId := range userContextDetails.Tenants {
			if id.Equals(assignedId) {
				return nil
			}
		}
	}

	logger.WithContext(ctx).
		WithError(ErrUserDoesNotHaveTenantAccess).
		Errorf("Tenant access check failed for tenantId %q", tenantId.String())
	return ErrUserDoesNotHaveTenantAccess
}

// HasAccessToTenant validates that the user has access to the tenant, and that the tenant is valid (the tenant exists).
//
// Returns nil if the tenant is accessible.  Returns ErrUserDoesNotHaveTenantAccess if the tenant is not
// accessible.  Returns ErrTenantDoesNotExist when the user has ACCESS_ALL_TENANTS and the tenant does not exist.
func HasAccessToTenant(ctx context.Context, tenantId types.UUID) error {

	//check if this is an PermissionAccessAllTenants account
	allTenants, _ := HasAccessAllTenants(ctx)

	if allTenants {
		return ValidateTenant(ctx, tenantId)
	}

	return HasTenant(ctx, tenantId)
}

// THIS IS NOT AN ACCESS CONTROL CHECK!
//
// ValidateTenant Checks if the tenant id is a valid tenant. See securitytest.MockedTenantValidation for test mocking
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
