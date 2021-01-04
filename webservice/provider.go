//go:generate mockery --name AuthenticationProvider --structname MockAuthenticationProvider --filename mock_AuthenticationProvider.go --inpackage
//go:generate mockery --name DocumentationProvider --structname MockDocumentationProvider --filename mock_DocumentationProvider.go --inpackage
//go:generate mockery --name ServiceProvider --structname MockServiceProvider --filename mock_ServiceProvider.go --inpackage
//go:generate mockery --name HttpHandler --structname MockHttpHandler --filename mock_HttpHandler.go --inpackage

package webservice

import (
	"github.com/emicklei/go-restful"
	"net/http"
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
	Path string
	File string
}

func (a StaticAlias) Alias(path string) string {
	p := strings.TrimSuffix(path, a.Path)
	return p + a.File
}
