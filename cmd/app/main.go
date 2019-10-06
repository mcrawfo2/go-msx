package main

import (
	"cto-github.cisco.com/NFV-BU/go-msx/app"
	"cto-github.cisco.com/NFV-BU/go-msx/support/log"
)

var logger = log.NewLogger("someservice")

func init() {
	app.OnEvent(app.EventStart, app.PhaseDuring, dumpConfiguration)
}

func dumpConfiguration() {
	cfg := app.Config()
	quiet, _ := cfg.BoolOr("cli.flag.quiet", false)
	if !quiet {
		logger.Info("Dumping application configuration")
		cfg.Each(func(name, value string) {
			logger.Infof("%s: %s", name, value)
		})
	}
}

func main() {
	rootCmd := app.FindCommand()
	rootCmd.PersistentFlags().Bool("quiet", false, "Be quiet")
	app.Run("someservice")
}
