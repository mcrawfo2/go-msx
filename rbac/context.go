package rbac

import "context"

type contextKey int

const (
	contextKeyTenantHierarchy contextKey = iota
)

func ContextWithTenantHierarchy(ctx context.Context, api TenantHierarchyApi) context.Context {
	return context.WithValue(ctx, contextKeyTenantHierarchy, api)
}

func TenantHierarchyFromContext(ctx context.Context) TenantHierarchyApi {
	api, _ := ctx.Value(contextKeyTenantHierarchy).(TenantHierarchyApi)
	return api
}
