// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package rbac

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

var testRootTenantId = types.MustParseUUID("215927cd-95b9-4e21-b29f-ef9bdaff9cbf")
var testChildTenantId = types.MustParseUUID("215927cd-95b9-4e21-b29f-ef9bdaff9cbd")

func newMockTenantHierarchyLoader(ctx context.Context) (*MockTenantHierarchyApi, context.Context) {
	loader := new(MockTenantHierarchyApi)

	loader.
		On("Root", mock.AnythingOfType("*context.valueCtx")).
		Return(testRootTenantId, nil)

	loader.
		On("Parent", mock.AnythingOfType("*context.valueCtx"), testRootTenantId).
		Return(nil, nil)

	loader.
		On("Ancestors", mock.AnythingOfType("*context.valueCtx"), testRootTenantId).
		Return(nil, nil)

	return loader, ContextWithTenantHierarchyLoader(ctx, loader)
}

func TestTenantHierarchyCache_Parent(t *testing.T) {
	tests := []struct {
		name       string
		tenantId   types.UUID
		loader     func(api *MockTenantHierarchyApi)
		wantResult types.UUID
		wantErr    bool
	}{
		{
			name:       "Root",
			tenantId:   testRootTenantId,
			wantResult: nil,
		},
		{
			name:     "Child",
			tenantId: testChildTenantId,
			loader: func(api *MockTenantHierarchyApi) {
				api.
					On("Ancestors", mock.AnythingOfType("*context.valueCtx"), testChildTenantId).
					Return([]types.UUID{testRootTenantId}, nil)
			},
			wantResult: testRootTenantId,
		},
		{
			name:     "ChildFailure",
			tenantId: testChildTenantId,
			loader: func(api *MockTenantHierarchyApi) {
				api.
					On("Ancestors", mock.AnythingOfType("*context.valueCtx"), testChildTenantId).
					Return(nil, errTenantNotFoundInCache)

			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader, ctx := newMockTenantHierarchyLoader(context.Background())

			if tt.loader != nil {
				tt.loader(loader)
			}

			cache, err := newTenantHierarchyCache(ctx)
			assert.NoError(t, err)

			gotResult, err := cache.Parent(ctx, tt.tenantId)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Parent() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestTenantHierarchyCache_Root(t *testing.T) {
	tests := []struct {
		name       string
		loader     func(api *MockTenantHierarchyApi)
		wantResult types.UUID
		wantErr    bool
	}{
		{
			name:       "Root",
			wantResult: testRootTenantId,
		},
		{
			name: "RootFailure",
			loader: func(api *MockTenantHierarchyApi) {
				api.ExpectedCalls = nil
				api.
					On("Root", mock.AnythingOfType("*context.valueCtx")).
					Return(nil, errTenantNotFoundInCache)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader, ctx := newMockTenantHierarchyLoader(context.Background())

			if tt.loader != nil {
				tt.loader(loader)
			}

			cache, err := newTenantHierarchyCache(ctx)
			assert.NoError(t, err)

			gotResult, err := cache.Root(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Root() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Root() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestTenantHierarchyCache_Ancestors(t *testing.T) {
	tests := []struct {
		name       string
		tenantId   types.UUID
		loader     func(api *MockTenantHierarchyApi)
		wantResult []types.UUID
		wantErr    bool
	}{
		{
			name:       "Root",
			tenantId:   testRootTenantId,
			wantResult: nil,
		},
		{
			name:     "Child",
			tenantId: testChildTenantId,
			loader: func(api *MockTenantHierarchyApi) {
				api.
					On("Ancestors", mock.AnythingOfType("*context.valueCtx"), testChildTenantId).
					Return([]types.UUID{testRootTenantId}, nil)
			},
			wantResult: []types.UUID{testRootTenantId},
		},
		{
			name:     "ChildFailure",
			tenantId: testChildTenantId,
			loader: func(api *MockTenantHierarchyApi) {
				api.
					On("Ancestors", mock.AnythingOfType("*context.valueCtx"), testChildTenantId).
					Return(nil, errTenantNotFoundInCache)

			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader, ctx := newMockTenantHierarchyLoader(context.Background())

			if tt.loader != nil {
				tt.loader(loader)
			}

			cache, err := newTenantHierarchyCache(ctx)
			assert.NoError(t, err)

			gotResult, err := cache.Ancestors(ctx, tt.tenantId)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ancestors() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Error("Ancestors() unexpected result")
				t.Errorf(testhelpers.Diff(tt.wantResult, gotResult))
			}
		})
	}
}
