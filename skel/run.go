// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package skel

import (
	"encoding/json"
	"os"
	"path/filepath"

	"cto-github.cisco.com/NFV-BU/go-msx/cli"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
)

// skel program is likely started by ../cmd/skel/skel.go

const appName = "skel"
const projectConfigFileName = ".skel.json"
const generateConfigFileName = "generate.json"

var logger = log.NewLogger("msx.skel")

var TitlingLanguage = language.English

// templates, loaded by provideStaticFiles
var staticFiles map[string]*staticFilesFile

func init() {
	var err error
	staticFiles, err = provideStaticFiles() // load the templates
	if err != nil {
		panic(err.Error())
	}
	rootCmd := cli.RootCmd()
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return GenerateSkeleton(args)
	}
	rootCmd.PersistentPreRunE = configure
}

func configure(cmd *cobra.Command, _ []string) error {
	if cmd.Use == "version" {
		return nil
	}

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

	// Configure a new project via the survey menus if no project was found
	return ConfigureInteractive()
}

func loadConfig(configFile string) error {
	bytes, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	skeletonConfig.Trunk = ""

	err = json.Unmarshal(bytes, &skeletonConfig)
	if err != nil {
		return err
	}

	if skeletonConfig.Trunk == "" {
		// default from before the field was added
		skeletonConfig.Trunk = "master"
	}

	return nil
}

// loadProjectConfig finds the project config file in the dir where skel is run, or above it
// in containing directories
func loadProjectConfig() (bool, error) {
	here, err := os.Getwd()
	if err != nil {
		return false, err
	}

	configFile := ""
	prev := ""
	for here != prev {
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
		prev = here
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

var buildNumber int

func Run(build int) {
	buildNumber = build
	log.SetLoggerLevel("msx.config", logrus.ErrorLevel)
	log.SetLoggerLevel("msx.config.pflagprovider", logrus.ErrorLevel)
	cli.Run(appName)
}
