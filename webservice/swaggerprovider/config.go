// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package swaggerprovider

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
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

type DocumentationConfig struct {
	Enabled     bool   `config:"default=false"`
	ApiPath     string `config:"default=/apidocs.json"`
	SwaggerPath string `config:"default=/swagger-resources"`
	Security    DocumentationSecurityConfig
	Ui          DocumentationUiConfig
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
