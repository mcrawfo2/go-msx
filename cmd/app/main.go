package main

import (
	"cto-github.cisco.com/NFV-BU/go-msx/app"
	"cto-github.cisco.com/NFV-BU/go-msx/discovery"
	"cto-github.cisco.com/NFV-BU/go-msx/support/log"
)

var logger = log.NewLogger("someservice")

func init() {
	app.OnEvent(app.EventStart, app.PhaseDuring, dumpConfiguration)
	app.OnEvent(app.EventReady, app.PhaseDuring, findUserManagement)
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

func findUserManagement() {
	if instances, err := discovery.Discover("usermanagementservice", true); err != nil {
		logger.Error(err)
	} else if len(instances) == 0 {
		logger.Error("No healthy instances of usermanagementservice found")
	} else {
		instance := instances.SelectRandom()
		logger.Info(instance)
	}

}

func main() {
	rootCmd := app.FindCommand()
	rootCmd.PersistentFlags().Bool("quiet", false, "Be quiet")
	app.Run("someservice")
}
