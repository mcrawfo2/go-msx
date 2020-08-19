package webservice

import (
	"github.com/emicklei/go-restful"
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

type StaticAlias struct {
	Path string
	File string
}

func (a StaticAlias) Alias(path string) string {
	p := strings.TrimSuffix(path, a.Path)
	return p + a.File
}
