package skel

import (
	"cto-github.cisco.com/NFV-BU/go-msx/cli"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const appName = "skel"

var logger = log.NewLogger("msx.skel")

func init() {
	rootCmd := cli.RootCmd()
	rootCmd.Flags().Bool("list", false, "List available build targets")
	rootCmd.PersistentFlags().StringArray("config", []string{"build.yml"}, "Specify one or more build config files")
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return GenerateSkeleton(args)
	}
	rootCmd.PersistentPreRunE = loadConfig
}

func loadConfig(cmd *cobra.Command, args []string) error {
	_, err := cmd.Root().PersistentFlags().GetStringArray("config")
	if err != nil {
		return err
	}
	//return LoadBuildConfig(context.Background(), configFiles)
	return nil
}

func Run() {
	log.SetLoggerLevel("msx.config", logrus.ErrorLevel)
	log.SetLoggerLevel("msx.config.pflagprovider", logrus.ErrorLevel)
	cli.Run(appName)
}
