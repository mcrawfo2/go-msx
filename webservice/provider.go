// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

//go:generate mockery --name AuthenticationProvider --structname MockAuthenticationProvider --filename mock_AuthenticationProvider.go --inpackage
//go:generate mockery --name DocumentationProvider --structname MockDocumentationProvider --filename mock_DocumentationProvider.go --inpackage
//go:generate mockery --name ServiceProvider --structname MockServiceProvider --filename mock_ServiceProvider.go --inpackage
//go:generate mockery --name HttpHandler --structname MockHttpHandler --filename mock_HttpHandler.go --inpackage

package webservice

import (
	"github.com/emicklei/go-restful"
	"net/http"
	"path"
	"strings"
)

type ServiceProvider interface {
	Actuate(webService *restful.WebService) error
	EndpointName() string
}

type DocumentationProvider interface {
	Actuate(container *restful.Container, webService *restful.WebService) error
}

type AuthenticationProvider interface {
	// Ensures user is logged in
	Authenticate(request *restful.Request) error
}

type HttpHandler interface {
	http.Handler
}

type StaticAlias struct {
	ContextPath string
	Path        string
	File        string
}

func (a StaticAlias) Alias(originalPath string) string {
	patternParts := strings.Split(strings.Trim(a.Path, "/"), "/")

	originalPath = strings.TrimPrefix(originalPath, a.ContextPath)
	originalPathParts := strings.Split(strings.Trim(originalPath, "/"), "/")

	// consume originalPathParts
	originalIndex := 0
	for patternIndex := 0; patternIndex < len(patternParts); patternIndex++ {
		patternPart := patternParts[patternIndex]

		switch {
		case len(patternPart) == 0:
			continue

		case patternPart[0] == '{':
			if strings.Index(patternPart, ":*") != -1 {
				originalIndex = len(originalPathParts)
			}

		default:
			originalIndex++
		}
	}

	originalPathParts = originalPathParts[originalIndex:]
	subPath := path.Join(originalPathParts...)
	return path.Join(a.ContextPath, subPath, a.File)
}
