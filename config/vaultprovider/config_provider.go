package vaultprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/support/config"
	"cto-github.cisco.com/NFV-BU/go-msx/support/log"
	"cto-github.cisco.com/NFV-BU/go-msx/support/vault"
	"fmt"
	"github.com/pkg/errors"
)

const (
	configRootVaultConfigProvider = "spring.cloud.vault.generic"
	configKeyAppName              = "spring.application.name"
)

var logger = log.NewLogger("msx.config.vaultprovider")

type ConfigProviderConfig struct {
	Enabled          bool   `config:"default=false"`
	Backend          string `config:"default=secret"`
	ProfileSeparator string `config:"default=/"`
	DefaultContext   string `config:"default=defaultapplication"`
}

type ConfigProvider struct {
	sourceConfig *ConfigProviderConfig
	appContext   string
	connection   *vault.Connection
}

func (f *ConfigProvider) Load(ctx context.Context) (settings map[string]string, err error) {
	settings = make(map[string]string)

	// load keys from default context
	var vaultPath = fmt.Sprintf("%s%s%s", f.sourceConfig.Backend, f.sourceConfig.ProfileSeparator, f.sourceConfig.DefaultContext)
	logger.Infof("Loading configuration from vault (%s): %s)", f.connection.Host(), vaultPath)
	var defaultSettings map[string]string
	if defaultSettings, err = f.connection.ListSecrets(ctx, vaultPath); err != nil {
		return nil, errors.Wrap(err, "Failed to load configuration from vault")
	}

	for k, v := range defaultSettings {
		settings[config.NormalizeKey(k)] = v
	}

	// load keys from application context
	vaultPath = fmt.Sprintf("%s%s%s", f.sourceConfig.Backend, f.sourceConfig.ProfileSeparator, f.appContext)
	logger.Infof("Loading configuration from vault (%s): %s", f.connection.Host(), vaultPath)
	var appSettings map[string]string
	if appSettings, err = f.connection.ListSecrets(ctx, vaultPath); err != nil {
		return nil, errors.Wrap(err, "Failed to load configuration from vault")
	}

	for k, v := range appSettings {
		settings[config.NormalizeKey(k)] = v
	}

	return settings, nil
}

func NewConfigProviderFromConfig(cfg *config.Config) config.Provider {
	var sourceConfig = &ConfigProviderConfig{}
	var err = cfg.Populate(sourceConfig, configRootVaultConfigProvider)
	if err != nil {
		logger.Warn(err.Error())
		return nil
	}

	if !sourceConfig.Enabled {
		logger.Warn("Vault configuration source disabled")
		return nil
	}

	var appContext string
	if appContext, err = cfg.String(configKeyAppName); err != nil {
		logger.Warn(err.Error())
		return nil
	}

	if err = vault.ConfigurePool(cfg); err != nil {
		logger.Warn(err.Error())
		return nil
	}

	return &ConfigProvider{
		sourceConfig: sourceConfig,
		appContext:   appContext,
		connection:   vault.Pool().Connection(),
	}
}
