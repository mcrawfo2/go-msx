package rbac

import "context"

type contextKey string

const (
	contextKeyTenantHierarchy = contextKey("TenantHierarchy")
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
