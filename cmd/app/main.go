package main

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/lifecycle"
	"cto-github.cisco.com/NFV-BU/go-msx/support/log"
)

func init() {
	lifecycle.OnEvent(lifecycle.EventInit, lifecycle.PhaseAfter, func() {
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
