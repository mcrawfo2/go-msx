package consul

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/support/config"
	"fmt"
	"github.com/pkg/errors"
)

type ConfigProviderConfig struct {
	Prefix         string `properties:"prefix,default=userviceconfiguration"`
	DefaultContext string `properties:"defaultContext,default=defaultapplication"`
}

type ConfigProvider struct {
	sourceConfig *ConfigProviderConfig
	appContext   string
	connection   *Connection
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

func NewConsulSource(cfg *config.Config) config.Provider {
	var sourceConfig = &ConfigProviderConfig{}
	var err = cfg.Populate(sourceConfig, "spring.cloud.consul.config")
	if err != nil {
		logger.Warn(err.Error())
		return nil
	}

	var appContext string
	if appContext, err = cfg.String("spring.app.name"); err != nil {
		logger.Warn(err.Error())
		return nil
	}

	var conn *Connection
	if conn, err = NewConnection(cfg); err != nil {
		logger.Warn(err.Error())
		return nil
	}

	return &ConfigProvider{
		sourceConfig: sourceConfig,
		appContext:   appContext,
		connection:   conn,
	}
}
