package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/cassandra/migrate"
	"cto-github.cisco.com/NFV-BU/go-msx/cli"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/config/cobraprovider"
	"github.com/spf13/cobra"
)

const (
	configKeyRedisEnable           = "spring.redis.enable"
	configKeyKafkaEnable           = "spring.cloud.stream.kafka.binder.enabled"
	configKeyConsulDiscoveryEnable = "spring.cloud.consul.discovery.enabled"
	configKeyServerEnable          = "server.enabled"

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
					"info.app.name": cmd.Use,
				})),
			}, nil
		})
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return application.Run(CommandRoot)
	}

	if _, err := AddCommand(CommandMigrate, "Migrate database schema", migrate.Migrate, commandMigrate); err != nil {
		cli.Fatal(err)
	}

	// TODO: Populate
}

func AddCommand(path, brief string, command Observer, init Observer) (cmd *cobra.Command, err error) {
	cmd, err = cli.AddCommand(path, brief, func(args []string) error {
		OnEvent(EventCommand, path, init)
		OnEvent(EventReady, PhaseAfter, func(ctx context.Context) error {
			logger.Infof("Executing command: %s", cmd.Use)
			if err := command(ctx); err != nil {
				logger.Errorf("Command %s returned error: %v", cmd.Use, err)
				cli.SetExitCode(1)
			}
			return application.Stop()
		})

		return application.Run(path)
	})
	return cmd, err
}

func Run(appName string) {
	cli.Run(appName)
}

func Noop(context.Context) error {
	return nil
}

func commandMigrate(context.Context) error {
	OverrideConfig(map[string]string{
		configKeyRedisEnable:           "false",
		configKeyKafkaEnable:           "false",
		configKeyConsulDiscoveryEnable: "false",
		configKeyServerEnable:          "false",
	})

	OnEvent(EventStart, PhaseBefore, setContextMigrationManifest)
	return nil
}

func setContextMigrationManifest(ctx context.Context) error {
	manifest, err := migrate.NewManifest(config.FromContext(ctx))
	if err != nil {
		return err
	}

	contextInjectors.Register(func(ctx context.Context) context.Context {
		return migrate.ContextWithManifest(ctx, manifest)
	})
	return nil
}
