package envprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/adminprovider"
	"github.com/emicklei/go-restful/v3"
	"strings"
)

const (
	endpointName = "env"
)

var logger = log.NewLogger("msx.webservice.envprovider")

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

	for _, cache := range config.FromContext(req.Request.Context()).Caches() {
		entries, err := cache.Load(req.Request.Context())
		if err != nil {
			logger.
				WithContext(req.Request.Context()).
				WithError(err).
				Errorf("Failed to load config properties from provider %q", cache.Description())
			continue
		}

		propertySource := PropertySource{
			Name:       cache.Description(),
			Properties: make(map[string]Property),
		}

		for _, entry := range entries {
			if h.isSecret(entry.Name) {
				propertySource.Properties[entry.Name] = Property{"*****"}
			} else {
				propertySource.Properties[entry.Name] = Property{entry.Value}
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
