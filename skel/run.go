package skel

import (
	"cto-github.cisco.com/NFV-BU/go-msx/cli"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
)

const appName = "skel"
const projectConfigFileName = ".skel.json"
const generateConfigFileName = "generate.json"

var logger = log.NewLogger("msx.skel")

func init() {
	rootCmd := cli.RootCmd()
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return GenerateSkeleton(args)
	}
	rootCmd.PersistentPreRunE = configure
}

func configure(cmd *cobra.Command, args []string) error {
	if loaded, err := loadProjectConfig(); err != nil {
		return err
	} else if loaded {
		return nil
	}

	if loaded, err := loadGenerateConfig(); err != nil {
		return err
	} else if loaded {
		return nil
	}

	// Configure a new project if no project was found
	return ConfigureInteractive(args)
}

func loadConfig (configFile string) error {
	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &skeletonConfig)
	if err != nil {
		return err
	}

	return nil
}

func loadProjectConfig() (bool, error) {
	here, err := os.Getwd()
	if err != nil {
		return false, err
	}

	configFile := ""
	for here != "/" {
		hereFile := filepath.Join(here, projectConfigFileName)
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

	err = loadConfig(configFile)
	if err != nil {
		return false, err
	}

	skeletonConfig.TargetDir = filepath.Dir(configFile)
	skeletonConfig.TargetParent = filepath.Dir(skeletonConfig.TargetDir)

	return true, nil
}

func loadGenerateConfig() (bool, error) {
	stat, err := os.Stat(generateConfigFileName)
	if err != nil && !os.IsNotExist(err) {
		return false, err
	} else if err != nil {
		return false, nil
	} else if stat.IsDir() {
		return false, nil
	}

	err = loadConfig(generateConfigFileName)
	if err != nil {
		return false, err
	}

	return true, nil
}

func Run() {
	log.SetLoggerLevel("msx.config", logrus.ErrorLevel)
	log.SetLoggerLevel("msx.config.pflagprovider", logrus.ErrorLevel)
	cli.Run(appName)
}
