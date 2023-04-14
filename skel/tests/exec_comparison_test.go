// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package tests

import (
	"bytes"
	"cto-github.cisco.com/NFV-BU/go-msx/skel/tests/txtartools/clienv"
	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"github.com/stretchr/testify/require"
	"gopkg.in/pipe.v2"
	"os"
	"path/filepath"
	"testing"
)

type ComparisonExecutor struct{}

func (c *ComparisonExecutor) PatchedTxtAr(t *testing.T, e TestWorkspace, test TargetTest) []byte {
	sourceBytes, err := os.ReadFile(e.Before())
	if err != nil {
		t.Fatalf("Failed parsing golden before %s: %s", test.Name, err)
	}

	patchBytes, err := os.ReadFile(e.AfterDiff(test))
	if err != nil {
		t.Fatalf("Failed reading after patch %s: %s", test.Name, err)
	}

	patches, _, err := gitdiff.Parse(bytes.NewReader(patchBytes))
	if err != nil {
		t.Fatalf("Missing patch file %s: %s", test.Name, err)
	}

	if len(patches) == 0 {
		return sourceBytes
	}

	dest := new(bytes.Buffer)

	err = gitdiff.Apply(dest, bytes.NewReader(sourceBytes), patches[0])
	if err != nil {
		t.Fatalf("Failed applying golden patch file %s: %s", test.Name, err)
	}

	return dest.Bytes()
}

func (c *ComparisonExecutor) Test(t *testing.T, e TestWorkspace, test TargetTest) {
	name := test.Name

	if test.Disabled {
		t.Skipf("Skipping: Disabled: %s of command: %s %s ⏭", name, test.CommandName(), test.CommandArgs())
		return
	}

	goldenAfter := e.After(test)

	// check if we can use patch-compressed txtar
	var goldenAfterBuffer *bytes.Buffer
	if test.SpecialRun == OrdinaryRun &&
		!test.NoRootBefore &&
		test.BeforeFunction == nil {
		goldenAfterBytes := c.PatchedTxtAr(t, e, test)
		goldenAfterBuffer = bytes.NewBuffer(goldenAfterBytes)
	}

	if goldenAfterBuffer == nil {
		// uncompressed txtar
		switch test.SpecialRun {
		case OrdinaryRun, SpecRunStdout:
			// we need a golden file
			_, err := os.Stat(goldenAfter)
			if os.IsNotExist(err) {
				t.Skipf("Skipping %s: No golden result %q found.   ⏭", test.Name, goldenAfter)
			}
			require.NoError(t, err, "Failed to locate golden result: %s", test.Name)
		}
	}

	t.Logf("Running: test %s of command: %s %v", name, test.CommandName(), test.CommandArgs())

	allArgs := append([]string{"--allow-dirty"}, test.Args...)

	var runIt pipe.Pipe
	switch test.SpecialRun {
	case OrdinaryRun: // it's a normal test
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
			// Perform the generation
			pipe.Exec("skel", allArgs...),
		)

		// Compare the result with the golden expectation
		if goldenAfterBuffer == nil {
			// stock txtar file
			pipes = append(pipes,
				pipe.SetEnvVar(clienv.EnvCmp, test.Globs()),
				pipe.Exec("txtarcmp", goldenAfter, "."))
		} else {
			// patch-generated txtar
			pipes = append(pipes,
				pipe.ChDir(e.TestDir),
				pipe.Line(
					pipe.SetEnvVar(clienv.EnvCmp, test.Globs()),
					pipe.Read(goldenAfterBuffer),
					pipe.Exec("txtarcmp", "-", "."),
				))
		}

		runIt = pipe.Script(pipes...)

	case SpecRunStdout:
		currentStdout := filepath.Join(e.TestDir, "stdout")

		runIt = pipe.Script(
			pipe.ChDir(e.TestDir), // parent of service directory
			pipe.Line(
				// Perform the generation
				pipe.Exec("skel", allArgs...),
				pipe.WriteFile(currentStdout, 0644),
			),
			// Compare the result with the golden expectation
			pipe.SetEnvVar(clienv.EnvCmp, test.Globs()),
			pipe.Exec("txtarcmp", goldenAfter, "."),
		)

	case SpecRunPipe: // it's a custom pipe test
		runIt = pipe.Script(
			pipe.ChDir(e.ProjectDir),
			pipe.SetEnvVar(clienv.EnvCmp, test.Globs()),
			test.SRPipe,
		)

	case SpecRunFunction: // it's a go function test
		ok := test.SRFunction()
		if !ok {
			require.True(t, ok, "Failed running test %s", name)
		}
		return
	}

	output, err := pipe.CombinedOutput(runIt)
	require.NoError(t, err, "Failed running test %s: %s\n%s", name, err, output)
}
