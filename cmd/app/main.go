package main

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	_ "cto-github.cisco.com/NFV-BU/go-msx/support/consul"
	_ "cto-github.cisco.com/NFV-BU/go-msx/support/vault"
	"cto-github.cisco.com/NFV-BU/go-msx/lifecycle"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
)

func init() {
	lifecycle.OnEvent(lifecycle.EventInit, lifecycle.PhaseBefore, func() {
		config.SetStaticConfig(map[string]string{
			"spring.app.name": "app",
		})
	})
}

func main() {
	var logger = log.NewLogger("app")

	lifecycle.OnEvent(lifecycle.EventReady, lifecycle.PhaseDuring, func() {
		logger.Info("Dumping application configuration")
		cfg := config.Application()
		settings := cfg.Settings()
		for name, value := range settings {
			logger.Infof("%s: %s", name, value)
		}
	})

	lifecycle.Run()
}
