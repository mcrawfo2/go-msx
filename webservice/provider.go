package webservice

import (
	"context"
	"github.com/emicklei/go-restful"
)

type ServiceProvider interface {
	Actuate(webService *restful.WebService) error
}

type DocumentationProvider interface {
	Actuate(container *restful.Container, webService *restful.WebService) error
}

type SecurityProvider interface {
	// Ensures user is logged in
	Authentication(request *restful.Request) error
}

type InjectorFunction func(ctx context.Context) context.Context

var (
	documentationProvider DocumentationProvider
	actuators             []ServiceProvider
	securityProvider      SecurityProvider
)


func RegisterDocumentationProvider(provider DocumentationProvider) {
	if provider != nil {
		documentationProvider = provider
	}
}

func RegisterActuator(provider ServiceProvider) {
	if provider != nil {
		actuators = append(actuators, provider)
	}
}

func RegisterSecurityProvider(provider SecurityProvider) {
	if provider != nil {
		securityProvider = provider
	}
}
