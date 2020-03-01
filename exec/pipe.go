package exec

import (
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"gopkg.in/pipe.v2"
	"os"
	"strings"
)

var logger = log.NewLogger("msx.exec")

func RemoveAll(dir string) pipe.Pipe {
	return func(s *pipe.State) error {
		return os.RemoveAll(s.Path(dir))
	}
}

func Exec(name string, args []string, moreArgs ...[]string) pipe.Pipe {
	for _, moreArg := range moreArgs {
		args = append(args, moreArg...)
	}
	logger.Infof("cmd: %s %v", name, args)
	return pipe.Exec(name, args...)
}

func ExecSimple(command ...string) pipe.Pipe {
	return Exec(command[0], command[1:])
}

func Info(template string, args ...interface{}) pipe.Pipe {
	return func(s *pipe.State) error {
		logger.Infof(template, args...)
		return nil
	}
}

func ExecutePipes(pipes ...pipe.Pipe) error {
	for _, p := range pipes {
		if outputBytes, err := pipe.CombinedOutput(WithOutput(p)); err != nil {
			_, _ = os.Stderr.Write(outputBytes)
			return err
		}
	}

	return nil
}

func WithEnv(env map[string]string, p pipe.Pipe) pipe.Pipe {
	var pipes []pipe.Pipe
	for k, v := range env {
		k = strings.ToUpper(k)
		logger.Infof("env: %s=`%s`", k, v)
		pipes = append(pipes, pipe.SetEnvVar(k, v))
	}
	pipes = append(pipes, p)
	return pipe.Line(pipes...)
}

func WithDir(directory string, p pipe.Pipe) pipe.Pipe {
	if directory == "" {
		return p
	}

	return pipe.Line(
		pipe.ChDir(directory),
		p)
}

func WithOutput(p pipe.Pipe) pipe.Pipe {
	return pipe.Line(
		p,
		pipe.Write(os.Stdout))
}
