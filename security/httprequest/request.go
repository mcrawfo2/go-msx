package httprequest

import (
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

const headerAuthorization = "Authorization"
const headerValuePrefixBearer = "Bearer "

var logger = log.NewLogger("msx.security.httprequest")
var ErrNotFound = errors.New("Authorization not found in request")
var ErrInvalidHeader = errors.New("Authorization header not in 'Bearer <token>' format")

// Extract the JWT from the request
func ExtractToken(req *http.Request) (token string, err error) {
	authorizationHeaderValue := req.Header.Get(headerAuthorization)
	if authorizationHeaderValue == "" {
		return "", ErrNotFound
	}
	if !strings.HasPrefix(strings.ToLower(authorizationHeaderValue), strings.ToLower(headerValuePrefixBearer)) {
		return "", ErrInvalidHeader
	}

	parts := strings.SplitN(authorizationHeaderValue, " ", 2)
	return strings.TrimSpace(parts[1]), nil
}

// Inject token into request from the request context
func InjectToken(req *http.Request) {
	userContext := security.UserContextFromContext(req.Context())
	token := userContext.Token
	if token == "" {
		logger.Warn("Token required but not set")
		return
	}

	authorizationHeaderValue := headerValuePrefixBearer + token
	req.Header.Set(headerAuthorization, authorizationHeaderValue)
}
