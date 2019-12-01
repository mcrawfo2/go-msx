package security

import "context"

var (
	defaultUserContext = &UserContext{
		UserName: "anonymous",
		Roles:    nil,
		TenantId: "",
		Scopes:   nil,
		Token:    "",
	}
)

type UserContext struct {
	UserName string   `json:"user_name"`
	Roles    []string `json:"roles"`
	TenantId string   `json:"tenant_id"`
	Scopes   []string `json:"scope"`
	Token    string   `json:"-"`
}

func (c *UserContext) Clone() *UserContext {
	return &UserContext{
		UserName: c.UserName,
		Roles:    c.Roles[:],
		TenantId: c.TenantId,
		Scopes:   c.Scopes[:],
		Token:    c.Token,
	}
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
		return defaultUserContext.Clone()
	}
	return userContextInterface.(*UserContext)
}
