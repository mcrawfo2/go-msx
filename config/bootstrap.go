package config

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/lifecycle"
	"cto-github.cisco.com/NFV-BU/go-msx/support/config"
	"cto-github.cisco.com/NFV-BU/go-msx/support/log"
	"time"
)

var logger = log.NewLogger("msx.config")
var bootstrapConfig *config.Config
var bootstrapSources *Sources

var staticConfig = make(map[string]string)

func init() {
	lifecycle.OnEvent(lifecycle.EventConfigure, lifecycle.PhaseBefore, createBootstrap)
	lifecycle.OnEvent(lifecycle.EventConfigure, lifecycle.PhaseBefore, mustLoadBootstrapConfig)
}

func mustLoadBootstrapConfig() {
	ctx, cancel := context.WithTimeout(lifecycle.Context(), time.Second*5)
	defer cancel()

	cfg := Bootstrap()
	if err := cfg.Load(ctx); err != nil {
		logger.Error(err)
		lifecycle.Shutdown()
	}
}

func createBootstrap() {
	bootstrapSources = &Sources{
		Defaults:      newDefaultsProvider(),
		BootstrapFile: newBootstrapProvider(),
		Environment:   newEnvironmentProvider(),
		Static:        newStaticProvider(staticConfig),
	}

	bootstrapConfig = config.NewConfig(bootstrapSources.Providers()...)
}

func Bootstrap() *config.Config {
	if bootstrapConfig == nil {
		createBootstrap()
	}

	return bootstrapConfig
}

func SetStaticConfig(static map[string]string) {
	if static != nil {
		staticConfig = static
	}
}
