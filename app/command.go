package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/cli"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/config/cobraprovider"
	"github.com/spf13/cobra"
)

func init() {
	// Configure the root command
	cmd := cli.RootCmd()

	cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		RegisterProviderFactory(SourceCommandLine, func(cfg *config.Config) (config.Provider, error) {
			return config.NewCachedLoader(cobraprovider.NewCobraSource(cmd, "cli.flag.")), nil
		})
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return application.Run()
	}
}

type CommandFunc func(ctx context.Context) error

func AddCommand(path, brief string, simpleFunc CommandFunc) (cmd *cobra.Command, err error) {
	cmd, err = cli.AddCommand(path, brief, func(args []string) error {
		OnEvent(EventReady, PhaseAfter, func(ctx context.Context) error {
			logger.Infof("Executing command: %s", cmd.Use)
			if err := simpleFunc(ctx); err != nil {
				logger.Errorf("Command %s returned error: %v", cmd.Use, err)
				cli.SetExitCode(1)
			}
			return application.Stop()
		})

		return application.Run()
	})
	return cmd, err
}

func Run(appName string) {
	cli.Run(appName)
}
