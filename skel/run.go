// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package skel

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"os"
	"path/filepath"
	"strings"

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
var incFiles []string
var excFiles []string

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
	rootCmd.PersistentFlags().StringSliceVarP(&incFiles, "include", "i",
		[]string{}, "eg: i=\"**/*.go,**/*.ts\"\nOnly output file operations matching these quoted doublestar patterns will be done\nAlt syntax {alt1,...} not supported\nIf you don't include a file that a subsequent generation step needs, it may fail")
	rootCmd.PersistentFlags().StringSliceVarP(&excFiles, "exclude", "e",
		[]string{}, "eg: e=\"**/*.mod,**/*.sum\"\nOutput file operations matching these doublestar quoted patterns will not be done\nAlt syntax {alt1,...} not supported\nIf you exclude a file that a subsequent generation step needs, it may fail")

	rootCmd.PersistentPreRunE = configure

}

func configure(cmd *cobra.Command, _ []string) error {

	if cmd.Use == "version" {
		return nil
	}

	logLevelName = strings.ToUpper(logLevelName)
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

	printErr := color.New(color.FgRed).PrintfFunc()
	dir := types.May(os.Getwd())

	if cmd != cli.RootCmd() {

		// checks and warnings before starting

		// warn if there are existing projects
		projs, _ := FindProjects(dir, 4) // 4 seems like a good cutoff
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

// GitCheckAsk determines if there are vulnerable modifications in the given dir
// if so, it asks whether to continue
func GitCheckAsk(dir string) (ok bool, err error) {

	printErr := color.New(color.FgRed).PrintfFunc()

	logger.Tracef("GitCheckAsk(%s)", dir)

	// warn if git is dirty
	r, err := git.PlainOpenWithOptions(dir, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		logger.Tracef("GitCheckAsk(%s) not a git repo: %s", dir, err)
		return true, nil
	}
	wt, err := r.Worktree()
	if err != nil {
		logger.Tracef("GitCheckAsk(%s) error wt: %s", dir, err)
		return true, nil
	}
	status, err := wt.Status()
	if err != nil {
		logger.Tracef("GitCheckAsk(%s) error status: %s", dir, err)
		return true, nil
	}
	if !status.IsClean() {
		printErr("\n\nThere are some uncommited modified files in this repo\n")
		printErr("You may want to commit them before running this command\n\n")
		keepCalm := false
		carryOn := &survey.Confirm{
			Message: "Continue anyway?",
		}
		err = survey.AskOne(carryOn, &keepCalm)
		if err != nil {
			logger.Tracef("GitCheckAsk(%s) error ask: %s", dir, err)
			return false, err
		}
		if !keepCalm {
			return false, nil
		}
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
