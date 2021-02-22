package rbac

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/integration/usermanagement"
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

func Test_ValidateTenant(t *testing.T) {

	rootTenantId, _ := types.NewUUID()
	tenantId, _ := types.NewUUID()
	tenantBadId, _ := types.NewUUID()

	//mock responses
	parentTenantId, _ := types.NewUUID()

	rootTenantIdMsxResponse := &integration.MsxResponse{
		StatusCode: 200,
		BodyString: rootTenantId.String(),
	}

	parentTenantIdMsxResponse := &integration.MsxResponse{
		StatusCode: 200,
		BodyString: parentTenantId.String(),
	}

	badTenantIdMsxResponse := &integration.MsxResponse{
		StatusCode: 200,
		BodyString: "",
	}

	//build modked context
	ctx := context.Background()
	userManagementApi := new(usermanagement.MockUserManagement)
	ctx = usermanagement.ContextWithIntegration(context.Background(), userManagementApi)
	setMockTokenDetailsProvider(ctx, []types.UUID{rootTenantId}, []string{})

	mockedUsermanagement := usermanagement.IntegrationFromContext(ctx).(*usermanagement.MockUserManagement)
	mockedUsermanagement.On("GetTenantHierarchyRoot").Return(rootTenantIdMsxResponse, nil).Maybe()
	mockedUsermanagement.On("GetTenantHierarchyParent", tenantId).Return(parentTenantIdMsxResponse, nil).Maybe()

	//validate true
	validateErr := ValidateTenant(ctx, tenantId)
	assert.Nil(t, validateErr)

	//validate false
	mockedUsermanagement.On("GetTenantHierarchyParent", tenantBadId).Return(badTenantIdMsxResponse, nil).Maybe()
	validateErr = ValidateTenant(ctx, tenantBadId)
	assert.Error(t, validateErr)

	mockedUsermanagement.AssertExpectations(t)
}
