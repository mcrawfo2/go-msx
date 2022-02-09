// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package swaggerprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"strings"
)

const (
	configRootDocumentation = "swagger"
	configKeyAppName        = "info.app.name"
	configKeyAppDescription = "info.app.description"
	configKeyBuildVersion   = "info.build.version"
	configKeyDisplayName    = "info.app.attributes.display-name"
)

type DocumentationSecuritySsoConfig struct {
	BaseUrl       string `config:"default=http://localhost:9103/idm"`
	TokenPath     string `config:"default=/v2/token"`
	AuthorizePath string `config:"default=/v2/authorize"`
	ClientId      string `config:"default=nfv-client"`
	ClientSecret  string `config:"default="`
}

type DocumentationSecurityConfig struct {
	Sso DocumentationSecuritySsoConfig
}

type DocumentationUiConfig struct {
	Enabled  bool   `config:"default=true"`
	Endpoint string `config:"default=/swagger"`
	View     string `config:"default=/swagger-ui.html"`
}

type DocumentationServerConfig struct {
	Host        string `config:"default=${server.host}"`
	Port        int    `config:"default=${server.port}"`
	ContextPath string `config:"default=${service.context-path}"`
}

type DocumentationConfig struct {
	Enabled     bool   `config:"default=false"`
	ApiPath     string `config:"default=/apidocs.json"`
	ApiYamlPath string `config:"default=/apidocs.yml"`
	SwaggerPath string `config:"default=/swagger-resources"`
	Version     string `config:"default=2.0"`
	Security    DocumentationSecurityConfig
	Ui          DocumentationUiConfig
	Server      DocumentationServerConfig
}

func DocumentationConfigFromConfig(cfg *config.Config) (*DocumentationConfig, error) {
	var documentationConfig DocumentationConfig
	if err := cfg.Populate(&documentationConfig, configRootDocumentation); err != nil {
		return nil, err
	}

	if !strings.HasPrefix(documentationConfig.SwaggerPath, "/") {
		documentationConfig.SwaggerPath = "/" + documentationConfig.SwaggerPath
	}

	if strings.HasSuffix(documentationConfig.SwaggerPath, "/") {
		documentationConfig.SwaggerPath = strings.TrimSuffix(documentationConfig.SwaggerPath, "/")
	}

	return &documentationConfig, nil
}

type AppInfo struct {
	Name        string
	DisplayName string
	Description string
	Version     string
}

func AppInfoFromConfig(cfg *config.Config) (*AppInfo, error) {
	var appInfo AppInfo
	var err error
	if appInfo.Name, err = cfg.String(configKeyAppName); err != nil {
		return nil, err
	}
	if appInfo.Description, err = cfg.String(configKeyAppDescription); err != nil {
		return nil, err
	}
	if appInfo.Version, err = cfg.String(configKeyBuildVersion); err != nil {
		return nil, err
	}
	if appInfo.DisplayName, err = cfg.String(configKeyDisplayName); err != nil {
		return nil, err
	}
	return &appInfo, nil
}

func RegisterSwaggerProvider(ctx context.Context) error {
	server := webservice.WebServerFromContext(ctx)
	if server == nil {
		return nil
	}

	cfg, err := DocumentationConfigFromConfig(config.MustFromContext(ctx))
	if err != nil {
		return err
	}

	if !cfg.Enabled {
		return ErrDisabled
	}

	appInfo, err := AppInfoFromConfig(config.MustFromContext(ctx))
	if err != nil {
		return err
	}

	oapiVersion, err := types.NewVersion(cfg.Version)
	if err != nil {
		return err
	}

	var provider webservice.DocumentationProvider
	switch oapiVersion[0] {
	case 2:
		provider = &SwaggerProvider{
			ctx:     ctx,
			cfg:     cfg,
			appInfo: appInfo,
		}

	case 3:
		provider = &OpenApiProvider{
			ctx:       ctx,
			cfg:       cfg,
			appInfo:   appInfo,
			reflector: &webservice.Reflector,
			spec:      webservice.Reflector.Spec,
		}
	}

	server.AddDocumentationProvider(provider)
	specProvider = provider.(SpecProvider)
	return nil
}
