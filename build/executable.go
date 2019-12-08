package build

import (
	"cto-github.cisco.com/NFV-BU/go-msx/exec"
	"os"
	"path"
	"strings"
)

func init() {
	AddTarget("build-executable", "Build the binary executable", BuildExecutable)
	AddTarget("build-debug-executable", "Build the binary debug executable", BuildDebugExecutable)
}

func BuildDebugExecutable(args []string) error {
	buildArgs := []string{
		"build",
		"-o", path.Join(BuildConfig.OutputBinaryPath(), BuildConfig.App.Name+"-debug"),
		`-gcflags="all=-N -l"`,
	}

	builderFlags := strings.Fields(os.Getenv("BUILDER_FLAGS"))

	sourceFile := strings.Fields(path.Join("cmd", BuildConfig.Executable.Cmd, "main.go"))

	return exec.ExecutePipes(exec.Exec(
		"go",
		buildArgs,
		builderFlags,
		sourceFile))
}

func BuildExecutable(args []string) error {
	buildArgs := []string{
		"build",
		"-o", path.Join(BuildConfig.OutputBinaryPath(), BuildConfig.App.Name),
		"-buildmode=pie",
	}

	builderFlags := strings.Fields(os.Getenv("BUILDER_FLAGS"))

	sourceFile := strings.Fields(path.Join("cmd", BuildConfig.Executable.Cmd, "main.go"))

	return exec.ExecutePipes(exec.Exec(
		"go",
		buildArgs,
		builderFlags,
		sourceFile))
}
