package skel

import (
	"github.com/AlecAivazis/survey/v2"
	"path"
	"strconv"
)

type SkeletonConfig struct {
	AppName        string `survey:"appName"`
	AppDisplayName string `survey:"appDisplayName"`
	AppDescription string `survey:"appDescription"`
	ServerPort     int    `survey:"serverPort"`
}

func (c SkeletonConfig) TargetDirectory() string {
	return path.Join("/Users/mcrawfo2/vms-3.1/demos", c.AppName)
}

var skeletonConfig = &SkeletonConfig{
	AppName:        "someservice",
	AppDisplayName: "Some Microservice",
	AppDescription: "Does Something",
	ServerPort:     9999,
}

var surveyQuestions = []*survey.Question{
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
