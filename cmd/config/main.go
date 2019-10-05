package main

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/support/log"
	"time"
)

func main() {
	var logger = log.StandardLogger()

	config.SetStaticConfig(map[string]string{
		"application.name": "config",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 90 * time.Second)
	defer cancel()

	if err := config.Bootstrap().Load(ctx); err != nil {
		logger.Fatal(err)
	}

	config.RegisterRemoteConfigProviders()
	if err := config.Application().Load(ctx); err != nil {
		logger.Fatal(err)
	}

	logger.Info("Dumping application configuration")
	config.Application().Each(func(name, value string) {
		logger.Infof("%s: %s", name, value)
	})
}
