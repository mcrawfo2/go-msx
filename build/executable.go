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
		"-o", path.Join(BuildConfig.OutputBinaryPath(), BuildConfig.App.Name + "-debug"),
		`-gcflags="all=-N -l"`,
	}

	if builderFlags := os.Getenv("BUILDER_FLAGS"); builderFlags != "" {
		buildArgs = append(buildArgs, strings.Fields(builderFlags)...)
	}

	buildArgs = append(buildArgs, path.Join("cmd", BuildConfig.Executable.Cmd, "main.go"))

	return exec.Execute("go", buildArgs...)
}

func BuildExecutable(args []string) error {
	buildArgs := []string{
		"build",
		"-o", path.Join(BuildConfig.OutputBinaryPath(), BuildConfig.App.Name),
		"-buildmode=pie",
	}

	if builderFlags := os.Getenv("BUILDER_FLAGS"); builderFlags != "" {
		buildArgs = append(buildArgs, strings.Fields(builderFlags)...)
	}

	buildArgs = append(buildArgs, path.Join("cmd", BuildConfig.Executable.Cmd, "main.go"))

	return exec.Execute("go", buildArgs...)
}
