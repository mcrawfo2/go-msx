package build

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/cli"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd := cli.RootCmd()
	rootCmd.Flags().Bool("list", false, "List available build targets")
	rootCmd.PersistentFlags().StringArray("config", []string{"build.yml"}, "Specify one or more build config files")
}

func loadConfig(cmd *cobra.Command, args []string) error {
	configFiles, err := cmd.Root().PersistentFlags().GetStringArray("config")
	if err != nil {
		return err
	}
	return LoadBuildConfig(context.Background(), configFiles)
}

func Run() {
	// TODO: 1.0.0 : Error, Exit
	logger.Warn("NOTE: cto-github.cisco.com/NFV-BU/go-msx/build package is deprecated.")
	logger.Warn("NOTE: please switch to cto-github.cisco.com/NFV-BU/go-msx-build/pkg")

	log.SetLoggerLevel("msx.config", logrus.ErrorLevel)
	log.SetLoggerLevel("msx.config.pflagprovider", logrus.ErrorLevel)
	cli.Run("build")
}
