package security

type TokenProvider interface {
	SecurityContextFromToken(token string) (*UserContext, error)
	TokenFromSecurityContext(userContext *UserContext) (token string, err error)
}

var tokenProvider TokenProvider

func SetTokenProvider(provider TokenProvider) {
	if provider != nil {
		tokenProvider = provider
	}
}

func NewUserContextFromToken(token string) (userContext *UserContext, err error) {
	return tokenProvider.SecurityContextFromToken(token)
}
