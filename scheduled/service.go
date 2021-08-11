package scheduled

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
	"github.com/thejerf/abtime"
	"sync/atomic"
)

var logger = log.NewLogger("msx.scheduled")

var errTaskExists = errors.New("Task already exists")
var errTaskNotConfigured = errors.New("Task not configured")

var taskCounter uint64 = 0

type SchedulerServiceApi interface {
	Schedule(taskName string, action types.ActionFunc) error
	Run(ctx context.Context) error
}

type schedulerService struct {
	cfg     *TasksConfig
	ctx     context.Context
	clock   abtime.AbstractTime
	tasks   map[string]scheduledTask
	started bool
}

func (s *schedulerService) Schedule(taskName string, action types.ActionFunc) error {
	taskKey := config.NormalizeKey(taskName)
	if _, ok := s.tasks[taskKey]; ok {
		return errors.Wrap(errTaskExists, taskName)
	}

	cfg, ok := s.cfg.Tasks[taskKey]
	if !ok {
		return errors.Wrap(errTaskNotConfigured, taskName)
	}

	s.tasks[taskKey] = scheduledTask{
		name:   taskName,
		cfg:    cfg,
		action: action,
		index:  taskCounter,
	}

	atomic.AddUint64(&taskCounter, 1)

	s.startTask(taskKey)

	return nil
}

func (s *schedulerService) startTask(key string) {
	if !s.started {
		// all tasks will be scheduled at start time
		return
	}

	task := s.tasks[key]

	task.worker = types.NewWorker(s.ctx)
	task.worker.Schedule(s.timer(key, task))

	s.tasks[key] = task

	return
}

func (s *schedulerService) timer(taskKey string, t scheduledTask) types.ActionFunc {
	operationName := config.PrefixWithName("scheduled.tasks", taskKey)

	return func(ctx context.Context) error {
		ticker := s.clock.NewTicker(*t.cfg.FixedInterval, int(t.index))
		defer func() { ticker.Stop() }()

		for {
			select {
			case <-ticker.Channel():
				trace.BackgroundOperation(ctx, operationName, func(ctx context.Context) error {
					logger.WithContext(ctx).Debugf("Scheduled task %q tick", t.name)
					err := t.action(ctx)
					if err != nil {
						logger.WithContext(ctx).WithError(err).Errorf("Scheduled Task %q failed.", t.name)
					} else {
						logger.WithContext(ctx).Debugf("Scheduled Task %q completed.", t.name)
					}
					return nil
				})

			case <-ctx.Done():
				logger.WithContext(ctx).WithError(ctx.Err()).Errorf("Scheduled Task %q Timer stopped.", t.name)
				return nil
			}
		}
	}
}

func (s *schedulerService) Run(ctx context.Context) error {
	s.started = true
	s.ctx = ctx
	for taskKey := range s.tasks {
		s.startTask(taskKey)
	}
	return nil
}

func NewSchedulerService(ctx context.Context) (SchedulerServiceApi, error) {
	service := SchedulerServiceFromContext(ctx)
	if service == nil {
		cfg, err := newTasksConfig(ctx)
		if err != nil {
			return nil, err
		}

		service = &schedulerService{
			cfg:   cfg,
			clock: types.NewClock(ctx),
			tasks: make(map[string]scheduledTask),
		}
	}

	return service, nil
}
