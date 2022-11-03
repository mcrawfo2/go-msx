// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package build

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/cli"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"runtime"
	"strings"
)

var logLevel logrus.Level
var logLevelName string

func init() {
	rootCmd := cli.RootCmd()
	rootCmd.Flags().Bool("list", false, "List available build targets")
	rootCmd.PersistentFlags().StringArray("config", []string{"build.yml"}, "Specify one or more build config files")
	rootCmd.PersistentFlags().StringVarP(&logLevelName, "loglevel", "l",
		"INFO", "Set logging level: TRACE, DEBUG, INFO, WARN, ERROR or FATAL")
	rootCmd.PersistentPreRunE = handleGlobalFlags
}

func handleGlobalFlags(cmd *cobra.Command, args []string) error {
	logLevelName = strings.ToUpper(logLevelName)
	if log.CheckLevel(logLevelName) != nil {
		logger.Fatalf("invalid log level: %s", logLevelName)
	}
	logLevel = log.LevelFromName(logLevelName)
	log.SetLoggerLevel("build", logLevel)
	log.SetLoggerLevel("build.maven", logLevel)
	log.SetLoggerLevel("build.maven.platform", logLevel)

	return nil
}

func loadConfig(cmd *cobra.Command, args []string) error {
	configFiles, err := cmd.Root().PersistentFlags().GetStringArray("config")
	if err != nil {
		return err
	}

	if err = LoadBuildConfig(context.Background(), cmd, configFiles); err != nil {
		return err
	}

	currentGoVersion := getGoVersion()

	logger.Infof("Module: %s", BuildConfig.Module.ModulePath)
	logger.Infof("Required Go version: %s", BuildConfig.Module.MinGoVersion)
	logger.Infof("Current Go version: %s", currentGoVersion)

	if !strings.HasSuffix(BuildConfig.Module.ModulePath, "go-msx") {
		logger.Error("NOTE: cto-github.cisco.com/NFV-BU/go-msx/build package is deprecated.")
		logger.Fatal("NOTE: please switch to cto-github.cisco.com/NFV-BU/go-msx-build/pkg")
	}

	return nil
}

func getGoVersion() string {
	return strings.TrimPrefix(runtime.Version(), "go")
}

func Run() {
	log.SetLoggerLevel("msx.config", logrus.ErrorLevel)
	log.SetLoggerLevel("msx.config.pflagprovider", logrus.ErrorLevel)
	cli.Run("build")
}
