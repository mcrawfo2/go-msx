package skel

import (
	"cto-github.cisco.com/NFV-BU/go-msx/cli"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
)

const appName = "skel"
const configFileName = ".skel.json"

var logger = log.NewLogger("msx.skel")

var generateBeat bool

func init() {
	rootCmd := cli.RootCmd()
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		switch skeletonConfig.Generator {
		case "app":
			return GenerateSkeletonApp(args)
		case "beat":
			return GenerateSkeletonBeat(args)
		}
		return errors.Errorf("Unknown generator %q", skeletonConfig.Generator)
	}
	rootCmd.PersistentPreRunE = configure
	rootCmd.Flags().BoolVar(&generateBeat, "beat", false, "Generate a beat app skeleton")
}

func configure(cmd *cobra.Command, args []string) error {
	if loaded, err := loadConfig(); err != nil {
		return err
	} else if loaded {
		return nil
	}

	// Configure a new project if no project was found
	return ConfigureInteractive(args)
}

func loadConfig() (bool, error) {
	here, err := os.Getwd()
	if err != nil {
		return false, err
	}

	configFile := ""
	for here != "/" {
		hereFile := filepath.Join(here, configFileName)
		stat, err := os.Stat(hereFile)
		if err != nil && !os.IsNotExist(err) {
			return false, err
		} else if err == nil && !stat.IsDir() {
			configFile = hereFile
			break
		} else {
			err = nil
		}
		here = filepath.Dir(here)
	}

	if configFile == "" {
		return false, nil
	}

	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal(bytes, &skeletonConfig)
	if err != nil {
		return false, err
	}

	skeletonConfig.TargetParent = filepath.Dir(filepath.Dir(configFile))

	return true, nil
}

func Run() {
	log.SetLoggerLevel("msx.config", logrus.ErrorLevel)
	log.SetLoggerLevel("msx.config.pflagprovider", logrus.ErrorLevel)
	cli.Run(appName)
}
