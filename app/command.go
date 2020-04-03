package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/cassandra/migrate"
	"cto-github.cisco.com/NFV-BU/go-msx/cli"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/config/cobraprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/populate"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

const (
	configKeyRedisEnable           = "spring.redis.enable"
	configKeyKafkaEnable           = "spring.cloud.stream.kafka.binder.enabled"
	configKeyConsulDiscoveryEnable = "spring.cloud.consul.discovery.enabled"
	configKeyServerEnable          = "server.enabled"
	configKeyLeaderEnable          = "consul.leader.election.enabled"

	CommandRoot     = ""
	CommandMigrate  = "migrate"
	CommandPopulate = "populate"
)

func init() {
	// Configure the root command
	cmd := cli.RootCmd()

	cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		RegisterProviderFactory(SourceCommandLine, func(name string, cfg *config.Config) ([]config.Provider, error) {
			return []config.Provider{
				config.NewCachedLoader(
					cobraprovider.NewCobraSource(name, cmd, "cli.flag."),
				),
				config.NewCachedLoader(config.NewStatic("Built-In", map[string]string{
					"spring.application.name": strings.Fields(cli.RootCmd().Use)[0],
					"info.app.name":           strings.Fields(cli.RootCmd().Use)[0],
				})),
			}, nil
		})
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return application.Run(CommandRoot)
	}

	if _, err := AddCommand(CommandMigrate, "Migrate database schema", migrate.Migrate, commandMigrateInit); err != nil {
		cli.Fatal(err)
	} else {
	}

	if populateCommand, err := AddCommand(CommandPopulate, "Populate data", populate.Populate, commandPopulateInit); err != nil {
		cli.Fatal(err)
	} else {
		populate.CustomizeCommand(populateCommand)
	}
}

func AddCommand(path, brief string, command CommandObserver, init Observer) (cmd *cobra.Command, err error) {
	return cli.AddCommand(path, brief, func(args []string) error {
		OnEvent(EventCommand, path, init)
		OnEvent(EventReady, PhaseAfter, func(ctx context.Context) error {
			logger.Infof("Executing command: %s", strings.Fields(cmd.Use)[0])
			if err := command(ctx, args); err != nil {
				logger.Errorf("Command %s returned error: %v", strings.Fields(cmd.Use)[0], err)
				cli.SetExitCode(1)
			}
			return application.Stop()
		})

		return application.Run(path)
	})
}

func Run(appName string) {
	// Convert environment variable POPULATE into migrate command
	if strings.ToLower(os.Getenv("POPULATE")) == "database" {
		os.Args = append(os.Args, "migrate")
	}

	cli.Run(appName)
}

func Noop(context.Context) error {
	return nil
}

func commandMigrateInit(context.Context) error {
	OverrideConfig(map[string]string{
		configKeyRedisEnable:           "false",
		configKeyKafkaEnable:           "false",
		configKeyConsulDiscoveryEnable: "false",
		configKeyServerEnable:          "false",
		configKeyLeaderEnable:          "false",
	})

	OnEvent(EventStart, PhaseBefore, func(ctx context.Context) error {
		manifest, err := migrate.NewManifest(config.FromContext(ctx))
		if err != nil {
			return err
		}

		contextInjectors.Register(func(ctx context.Context) context.Context {
			return migrate.ContextWithManifest(ctx, manifest)
		})
		return nil
	})
	return nil
}

func commandPopulateInit(context.Context) error {
	OverrideConfig(map[string]string{
		configKeyConsulDiscoveryEnable: "false",
		configKeyServerEnable:          "false",
		configKeyLeaderEnable:          "false",
	})

	return nil
}
