package security

import "context"

type securityContextKey int

const (
	contextDefaultUserName                    = "system"
	contextKeyUserContext  securityContextKey = iota
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

func UserNameFromContext(ctx context.Context) string {
	user := UserContextFromContext(ctx)
	userName := contextDefaultUserName
	if user != nil {
		userName = user.UserName
	}
	return userName
}
