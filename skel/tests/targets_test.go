// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

//go:build integration

package tests

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel"
	"cto-github.cisco.com/NFV-BU/go-msx/skel/tests/txtartools/clienv"
	"flag"
	"github.com/stretchr/testify/require"
	"gopkg.in/pipe.v2"
	"os"
	"path"
	"testing"
)

// TestExecutor executes the test and either compares the generated results or stores them in the golden config
type TestExecutor interface {
	Test(t *testing.T, e TestWorkspace, d TargetTest)
}

type SkelTestSuite struct {
	Execution TestWorkspace
}

func (s *SkelTestSuite) SetupSuite(t *testing.T) {
	var err error

	s.Execution.Global.TestsRoot, err = os.Getwd()
	require.NoError(t, err)

	// create the suite folder
	s.Execution.Global.WorkDir, err = os.MkdirTemp("", "skel-integration-test-")
	require.NoError(t, err)
	t.Logf("Creating test suite folder %s", s.Execution.Global.WorkDir)

	s.Execution.Global.BinDir = path.Join(s.Execution.Global.WorkDir, "bin")
	require.NoError(t, os.MkdirAll(s.Execution.Global.BinDir, 0755))

	// deploy skel, txtar, testscript, mockery into the BinDir
	output, err := pipe.CombinedOutput(pipe.Script(
		pipe.SetEnvVar(envVarGOBIN, s.Execution.Global.BinDir),
		pipe.Exec("go", "install", "github.com/vektra/mockery/v2@v2.21.1"),
		pipe.Exec("go", "install", "cto-github.cisco.com/NFV-BU/go-msx/skel/tests/txtartools/txtargen"),
		pipe.Exec("go", "install", "cto-github.cisco.com/NFV-BU/go-msx/skel/tests/txtartools/txtarwrap"),
		pipe.Exec("go", "install", "cto-github.cisco.com/NFV-BU/go-msx/skel/tests/txtartools/txtarunwrap"),
		pipe.Exec("go", "install", "cto-github.cisco.com/NFV-BU/go-msx/skel/tests/txtartools/txtarcmp"),
		pipe.Exec("go", "install", "cto-github.cisco.com/NFV-BU/go-msx/cmd/skel"),
		pipe.Exec("go", "install", "github.com/rogpeppe/go-internal/cmd/testscript@v1.10.0"),
	))

	require.NoError(t, err, "Failed to install tools: %s\n%s", err, output)

	err = os.Setenv("PATH", s.Execution.Path())
	require.NoError(t, err, "Failed to set executable path: %s", err)
}

func (s *SkelTestSuite) TearDownSuite(t *testing.T) {
	if "" != os.Getenv("WORKSPACE") {
		return
	}

	err := os.RemoveAll(s.Execution.WorkDir)
	var err error
	if err != nil {
		t.Errorf("Failed to remove work directory %q", s.Execution.Global.WorkDir)
	}
}

func (s *SkelTestSuite) BeforeTest(t *testing.T, testName string) (execution TestWorkspace) {
	// set Execution
	execution.Global = s.Execution.Global
	execution.TestDir = path.Join(execution.Global.WorkDir, testName)
	execution.ProjectDir = path.Join(execution.TestDir, "someservice")

	// create the per-test folder
	require.NoError(t, os.MkdirAll(execution.ProjectDir, 0755))

	return execution
}

var generateGolden bool
var noOverwrite bool
var noParallel bool
var overrideAfter string

func init() {
	flag.BoolVar(&generateGolden, "skel.generate-golden", false,
		"Create new golden expectation files")
	flag.BoolVar(&noOverwrite, "skel.no-overwrite", false,
		"Do not overwrite existing golden expectation files")
	flag.BoolVar(&noParallel, "skel.no-parallel", false,
		"Do not run tests in parallel")
	flag.StringVar(&overrideAfter, "skel.golden-archive", "",
		"Specify a custom txtar output file")
}

