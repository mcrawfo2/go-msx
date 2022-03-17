//go:generate mockery --inpackage --name=UpperCamelSingularTimerApi --structname=MockUpperCamelSingularTimer --filename mock_timer_lowersingular.go
package lowerplural

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/app"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/leader"
	"cto-github.cisco.com/NFV-BU/go-msx/scheduled"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

const (
	contextKeyUpperCamelSingularTimer = contextKey("lowerCamelSingularTimer")
	configRootUpperCamelSingularTimer = "${app.name}.lowersingular.timer"
	taskNameUpperCamelSingular        = "lowersingular"
)

type UpperCamelSingularTimerApi interface {
	Run(ctx context.Context) error
}

type lowerCamelSingularTimer struct {
	cfg *lowerCamelSingularTimerConfig
}

type lowerCamelSingularTimerConfig struct {
}

func newUpperCamelSingularTimerConfig(ctx context.Context) (*lowerCamelSingularTimerConfig, error) {
	var cfg lowerCamelSingularTimerConfig
	if err := config.FromContext(ctx).Populate(&cfg, configRootUpperCamelSingularTimer); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (t *lowerCamelSingularTimer) Run(ctx context.Context) error {
	// TODO : Implement
	return nil
}

func newUpperCamelSingularTimer(ctx context.Context) (UpperCamelSingularTimerApi, error) {
	timer := UpperCamelSingularTimerFromContext(ctx)
	if timer == nil {
		cfg, err := newUpperCamelSingularTimerConfig(ctx)
		if err != nil {
			return nil, err
		}

		timer = &lowerCamelSingularTimer{
			cfg: cfg,
		}
	}
	return timer, nil
}

func UpperCamelSingularTimerFromContext(ctx context.Context) UpperCamelSingularTimerApi {
	value, _ := ctx.Value(contextKeyUpperCamelSingularTimer).(UpperCamelSingularTimerApi)
	return value
}

func ContextWithUpperCamelSingularTimer(ctx context.Context, timer UpperCamelSingularTimerApi) context.Context {
	return context.WithValue(ctx, contextKeyUpperCamelSingularTimer, timer)
}

func init() {
	var timer UpperCamelSingularTimerApi

	app.OnRootEvent(app.EventStart, app.PhaseDuring, func(ctx context.Context) (err error) {
		timer, err = newUpperCamelSingularTimer(ctx)
		if err != nil {
			return err
		}
		return
	})

	app.OnRootEvent(app.EventStart, app.PhaseAfter, func(ctx context.Context) error {
		operation := types.
			NewOperation(timer.Run).
			WithDecorator(leader.MasterLeaderDecorator)

		return scheduled.ScheduleTask(ctx, taskNameUpperCamelSingular, operation.Run)
	})
}
