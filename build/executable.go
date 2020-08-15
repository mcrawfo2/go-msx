package build

import (
	"cto-github.cisco.com/NFV-BU/go-msx/exec"
	"github.com/pkg/errors"
	"os"
	"path"
	"strings"
)

func init() {
	AddTarget("build-executable", "Build the binary executable", BuildExecutable)
	AddTarget("build-debug-executable", "Build the binary debug executable", BuildDebugExecutable)
}

func BuildDebugExecutable(args []string) error {
	if BuildConfig.App.Name == "" {
		return errors.New("App name not specified.  Please provide app.name in app config file")
	}
	if BuildConfig.Executable.Cmd == "" {
		return errors.New("Entrypoint not specified.  Please provide executable.cmd in build.yml")
	}

	buildArgs := []string{
		"build",
		"-o", path.Join(BuildConfig.OutputBinaryPath(), BuildConfig.App.Name+"-debug"),
		`-gcflags="all=-N -l"`,
	}

	builderFlags := strings.Fields(os.Getenv("BUILDER_FLAGS"))

	sourceFile := strings.Fields(path.Join("cmd", BuildConfig.Executable.Cmd, "main.go"))

	return exec.ExecutePipes(
		exec.Exec("go",
			buildArgs,
			builderFlags,
			sourceFile))
}

func BuildExecutable(args []string) error {
	if BuildConfig.App.Name == "" {
		return errors.New("App name not specified.  Please provide app.name in app config file")
	}
	if BuildConfig.Executable.Cmd == "" {
		return errors.New("Entrypoint not specified.  Please provide executable.cmd in build.yml")
	}

	buildArgs := []string{
		"build",
		"-o", path.Join(BuildConfig.OutputBinaryPath(), BuildConfig.App.Name),
	}

	builderFlags := strings.Fields(os.Getenv("BUILDER_FLAGS"))

	sourceFile := strings.Fields(path.Join("cmd", BuildConfig.Executable.Cmd, "main.go"))

	return exec.ExecutePipes(
		exec.WithEnv(BuildConfig.Go.Environment(),
			exec.Exec(
				"go",
				buildArgs,
				builderFlags,
				sourceFile)))
}
