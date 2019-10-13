package main

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/app"
	"cto-github.cisco.com/NFV-BU/go-msx/cassandra"
	"cto-github.cisco.com/NFV-BU/go-msx/cli"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/discovery"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
)

const AppName = "someservice"

var logger = log.NewLogger(AppName)

func init() {
	app.OnEvent(app.EventStart, app.PhaseDuring, dumpConfiguration)
	app.OnEvent(app.EventReady, app.PhaseDuring, findUserManagement)
	app.OnEvent(app.EventReady, app.PhaseDuring, listGauges)
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

func listGauges(ctx context.Context) error {
	cassandraPool := cassandra.PoolFromContext(ctx)
	if cassandraPool == nil {
		return errors.New("Cassandra connection pool not found")
	}

	return cassandraPool.WithSession(listGaugesFromSession)
}

func listGaugesFromSession(session *gocql.Session) error {
	var serviceType, deviceType, deviceSubType, beatType *string
	if err := session.Query(`SELECT servicetype, devicetype, devicesubtype, beattype FROM gauges LIMIT 1 ALLOW FILTERING`).
			Consistency(gocql.One).
			Scan(&serviceType, &deviceType, &deviceSubType, &beatType); err != nil {
		logger.Error(err)
	} else {
		logger.Infof("Found gauges: serviceType=%s deviceType=%s deviceSubType=%s beatType=%s",
			*serviceType, *deviceType, *deviceSubType, *beatType)
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
