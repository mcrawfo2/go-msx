package security

import "context"

type TokenProvider interface {
	UserContextFromToken(ctx context.Context, token string) (*UserContext, error)
}

var tokenProvider TokenProvider

func SetTokenProvider(provider TokenProvider) {
	if provider != nil {
		tokenProvider = provider
	}
}

func NewUserContextFromToken(ctx context.Context, token string) (userContext *UserContext, err error) {
	return tokenProvider.UserContextFromToken(ctx, token)
}
