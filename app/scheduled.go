// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package app

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/scheduled"
)

func init() {
	OnEvent(EventStart, PhaseBefore, func(ctx context.Context) error {
		schedulerService, err := scheduled.NewSchedulerService(ctx)
		if err != nil {
			return err
		}

		RegisterContextInjector(func(ctx context.Context) context.Context {
			return scheduled.ContextWithSchedulerService(ctx, schedulerService)
		})

		OnEvent(EventStart, PhaseAfter, schedulerService.Run)

		return nil
	})
}
