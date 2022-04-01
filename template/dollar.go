// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package template

import (
	"bytes"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"regexp"
	"strings"
)

type DollarRenderOptions struct {
	Variables map[string]string
}

func NewDollarRenderOptions() DollarRenderOptions {
	return DollarRenderOptions{
		Variables: map[string]string{},
	}
}

// DollarTemplate is a simplistic ES template renderer (variable names only, no expressions).
// To inject a variable, use `${variable}` syntax in the template.
type DollarTemplate []byte

func (t DollarTemplate) RenderBytes(options DollarRenderOptions) []byte {
	// Load the source
	contents := append([]byte{}, t[:]...)

	// Substitute variables
	variableInstanceRegex := regexp.MustCompile(`\${([^}]+)}`)
	substituted := make(types.StringSet)
	for _, variableInstance := range variableInstanceRegex.FindAllSubmatch(contents, -1) {
		variableName := strings.TrimSpace(string(variableInstance[1]))
		if substituted.Contains(variableName) {
			continue
		}

		variableValue, ok := options.Variables[variableName]
		variableKey := bytes.NewBufferString("${" + variableName + "}").Bytes()
		if ok {
			contents = bytes.ReplaceAll(contents, variableKey, []byte(variableValue))
		}

		substituted.Add(variableName)
	}

	return contents
}

func (t DollarTemplate) RenderString(options DollarRenderOptions) string {
	// Load the source
	contents := string(t)

	// Substitute variables
	variableInstanceRegex := regexp.MustCompile(`\${([^}]+)}`)
	substituted := make(types.StringSet)
	for _, variableInstance := range variableInstanceRegex.FindAllStringSubmatch(contents, -1) {
		variableName := strings.TrimSpace(variableInstance[1])
		if substituted.Contains(variableName) {
			continue
		}

		variableValue, ok := options.Variables[variableName]
		variableKey := "${" + variableName + "}"
		if ok {
			contents = strings.ReplaceAll(contents, variableKey, variableValue)
		}

		substituted.Add(variableName)
	}

	return contents
}
