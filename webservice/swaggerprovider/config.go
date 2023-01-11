// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package swaggerprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/schema"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"github.com/pkg/errors"
	"strings"
)

const (
	configRootDocumentation = "swagger"
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
	Enabled     bool     `config:"default=true"`
	Endpoint    string   `config:"default=/swagger"`
	StaticFiles string   `config:"default=/swagger/ui"`
	StaticView  string   `config:"default=/swagger/static"`
	View        string   `config:"default=/swagger-ui.html"`
	RootFiles   []string `config:"default=swagger-sso-redirect.html"`
}

type DocumentationServerConfig struct {
	TlsEnabled  bool   `config:"default=${server.tls.enabled:false}"`
	Host        string `config:"default=${server.host:localhost}"`
	Port        int    `config:"default=${server.port:0}"`
	ContextPath string `config:"default=${server.context-path:/${info.app.name}}"`
}

func (c DocumentationServerConfig) Scheme() string {
	switch c.TlsEnabled {
	case false:
		return "http"
	default:
		return "https"
	}
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

	appInfo, err := schema.AppInfoFromConfig(config.MustFromContext(ctx))
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
		provider = NewSwaggerProvider(ctx, cfg, appInfo)

	case 3:
		provider = NewOpenApiProvider(ctx, cfg, appInfo)

	default:
		return errors.Errorf("Unknown OpenApi version: %q", cfg.Version)
	}

	server.AddDocumentationProvider(provider)
	specProvider = provider.(SpecProvider)
	return nil
}
