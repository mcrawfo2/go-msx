// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package asyncapiprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"github.com/pkg/errors"
	"strings"
)

const configRootDocumentation = "asyncapi"

var (
	ErrDisabled = errors.New("AsyncApi disabled")
)

type DocumentationServerConfig struct {
	Host        string `config:"default=${server.host}"`
	Port        int    `config:"default=${server.port}"`
	ContextPath string `config:"default=${server.context-path}"`
}

type DocumentationUiConfig struct {
	Enabled  bool   `config:"default=true"`
	Endpoint string `config:"default=/asyncapi/studio"`
}

type DocumentationResourcesConfig struct {
	Path         string `config:"default=/asyncapi/resources"`
	YamlSpecPath string `config:"default=/asyncapi.yaml"`
	YamlSpecFile string `config:"default=/api/asyncapi.yaml"`
}

type DocumentationConfig struct {
	Enabled   bool   `config:"default=true"`
	Source    string `config:"default=registry"`
	Resources DocumentationResourcesConfig
	Ui        DocumentationUiConfig
	Server    DocumentationServerConfig
}

func NewDocumentationConfig(ctx context.Context) (*DocumentationConfig, error) {
	var documentationConfig DocumentationConfig
	if err := config.FromContext(ctx).Populate(&documentationConfig, configRootDocumentation); err != nil {
		return nil, err
	}

	if !strings.HasPrefix(documentationConfig.Resources.Path, "/") {
		documentationConfig.Resources.Path = "/" + documentationConfig.Resources.Path
	}

	if strings.HasSuffix(documentationConfig.Resources.Path, "/") {
		documentationConfig.Resources.Path = strings.TrimSuffix(documentationConfig.Resources.Path, "/")
	}

	return &documentationConfig, nil
}

func RegisterProvider(ctx context.Context) error {
	server := webservice.WebServerFromContext(ctx)
	if server == nil {
		return nil
	}

	provider, err := NewStudioProvider(ctx)
	if err != nil {
		return err
	}

	server.AddDocumentationProvider(provider)
	return nil
}
