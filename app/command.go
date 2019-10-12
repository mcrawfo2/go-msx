package app

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/config/cobraprovider"
	"github.com/spf13/cobra"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:                "app",
	FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
	RunE: func(cmd *cobra.Command, args []string) error {
		RegisterProviderFactory(SourceCommandLine, func(cfg *config.Config) (config.Provider, error) {
			return config.NewCachedLoader(cobraprovider.NewCobraSource(cmd, "cli.flag.")), nil
		})
		return Lifecycle()
	},
}

func FindCommand(path ...string) *cobra.Command {
	var next *cobra.Command
	here := rootCmd
	for _, pathPart := range path {
		hereCommands := here.Commands()

		next = nil
		for _, hereCommand := range hereCommands {
			if hereCommand.Use == pathPart || strings.HasPrefix(hereCommand.Use, pathPart+" ") {
				next = hereCommand
				break
			}
		}

		here = next
		if here == nil {
			break
		}
	}

	return here
}

func Run(appName string) {
	rootCmd.Use = appName
	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err)
	}
}
