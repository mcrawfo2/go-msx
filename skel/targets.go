// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package skel

import (
	"cto-github.cisco.com/NFV-BU/go-msx/cli"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var allTargets = make(map[string]cli.CommandFunc)

func AddTarget(name, description string, fn cli.CommandFunc) *cobra.Command {
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

	cmd, err := cli.AddCommand(name, description, wrapper)
	if err != nil {
		panic(err.Error())
	}

	return cmd
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
