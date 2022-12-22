// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package securitytest

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/integration/usermanagement"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/rbac"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"net/http"
)

var logger = log.NewPackageLogger()

// MockedTenantValidation test helper to handle mock pass/fail rbac.ValidateTenant method
func MockedTenantValidation(ctx context.Context, pass bool) (context.Context, *usermanagement.MockUserManagement) {
	ctx, mockedUserManagementIntegration := TenantHierarchyInjector(ctx)

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

func TenantHierarchyInjector(ctx context.Context) (context.Context, *usermanagement.MockUserManagement) {
	mockedUserManagementIntegration := new(usermanagement.MockUserManagement)
	return usermanagement.ContextWithIntegration(ctx, mockedUserManagementIntegration), mockedUserManagementIntegration
}

type MockTenantHierarchyNode struct {
	TenantId types.UUID
	Children []*MockTenantHierarchyNode
}

func (n *MockTenantHierarchyNode) DescendantIds() []types.UUID {
	var result []types.UUID
	for _, child := range n.Children {
		result = append(result, child.TenantId)
		result = append(result, child.DescendantIds()...)
	}
	return result
}

type MockTenantHierarchyVisitor struct {
	Ancestors []*MockTenantHierarchyNode
	Node      *MockTenantHierarchyNode
}

func (v MockTenantHierarchyVisitor) AncestorIds() []types.UUID {
	var result []types.UUID
	for i := len(v.Ancestors) - 1; i >= 0; i-- {
		result = append(result, v.Ancestors[i].TenantId)
	}
	return result
}

func (v MockTenantHierarchyVisitor) Visit(ctx context.Context) {
	TenantHierarchyDescendantsInjector(ctx, v.Node.TenantId, v.Node.DescendantIds())
	ancestorIds := v.AncestorIds()
	TenantHierarchyAncestorsInjector(ctx, v.Node.TenantId, ancestorIds)
	if len(ancestorIds) > 0 {
		TenantHierarchyParentInjector(ctx, v.Node.TenantId, ancestorIds[0])
	} else {
		TenantHierarchyParentInjector(ctx, v.Node.TenantId, nil)
	}

	ancestors := append([]*MockTenantHierarchyNode{}, v.Ancestors...)
	ancestors = append(ancestors, v.Node)
	for _, node := range v.Node.Children {
		visitor := MockTenantHierarchyVisitor{
			Ancestors: ancestors,
			Node:      node,
		}
		visitor.Visit(ctx)
	}
}

func TenantHierarchyNodesInjector(ctx context.Context, root *MockTenantHierarchyNode) context.Context {
	mockUserManagement, ok := usermanagement.IntegrationFromContext(ctx).(*usermanagement.MockUserManagement)
	if !ok {
		ctx, mockUserManagement = TenantHierarchyInjector(ctx)
	}

	testhelpers.UnsetMockExpectedCallsByMethod(&mockUserManagement.Mock, "GetTenantHierarchyDescendants")
	testhelpers.UnsetMockExpectedCallsByMethod(&mockUserManagement.Mock, "GetTenantHierarchyAncestors")
	testhelpers.UnsetMockExpectedCallsByMethod(&mockUserManagement.Mock, "GetTenantHierarchyParent")
	testhelpers.UnsetMockExpectedCallsByMethod(&mockUserManagement.Mock, "GetTenantHierarchyRoot")

	TenantHierarchyRootInjector(ctx, root.TenantId)
	visitor := MockTenantHierarchyVisitor{
		Ancestors: nil,
		Node:      root,
	}
	visitor.Visit(ctx)

	return ctx
}

func TenantHierarchyDescendantsInjector(ctx context.Context, tenantId types.UUID, descendants []types.UUID) context.Context {
	mockUserManagement, ok := usermanagement.IntegrationFromContext(ctx).(*usermanagement.MockUserManagement)
	if !ok {
		ctx, mockUserManagement = TenantHierarchyInjector(ctx)
	}

	mockUserManagement.
		On("GetTenantHierarchyDescendants", tenantId).
		Run(func(args mock.Arguments) {
			logger.WithContext(ctx).WithField("descendants", descendants).Infof("%+v", args)
		}).
		Return(nil, descendants, nil)
	return ctx
}

func TenantHierarchyAncestorsInjector(ctx context.Context, tenantId types.UUID, ancestors []types.UUID) context.Context {
	mockUserManagement, ok := usermanagement.IntegrationFromContext(ctx).(*usermanagement.MockUserManagement)
	if !ok {
		ctx, mockUserManagement = TenantHierarchyInjector(ctx)
	}

	mockUserManagement.
		On("GetTenantHierarchyAncestors", tenantId).
		Run(func(args mock.Arguments) {
			logger.WithContext(ctx).WithField("ancestors", ancestors).Infof("GetTenantHierarchyAncestors: %+v", args)
		}).
		Return(nil, ancestors, nil)
	return ctx
}

func TenantHierarchyRootInjector(ctx context.Context, rootTenantId types.UUID) context.Context {
	mockUserManagement, ok := usermanagement.IntegrationFromContext(ctx).(*usermanagement.MockUserManagement)
	if !ok {
		ctx, mockUserManagement = TenantHierarchyInjector(ctx)
	}

	response := &integration.MsxResponse{
		StatusCode: http.StatusOK,
		BodyString: rootTenantId.String(),
	}

	mockUserManagement.
		On("GetTenantHierarchyRoot").
		Run(func(args mock.Arguments) {
			logger.WithContext(ctx).WithField("root", rootTenantId).Infof("%+v", args)
		}).
		Return(response, nil)

	return ctx
}

func TenantHierarchyParentInjector(ctx context.Context, tenantId, parentTenantId types.UUID) context.Context {
	mockUserManagement, ok := usermanagement.IntegrationFromContext(ctx).(*usermanagement.MockUserManagement)
	if !ok {
		ctx, mockUserManagement = TenantHierarchyInjector(ctx)
	}

	if parentTenantId == nil {
		response2 := &integration.MsxResponse{
			StatusCode: http.StatusOK,
		}

		mockUserManagement.
			On("GetTenantHierarchyParent", tenantId).
			Run(func(args mock.Arguments) {
				logger.WithContext(ctx).WithField("parent", "none").Infof("%+v", args)
			}).
			Return(response2, nil)

		return ctx
	}

	response := &integration.MsxResponse{
		StatusCode: http.StatusOK,
		BodyString: parentTenantId.String(),
	}

	mockUserManagement.On("GetTenantHierarchyParent", tenantId).
		Run(func(args mock.Arguments) {
			logger.WithContext(ctx).WithField("parent", parentTenantId).Infof("%+v", args)
		}).
		Return(response, nil)
	return ctx
}
