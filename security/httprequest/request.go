// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package httprequest

import (
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strings"
)

const headerAuthorization = "Authorization"
const headerValuePrefixBearer = "Bearer "
const headerSslCert = "X-Ssl-Cert"

var logger = log.NewLogger("msx.security.httprequest")
var ErrNotFound = errors.New("Authorization not found in request")
var ErrInvalidHeader = errors.New("Authorization header not in 'Bearer <token>' format")

// ExtractToken xtracts the JWT from the request
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

// InjectToken into request from the request context
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

// ExtractCertificate returns the PEM-encoded x509 certificate from the request
func ExtractCertificate(req *http.Request) (certificate string, err error) {
	certificates := req.Header[headerSslCert]
	if len(certificates) == 0 {
		return "", ErrNotFound
	}
	// Reverse the nginx encoding
	return url.QueryUnescape(certificates[0])
}
