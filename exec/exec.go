package exec

import (
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"os"
	"os/exec"
)

var logger = log.NewLogger("msx.exec")

func execute(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func ExecuteIn(dir string, name string, args ...string) error {
	return execute(dir, name, args...)
}

func Execute(name string, args ...string) error {
	return ExecuteIn("", name, args...)
}
