package securitytest

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/integration/usermanagement"
	"cto-github.cisco.com/NFV-BU/go-msx/rbac"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockedTenantValidation test helper to handle mock pass/fail rbac.ValidateTenant method
func MockedTenantValidation(ctx context.Context, pass bool) (context.Context, *usermanagement.MockUserManagement) {
	mockedUserManagementIntegration := new(usermanagement.MockUserManagement)
	ctx = usermanagement.ContextWithIntegration(ctx, mockedUserManagementIntegration)

	rootUUID, _ := uuid.NewRandom()
	var mockedRootResponse = integration.MsxResponse{
		StatusCode: 200,
		BodyString: rootUUID.String(),
	}

	parentUUID, _ := uuid.NewRandom()
	var mockedParentResponse = integration.MsxResponse{
		StatusCode: 200,
		BodyString: parentUUID.String(),
	}

	//mocking conditions
	mockedUserManagementIntegration.On("GetTenantHierarchyRoot").Return(&mockedRootResponse, nil)
	if pass {
		mockedParentResponse.BodyString = parentUUID.String()
		mockedUserManagementIntegration.On("GetTenantHierarchyParent", mock.Anything).Return(&mockedParentResponse, nil)
	} else {
		mockedUserManagementIntegration.On("GetTenantHierarchyParent", mock.Anything).Return(nil, rbac.ErrTenantDoesNotExist)
	}
	return ctx, mockedUserManagementIntegration
}
