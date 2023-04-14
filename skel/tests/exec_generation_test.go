// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package tests

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel/tests/txtartools/clienv"
	"gopkg.in/pipe.v2"
	"os"
	"path/filepath"
	"testing"
)

const (
	ignore      = "**/.git/** **/go.sum generate.json" // glob of files we will leave out of golden sets and tests
	envVarGOBIN = "GOBIN"
)

type GenerationExecutor struct {
	NoOverwrite   bool
	OverrideAfter string
}

func (g *GenerationExecutor) FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return true
}

func (g *GenerationExecutor) Test(t *testing.T, e TestWorkspace, test TargetTest) {
	name := test.Name

	if test.Disabled {
		t.Skipf("Skipping: Disabled: %s of command: %s %s ⏭", name, test.CommandName(), test.CommandArgs())
	}

	switch test.SpecialBuild {
	case SpecBuildNone:
		t.Skipf("Skipping    ⏭  special, no build required: test %s of command: %s %s", test.Name, test.CommandName(), test.CommandArgs())
	}

	goldenAfter := e.After(test)
	if g.OverrideAfter != "" {
		goldenAfter = filepath.Join(e.Fixtures(), g.OverrideAfter)
	}
	if g.NoOverwrite && g.FileExists(goldenAfter) {
		t.Skipf("Skipping    ⏭  Nooverwrite set & exists: %s for test %s of command: %s %s", goldenAfter, test.Name, test.CommandName(), test.CommandArgs())
	}

	relname, _ := filepath.Rel(e.Fixtures(), goldenAfter)
	t.Logf("Making test 🚧 skel %s, args:%v in 📁 %s", test.CommandName(), test.CommandArgs(), relname)

	var makeIt pipe.Pipe

	switch test.SpecialBuild {
	case SpecBuildScript:
		makeIt = pipe.Script(
			pipe.ChDir(e.TestDir),
			pipe.Exec(test.SBScript, test.Args...),
		)

	case SpecBuildStdout:
		currentStdout := filepath.Join(e.TestDir, "stdout")

		makeIt = pipe.Script(
			pipe.ChDir(e.TestDir), // parent of service directory
			pipe.Line(
				pipe.Exec("skel", test.Args...),
				pipe.WriteFile(currentStdout, 0644),
			),
			pipe.SetEnvVar(clienv.EnvIgnore, ignore),
			pipe.Line(
				pipe.Exec("txtarwrap", "."),
				pipe.WriteFile(goldenAfter, 0644),
			),
		)

	case OrdinaryBuild:
		var pipes []pipe.Pipe

		if !test.NoRootBefore {
			pipes = append(pipes,
				pipe.ChDir(e.TestDir), // parent of service directory
				pipe.Exec("txtarunwrap", e.Before(), "."),
			)
		}

		if test.BeforeFunction != nil {
			test.BeforeFunction(t, e, test)
		}

		skelArgs := append([]string{"--allow-dirty"}, test.Args...)

		if test.RunInTestDir {
			pipes = append(pipes,
				pipe.ChDir(e.TestDir),
			)
		} else {
			pipes = append(pipes,
				pipe.ChDir(e.ProjectDir),
				pipe.Exec("go", "get", "./..."),
			)
		}

		pipes = append(pipes,
			pipe.Exec("skel", skelArgs...),
			pipe.SetEnvVar(clienv.EnvIgnore, ignore),
			pipe.Line(
				pipe.Exec("txtarwrap", "."),
				pipe.WriteFile(goldenAfter, 0644),
			),
		)

		makeIt = pipe.Script(pipes...)
	}

	barf, err := pipe.CombinedOutput(makeIt)
	if err != nil {
		t.Fatalf("Failed making test %s: %s\n%s\n", test.Name, err, barf)
	}
}
