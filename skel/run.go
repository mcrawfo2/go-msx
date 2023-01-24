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
	"github.com/pkg/errors"
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

var ErrUserCancel = errors.New("user cancelled skel run")
var ErrNoProjects = errors.New("no projects found")

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

	// set the log level
	logLevelName = strings.ToUpper(logLevelName)
	if log.CheckLevel(logLevelName) != nil {
		logger.Fatalf("invalid log level: %s", logLevelName)
	}
	logLevel = log.LevelFromName(logLevelName)
	logger.SetLevel(logLevel)
	logger.Debugf("Log level set to %s (%d)", log.LoggerLevel(logLevel).Name(), logLevel)

	// compute some flags:
	//   root -- are we running the root command?
	//   project -- are we in a project subdir?
	//   subProjs -- are there project subdirs in the dir?
	//   dirtyGit -- are there uncommitted files in this git repo?
	// then load the project or generator config

	root := cmd == cli.RootCmd()
	project := false

	loaded, err := loadProjectConfig()
	if err != nil {
		return err
	}
	if loaded {
		project = true
		logger.Info("Loaded project config")
	}

	if !loaded {
		loaded, err := loadGenerateConfig()
		if err != nil {
			return err
		}
		if loaded {
			project = false
			logger.Info("Loaded generator config")
		}
	}

	subProjs := false
	var projs []DirName
	dir := types.May(os.Getwd())
	if !loaded { // FindProjects can be expensive
		projs, _ = FindProjects(dir, 4) // 4 seems like a good cutoff
		subProjs = len(projs) > 0
	}

	logger.Debugf("configure: root:%t, project:%t, #projects:%d",
		root, project, len(projs))

	// flags now set

	printErr := color.New(color.FgRed).PrintfFunc()
	printInfo := color.New(color.FgBlue).PrintfFunc()

	if subProjs {
		printErr("We found %d possible project(s) in this dir:\n", len(projs))
		for _, proj := range projs {
			fmt.Println("  " + proj + "\n")
		}
		printErr("Please switch to one of these folders first.\n")
		os.Exit(1)
	}

	if project { // only makes sense to check for dirty git if we are in a project
		ok, err := GitCheckAsk(dir)
		if err != nil {
			logger.WithError(err).Errorf("configure(%s) error ask:", dir)
			return err
		}

		if !ok {
			logger.WithError(ErrUserCancel).Errorf("configure(%s) user cancels, uncommitted files", dir)
			os.Exit(1)
		}
	}

	// root cmd in empty dir generates from scratch using menus
	if root && !project {
		return ConfigureInteractive()
	}

	// root cmd in dir containing .skel.json (project dir) uses those settings,
	// explains it will regenerate and confirms
	if root && project {
		printErr("\n\nThere is already a project in this directory\n")
		printInfo("Did you, perhaps, mean to run a skel subcommand?\n")
		printErr("Do you want to continue, which will regenerate the project, and may overwrite files?\n\n")
		keepGoing := false
		contQ := &survey.Confirm{
			Message: "Continue and regenerate?",
		}
		err = survey.AskOne(contQ, &keepGoing)
		if err != nil {
			logger.WithError(err).Errorf("configure(%s) error ask:", dir)
			return err
		}
		if !keepGoing {
			logger.WithError(ErrUserCancel).Errorf("configure(%s) user cancels, no regen:", dir)
			os.Exit(1)
		}

		skeletonConfig.noOverwrite = false // ovewriting still allowed

		return nil // regenerate using the loaded settings
	}

	// !root so should be in a project dir, if not, we error out
	if !project {
		printErr("No projects found in parent or child folders.  Please create a project first using `skel`.\n")
		logger.WithError(ErrNoProjects).Errorf("Non-root with no projects: %s", dir)
		os.Exit(1)
	}

	// subcommands do not overwrite files
	skeletonConfig.noOverwrite = true // no overwriting allowed

	return nil
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

// GitDirtyCheck determines if there are vulnerable modifications in the given directory
func GitDirtyCheck(dir string) (dirty bool, err error) {
	r, err := git.PlainOpenWithOptions(dir, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		logger.Tracef("GitDirtyCheck(%s) not a git repo: %s", dir, err)
		return true, nil
	}
	wt, err := r.Worktree()
	if err != nil {
		logger.Tracef("GitDirtyCheck(%s) error wt: %s", dir, err)
		return true, nil
	}
	status, err := wt.Status()
	if err != nil {
		logger.Tracef("GitDirtyCheck(%s) error status: %s", dir, err)
		return true, nil
	}
	if !status.IsClean() {
		return true, nil
	}
	return false, nil
}

// GitCheckAsk determines if there are vulnerable modifications in the given directory
// and if so, it asks whether to continue and returns an ok flag if so directed
func GitCheckAsk(dir string) (ok bool, err error) {

	logger.Tracef("GitCheckAsk(%s)", dir)
	dirty, err := GitDirtyCheck(dir)
	if err != nil {
		logger.WithError(err).Debugf("Config(%s) error ask:", dir)
		return false, err
	}
	if !dirty {
		return true, nil
	}

	printErr := color.New(color.FgRed).PrintfFunc()

	printErr("\n\nThere are some uncommited modified files in this repo\n")
	printErr("You may want to commit them before running this command\n\n")
	keepGoing := false
	contQ := &survey.Confirm{
		Message: "Continue anyway?",
	}
	err = survey.AskOne(contQ, &keepGoing)
	if err != nil {
		logger.WithError(err).Debugf("GitCheckAsk(%s) error ask:", dir)
		return false, err
	}
	if !keepGoing {
		return false, nil
	}
	return true, nil
}

var buildNumber int

// Run the skel command. This is the entry point for the skel command line
// it is called from cmd/skel.go
func Run(build int) {
	buildNumber = build
	log.SetFormat(log.LogFormatLogFmt)
	log.SetLoggerLevel("msx.config", logrus.ErrorLevel)
	log.SetLoggerLevel("msx.config.pflagprovider", logrus.ErrorLevel)

	cli.Run(appName)
}
