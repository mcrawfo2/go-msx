package scheduled

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/thejerf/abtime"
	"sync/atomic"
)

var logger = log.NewLogger("msx.scheduled")

var errTaskExists = errors.New("Task already exists")
var errTaskNotConfigured = errors.New("Task not configured")

var taskCounter uint64 = 0

var cronParserOpts = cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow

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
	task.worker.Schedule(s.looper(key, task))

	s.tasks[key] = task

	return
}

func (s *schedulerService) timer(t scheduledTask, first bool) (abtime.Timer, error) {
	switch {
	case t.cfg.InitialDelay != nil && first:
		return s.clock.NewTimer(*t.cfg.InitialDelay, int(t.index)), nil

	case t.cfg.FixedInterval != nil:
		return s.clock.NewTimer(*t.cfg.FixedInterval, int(t.index)), nil

	case t.cfg.FixedDelay != nil:
		return s.clock.NewTimer(*t.cfg.FixedDelay, int(t.index)), nil

	case t.cfg.CronExpression != nil:
		parser := cron.NewParser(cronParserOpts)
		schedule, err := parser.Parse(*t.cfg.CronExpression)
		if err != nil {
			return nil, err
		}

		nextTime := schedule.Next(s.clock.Now())
		delay := nextTime.Sub(s.clock.Now())
		if delay < 1 {
			delay = 1
		}
		return s.clock.NewTimer(delay, int(t.index)), nil

	default:
		return nil, errSingleSchedule
	}
}

func (s *schedulerService) looper(taskKey string, t scheduledTask) types.ActionFunc {
	operationName := config.PrefixWithName("scheduled.tasks", taskKey)

	return func(ctx context.Context) error {
		timer, err := s.timer(t, true)
		if err != nil {
			return err
		}
		defer func() { timer.Stop() }()

		for {
			select {
			case <-timer.Channel():
				action := types.NewOperation(t.action).
					WithDecorator(s.taskDecorator(t)).
					Run

				if t.cfg.FixedDelay != nil {
					trace.ForegroundOperation(ctx, operationName, action)
				} else {
					trace.BackgroundOperation(ctx, operationName, action)
				}

				timer, err = s.timer(t, false)
				if err != nil {
					return err
				}

			case <-ctx.Done():
				logger.WithContext(ctx).WithError(ctx.Err()).Errorf("Scheduled Task %q Timer stopped.", t.name)
				return nil
			}
		}
	}
}

func (s *schedulerService) taskDecorator(t scheduledTask) types.ActionFuncDecorator {
	return func(action types.ActionFunc) types.ActionFunc {
		return func(ctx context.Context) error {
			logger.WithContext(ctx).Debugf("Scheduled task %q tick", t.name)
			err := t.action(ctx)
			if err != nil {
				logger.WithContext(ctx).WithError(err).Errorf("Scheduled Task %q failed.", t.name)
			} else {
				logger.WithContext(ctx).Debugf("Scheduled Task %q completed.", t.name)
			}
			return nil
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
