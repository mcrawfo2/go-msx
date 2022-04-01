// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"github.com/elastic/go-seccomp-bpf"
	"github.com/pkg/errors"
	"runtime"
)

const configRootSecComp = "seccomp"

func init() {
	OnEvent(EventStart, PhaseBefore, applySecCompProfile)
}

func applySecCompProfile(ctx context.Context) error {
	if !seccomp.Supported() {
		if runtime.GOOS == "linux" {
			logger.WithContext(ctx).Warn("SecComp is not supported on this installation of linux.")
		} else {
			logger.WithContext(ctx).Warnf("SecComp is not supported on %s.", runtime.GOOS)
		}
		return nil
	}

	cfg := config.FromContext(ctx)
	if cfg == nil {
		return config.ErrNotLoaded
	}

	enabled, err := cfg.BoolOr("seccomp.enabled", false)
	if err != nil {
		return err
	}

	if !enabled {
		logger.WithContext(ctx).Warn("SecComp is disabled.")
		return nil
	}

	noNewPrivs, err := cfg.BoolOr("seccomp.no-new-privs", true)
	if err != nil {
		return err
	}

	var policy seccomp.Policy
	if err = cfg.Populate(&policy, configRootSecComp); err != nil {
		logger.WithContext(ctx).WithError(err).Warn("SecComp policy not defined or invalid.")
		return nil
	}

	// Create a filter based on config.
	filter := seccomp.Filter{
		NoNewPrivs: noNewPrivs,
		Flag:       seccomp.FilterFlagTSync,
		Policy:     policy,
	}

	err = seccomp.LoadFilter(filter)
	if err != nil {
		err = errors.Wrap(err, "Failed to activate seccomp")
	}

	return err
}
