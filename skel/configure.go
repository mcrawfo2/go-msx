// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package skel

import (
	"errors"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"golang.org/x/text/cases"
)

var ErrUserCancelled = errors.New("user cancelled")

type SkeletonConfig struct {
	Archetype         string `survey:"generator" json:"generator"`
	TargetParent      string `survey:"targetParent" json:"targetParent"`
	TargetDir         string `json:"-"`
	AppName           string `survey:"appName" json:"appName"`
	AppUUID           string `survey:"appUUID" json:"appUUID"`
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
	SlackChannel      string `survey:"slackChannel" json:"slackChannel"`
	Trunk             string `survey:"trunk" json:"trunk"`
	ImageFile         string `survey:"imageFile" json:"imageFile"`
}

func (c SkeletonConfig) TargetDirectory() string {

	if c.TargetDir != "" {
		return c.TargetDir
	}

	switch c.Archetype {
	case archetypeKeyApp:
		return path.Join(c.TargetParent, c.AppName)
	case archetypeKeyBeat:
		return path.Join(c.TargetParent, c.BeatProtocol+"beat")
	case archetypeKeyServicePack:
		return path.Join(c.TargetParent, c.AppName)
	case archetypeKeySPUI:
		return path.Join(c.TargetParent, c.AppName+"-ui")
	}
	return path.Join(c.TargetParent, c.AppName)

}

func (c SkeletonConfig) AppMigrateVersion() string {
	return "V" + strings.ReplaceAll(c.AppVersion, ".", "_")
}

func (c SkeletonConfig) AppPackageUrl() string {
	return path.Join("cto-github.cisco.com", "NFV-BU", c.AppName)
}

func (c SkeletonConfig) ApiPackageUrl() string {
	return path.Join(c.AppPackageUrl(), "pkg", "api")
}

func (c SkeletonConfig) AppShortName() string {
	noService := strings.TrimSuffix(c.AppName, "service")
	noBeat := strings.TrimSuffix(noService, "beat")
	noUi := strings.TrimSuffix(noBeat, "-ui")
	return noUi
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
	TargetParent:      func() string { wd, _ := os.Getwd(); return wd }(),
	AppName:           "someservice",
	AppDisplayName:    "Some Microservice",
	AppDescription:    "Does Something",
	AppVersion:        "5.0.0",
	DeploymentGroup:   "something",
	ServerPort:        9999,
	ServerContextPath: "/some",
	Repository:        "cockroach",
	BeatProtocol:      "",
	ServiceType:       "",
	SlackChannel:      "go-msx-build",
	Trunk:             "main",
	ImageFile:         "msx.png",
}

func Config() *SkeletonConfig {
	return skeletonConfig
}

var archetypeSurveyQuestions = map[string][]*survey.Question{
	archetypeKeyApp: {
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
		{
			Name: "slackChannel",
			Prompt: &survey.Input{
				Message: "Build notifications slack channel:",
				Default: skeletonConfig.SlackChannel,
			},
			Validate: survey.Required,
		},
		{
			Name: "trunk",
			Prompt: &survey.Input{
				Message: "Primary branch name:",
				Default: skeletonConfig.Trunk,
			},
			Validate: survey.Required,
		},
	},
	archetypeKeyBeat: {
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
		{
			Name: "slackChannel",
			Prompt: &survey.Input{
				Message: "Build notifications slack channel:",
				Default: skeletonConfig.SlackChannel,
			},
			Validate: survey.Required,
		},
		{
			Name: "trunk",
			Prompt: &survey.Input{
				Message: "Primary branch name:",
				Default: skeletonConfig.Trunk,
			},
			Validate: survey.Required,
		},
	},
	archetypeKeyServicePack: {
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
				Message: "Service Pack Name:",
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
				Message: "Catalog Service Type:",
				Default: skeletonConfig.ServiceType,
			},
			Validate: survey.Required,
		},
		{
			Name: "slackChannel",
			Prompt: &survey.Input{
				Message: "Build notifications slack channel:",
				Default: skeletonConfig.SlackChannel,
			},
			Validate: survey.Required,
		},
		{
			Name: "trunk",
			Prompt: &survey.Input{
				Message: "Primary branch name:",
				Default: skeletonConfig.Trunk,
			},
			Validate: survey.Required,
		},
	},
	archetypeKeySPUI: {
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
			Name: "serviceType",
			Prompt: &survey.Input{
				Message: "Catalog Service Type:",
				Default: skeletonConfig.ServiceType,
			},
			Validate: survey.Required,
		},
		{
			Name: "slackChannel",
			Prompt: &survey.Input{
				Message: "Build notifications slack channel:",
				Default: skeletonConfig.SlackChannel,
			},
			Validate: survey.Required,
		},
		{
			Name: "trunk",
			Prompt: &survey.Input{
				Message: "Primary branch name:",
				Default: skeletonConfig.Trunk,
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

// ConfigureInteractive is the only entry point to the menu UI
func ConfigureInteractive() error {

	var archetypeIndex int
	err := survey.Ask(archetypeQuestions, &archetypeIndex) // determine the archetype
	if err != nil {
		return err
	}

	// Ask for the TargetParent and compute the target from it and the archetype
	targetPQ := &survey.Input{
		Message: "Project Parent Directory:",
		Default: skeletonConfig.TargetParent,
	}
	err = survey.AskOne(targetPQ, &skeletonConfig.TargetParent, survey.WithValidator(survey.Required))
	if err != nil {
		return err
	}

	target := skeletonConfig.TargetDirectory()
	logger.Debugf("Target directory: %s", target)

	carryOn, err := GitCheckAsk(target)
	if err != nil || carryOn == false {
		return ErrUserCancelled
	}

	// Configure the archetype using the questions for it
	skeletonConfig.Archetype = archetypes.Key(archetypeIndex)

	var questions = archetypeSurveyQuestions[skeletonConfig.Archetype]
	err = survey.Ask(questions, skeletonConfig)
	if err != nil {
		return err
	}

	// Post-Process answers
	caser := cases.Title(TitlingLanguage)
	switch skeletonConfig.Archetype {
	case archetypeKeyApp:
		skeletonConfig.KubernetesGroup = "platformms"
		skeletonConfig.DeploymentGroup = skeletonConfig.AppName

	case archetypeKeyBeat:
		skeletonConfig.BeatProtocol = strings.ToLower(skeletonConfig.BeatProtocol)
		skeletonConfig.AppName = skeletonConfig.BeatProtocol + "beat"
		skeletonConfig.AppDescription = "Probes " + skeletonConfig.BeatProtocol
		skeletonConfig.AppDisplayName = caser.String(skeletonConfig.AppName)
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
