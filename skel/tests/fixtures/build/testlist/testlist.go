// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

// Package testlist contains a central list of skel commands and tests
package testlist

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel/tests/txtartools/clienv"
	"fmt"
	"gopkg.in/pipe.v2"
	"os"
	"os/exec"
	"path/filepath"
)

// FIXT obtained from the FIXT environment variable during init, it should be the path to the fixtures directory
var FIXT string

type SpecialRunType int

const (
	OrdinaryRun SpecialRunType = iota
	SpecRunPipe
	SpecRunFunction
)

type SpecialBuildType int

const (
	OrdinaryBuild   SpecialBuildType = iota // ordinary in fact :|
	SpecBuildNone                           // no build needed
	SpecBuildScript                         // build using a bash script
)

type Test struct {
	Command      string           // the command to run for the test
	Args         string           // the arguments to pass to the command
	SpecialBuild SpecialBuildType // test requires a special build
	SpecialRun   SpecialRunType   // the test is special
	SRPipe       pipe.Pipe        // the pipe if the special run is a pipe
	SRFunction   func() (ok bool) // the function if the special run is a function
	SBScript     string           // the script if the special build is a script
	Disabled     bool             // if true, the test is disabled, probably it needs to be fixed :(
	CmpGlobs     string           // list of globs to steer testing generation
	NoHelp       bool             // if true, the command does not need its help tested
	ExtraBefores []string         // extra files to make available before the test
}

// Tests is the list of tests to run
var Tests map[string]Test

func init() { // set defaults and help test, which is special :|, so very special

	wd, _ := os.Getwd()

	FIXT = os.Getenv("FIXT")
	if len(FIXT) == 0 {
		tests, err := trimPathBackTo("tests", wd)
		if err != nil {
			panic(fmt.Sprintf("Could not find tests directory in path: %s", wd))
		}
		FIXT = filepath.Join(tests, "fixtures")
		err = os.Setenv("FIXT", FIXT)
		if err != nil {
			panic(fmt.Sprintf("Could not find set fixt: %s", err))
		}
		fmt.Printf("🔻 Set FIXT to: %s\n", FIXT)
	}
	fmt.Printf("🔻 FIXT is: %s\n", os.Getenv("FIXT"))
	if !FileExists(FIXT) {
		panic(fmt.Sprintf("FIXT env var (%s) does not point to a directory that exists", FIXT))
	}

	var versionCheck = pipe.Line(
		pipe.Exec("skel", "version"),
		pipe.Exec("grep", `"Current build: \d*"`),
	)

	var completionCheck = pipe.Script(
		pipe.Exec("testscript", filepath.Join(FIXT, "final", "completions-test.txtar")),
	)

	// Tests uses names of the tests as the keys
	Tests = map[string]Test{
		"add-go-msx-dependency":         {Command: "add-go-msx-dependency"},
		"generate-app":                  {Command: "generate-app"},
		"generate-build":                {Command: "generate-build"},
		"generate-certificate":          {Command: "generate-certificate", CmpGlobs: "notsame:local/server.{crt,key} " + clienv.DefaultCmpGlob},
		"generate-channel":              {Command: "generate-channel", Args: "weasels"},
		"generate-channel-publisher":    {Command: "generate-channel-publisher", Args: "weasels"},
		"generate-channel-subscriber":   {Command: "generate-channel-subscriber", Args: "weasels"},
		"generate-deployment-variables": {Command: "generate-deployment-variables"},
		"generate-dockerfile":           {Command: "generate-dockerfile"},
		"generate-domain-beats":         {Command: "generate-domain-beats", Args: "toad"},
		"generate-domain-openapi":       {Command: "generate-domain-openapi", Args: "toad", Disabled: true},
		"generate-domain-system":        {Command: "generate-domain-system", Args: "toadhall"},
		"generate-domain-tenant":        {Command: "generate-domain-tenant", Args: "toad"},
		"generate-git":                  {Command: "generate-git"},
		"generate-github":               {Command: "generate-github"},
		"generate-goland":               {Command: "generate-goland"},
		"generate-harness":              {Command: "generate-harness"},
		"generate-jenkins":              {Command: "generate-jenkins"},
		"generate-kubernetes":           {Command: "generate-kubernetes"},
		"generate-local":                {Command: "generate-local"},
		"generate-manifest":             {Command: "generate-manifest"},
		"generate-migrate":              {Command: "generate-migrate"},
		"generate-service-pack":         {Command: "generate-service-pack"},
		"generate-skel-json":            {Command: "generate-skel-json"},
		"generate-spui":                 {Command: "generate-spui"},
		"generate-test":                 {Command: "generate-test"},
		"generate-timer":                {Command: "generate-timer", Args: "wabbit"},
		"generate-topic-publisher":      {Command: "generate-topic-publisher", Args: "weasels"},
		"generate-topic-subscriber":     {Command: "generate-topic-subscriber", Args: "weasels", CmpGlobs: "exists:go.sum " + clienv.DefaultCmpGlob},
		"generate-vscode":               {Command: "generate-vscode"},
		"generate-webservices":          {Command: "generate-webservices", Disabled: true},
		"version":                       {Command: "version", SpecialBuild: SpecBuildNone, SpecialRun: SpecRunPipe, SRPipe: versionCheck},
		"generate-channel-asyncapi": {Command: "generate-channel-asyncapi",
			Args: "api/asyncapi.yaml -a", ExtraBefores: []string{"api/asyncapi.yaml"}},
		"generate-domain": {Command: "generate-domain", Args: "toad",
			CmpGlobs: "exists:internal/toads/payloads_toad.go same:**"}, // snippets are emitted in random order
		"completions": {Command: "completion",
			SpecialBuild: SpecBuildScript, SBScript: filepath.Join(FIXT, "build", "make-completions-test.sh"),
			SpecialRun: SpecRunPipe, SRPipe: completionCheck},
		// help test is added in init below
	}

	// set the defaul cmp operation
	for k, v := range Tests {
		if v.CmpGlobs == "" {
			v.CmpGlobs = clienv.DefaultCmpGlob
		}
		Tests[k] = v
	}

	// add the help test
	var needhelp []string
	for k, v := range Tests {
		if !v.NoHelp {
			needhelp = append(needhelp, k)
		}
	}

	thf := func() (ok bool) {
		ok = true
		for _, needy := range needhelp {
			cmd := exec.Command("skel", "help", needy)
			err := cmd.Run()
			if err != nil {
				fmt.Printf("Help for %s failed ❌\n", needy)
				ok = false
			}
		}
		return ok
	}

	Tests["help"] = Test{Command: "help", SpecialBuild: SpecBuildNone, SpecialRun: SpecRunFunction, SRFunction: thf}

	// end of help test

}

func Test2Filename(kind, name, suffix string) string {
	return filepath.Join(FIXT, kind, name+suffix)
}

func FileExists(name string) bool {
	_, err := os.Stat(name)
	if err == os.ErrNotExist || err != nil {
		return false
	}
	return true
}

// trimPathTo trims the path back to the given dir by removing elements from the end
func trimPathBackTo(path, dir string) (base string, err error) {

	for {
		if filepath.Base(path) == dir {
			return path, nil
		}
		path = filepath.Dir(path)
		if path == "." {
			return "", fmt.Errorf("could not find %s in %s", dir, path)
		}
	}

}
