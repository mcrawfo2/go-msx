// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package rbac

import (
	"context"
	"testing"

	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/integration/usermanagement"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/stretchr/testify/assert"
)

func TestHasTenant_Explicit(t *testing.T) {
	tenantId, _ := types.NewUUID()
	tenantBadId, _ := types.NewUUID()

	userManagementApi := new(usermanagement.MockUserManagement)
	ctx := usermanagement.ContextWithIntegration(context.Background(), userManagementApi)
	setMockTokenDetailsProvider(ctx, []types.UUID{tenantId}, []string{})
	mockedUsermanagement := usermanagement.IntegrationFromContext(ctx).(*usermanagement.MockUserManagement)

	var err error
	err = HasTenant(ctx, tenantId)
	assert.NoError(t, err)

	mockedUsermanagement.On("GetTenantHierarchyAncestors", tenantBadId).Return(nil, []types.UUID{}, nil)
	err = HasTenant(ctx, tenantBadId)
	assert.Error(t, err)
	assert.Equal(t, ErrUserDoesNotHaveTenantAccess, err)

	rootTenantId, _ := types.NewUUID()
	childTenantId, _ := types.NewUUID()
	ancestors := []types.UUID{rootTenantId, tenantId}
	mockedUsermanagement.On("GetTenantHierarchyAncestors", childTenantId).Return(nil, ancestors, nil)
	err = HasAccessToTenant(ctx, childTenantId)
	assert.NoError(t, err)
}

func TestHasTenant_Implicit(t *testing.T) {
	tenantId, _ := types.NewUUID()
	ctx := context.Background()
	setMockTokenDetailsProvider(ctx, []types.UUID{}, []string{PermissionAccessAllTenants})

	var err error
	err = HasTenant(ctx, tenantId)
	assert.NoError(t, err)
}

func TestHasAccessToTenant(t *testing.T) {
	tenantId, _ := types.NewUUID()
	rootTenantId, _ := types.NewUUID()

	//build mocked context
	userManagementApi := new(usermanagement.MockUserManagement)
	ctx := usermanagement.ContextWithIntegration(context.Background(), userManagementApi)
	setMockTokenDetailsProvider(ctx, []types.UUID{tenantId}, []string{})
	mockedUsermanagement := usermanagement.IntegrationFromContext(ctx).(*usermanagement.MockUserManagement)

	//validate
	err := HasAccessToTenant(ctx, tenantId)
	assert.NoError(t, err)

	tenantBadId, _ := types.NewUUID()
	mockedUsermanagement.On("GetTenantHierarchyAncestors", tenantBadId).Return(nil, []types.UUID{rootTenantId}, nil)
	notValid := HasAccessToTenant(ctx, tenantBadId)
	assert.Error(t, notValid, ErrTenantDoesNotExist)

	mockedUsermanagement.AssertExpectations(t)
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

	//build mocked context
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
