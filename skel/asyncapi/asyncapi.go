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
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"path"
)

var logger = log.NewLogger("skel.asyncapi")

func init() {
	cmd := skel.AddTarget("generate-streams-asyncapi", "Create stream from AsyncApi 2.4 manifest", GenerateStreamAsyncApi)
	cmd.Use = "generate-streams-asyncapi <document> [<channel-name> | --all]"
	cmd.Args = cobra.RangeArgs(1, 2)
	cmd.Flags().BoolP("all", "a", false, "Generate streams for all channels in document")
}

func GenerateStreamAsyncApi(args []string) error {
	doc := args[0]

	channels := []string{}
	if len(args) == 2 {
		channels = append(channels, args[1])
	}

	spec, err := loadSpec(doc)
	if err != nil {
		return errors.Wrap(err, "Failed to load AsyncApi spec")
	}

	return NewGenerator(spec).GenerateChannels(channels...)
}

func loadSpec(filename string) (spec asyncapi.Spec, err error) {
	specBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return
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
