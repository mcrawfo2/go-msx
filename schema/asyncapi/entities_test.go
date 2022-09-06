// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package asyncapi

import (
	"encoding/json"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"testing"
)

func TestParse(t *testing.T) {
	yamlBytes, _ := ioutil.ReadFile("testdata/compliance-asyncapi.yaml")
	jsonBytes, _ := yaml.YAMLToJSON(yamlBytes)

	var spec Spec
	err := json.Unmarshal(jsonBytes, &spec)
	if err != nil {
		t.Error(err.Error())
	}
}
