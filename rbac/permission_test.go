package rbac

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func mockUserContextDetails(tenants []types.UUID, permissions []string) *security.UserContextDetails {
	return &security.UserContextDetails{
		Permissions: permissions,
		Tenants:     tenants,
	}
}

func setMockTokenDetailsProvider(ctx context.Context, tenants []types.UUID, permissions []string) {
	userContextDetails := mockUserContextDetails(tenants, permissions)
	mockTokenDetailsProvider := new(security.MockTokenDetailsProvider)
	mockTokenDetailsProvider.
		On("TokenDetails", ctx).
		Return(userContextDetails, nil)
	security.SetTokenDetailsProvider(mockTokenDetailsProvider)
}

func TestHasPermission_Explicit(t *testing.T) {
	ctx := context.Background()
	setMockTokenDetailsProvider(ctx, nil, []string{"SOME_PERMISSION"})

	var err error
	err = HasPermission(ctx, []string{"SOME_PERMISSION"})
	assert.NoError(t, err)

	err = HasPermission(ctx, []string{"OTHER_PERMISSION"})
	assert.Error(t, err)
	assert.Equal(t, ErrUserDosNotHavePermission, err)
}

func TestHasPermission_Implicit(t *testing.T) {
	ctx := context.Background()
	setMockTokenDetailsProvider(ctx, nil, []string{PermissionIsApiAdmin})

	var err error
	err = HasPermission(ctx, []string{"SOME_PERMISSION"})
	assert.NoError(t, err)

	err = HasPermission(ctx, []string{"OTHER_PERMISSION"})
	assert.NoError(t, err)
}

func TestHasAccessAllTenants(t *testing.T) {
	ctx := context.Background()
	setMockTokenDetailsProvider(ctx, nil, []string{PermissionAccessAllTenants})

	access, err := HasAccessAllTenants(ctx)
	assert.NoError(t, err)
	assert.True(t, access)
}

func TestHasAccessAllTenants_False(t *testing.T) {
	ctx := context.Background()
	setMockTokenDetailsProvider(ctx, nil, []string{PermissionAccessAllTenants})

	access, err := HasAccessAllTenants(ctx)
	assert.NoError(t, err)
	assert.True(t, access)
}
