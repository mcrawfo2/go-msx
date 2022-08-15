// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package rbac

import "context"

type contextKey string

const (
	contextKeyTenantHierarchy       = contextKey("TenantHierarchy")
	contextKeyTenantHierarchyLoader = contextKey("TenantHierarchyLoader")
)

func ContextWithTenantHierarchy(ctx context.Context, api TenantHierarchyApi) context.Context {
	return context.WithValue(ctx, contextKeyTenantHierarchy, api)
}

func TenantHierarchyFromContext(ctx context.Context) TenantHierarchyApi {
	api, _ := ctx.Value(contextKeyTenantHierarchy).(TenantHierarchyApi)
	return api
}

func ContextWithTenantHierarchyLoader(ctx context.Context, api TenantHierarchyApi) context.Context {
	return context.WithValue(ctx, contextKeyTenantHierarchyLoader, api)
}

func TenantHierarchyLoaderFromContext(ctx context.Context) TenantHierarchyApi {
	api, _ := ctx.Value(contextKeyTenantHierarchyLoader).(TenantHierarchyApi)
	return api
}
