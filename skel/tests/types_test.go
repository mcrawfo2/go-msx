// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package tests

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel/tests/txtartools/clienv"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/pipe.v2"
	"os"
	"path"
	"path/filepath"
	"testing"
)

// TestWorkspace describes the environment for the test run
type TestWorkspace struct {
	// Global is shared across all tests
	Global struct {
		TestsRoot string // root integration tests folder
		WorkDir   string // parent of where the tests are run
		BinDir    string // where the skel and test binaries are installed
	}
	TestDir    string // single integration test folder
	ProjectDir string // where the project is generated
}

func (e TestWorkspace) Fixtures() string {
	return filepath.Join(e.Global.TestsRoot, "fixtures")
}

func (e TestWorkspace) BeforeFixtures() string {
	return filepath.Join(e.Fixtures(), "before")
}

func (e TestWorkspace) AfterFixtures() string {
	return filepath.Join(e.Fixtures(), "golden")
}

func (e TestWorkspace) Before() string {
	return filepath.Join(e.BeforeFixtures(), "plain-root.txtar")
}

func (e TestWorkspace) WriteTestJson(testSubPath string, value any) error {
	outFileName := filepath.Join(e.TestDir, testSubPath)
	outDir := filepath.Dir(outFileName)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(value, "", "    ")
	if err != nil {
		return err
	}

	if err = os.WriteFile(outFileName, data, 0644); err != nil {
		return err
	}

	return nil
}

func (e TestWorkspace) archiveName(test TargetTest) string {
	return test.Name + "-test.txtar"
}

func (e TestWorkspace) After(test TargetTest) string {
	return path.Join(e.AfterFixtures(), e.archiveName(test))
}

func (e TestWorkspace) diffName(test TargetTest) string {
	return test.Name + "-test.patch"
}

func (e TestWorkspace) AfterDiff(test TargetTest) string {
	return path.Join(e.AfterFixtures(), e.diffName(test))
}

func (e TestWorkspace) Path() string {
	return fmt.Sprintf("%s%c%s", // include all our binaries
		e.Global.BinDir,
		os.PathListSeparator,
		os.Getenv("PATH"))
}

func (e TestWorkspace) WriteBeforeFiles(test TargetTest, befores []string) error {
	for _, extraBefore := range befores {
		// Read the source file
		contents, err := os.ReadFile(filepath.Join(e.BeforeFixtures(), extraBefore))
		if err != nil {
			return errors.Wrapf(err, "Failed loading extra before %q for test %s", extraBefore, test.Name)
		}

		// Write the output file
		targetFile := filepath.Join(e.ProjectDir, extraBefore)
		targetDir := filepath.Dir(targetFile)
		err = os.MkdirAll(targetDir, 0755)
		if err != nil {
			return errors.Wrapf(err, "Failed to create directory for %s for test %s", extraBefore, test.Name)
		}

		err = os.WriteFile(targetFile, contents, 0644)
		if err != nil {
			return errors.Wrapf(err, "Failed to write file for %s for test %s", extraBefore, test.Name)
		}
	}

	return nil
}

type SpecialRunType int

const (
	OrdinaryRun SpecialRunType = iota
	SpecRunPipe
	SpecRunFunction
	SpecRunStdout
)

type SpecialBuildType int

const (
	OrdinaryBuild   SpecialBuildType = iota // ordinary in fact :|
	SpecBuildNone                           // no build needed
	SpecBuildScript                         // build using a bash script
	SpecBuildStdout                         // build from script stdout
)

type TargetTest struct {
	Name     string // the name of the test
	Disabled bool   // if true, the test is disabled, probably it needs to be fixed :(

	Args []string // the arguments to pass to the command

	NoRootBefore   bool                                                 // do not make root files available before the test
	BeforeFunction func(t *testing.T, e TestWorkspace, test TargetTest) // execute some arbitrary code before the test

	// TODO: Remove SBScript, SpecBuildScript
	SpecialBuild SpecialBuildType // test requires a special build
	SBScript     string           // the script if the special build is a script

	// TODO: Remove SpecialRun, SRFunction, SpecRunFunction
	SpecialRun   SpecialRunType   // the test is special
	SRPipe       pipe.Pipe        // the pipe if the special run is a pipe
	SRFunction   func() (ok bool) // the function if the special run is a function
	CmpGlobs     string           // list of globs to steer testing generation
	RunInTestDir bool             // alternative run directory
}

func (t TargetTest) CommandName() string {
	if len(t.Args) == 0 {
		return t.Name
	}
	return t.Args[0]
}

func (t TargetTest) CommandArgs() []string {
	if len(t.Args) == 0 {
		return nil
	}
	return t.Args[1:]
}

func (t TargetTest) Globs() string {
	if t.CmpGlobs == "" {
		return clienv.DefaultCmpGlob
	}
	return t.CmpGlobs
}
