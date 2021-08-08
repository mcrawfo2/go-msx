package rbac

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/repository"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

type TenantAccess struct {
	// ACCESS_ALL_TENANTS
	Unfiltered bool
	// Explicit Tenant Association
	TenantIds []types.UUID
	// Ancestor Tenants of TenantIds
	AncestorIds []types.UUID
}

func (t TenantAccess) ValidateTenantAccess(tenantId types.UUID) error {
	if t.Unfiltered {
		// User has access to all tenants
		return nil
	}

	// User has access to a set of tenants
	for _, hasTenantUuid := range t.TenantIds {
		if hasTenantUuid.Equals(tenantId) {
			return nil
		}
	}

	// No access to this task
	return repository.ErrNotFound
}

func (t TenantAccess) WithAncestors(ctx context.Context) (TenantAccess, error) {
	if t.Unfiltered {
		return t, nil
	}

	ancestorIds := make(map[string]types.UUID)

	hierarchy, err := GetTenantHierarchyApi(ctx)
	if err != nil {
		return TenantAccess{}, err
	}

	for _, tenantId := range t.TenantIds {
		tenantAncestors, err := hierarchy.Ancestors(ctx, tenantId)
		if err != nil {
			return TenantAccess{}, err
		}

		for _, tenantAncestorId := range tenantAncestors {
			ancestorIds[tenantAncestorId.String()] = tenantAncestorId
		}
	}

	for _, ancestorId := range ancestorIds {
		t.AncestorIds = append(t.AncestorIds, ancestorId)
	}

	return t, nil
}

func (t TenantAccess) ValidateAncestorAccess(tenantId types.UUID) error {
	if t.Unfiltered {
		// User has access to all tenants
		return nil
	}

	// User has access to a set of tenants
	for _, hasAncestorUuid := range t.AncestorIds {
		if hasAncestorUuid.Equals(tenantId) {
			return nil
		}
	}

	// No access to this task
	return repository.ErrNotFound
}

type TenantAccessWithAncestors struct {
	TenantAccess
	Ancestors []types.UUID
}

func NewTenantAccess(ctx context.Context) (result TenantAccess, err error) {
	if accessAllTenants, err := HasAccessAllTenants(ctx); err != nil {
		return result, err
	} else if accessAllTenants {
		result.Unfiltered = true
	}

	if !result.Unfiltered {
		userContextDetails, err := security.NewUserContextDetails(ctx)
		if err != nil {
			return TenantAccess{}, err
		}

		result.TenantIds = userContextDetails.Tenants
	}

	return
}

func NewTenantAccessForTenant(ctx context.Context, tenantId types.UUID) (result TenantAccess, err error) {
	if accessAllTenants, err := HasAccessAllTenants(ctx); err != nil {
		return result, err
	} else if accessAllTenants {
		result.Unfiltered = false
		result.TenantIds = []types.UUID{tenantId}
		return result, nil
	}

	userContextDetails, err := security.NewUserContextDetails(ctx)
	if err != nil {
		return TenantAccess{}, err
	}
	result.TenantIds = userContextDetails.Tenants

	err = result.ValidateTenantAccess(tenantId)
	if err != nil {
		result.TenantIds = nil
	} else {
		result.TenantIds = []types.UUID{tenantId}
	}

	return result, err
}

func NewTenantAccessForOptionalTenant(ctx context.Context, tenantId *types.UUID) (result TenantAccess, err error) {
	if tenantId == nil {
		return NewTenantAccess(ctx)
	} else {
		return NewTenantAccessForTenant(ctx, *tenantId)
	}
}

func NewTenantAccessForTenantSlice(ctx context.Context, tenantIds []types.UUID) (result TenantAccess, err error) {
	if accessAllTenants, err := HasAccessAllTenants(ctx); err != nil {
		return TenantAccess{}, err
	} else if accessAllTenants {
		result.TenantIds = tenantIds[:]
		return result, nil
	}

	userContextDetails, err := security.NewUserContextDetails(ctx)
	if err != nil {
		return TenantAccess{}, err
	}

	userTenantIds := make(types.UUIDSet)
	userTenantIds.Add(userContextDetails.Tenants...)

	for _, tenantId := range tenantIds {
		if !userTenantIds.Contains(tenantId) {
			return TenantAccess{}, ErrTenantDoesNotExist
		}
	}

	result.TenantIds = tenantIds[:]
	return result, err
}
