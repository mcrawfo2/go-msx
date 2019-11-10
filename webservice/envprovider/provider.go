package envprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/adminprovider"
	"encoding/json"
	"github.com/emicklei/go-restful"
)

const (
	providerName = "env"
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

func (h Provider) Actuate(webService *restful.WebService) error {
	webService.Consumes(restful.MIME_JSON)
	webService.Produces(restful.MIME_JSON)

	webService.Path(webService.RootPath() + "/admin/" + providerName)

	// Unsecured routes for info
	webService.Route(webService.GET("").
		Operation("admin.env").
		To(h.report).
		Do(webservice.Returns200))

	return nil
}

func (h Provider) report(req *restful.Request, resp *restful.Response) {
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
				propertySource.Properties[k] = Property{v}
			}
		}

		report.PropertySources = append(report.PropertySources, propertySource)
	}

	bodyBytes, _ := json.Marshal(report)

	resp.Header().Set("Expires", "0")
	resp.Header().Set("X-Frame-Options", "SAMEORIGIN")
	resp.Header().Set("Pragma", "no-cache")
	resp.Header().Set("X-Content-Type-Options", "nosniff")
	resp.Header().Set("X-XSS-Protection", "1; mode=block")
	resp.Header().Set("Content-Type", "application/vnd.spring-boot.actuator.v2+json")
	resp.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")

	resp.WriteHeader(200)

	_, _ = resp.Write(bodyBytes)
}

func RegisterProvider(ctx context.Context) error {
	server := webservice.WebServerFromContext(ctx)
	if server != nil {
		server.RegisterActuator(new(Provider))
		adminprovider.RegisterLink(providerName, providerName, false)
	}
	return nil
}
