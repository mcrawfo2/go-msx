// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package scheduled

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
)

var errSchedulerServiceNotAvailable = errors.New("Scheduler Service not available")

type scheduledTask struct {
	index  uint64
	name   string
	cfg    TaskConfig
	action types.ActionFunc
	worker *types.Worker
}

func ScheduleTask(ctx context.Context, taskName string, action types.ActionFunc) error {
	service := SchedulerServiceFromContext(ctx)
	if service == nil {
		return errors.Wrap(errSchedulerServiceNotAvailable, "Failed to schedule task")
	}

	return service.Schedule(taskName, action)
}