func TestSkelTargets(t *testing.T) {
	var executor TestExecutor
	if generateGolden {
		executor = &GenerationExecutor{
			NoOverwrite:   noOverwrite,
			OverrideAfter: overrideAfter,
		}
	} else {
		executor = &ComparisonExecutor{}
	}

	testSuite := new(SkelTestSuite)

	testSuite.SetupSuite(t)
	t.Cleanup(func() { testSuite.TearDownSuite(t) })

	tests := []struct {
		Name string
		Test TargetTest
	}{
		{
			Name: "archetype-app",
			Test: TargetTest{
				NoRootBefore: true,
				RunInTestDir: true,
				BeforeFunction: func(t *testing.T, e TestWorkspace, test TargetTest) {
					require.NoError(t, e.WriteTestJson(
						skel.GenerateConfigFileName,
						skel.SkeletonConfig{
							Archetype:         "app",
							TargetParent:      e.TestDir,
							AppName:           "someservice",
							AppUUID:           "",
							AppDisplayName:    "Some Microservice",
							AppDescription:    "Does Something",
							ServerPort:        9999,
							DebugPort:         40000,
							ServerContextPath: "/some",
							AppVersion:        "5.0.0",
							Repository:        "cockroach",
							BeatProtocol:      "",
							ServiceType:       "",
							DeploymentGroup:   "someservice",
							KubernetesGroup:   "platformms",
							SlackChannel:      "go-msx-build",
							Trunk:             "main",
							ImageFile:         "msx.png",
							ScmHost:           "cto-github.cisco.com",
							ScmOrganization:   "NFV-BU",
						}))
				},
				CmpGlobs: clienv.DefaultCmpGlob +
					` ignorelines:someservice/.skel.json:.*"targetParent".*"`,
			},
		},
		{
			Name: "add-go-msx-dependency",
			Test: TargetTest{
				Args: []string{"add-go-msx-dependency"},
			},
		},
		{
			Name: "generate-app",
			Test: TargetTest{
				Args: []string{"generate-app"},
			},
		},
		{
			Name: "generate-build",
			Test: TargetTest{
				Args: []string{"generate-build"},
			},
		},
		{
			Name: "generate-certificate",
			Test: TargetTest{
				Args:     []string{"generate-certificate"},
				CmpGlobs: "notsame:local/server.{crt,key} " + clienv.DefaultCmpGlob,
			},
		},
		{
			Name: "generate-channel",
			Test: TargetTest{
				Args: []string{"generate-channel", "weasels"},
			},
		},
		{
			Name: "generate-channel-asyncapi",
			Test: TargetTest{
				Args: []string{"generate-channel-asyncapi", "api/asyncapi.yaml", "-a"},
				BeforeFunction: func(t *testing.T, e TestWorkspace, test TargetTest) {
					require.NoError(t, e.WriteBeforeFiles(test, []string{
						"api/asyncapi.yaml",
					}))
				},
			},
		},
		{
			Name: "generate-channel-publisher",
			Test: TargetTest{
				Args: []string{"generate-channel-publisher", "weasels"},
			},
		},
		{
			Name: "generate-channel-subscriber",
			Test: TargetTest{
				Args: []string{"generate-channel-subscriber", "weasels"},
			},
		},
		{
			Name: "generate-deployment-variables",
			Test: TargetTest{
				Args: []string{"generate-deployment-variables"},
			},
		},
		{
			Name: "generate-dockerfile",
			Test: TargetTest{
				Args: []string{"generate-dockerfile"},
			},
		},
		{
			Name: "generate-domain",
			Test: TargetTest{
				Args:     []string{"generate-domain", "toad"},
				CmpGlobs: "exists:internal/toads/payloads_toad.go same:**",
			},
		},
		{
			Name: "generate-domain-beats",
			Test: TargetTest{
				Args: []string{"generate-domain-beats", "toad"},
			},
		},
		{
			Name: "generate-domain-openapi",
			Test: TargetTest{
				Args:     []string{"generate-domain-openapi", "toad"},
				Disabled: true,
			},
		},
		{
			Name: "generate-domain-system",
			Test: TargetTest{
				Args: []string{"generate-domain-system", "toadhall"},
			},
		},
		{
			Name: "generate-domain-tenant",
			Test: TargetTest{
				Args: []string{"generate-domain-tenant", "toad"},
			},
		},
		{
			Name: "generate-git",
			Test: TargetTest{
				Args: []string{"generate-git"},
			},
		},
		{
			Name: "generate-github",
			Test: TargetTest{
				Args: []string{"generate-github"},
			},
		},
		{
			Name: "generate-goland",
			Test: TargetTest{
				Args: []string{"generate-goland"},
			},
		},
		{
			Name: "generate-harness",
			Test: TargetTest{
				Args: []string{"generate-harness"},
			},
		},
		{
			Name: "generate-jenkins",
			Test: TargetTest{
				Args: []string{"generate-jenkins"},
			},
		},
		{
			Name: "generate-kubernetes",
			Test: TargetTest{
				Args: []string{"generate-kubernetes"},
			},
		},
		{
			Name: "generate-local",
			Test: TargetTest{
				Args: []string{"generate-local"},
			},
		},
		{
			Name: "generate-manifest",
			Test: TargetTest{
				Args: []string{"generate-manifest"},
			},
		},
		{
			Name: "generate-migrate",
			Test: TargetTest{
				Args: []string{"generate-migrate"},
			},
		},
		{
			Name: "generate-service-pack",
			Test: TargetTest{
				Args: []string{"generate-service-pack"},
			},
		},
		{
			Name: "generate-skel-json",
			Test: TargetTest{
				Args: []string{"generate-skel-json"},
				CmpGlobs: clienv.DefaultCmpGlob +
					` ignorelines:.skel.json:.*"targetParent".*"`,
			},
		},
		{
			Name: "generate-spui",
			Test: TargetTest{
				Args: []string{"generate-spui"},
			},
		},
		{
			Name: "generate-test",
			Test: TargetTest{
				Args: []string{"generate-test"},
			},
		},
		{
			Name: "generate-timer",
			Test: TargetTest{
				Args: []string{"generate-timer", "wabbit"},
			},
		},
		{
			Name: "generate-topic-publisher",
			Test: TargetTest{
				Args: []string{"generate-topic-publisher", "weasels"},
			},
		},
		{
			Name: "generate-topic-subscriber",
			Test: TargetTest{
				Args:     []string{"generate-topic-subscriber", "weasels"},
				CmpGlobs: "exists:go.sum " + clienv.DefaultCmpGlob,
			},
		},
		{
			Name: "generate-vscode",
			Test: TargetTest{
				Args: []string{"generate-vscode"},
			},
		},
		{
			Name: "generate-webservices",
			Test: TargetTest{
				Args:     []string{"generate-webservices"},
				Disabled: true,
			},
		},
		{
			Name: "completion-bash",
			Test: TargetTest{
				Args:         []string{"completion", "bash"},
				SpecialBuild: SpecBuildStdout,
				SpecialRun:   SpecRunStdout,
			},
		},
		{
			Name: "completion-zsh",
			Test: TargetTest{
				Args:         []string{"completion", "zsh"},
				SpecialBuild: SpecBuildStdout,
				SpecialRun:   SpecRunStdout,
			},
		},
		{
			Name: "version",
			Test: TargetTest{
				Args:         []string{"version"},
				SpecialBuild: SpecBuildNone,
				SpecialRun:   SpecRunPipe,
				SRPipe: pipe.Line(
					pipe.Exec("skel", "version"),
					pipe.Exec("grep", `Current build: \d*`),
				),
			},
		},
		{
			Name: "help",
			Test: TargetTest{
				Args:         []string{"help"},
				SpecialBuild: SpecBuildNone,
				SpecialRun:   SpecRunPipe,
				SRPipe: pipe.Line(
					pipe.Exec("skel", "help"),
				),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			if !noParallel {
				t.Parallel()
			}

			tt.Test.Name = tt.Name
			execution := testSuite.BeforeTest(t, tt.Name)
			executor.Test(t, execution, tt.Test)
		})
	}
}
