package consulprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/consul"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/retry"
	"fmt"
	"github.com/pkg/errors"
	"time"
)

const (
	configRootConsulConfigProvider = "spring.cloud.consul.config"
	configKeyAppName               = "spring.application.name"
)

var logger = log.NewLogger("msx.config.consulprovider")

type ConfigProviderConfig struct {
	Enabled        bool   `config:"default=false"`
	Prefix         string `config:"default=userviceconfiguration"`
	DefaultContext string `config:"default=defaultapplication"`
	Pool           bool   `config:"default=false"`
}

type ConfigProvider struct {
	name         string
	sourceConfig *ConfigProviderConfig
	contextPath  string
	connection   *consul.Connection
}

func (f *ConfigProvider) Description() string {
	return fmt.Sprintf("%s: [%s]", f.name, f.ContextPath())
}

func (f *ConfigProvider) ContextPath() string {
	return fmt.Sprintf("%s/%s", f.sourceConfig.Prefix, f.contextPath)
}

func (f *ConfigProvider) Load(ctx context.Context) (settings map[string]string, err error) {
	settings = make(map[string]string)

	// load keys from default context
	var consulPrefix = f.ContextPath()
	logger.Infof("Loading configuration from consul (%s): %s)", f.connection.Host(), consulPrefix)
	var defaultSettings map[string]string

	err = retry.Retry{
		Attempts: 10,
		Delay:    3 * time.Second,
		BackOff:  0.0,
		Linear:   true,
		Context:  ctx,
	}.Retry(func() error {
		if ctx.Err() != nil {
			return &retry.PermanentError{Cause: err}
		}
		if defaultSettings, err = f.connection.ListKeyValuePairs(ctx, consulPrefix); err != nil {
			return errors.Wrap(err, "Failed to load configuration from consul")
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	for k, v := range defaultSettings {
		settings[config.NormalizeKey(k)] = v
	}

	return settings, nil
}

func NewConfigProvidersFromConfig(name string, cfg *config.Config) ([]config.Provider, error) {
	var providerConfig = &ConfigProviderConfig{}
	var err = cfg.Populate(providerConfig, configRootConsulConfigProvider)
	if err != nil {
		return nil, err
	}

	if !providerConfig.Enabled {
		logger.Warn("Consul configuration source disabled")
		return nil, nil
	}

	var appContext string
	if appContext, err = cfg.String(configKeyAppName); err != nil {
		return nil, err
	}

	var conn *consul.Connection
	if providerConfig.Pool {
		if err = consul.ConfigurePool(cfg); err != nil {
			return nil, err
		}
		conn = consul.Pool().Connection()
	} else if conn, err = consul.NewConnectionFromConfig(cfg); err != nil {
		return nil, err
	}

	return []config.Provider{
		&ConfigProvider{
			name:         name,
			sourceConfig: providerConfig,
			contextPath:  providerConfig.DefaultContext,
			connection:   conn,
		},
		&ConfigProvider{
			name:         name,
			sourceConfig: providerConfig,
			contextPath:  appContext,
			connection:   conn,
		},
	}, nil
}
