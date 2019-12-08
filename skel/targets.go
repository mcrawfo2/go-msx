package skel

import (
	"cto-github.cisco.com/NFV-BU/go-msx/cli"
	"github.com/pkg/errors"
)

var allTargets = make(map[string]cli.CommandFunc)

func AddTarget(name, description string, fn cli.CommandFunc) {
	wrapper := func(args []string) error {
		logger.Infof("Executing target '%s': %s", name, description)
		err := fn(args)
		if err != nil {
			logger.Infof("Target failed: '%s': %s", name, err.Error())
		} else {
			logger.Infof("Target succeeded: '%s'", name)
		}
		return err
	}

	allTargets[name] = fn

	_, err := cli.AddCommand(name, description, wrapper)
	if err != nil {
		panic(err.Error())
	}
}

func ExecTargets(targets ...string) error {
	for _, target := range targets {
		fn, ok := allTargets[target]
		if !ok {
			return errors.Errorf("Target not found: %s", target)
		}
		if err := fn(nil); err != nil {
			return err
		}
	}
	return nil
}
