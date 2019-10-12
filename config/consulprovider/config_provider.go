package consulprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/consul"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"fmt"
	"github.com/pkg/errors"
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
	sourceConfig *ConfigProviderConfig
	appContext   string
	connection   *consul.Connection
}

func (f *ConfigProvider) Load(ctx context.Context) (settings map[string]string, err error) {
	settings = make(map[string]string)

	// load keys from default context
	var consulPrefix = fmt.Sprintf("%s/%s", f.sourceConfig.Prefix, f.sourceConfig.DefaultContext)
	logger.Infof("Loading configuration from consul (%s): %s)", f.connection.Host(), consulPrefix)
	var defaultSettings map[string]string
	if defaultSettings, err = f.connection.ListKeyValuePairs(ctx, consulPrefix); err != nil {
		return nil, errors.Wrap(err, "Failed to load configuration from consul")
	}

	for k, v := range defaultSettings {
		settings[config.NormalizeKey(k)] = v
	}

	// load keys from application context
	consulPrefix = fmt.Sprintf("%s/%s", f.sourceConfig.Prefix, f.appContext)
	logger.Infof("Loading configuration from consul (%s): %s", f.connection.Host(), consulPrefix)
	var appSettings map[string]string
	if appSettings, err = f.connection.ListKeyValuePairs(ctx, consulPrefix); err != nil {
		return nil, errors.Wrap(err, "Failed to load configuration from consul")
	}

	for k, v := range appSettings {
		settings[config.NormalizeKey(k)] = v
	}

	return settings, nil
}

func NewConfigProviderFromConfig(cfg *config.Config) (config.Provider, error) {
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

	return &ConfigProvider{
		sourceConfig: providerConfig,
		appContext:   appContext,
		connection:   conn,
	}, nil
}
