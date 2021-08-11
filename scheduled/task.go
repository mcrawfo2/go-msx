package scheduled

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
)

var errSchedulerServiceNotAvailable = errors.New("Scheduler Service not available")

type scheduledTask struct {
	index         uint64
	name          string
	cfg           TaskConfig
	action        types.ActionFunc
	worker        *types.Worker
}

func ScheduleTask(ctx context.Context, taskName string, action types.ActionFunc) error {
	service  := SchedulerServiceFromContext(ctx)
	if service == nil {
		return errors.Wrap(errSchedulerServiceNotAvailable, "Failed to schedule task")
	}

	return service.Schedule(taskName, action)
}
