// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

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

func AddNode(path, brief string) (*cobra.Command, error) {
	path = strings.TrimSpace(path)
	if len(path) == 0 {
		return nil, errors.New("Missing command path")
	}

	pathParts := strings.Split(path, " ")
	parentPath := pathParts[:len(pathParts)-1]
	commandName := pathParts[len(pathParts)-1]

	existing := FindCommand(pathParts...)
	if existing != nil {
		return nil, errors.New("Command node alread exists")
	}

	parent := FindCommand(parentPath...)
	if parent == nil {
		return nil, errors.New("Could not find parent command")
	}

	cmd := &cobra.Command{
		Use:                commandName,
		Short:              brief,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
	}

	parent.AddCommand(cmd)

	return cmd, nil

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

func GetExitCode() int {
	return exitCode
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
		logger.Errorf("%+v", err)
	}
	Exit()
}

func Fatal(err error) {
	if err != nil {
		logger.Errorf("%+v", err)
		SetExitCode(1)
		Exit()
	}
}
