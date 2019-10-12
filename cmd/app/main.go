package main

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/app"
	"cto-github.cisco.com/NFV-BU/go-msx/discovery"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"fmt"
	"github.com/pkg/errors"
)

var logger = log.NewLogger("someservice")

func init() {
	app.OnEvent(app.EventStart, app.PhaseDuring, dumpConfiguration)
	app.OnEvent(app.EventReady, app.PhaseDuring, findUserManagement)
}

func dumpConfiguration(ctx context.Context) error {
	cfg := app.Config()
	quiet, _ := cfg.BoolOr("cli.flag.quiet", false)
	if !quiet {
		logger.Info("Dumping application configuration")
		cfg.Each(func(name, value string) {
			logger.Infof("%s: %s", name, value)
		})
	}
	return nil
}

func findUserManagement(ctx context.Context) error {
	serviceName := "usermanagementservice"
	if instances, err := discovery.Discover(serviceName, true); err != nil {
		return err
	} else if len(instances) == 0 {
		return errors.New(fmt.Sprintf("No healthy instances of %s found", serviceName))
	} else {
		instance := instances.SelectRandom()
		logger.Info(instance)
	}
	return nil
}

func main() {
	rootCmd := app.FindCommand()
	rootCmd.PersistentFlags().Bool("quiet", false, "Be quiet")
	app.Run("someservice")
}
