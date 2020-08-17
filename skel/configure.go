package skel

import (
	"github.com/AlecAivazis/survey/v2"
	"os"
	"path"
	"strconv"
	"strings"
)

type SkeletonConfig struct {
	Archetype         string `survey:"generator" json:"generator"`
	TargetParent      string `survey:"targetParent" json:"-"`
	TargetDir         string `json:"-"`
	AppName           string `survey:"appName" json:"appName"`
	AppDisplayName    string `survey:"appDisplayName" json:"appDisplayName"`
	AppDescription    string `survey:"appDescription" json:"appDescription"`
	ServerPort        int    `survey:"serverPort" json:"serverPort"`
	ServerContextPath string `survey:"serverContextPath" json:"serverContextPath"`
	AppVersion        string `survey:"appVersion" json:"appVersion"`
	Repository        string `survey:"repository" json:"repository"`
	BeatProtocol      string `survey:"protocol" json:"protocol"`
	ServiceType       string `survey:"serviceType" json:"serviceType"`
	DeploymentGroup   string `survey:"deploymentGroup" json:"deploymentGroup"`
	KubernetesGroup   string `json:"kubernetesGroup"`
}

func (c SkeletonConfig) TargetDirectory() string {
	if c.TargetDir == "" {
		return path.Join(c.TargetParent, c.AppName)
	}
	return c.TargetDir
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
	Archetype:         archetypeKeyApp,
	TargetParent:      path.Join(os.Getenv("HOME"), "msx"),
	AppName:           "someservice",
	AppDisplayName:    "Some Microservice",
	AppDescription:    "Does Something",
	AppVersion:        "3.10.0",
	DeploymentGroup:   "something",
	ServerPort:        9999,
	ServerContextPath: "/some",
	Repository:        "cassandra",
	BeatProtocol:      "",
	ServiceType:       "",
}

var archetypeSurveyQuestions = map[string][]*survey.Question{
	archetypeKeyApp: {
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
			Validate: survey.Required,
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
	},
	archetypeKeyBeat: {
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
			Validate: survey.Required,
		},
		{
			Name: "protocol",
			Prompt: &survey.Input{
				Message: "Protocol:",
				Default: "icmp",
			},
			Validate: survey.Required,
		},
	},
	archetypeKeyServicePack: {
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
			Validate: survey.Required,
		},
		{
			Name: "deploymentGroup",
			Prompt: &survey.Input{
				Message: "Service Pack Name: ",
				Default: skeletonConfig.ServiceType,
			},
			Validate: survey.Required,
		},
		{
			Name: "appName",
			Prompt: &survey.Input{
				Message: "Microservice name:",
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
		{
			Name: "serviceType",
			Prompt: &survey.Input{
				Message: "Catalog Service Type: ",
				Default: skeletonConfig.ServiceType,
			},
			Validate: survey.Required,
		},
	},
}

var archetypeQuestions = []*survey.Question{
	{
		Name: "generator",
		Prompt: &survey.Select{
			Message: "Generate archetype:",
			Options: archetypes.DisplayNames(),
			Default: 0,
		},
	},
}

func ConfigureInteractive(args []string) error {
	var archetypeIndex int
	err := survey.Ask(archetypeQuestions, &archetypeIndex)
	if err != nil {
		return err
	}

	// Configure the archetype
	skeletonConfig.Archetype = archetypes.Key(archetypeIndex)
	var questions = archetypeSurveyQuestions[skeletonConfig.Archetype]
	err = survey.Ask(questions, skeletonConfig)
	if err != nil {
		return err
	}

	// Post-Process answers
	switch skeletonConfig.Archetype {
	case archetypeKeyApp:
		skeletonConfig.KubernetesGroup = "platformms"
		skeletonConfig.DeploymentGroup = skeletonConfig.AppName

	case archetypeKeyBeat:
		skeletonConfig.BeatProtocol = strings.ToLower(skeletonConfig.BeatProtocol)
		skeletonConfig.AppName = skeletonConfig.BeatProtocol + "beat"
		skeletonConfig.AppDescription = "Probes " + skeletonConfig.BeatProtocol
		skeletonConfig.AppDisplayName = strings.Title(skeletonConfig.AppName)
		skeletonConfig.ServerPort = 8080
		skeletonConfig.ServerContextPath = ""
		skeletonConfig.Repository = ""
		skeletonConfig.KubernetesGroup = "dataplatform"
		skeletonConfig.DeploymentGroup = skeletonConfig.AppName

	case archetypeKeyServicePack:
		skeletonConfig.KubernetesGroup = "servicepackms"
	}

	return nil
}
