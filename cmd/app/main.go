package main

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/lifecycle"
	"cto-github.cisco.com/NFV-BU/go-msx/support/log"
)

var logger = log.NewLogger("app")

func init() {
	lifecycle.OnEvent(lifecycle.EventInit, lifecycle.PhaseAfter, setStaticConfig)
	lifecycle.OnEvent(lifecycle.EventStart, lifecycle.PhaseDuring, dumpConfiguration)
}

func setStaticConfig() {
	config.SetStaticConfig(map[string]string{
		"spring.app.name": "app",
	})
}

func dumpConfiguration() {
	logger.Info("Dumping application configuration")
	config.Application().Each(func(name, value string) {
		logger.Infof("%s: %s", name, value)
	})
}

func main() {
	lifecycle.Run()
}
