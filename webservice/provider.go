package webservice

import (
	"github.com/emicklei/go-restful"
)

type ServiceProvider interface {
	Actuate(webService *restful.WebService) error
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
