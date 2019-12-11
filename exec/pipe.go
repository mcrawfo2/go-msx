package exec

import (
	"gopkg.in/pipe.v2"
	"os"
)

func RemoveAll(dir string) pipe.Pipe {
	return func(s *pipe.State) error {
		return os.RemoveAll(s.Path(dir))
	}
}

func Exec(name string, args []string, moreArgs ...[]string) pipe.Pipe {
	for _, moreArg := range moreArgs {
		args = append(args, moreArg...)
	}
	return pipe.Exec(name, args...)
}

func Info(template string, args ...interface{}) pipe.Pipe {
	return func(s *pipe.State) error {
		logger.Infof(template, args...)
		return nil
	}
}

func ExecutePipes(pipes ...pipe.Pipe) error {
	return ExecutePipesIn("", pipes...)
}

func ExecutePipesIn(directory string, pipes ...pipe.Pipe) error {
	for _, p := range pipes {
		linePipes := []pipe.Pipe{
			p,
			pipe.Write(os.Stdout),
		}

		if directory != "" {
			linePipes = append([]pipe.Pipe{pipe.ChDir(directory)}, linePipes...)
		}

		if outputBytes, err := pipe.CombinedOutput(pipe.Line(linePipes...)); err != nil {
			_, _ = os.Stderr.Write(outputBytes)
			return err
		}
	}
	return nil
}
