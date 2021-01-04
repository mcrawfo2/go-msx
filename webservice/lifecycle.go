package webservice

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"github.com/pkg/errors"
)

var (
	ErrDisabled = errors.New("Web server disabled")
	service     *WebServer
)

// Start initiates the server listening
func Start(ctx context.Context) error {
	if service == nil {
		return nil
	}
	return service.Serve(ctx)
}

// Stop terminates the listening web server
func Stop(ctx context.Context) error {
	if service == nil {
		return nil
	}
	return service.StopServing(ctx)
}

// ConfigureWebServer creates a new WebServer from the supplied configuration
func NewWebServerFromConfig(cfg *config.Config, ctx context.Context) (*WebServer, error) {
	var webServerConfig WebServerConfig
	if err := cfg.Populate(&webServerConfig, configRootWebServer); err != nil {
		return nil, err
	}

	if !webServerConfig.Enabled {
		return nil, ErrDisabled
	}

	actuatorConfig, err := NewManagementSecurityConfig(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read management security config")
	}

	return NewWebServer(&webServerConfig, actuatorConfig, ctx)
}

// ConfigureWebServer creates a new WebServer from the supplied configuration and
// stores it in the `service` global
func ConfigureWebServer(cfg *config.Config, ctx context.Context) (err error) {
	service, err = NewWebServerFromConfig(cfg, ctx)
	return
}
