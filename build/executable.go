package build

import (
	"cto-github.cisco.com/NFV-BU/go-msx/exec"
	"path"
)

func init() {
	AddTarget("build-executable", "Build the binary executable", BuildExecutable)
	AddTarget("build-debug-executable", "Build the binary debug executable", BuildDebugExecutable)
}

func BuildDebugExecutable(args []string) error {
	// TODO
	return exec.MustExecute("go", "build",
		"-o", path.Join(BuildConfig.OutputBinaryPath(), BuildConfig.App.Name),
		path.Join("cmd", BuildConfig.Executable.Cmd, "main.go"))
}

func BuildExecutable(args []string) error {
	return exec.MustExecute("go", "build",
		"-o", path.Join(BuildConfig.OutputBinaryPath(), BuildConfig.App.Name),
		path.Join("cmd", BuildConfig.Executable.Cmd, "main.go"))
}
