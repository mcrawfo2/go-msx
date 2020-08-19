package envprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/adminprovider"
	"github.com/emicklei/go-restful"
	"strings"
)

const (
	endpointName = "env"
)

type CachingProvider interface {
	Cache() map[string]string
}

type Property struct {
	Value interface{} `json:"value"`
}

type PropertySource struct {
	Name       string              `json:"name"`
	Properties map[string]Property `json:"properties"`
}

type Report struct {
	ActiveProfiles  []string         `json:"activeProfiles"`
	PropertySources []PropertySource `json:"propertySources"`
}

type Provider struct{}

func (h Provider) EndpointName() string {
	return endpointName
}

func (h Provider) Actuate(webService *restful.WebService) error {
	webService.Consumes(restful.MIME_JSON)
	webService.Produces(restful.MIME_JSON)

	webService.Path(webService.RootPath() + "/admin/" + endpointName)

	// Unsecured routes for info
	webService.Route(webService.GET("").
		Operation("admin.env").
		To(adminprovider.RawAdminController(h.report)).
		Do(webservice.Returns200))

	return nil
}

func (h Provider) report(req *restful.Request) (body interface{}, err error) {
	profile, err := config.FromContext(req.Request.Context()).StringOr("profile", "")
	profiles := []string{}
	if err != nil && profile != "" {
		profiles = []string{profile}
	}

	report := &Report{
		ActiveProfiles:  profiles,
		PropertySources: []PropertySource{},
	}

	for _, provider := range config.FromContext(req.Request.Context()).Providers {
		propertySource := PropertySource{
			Name: provider.Description(),
		}

		if cachingProvider, ok := provider.(CachingProvider); ok {
			propertySource.Properties = make(map[string]Property)
			for k, v := range cachingProvider.Cache() {
				if h.isSecret(k) {
					propertySource.Properties[k] = Property{"*****"}
				} else {
					propertySource.Properties[k] = Property{v}
				}
			}
		}

		report.PropertySources = append([]PropertySource{propertySource}, report.PropertySources...)
	}

	return report, nil
}

func (h Provider) isSecret(key string) bool {
	key = strings.ToLower(key)
	return strings.Contains(key, "secret") ||
		strings.Contains(key, "password") ||
		strings.Contains(key, "token") ||
		strings.Contains(key, "credentials")
}

func RegisterProvider(ctx context.Context) error {
	server := webservice.WebServerFromContext(ctx)
	if server != nil {
		server.RegisterActuator(new(Provider))
		adminprovider.RegisterLink(endpointName, endpointName, false)
	}
	return nil
}
