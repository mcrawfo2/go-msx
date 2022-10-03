// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package asyncapi

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel"
	"encoding/json"
	"github.com/AlecAivazis/survey/v2"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/iancoleman/strcase"
	"strings"
)

type GeneratorConfig struct {
	ChannelName string
	Messages    []string
	Multi       bool
	Publisher   bool
	Subscriber  bool
	Domain      string
	Deep        bool
}

func (c GeneratorConfig) ChannelMessageName(suffix string) string {
	channelNoTopic := channelShortName(c.ChannelName)
	channelCamel := strcase.ToCamel(strings.ToLower(channelNoTopic))
	messageName := channelCamel + strcase.ToCamel(suffix)
	return messageName
}

func (c GeneratorConfig) String() string {
	data, _ := json.Marshal(c)
	return string(data)
}

var generatorConfig GeneratorConfig

// configure asks for channel and message details
func configureInteractive() (err error) {
	skeletonConfig := skel.Config()

	// Channel

	if generatorConfig.ChannelName == "" {
		appShortName := skeletonConfig.AppShortName()
		appScreamingShortName := strings.ToUpper(appShortName)
		generatorConfig.ChannelName = appScreamingShortName + "_TOPIC"
	}

	err = survey.AskOne(
		&survey.Input{
			Message: "Channel name",
			Default: generatorConfig.ChannelName,
			Help:    "The name of the channel to generate.  Generally matches the Kafka Topic name.",
		},
		&generatorConfig.ChannelName)
	if err != nil {
		return err
	}

	// Publish / Subscribe
	var pubsub string
	err = survey.AskOne(
		&survey.Select{
			Message: "Will this channel publish or subscribe?:",
			Options: []string{"Publish", "Subscribe"},
		},
		&pubsub)
	if err != nil {
		return err
	}

	switch strings.ToLower(pubsub) {
	case OperationPublish:
		generatorConfig.Publisher = true
	case OperationSubscribe:
		generatorConfig.Subscriber = true
	}

	if len(generatorConfig.Messages) == 0 {
		if generatorConfig.Publisher && generatorConfig.Subscriber {
			generatorConfig.Messages = []string{
				skeletonConfig.AppShortName() + "Event",
			}
		} else if generatorConfig.Subscriber {
			generatorConfig.Messages = []string{
				skeletonConfig.AppShortName() + "Request",
			}
		} else if generatorConfig.Publisher {
			generatorConfig.Messages = []string{
				skeletonConfig.AppShortName() + "Response",
			}
		}
	}

	var messages = strings.Join(generatorConfig.Messages, "\n")
	err = survey.AskOne(
		&survey.Multiline{
			Message: "Enter message types, one per line:\n",
			Help:    "Message identifiers",
			Default: messages,
		},
		&messages)
	if err != nil {
		return err
	}
	messages = strings.TrimSpace(messages)
	generatorConfig.Messages = strings.Split(messages, "\n")

	generatorConfig.Multi = len(generatorConfig.Messages) > 1

	if generatorConfig.Subscriber {
		err = survey.AskOne(
			&survey.Input{
				Message: "Subscriber domain",
				Help:    "Connects the subscriber to the domain's application service",
			},
			&generatorConfig.Domain,
			survey.WithValidator(func(ans interface{}) error {
				return validation.Validate(ans, validation.Required)
			}))
		if err != nil {
			return err
		}
	}

	generatorConfig.Deep = true

	return nil
}

type OperationConfig struct {
	MessageName string
	Domain      string
}

var operationConfig OperationConfig

type SpecificationConfig struct {
	Document    string
	ChannelName string
	AllChannels bool
	Invert      bool
}

var specificationConfig SpecificationConfig

func configureInteractiveSpec() (err error) {
	document := "api/asyncapi.yaml"
	if specificationConfig.Document != "" {
		document = specificationConfig.Document
	}
	err = survey.AskOne(
		&survey.Input{
			Message: "Specification file or url:",
			Help:    "Location of local or remote AsyncApi specification from which to generate components",
			Default: document,
		},
		&specificationConfig.Document)
	if err != nil {
		return err
	}

	spec, err := loadSpec(specificationConfig.Document)
	if err != nil {
		return err
	}

	if specificationConfig.AllChannels {
		return nil
	}

	var channelNames []string
	for channelName := range spec.Channels {
		channelNames = append(channelNames, channelName)
	}

	var defaultChannel interface{}
	if specificationConfig.ChannelName != "" {
		defaultChannel = specificationConfig.ChannelName
	}
	err = survey.AskOne(
		&survey.Select{
			Message: "Channel for which to generate components:",
			Options: channelNames,
			Default: defaultChannel,
		},
		&specificationConfig.ChannelName)
	if err != nil {
		return err
	}

	return nil
}
