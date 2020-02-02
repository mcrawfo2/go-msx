package cli

import (
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

type CommandFunc func(args []string) error

var rootCmd = &cobra.Command{
	Use:                "",
	SilenceUsage:       true,
	SilenceErrors:      true,
	FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
}

var logger = log.NewLogger("msx.cli")

var exitCode = 0

func RootCmd() *cobra.Command {
	return rootCmd
}

func AddCommand(path, brief string, cmdFunc CommandFunc) (*cobra.Command, error) {
	path = strings.TrimSpace(path)
	if len(path) == 0 {
		return nil, errors.New("Missing command path")
	}

	pathParts := strings.Split(path, " ")
	parentPath := pathParts[:len(pathParts)-1]
	commandName := pathParts[len(pathParts)-1]

	parent := FindCommand(parentPath...)
	if parent == nil {
		return nil, errors.New("Could not find parent command")
	}

	cmd := &cobra.Command{
		Use:                commandName,
		Short:              brief,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmdFunc(args)
		},
	}

	parent.AddCommand(cmd)

	return cmd, nil
}

func FindCommand(path ...string) *cobra.Command {
	var next *cobra.Command
	here := rootCmd
	for _, pathPart := range path {
		hereCommands := here.Commands()

		next = nil
		for _, hereCommand := range hereCommands {
			if hereCommand.Use == pathPart || strings.HasPrefix(hereCommand.Use, pathPart+" ") {
				next = hereCommand
				break
			}
		}

		here = next
		if here == nil {
			break
		}
	}

	return here
}

func SetExitCode(exit int) {
	exitCode = exit
}

func Exit() {
	os.Exit(exitCode)
}

func Run(appName string) {
	RootCmd().Use = appName
	if err := RootCmd().Execute(); err != nil {
		if exitCode == 0 {
			Fatal(err)
		}
		logger.Error(err)
	}
	Exit()
}

func Fatal(err error) {
	if err != nil {
		logger.Error(err)
		SetExitCode(1)
		Exit()
	}
}
