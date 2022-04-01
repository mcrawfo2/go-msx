// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package scheduled

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/validate"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"time"
)

const (
	configRootScheduled        = "scheduled"
	configPrefixScheduledTasks = configRootScheduled + ".tasks"
)

var errSingleSchedule = errors.New("exactly one of fixed-internal, fixed-delay, cron-expression must be specified")

type TaskConfig struct {
	FixedInterval  *time.Duration `config:"optional"`
	FixedDelay     *time.Duration `config:"optional"`
	InitialDelay   *time.Duration `config:"optional"`
	CronExpression *string        `config:"optional"`
}

func (c *TaskConfig) Validate() error {
	return types.ErrorMap{
		"fixed-interval":  validation.Validate(&c.FixedInterval, validate.Iff(c.FixedInterval != nil, validation.Required)),
		"fixed-delay":     validation.Validate(&c.FixedDelay, validate.Iff(c.FixedDelay != nil, validation.Required)),
		"initial-delay":   validation.Validate(&c.InitialDelay, validate.Iff(c.InitialDelay != nil, validation.Required)),
		"cron-expression": validation.Validate(&c.CronExpression, validate.Iff(c.CronExpression != nil, validation.Required)),
		"schedule": validation.Validate(nil, validate.OneOf(errSingleSchedule,
			c.FixedInterval != nil,
			c.FixedDelay != nil,
			c.CronExpression != nil),
		),
	}
}

func NewTaskConfig(ctx context.Context, taskName string) (*TaskConfig, error) {
	configPrefix := config.NormalizeKey(config.PrefixWithName(configPrefixScheduledTasks, taskName))
	var cfg TaskConfig
	if err := config.MustFromContext(ctx).Populate(&cfg, configPrefix); err != nil {
		return nil, err
	}
	return &cfg, nil
}

type TasksConfig struct {
	Tasks map[string]TaskConfig
}

func newTasksConfig(ctx context.Context) (*TasksConfig, error) {
	var cfg TasksConfig
	if err := config.MustFromContext(ctx).Populate(&cfg, configRootScheduled); err != nil {
		return nil, err
	}
	return &cfg, nil
}
