package security

import "context"

type UserContext struct {
	UserName string   `json:"user_name"`
	Roles    []string `json:"roles"`
	TenantId string   `json:"tenant_id"`
	Scopes   []string `json:"scope"`
	Token    string   `json:"-"`
}

type securityContextKey int

const (
	contextKeyUserContext securityContextKey = iota
)

func ContextWithUserContext(ctx context.Context, userContext *UserContext) context.Context {
	return context.WithValue(ctx, contextKeyUserContext, userContext)
}

func UserContextFromContext(ctx context.Context) *UserContext {
	userContextInterface := ctx.Value(contextKeyUserContext)
	if userContextInterface == nil {
		return nil
	}
	return userContextInterface.(*UserContext)
}
