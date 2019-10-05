package main

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/support/consul"
	"cto-github.cisco.com/NFV-BU/go-msx/support/log"
	"cto-github.cisco.com/NFV-BU/go-msx/support/vault"
)

func main() {
	var logger = log.StandardLogger()

	config.SetStaticConfig(map[string]string{
		"application.name": "config",
	})

	if err := config.Bootstrap().Load(context.Background()); err != nil {
		logger.Fatal(err)
	}

	config.RegisterProviderFactory(config.SourceConsul, consul.NewConsulSource)
	config.RegisterProviderFactory(config.SourceVault, vault.NewVaultSource)

	cfg := config.Application()
	if err := cfg.Load(context.Background()); err != nil {
		logger.Fatal(err)
	}

	logger.Info("Dumping application configuration")
	settings := cfg.Settings()
	for name, value := range settings {
		logger.Infof("%s: %s", name, value)
	}
}
