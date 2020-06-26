package rbac

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHasTenant_Explicit(t *testing.T) {
	tenantId, _ := types.NewUUID()
	tenantBadId, _ := types.NewUUID()
	ctx := context.Background()
	setMockTokenDetailsProvider(ctx, []types.UUID{tenantId}, []string{})

	var err error
	err = HasTenant(ctx, tenantId)
	assert.NoError(t, err)

	err = HasTenant(ctx, tenantBadId)
	assert.Error(t, err)
	assert.Equal(t, ErrUserDoesNotHaveTenantAccess, err)
}

func TestHasTenant_Implicit(t *testing.T) {
	tenantId, _ := types.NewUUID()
	ctx := context.Background()
	setMockTokenDetailsProvider(ctx, []types.UUID{}, []string{PermissionAccessAllTenants})

	var err error
	err = HasTenant(ctx, tenantId)
	assert.NoError(t, err)
}
