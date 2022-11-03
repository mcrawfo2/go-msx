// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package skel

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
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
var logLevel logrus.Level
var logLevelName string

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
	rootCmd.PersistentFlags().StringVarP(&logLevelName, "loglevel", "l",
		"INFO", "Set logging level: TRACE, DEBUG, INFO, WARN, ERROR or FATAL")

	rootCmd.PersistentPreRunE = configure

}

func configure(cmd *cobra.Command, _ []string) error {

	if cmd.Use == "version" {
		return nil
	}

	if log.CheckLevel(logLevelName) != nil {
		logger.Fatalf("invalid log level: %s", logLevelName)
	}
	logLevel = log.LevelFromName(logLevelName)
	logger.SetLevel(logLevel)
	if logLevel <= log.InfoLevel {
		fmt.Printf("Log level set to %s (%d)\n",
			log.LoggerLevel(logLevel).Name(), logLevel)
	}
	logger.Printf("Log level set to %s (%d)",
		log.LoggerLevel(logLevel).Name(), logLevel)

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

	if cmd != cli.RootCmd() {
		printErr := color.New(color.FgRed).PrintfFunc()

		// checks and warnings before starting
		// warn if there are existing projects
		projs, _ := FindProjects(types.May(os.Getwd()), 4) // 4 seems like a good cutoff
		if len(projs) > 0 {
			printErr("We found %d possible project(s) in this dir:\n", len(projs))
			for _, proj := range projs {
				fmt.Println("  " + proj + "\n")
			}
			printErr("Please switch to one of these folders first.\n")
		} else {
			printErr("No projects found in parent or child folders.  Please create a project first using `skel`.\n")
		}
		os.Exit(1)
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

// Run the skel command. This is the entry point for the skel command line
// it is called from cmd/skel.go
func Run(build int) {
	buildNumber = build
	log.SetLoggerLevel("msx.config", logrus.ErrorLevel)
	log.SetLoggerLevel("msx.config.pflagprovider", logrus.ErrorLevel)

	cli.Run(appName)
}
