package skel

import (
	"github.com/AlecAivazis/survey/v2"
	"os"
	"path"
	"strconv"
	"strings"
)

type SkeletonConfig struct {
	TargetParent      string `survey:"targetParent" json:"-"`
	AppName           string `survey:"appName" json:"appName"`
	AppDisplayName    string `survey:"appDisplayName" json:"appDisplayName"`
	AppDescription    string `survey:"appDescription" json:"appDescription"`
	ServerPort        int    `survey:"serverPort" json:"serverPort"`
	ServerContextPath string `survey:"serverContextPath" json:"serverContextPath"`
	AppVersion        string `survey:"appVersion" json:"appVersion"`
	Repository        string `survey:"repository" json:"repository"`
}

func (c SkeletonConfig) TargetDirectory() string {
	return path.Join(c.TargetParent, c.AppName)
}

func (c SkeletonConfig) AppMigrateVersion() string {
	return "V" + strings.ReplaceAll(c.AppVersion, ".", "_")
}

func (c SkeletonConfig) AppPackageUrl() string {
	return path.Join("cto-github.cisco.com", "NFV-BU", c.AppName)
}

var skeletonConfig = &SkeletonConfig{
	TargetParent:      path.Join(os.Getenv("HOME"), "msx"),
	AppName:           "someservice",
	AppDisplayName:    "Some Microservice",
	AppDescription:    "Does Something",
	AppVersion:        "3.9.0",
	ServerPort:        9999,
	ServerContextPath: "/some",
	Repository:        "cassandra",
}

var surveyQuestions = []*survey.Question{
	{
		Name: "targetParent",
		Prompt: &survey.Input{
			Message: "Project Parent Directory:",
			Default: skeletonConfig.TargetParent,
		},
		Validate: survey.Required,
	},
	{
		Name: "appVersion",
		Prompt: &survey.Input{
			Message: "Version:",
			Default: skeletonConfig.AppVersion,
		},
	},
	{
		Name: "appName",
		Prompt: &survey.Input{
			Message: "App name:",
			Default: skeletonConfig.AppName,
		},
		Validate:  survey.Required,
		Transform: survey.ToLower,
	},
	{
		Name: "appDisplayName",
		Prompt: &survey.Input{
			Message: "App display name:",
			Default: skeletonConfig.AppDisplayName,
		},
		Validate: survey.Required,
	},
	{
		Name: "appDescription",
		Prompt: &survey.Input{
			Message: "App description:",
			Default: skeletonConfig.AppDescription,
		},
		Validate: survey.Required,
	},
	{
		Name: "serverPort",
		Prompt: &survey.Input{
			Message: "Web server port:",
			Default: strconv.Itoa(skeletonConfig.ServerPort),
		},
		Validate: survey.Required,
	},
	{
		Name: "serverContextPath",
		Prompt: &survey.Input{
			Message: "Web server context path:",
			Default: skeletonConfig.ServerContextPath,
		},
		Validate: survey.Required,
	},
	{
		Name: "repository",
		Prompt: &survey.Select{
			Message: "Repository:",
			Options: []string{
				"cassandra",
				"cockroach",
			},
			Default: skeletonConfig.Repository,
		},
		Validate: survey.Required,
	},
}

func ConfigureInteractive(args []string) error {
	return survey.Ask(surveyQuestions, skeletonConfig)
}
