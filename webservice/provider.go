package webservice

import "github.com/emicklei/go-restful"

type ServiceProvider interface {
	Actuate(webService *restful.WebService) error
}

type DocumentationProvider interface {
	Actuate(container *restful.Container, webService *restful.WebService) error
}

var (
	documentationProvider DocumentationProvider
	actuators             []ServiceProvider
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
