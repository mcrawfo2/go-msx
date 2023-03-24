// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package asyncapi

import (
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/schema/asyncapi"
	"cto-github.cisco.com/NFV-BU/go-msx/skel"
	"encoding/json"
	"github.com/ghodss/yaml"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

var logger = log.NewLogger("skel.asyncapi")

func init() {
	// Swiss Army Knife
	cmd := skel.AddTarget("generate-channel", "Create async channel", GenerateChannel)
	cmd.Flags().StringSliceVarP(&generatorConfig.Messages, "message", "m", nil, "Message name")
	cmd.Flags().BoolVarP(&generatorConfig.Publisher, "publisher", "p", false, "Generate channel and message publisher")
	cmd.Flags().BoolVarP(&generatorConfig.Subscriber, "subscriber", "s", false, "Generate channel and message subscriber")
	//cmd.Flags().StringVarP(&generatorConfig.Domain, "domain", "d", "", "Attach subscriber to domain's application service")

	// Publisher
	cmd = skel.AddTarget("generate-channel-publisher", "Create async channel publisher", GenerateChannelPublisher)
	cmd.Args = cobra.ExactArgs(1)
	cmd.Flags().StringVarP(&operationConfig.MessageName, "message", "m", "", "Message name to add to multi-message channel")

	// Subscriber
	cmd = skel.AddTarget("generate-channel-subscriber", "Create async channel subscriber", GenerateChannelSubscriber)
	cmd.Args = cobra.ExactArgs(1)
	cmd.Flags().StringVarP(&operationConfig.MessageName, "message", "m", "", "Message name to add to multi-message channel")
	cmd.Flags().StringVarP(&operationConfig.Domain, "domain", "d", "", "Attach subscriber to domain's application service")

	// Inverted Spec
	cmd = skel.AddTarget("generate-channel-asyncapi", "Create stream from AsyncApi 2.4 specification", GenerateChannelAsyncApi)
	cmd.Use = "generate-channel-asyncapi <document> [<channel-name> | --all] [--invert]"
	cmd.Args = cobra.RangeArgs(0, 2)
	cmd.Flags().BoolVarP(&specificationConfig.AllChannels, "all", "a", false, "Generate streams for all channels in document")
	cmd.Flags().BoolVarP(&specificationConfig.Invert, "invert", "i", false, "Swap publish and subscribe operations")

}

// GenerateChannel is the CLI entry point for generating generic asyncapi-enabled channel modules
func GenerateChannel(args []string) (err error) {
	if len(args) != 1 {
		// Incorrect channel name
		if err = configureInteractive(); err != nil {
			return err
		}
	} else {
		generatorConfig.ChannelName = args[0]
	}

	for i, message := range generatorConfig.Messages {
		generatorConfig.Messages[i] = strcase.ToLowerCamel(message)
		generatorConfig.Deep = true
	}

	return newTemplatedGenerator(generatorConfig).Generate()
}

func GenerateChannelPublisher(args []string) (err error) {
	cfg := GeneratorConfig{
		ChannelName: args[0],
		Multi:       operationConfig.MessageName != "",
		Publisher:   true,
		Deep:        true,
	}

	if operationConfig.MessageName == "" {
		cfg.Messages = []string{
			messageName(channelShortName(cfg.ChannelName), "Request"),
		}
	} else {
		cfg.Messages = []string{
			strcase.ToLowerCamel(operationConfig.MessageName),
		}
	}

	return newTemplatedGenerator(cfg).Generate()
}

func GenerateChannelSubscriber(args []string) (err error) {
	cfg := GeneratorConfig{
		ChannelName: args[0],
		Subscriber:  true,
		Domain:      operationConfig.Domain,
		Deep:        true,
	}

	if operationConfig.MessageName == "" {
		cfg.Messages = []string{
			messageName(channelShortName(cfg.ChannelName), "Response"),
		}
	} else {
		cfg.Messages = []string{
			strcase.ToLowerCamel(operationConfig.MessageName),
		}
		cfg.Multi = true
	}

	if cfg.Domain == "" {
		cfg.Domain = "Unknown"
	}

	return newTemplatedGenerator(cfg).Generate()
}

func GenerateChannelAsyncApi(args []string) error {
	if len(args) > 0 {
		specificationConfig.Document = args[0]
	}

	if len(args) > 1 {
		if specificationConfig.AllChannels {
			return errors.New("Cannot specify channel name(s) and --all")
		}
		specificationConfig.ChannelName = args[1]
	}

	hasChannels := specificationConfig.ChannelName != "" || specificationConfig.AllChannels
	if !hasChannels || specificationConfig.Document == "" {
		if err := configureInteractiveSpec(); err != nil {
			return err
		}
	}

	spec, err := loadSpec(specificationConfig.Document)
	if err != nil {
		return errors.Wrap(err, "Failed to load AsyncApi spec")
	}

	var channels []string
	if specificationConfig.ChannelName != "" {
		channels = []string{specificationConfig.ChannelName}
	}

	return NewGenerator(spec).
		WithInvert(specificationConfig.Invert).
		GenerateChannels(channels...)
}

func loadSpec(filename string) (spec asyncapi.Spec, err error) {
	var specBytes []byte
	if strings.HasPrefix(filename, "http://") || strings.HasPrefix(filename, "https://") {
		client := &http.Client{Transport: &http.Transport{}}
		req, _ := http.NewRequest("GET", filename, http.NoBody)
		var resp *http.Response
		resp, err = client.Do(req)
		if err != nil {
			err = errors.Wrap(err, "Failed to locate spec")
			return
		}
		defer resp.Body.Close()

		specBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			err = errors.Wrap(err, "Failed to download spec")
			return
		}
	} else {
		specBytes, err = ioutil.ReadFile(filename)
		if err != nil {
			err = errors.Wrap(err, "Failed to load spec")
			return
		}

	}

	switch path.Ext(filename) {
	case ".yml", ".yaml":
		specBytes, err = yaml.YAMLToJSON(specBytes)
		if err != nil {
			return
		}
		fallthrough
	case ".json":
		err = json.Unmarshal(specBytes, &spec)
	default:
		err = errors.Errorf("Unknown spec file format: %s", path.Ext(filename))
	}

	return
}
