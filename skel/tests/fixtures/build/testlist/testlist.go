// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

// Package testlist contains a central list of skel commands and tests
package testlist

type Test struct {
	Command      string // the command to run for the test
	Args         string // the arguments to pass to the command
	SpecialBuild bool   // if true, the test requires a special test and therfore a special build too
	Disabled     bool   // if true, the test is disabled, probably it needs to be fixed :(
}

// Tests uses names of the tests as the keys
var Tests = map[string]Test{
	"add-go-msx-dependency":         {Command: "add-go-msx-dependency", SpecialBuild: true},
	"completion":                    {Command: "completion", SpecialBuild: true},
	"generate-app":                  {Command: "generate-app"},
	"generate-build":                {Command: "generate-build"},
	"generate-certificate":          {Command: "generate-certificate", Disabled: true},
	"generate-channel":              {Command: "generate-channel", Args: "-d toad -m weasels", Disabled: true},
	"generate-channel-asyncapi":     {Command: "generate-channel-asyncapi", Args: "api/asyncapi.yaml", Disabled: true},
	"generate-channel-publisher":    {Command: "generate-channel-publisher", Args: "weasels", Disabled: true},
	"generate-channel-subscriber":   {Command: "generate-channel-subscriber", Args: "weasels", Disabled: true},
	"generate-deployment-variables": {Command: "generate-deployment-variables"},
	"generate-dockerfile":           {Command: "generate-dockerfile"},
	"generate-domain":               {Command: "generate-domain", Args: "toad", Disabled: true},
	"generate-domain-beats":         {Command: "generate-domain-beats", Args: "toad", Disabled: true},
	"generate-domain-openapi":       {Command: "generate-domain-openapi", Args: "toad", Disabled: true},
	"generate-domain-system":        {Command: "generate-domain-system", Args: "toadhall", Disabled: true},
	"generate-domain-tenant":        {Command: "generate-domain-tenant", Args: "toad", Disabled: true},
	"generate-git":                  {Command: "generate-git", Disabled: true},
	"generate-github":               {Command: "generate-github"},
	"generate-goland":               {Command: "generate-goland"},
	"generate-harness":              {Command: "generate-harness"},
	"generate-jenkins":              {Command: "generate-jenkins"},
	"generate-kubernetes":           {Command: "generate-kubernetes"},
	"generate-local":                {Command: "generate-local"},
	"generate-manifest":             {Command: "generate-manifest"},
	"generate-migrate":              {Command: "generate-migrate"},
	"generate-service-pack":         {Command: "generate-service-pack"},
	"generate-skel-json":            {Command: "generate-skel-json", Disabled: true},
	"generate-spui":                 {Command: "generate-spui", Disabled: true},
	"generate-test":                 {Command: "generate-test"},
	"generate-timer":                {Command: "generate-timer", Disabled: true},
	"generate-topic-publisher":      {Command: "generate-topic-publisher", Args: "weasels", Disabled: true},
	"generate-topic-subscriber":     {Command: "generate-topic-subscriber", Args: "weasels", Disabled: true},
	"generate-vscode":               {Command: "generate-vscode"},
	"generate-webservices":          {Command: "generate-webservices", Disabled: true},
	"help":                          {Command: "help", SpecialBuild: true},
	"version":                       {Command: "version", SpecialBuild: true},
}
