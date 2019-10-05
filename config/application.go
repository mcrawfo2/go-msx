package config

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/lifecycle"
	"cto-github.cisco.com/NFV-BU/go-msx/support/config"
	"cto-github.cisco.com/NFV-BU/go-msx/support/consul"
	"cto-github.cisco.com/NFV-BU/go-msx/support/vault"
	"time"
)

var applicationConfig *config.Config
var applicationSources *Sources

func init() {
	lifecycle.OnEvent(lifecycle.EventConfigure, lifecycle.PhaseBefore, registerConfigProviders)
	lifecycle.OnEvent(lifecycle.EventConfigure, lifecycle.PhaseDuring, createApplicationConfig)
	lifecycle.OnEvent(lifecycle.EventConfigure, lifecycle.PhaseAfter, mustLoadApplicationConfig)
	lifecycle.OnEvent(lifecycle.EventConfigure, lifecycle.PhaseAfter, watchApplicationConfig)
}

func registerConfigProviders() {
	RegisterProviderFactory(SourceConsul, consul.NewConsulSource)
	RegisterProviderFactory(SourceVault, vault.NewVaultSource)
}

func watchApplicationConfig() {
	cfg := Application()
	cfg.Watch(lifecycle.Context())
}

func mustLoadApplicationConfig() {
	ctx, cancel := context.WithTimeout(lifecycle.Context(), time.Second*15)
	defer cancel()

	cfg := Application()
	if err := cfg.Load(ctx); err != nil {
		logger.Error(err)
		lifecycle.Shutdown()
	}
}

func createApplicationConfig() {
	if bootstrapConfig == nil {
		createBootstrap()
	}

	// Load full config
	applicationSources = &Sources{
		Defaults:        bootstrapSources.Defaults,
		BootstrapFile:   bootstrapSources.BootstrapFile,
		ApplicationFile: newProvider(SourceApplication, bootstrapConfig),
		Consul:          newProvider(SourceConsul, bootstrapConfig),
		Vault:           newProvider(SourceVault, bootstrapConfig),
		Profile:         newProvider(SourceProfile, bootstrapConfig),
		Environment:     bootstrapSources.Environment,
		Static:          bootstrapSources.Static,
	}

	applicationConfig = config.NewConfig(applicationSources.Providers()...)
}

func Application() *config.Config {
	if applicationConfig == nil {
		createApplicationConfig()
	}

	return applicationConfig
}
