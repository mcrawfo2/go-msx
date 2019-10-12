package main

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/app"
	"cto-github.cisco.com/NFV-BU/go-msx/cli"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/discovery"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"fmt"
	"github.com/pkg/errors"
)

const AppName = "someservice"

var logger = log.NewLogger(AppName)

func init() {
	app.OnEvent(app.EventStart, app.PhaseDuring, dumpConfiguration)
	app.OnEvent(app.EventReady, app.PhaseDuring, findUserManagement)
}

func dumpConfiguration(ctx context.Context) error {
	cfg := config.FromContext(ctx)
	if cfg == nil {
		return errors.New("Failed to obtain application config")
	}
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
	logger.Infof("Discovering %s", integration.ServiceNameUserManagement)
	if instances, err := discovery.Discover(ctx, integration.ServiceNameUserManagement, true); err != nil && err != discovery.ErrDiscoveryProviderNotDefined {
		return err
	} else if err == discovery.ErrDiscoveryProviderNotDefined {
		// Do nothing, discovery providers are disabled
	} else if len(instances) == 0 {
		return errors.New(fmt.Sprintf("No healthy instances of %s found", integration.ServiceNameUserManagement))
	} else {
		instance := instances.SelectRandom()
		logger.Info(instance)
	}
	return nil
}

func migrate(ctx context.Context) error {
	logger.Info("Migrate activity here")
	return nil
}

func populate(ctx context.Context) error {
	logger.Info("Populate activity here")
	return errors.New("Population failed")
}

func main() {
	cli.RootCmd().PersistentFlags().Bool("quiet", false, "Be quiet")
	if _, err := app.AddCommand("migrate", "Migrate database schema", migrate); err != nil {
		cli.Fatal(err)
	}
	if _, err := app.AddCommand("populate", "Populate remote microservices", populate); err != nil {
		cli.Fatal(err)
	}
	app.Run(AppName)
}
