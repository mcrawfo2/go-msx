package vaultprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/consul"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/vault"
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
	Pool             bool   `config:"default=false"`
}

type ConfigProvider struct {
	name         string
	sourceConfig *ConfigProviderConfig
	appContext   string
	connection   *vault.Connection
}

func (f *ConfigProvider) Description() string {
	return fmt.Sprintf("%s: [%s, %s]", f.name, f.defaultContextPath(), f.applicationContextPath())
}

func (f *ConfigProvider) defaultContextPath() string {
	return fmt.Sprintf("%s%s%s", f.sourceConfig.Backend, f.sourceConfig.ProfileSeparator, f.sourceConfig.DefaultContext)
}

func (f *ConfigProvider) applicationContextPath() string {
	return fmt.Sprintf("%s%s%s", f.sourceConfig.Backend, f.sourceConfig.ProfileSeparator, f.appContext)
}

func (f *ConfigProvider) Load(ctx context.Context) (settings map[string]string, err error) {
	settings = make(map[string]string)

	// load keys from default context
	var vaultPath = f.defaultContextPath()
	logger.Infof("Loading configuration from vault (%s): %s)", f.connection.Host(), vaultPath)
	var defaultSettings map[string]string
	if defaultSettings, err = f.connection.ListSecrets(ctx, vaultPath); err != nil {
		return nil, errors.Wrap(err, "Failed to load configuration from vault")
	}

	for k, v := range defaultSettings {
		settings[config.NormalizeKey(k)] = v
	}

	// load keys from application context
	vaultPath = f.applicationContextPath()
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

func NewConfigProviderFromConfig(name string, cfg *config.Config) (config.Provider, error) {
	var providerConfig = &ConfigProviderConfig{}
	var err = cfg.Populate(providerConfig, configRootVaultConfigProvider)
	if err != nil {
		return nil, err
	}

	if !providerConfig.Enabled {
		logger.Warn("Vault configuration source disabled")
		return nil, nil
	}

	var appContext string
	if appContext, err = cfg.String(configKeyAppName); err != nil {
		return nil, err
	}

	var conn *vault.Connection
	if providerConfig.Pool {
		if err = consul.ConfigurePool(cfg); err != nil {
			return nil, err
		}
		conn = vault.Pool().Connection()
	} else if conn, err = vault.NewConnectionFromConfig(cfg); err != nil {
		return nil, err
	}

	return &ConfigProvider{
		name:         name,
		sourceConfig: providerConfig,
		appContext:   appContext,
		connection:   conn,
	}, nil
}
