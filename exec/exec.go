package exec

import (
	"bytes"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"os/exec"
)

var logger = log.NewLogger("msx.exec")

func Execute(name string, dir string, args ...string) (string, string, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	var stdErr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stdErr
	cmd.Dir = dir
	err := cmd.Run()

	return out.String(), stdErr.String(), err
}

func MustExecuteIn(dir string, name string, args ...string) error {
	stdout, stderr, err := Execute(name, dir, args...)
	if err != nil {
		if stdout != "" {
			logger.Warn(stdout)
		}
		if stderr != "" {
			logger.Error(stderr)
		}
	}
	return err
}

func MustExecute(name string, args ...string) error {
	return MustExecuteIn("", name, args...)
}
