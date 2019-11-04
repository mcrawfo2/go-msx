package swaggerprovider

import "cto-github.cisco.com/NFV-BU/go-msx/config"

const (
	configRootDocumentation = "server.swagger"
)

type DocumentationConfig struct {
	Enabled        bool   `config:"default=false"`
	Path           string `config:"default=thirdparty/swagger-ui/dist"`
	WebServicesUrl string
	ApiPath        string
	SwaggerPath    string
}

func DocumentationConfigFromConfig(cfg *config.Config) (*DocumentationConfig, error) {
	var documentationConfig DocumentationConfig
	if err := cfg.Populate(&documentationConfig, configRootDocumentation); err != nil {
		return nil, err
	}

	return &documentationConfig, nil
}
