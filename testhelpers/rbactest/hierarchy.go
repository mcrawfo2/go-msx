// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package rbactest

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/rbac"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/securitytest"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
)

type MockTenantHierarchy struct {
	RootId  types.UUID
	Parents map[string]string
}

func (m MockTenantHierarchy) Parent(ctx context.Context, tenantId types.UUID) (types.UUID, error) {
	if parentId, ok := m.Parents[tenantId.String()]; !ok {
		return nil, errors.Wrap(rbac.ErrTenantDoesNotExist, tenantId.String())
	} else {
		return types.MustParseUUID(parentId), nil
	}
}

func (m MockTenantHierarchy) Ancestors(ctx context.Context, tenantId types.UUID) ([]types.UUID, error) {
	var ancestors []types.UUID

	rootId, err := m.Root(ctx)
	if err != nil {
		return nil, err
	}

	if tenantId.Equals(rootId) {
		return nil, nil
	}

	for !tenantId.Equals(rootId) {
		parentId, err := m.Parent(ctx, tenantId)
		if err != nil {
			return nil, err
		}
		ancestors = append(ancestors, parentId)
		tenantId = parentId
	}

	return ancestors, nil
}

func (m MockTenantHierarchy) Root(ctx context.Context) (types.UUID, error) {
	if m.RootId == nil {
		return nil, errors.Wrap(rbac.ErrTenantDoesNotExist, "Root Tenant")
	} else {
		return m.RootId, nil
	}
}

func NewMockTenantHierarchy(parents map[string]string) *MockTenantHierarchy {
	if parents == nil {
		parents = map[string]string{
			securitytest.DefaultTenantId.String(): securitytest.DefaultRootTenantId.String(),
		}
	}
	return &MockTenantHierarchy{
		RootId:  securitytest.DefaultRootTenantId,
		Parents: parents,
	}
}
