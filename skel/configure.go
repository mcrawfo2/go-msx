package skel

import (
	"github.com/AlecAivazis/survey/v2"
	"os"
	"path"
	"strconv"
	"strings"
)

type SkeletonConfig struct {
	Generator         string `json:"generator"`
	TargetParent      string `survey:"targetParent" json:"-"`
	AppName           string `survey:"appName" json:"appName"`
	AppDisplayName    string `survey:"appDisplayName" json:"appDisplayName"`
	AppDescription    string `survey:"appDescription" json:"appDescription"`
	ServerPort        int    `survey:"serverPort" json:"serverPort"`
	ServerContextPath string `survey:"serverContextPath" json:"serverContextPath"`
	AppVersion        string `survey:"appVersion" json:"appVersion"`
	Repository        string `survey:"repository" json:"repository"`
	BeatProtocol      string `survey:"protocol" json:"protocol"`
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

func (c SkeletonConfig) RepositoryQueryFileExtension() string {
	queryFileExtension := "cql"
	if skeletonConfig.Repository == "cockroach" {
		queryFileExtension = "sql"
	}
	return queryFileExtension
}

var skeletonConfig = &SkeletonConfig{
	Generator:         "app",
	TargetParent:      path.Join(os.Getenv("HOME"), "msx"),
	AppName:           "someservice",
	AppDisplayName:    "Some Microservice",
	AppDescription:    "Does Something",
	AppVersion:        "3.10.0",
	ServerPort:        9999,
	ServerContextPath: "/some",
	Repository:        "cassandra",
	BeatProtocol:      "",
}

func appSurveyQuestions() []*survey.Question {
	return []*survey.Question{
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
}

func beatsSurveyQuestions() []*survey.Question {
	skeletonConfig.Generator = "beat"
	skeletonConfig.AppName = "heartbeat"
	skeletonConfig.AppDisplayName = "ICMP Probe"
	skeletonConfig.AppDescription = "MSX ICMP probe"
	skeletonConfig.ServerPort = 8080
	skeletonConfig.ServerContextPath = ""
	skeletonConfig.BeatProtocol = "icmp"
	skeletonConfig.Repository = ""

	return []*survey.Question{

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
			Name: "protocol",
			Prompt: &survey.Input{
				Message: "Network protocol:",
				Default: skeletonConfig.BeatProtocol,
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
	}
}
func ConfigureInteractive(args []string) error {
	questions := appSurveyQuestions()
	if generateBeat {
		questions = beatsSurveyQuestions()
	}

	return survey.Ask(questions, skeletonConfig)
}
