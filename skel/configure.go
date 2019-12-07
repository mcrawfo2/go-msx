package skel

import (
	"github.com/AlecAivazis/survey/v2"
	"os"
	"path"
	"strconv"
)

type SkeletonConfig struct {
	TargetParent   string `survey:"targetParent"`
	AppName        string `survey:"appName"`
	AppDisplayName string `survey:"appDisplayName"`
	AppDescription string `survey:"appDescription"`
	ServerPort     int    `survey:"serverPort"`
}

func (c SkeletonConfig) TargetDirectory() string {
	return path.Join(c.TargetParent, c.AppName)
}

var skeletonConfig = &SkeletonConfig{
	TargetParent:   path.Join(os.Getenv("HOME"), "Documents"),
	AppName:        "someservice",
	AppDisplayName: "Some Microservice",
	AppDescription: "Does Something",
	ServerPort:     9999,
}

var surveyQuestions = []*survey.Question{
	{
		Name:      "targetParent",
		Prompt:    &survey.Input{
			Message: "Project Parent Directory:",
			Default: skeletonConfig.TargetParent,
		},
		Validate:  survey.Required,
	},
	{
		Name:      "appName",
		Prompt:    &survey.Input{
			Message: "App name:",
			Default: skeletonConfig.AppName,
		},
		Validate:  survey.Required,
		Transform: survey.ToLower,
	},
	{
		Name:      "appDisplayName",
		Prompt:    &survey.Input{
			Message: "App display name:",
			Default: skeletonConfig.AppDisplayName,
		},
		Validate:  survey.Required,
	},
	{
		Name:      "appDescription",
		Prompt:    &survey.Input{
			Message: "App description:",
			Default: skeletonConfig.AppDescription,
		},
		Validate:  survey.Required,
	},
	{
		Name:     "serverPort",
		Prompt:   &survey.Input{
			Message: "Web server port:",
			Default: strconv.Itoa(skeletonConfig.ServerPort),
		},
		Validate: survey.Required,
	},
}

func ConfigureInteractive(args []string) error {
	return survey.Ask(surveyQuestions, skeletonConfig)
}
