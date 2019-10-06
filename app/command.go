package app

import (
	cobraConfig "cto-github.cisco.com/NFV-BU/go-msx/config/cobraprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/support/config"
	"github.com/spf13/cobra"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:                "app",
	FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
	Run: func(cmd *cobra.Command, args []string) {
		RegisterProviderFactory(SourceCommandLine, func(cfg *config.Config) config.Provider {
			return config.NewCachedLoader(cobraConfig.NewCobraSource(cmd, "cli.flag."))
		})
		Lifecycle()
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
