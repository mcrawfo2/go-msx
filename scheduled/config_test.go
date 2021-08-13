package scheduled

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"cto-github.cisco.com/NFV-BU/go-msx/validate"
	"reflect"
	"testing"
	"time"
)

func TestTaskConfig_Validate(t *testing.T) {
	duration := time.Duration(15) * time.Second
	cron := "0 15 10 15 * ?"
	type fields struct {
		FixedInterval  *time.Duration
		FixedDelay     *time.Duration
		InitialDelay   *time.Duration
		CronExpression *string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "NoSchedule",
			fields: fields{
				InitialDelay: &duration,
			},
			wantErr: true,
		},
		{
			name: "FixedInterval",
			fields: fields{
				FixedInterval: &duration,
			},
		},
		{
			name: "FixedDelay",
			fields: fields{
				FixedDelay: &duration,
			},
		},
		{
			name: "CronExpression",
			fields: fields{
				CronExpression: &cron,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &TaskConfig{
				FixedInterval:  tt.fields.FixedInterval,
				FixedDelay:     tt.fields.FixedDelay,
				InitialDelay:   tt.fields.InitialDelay,
				CronExpression: tt.fields.CronExpression,
			}
			if err := validate.Validate(c); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewTasksConfig(t *testing.T) {
	duration := time.Duration(15) * time.Second
	cron := "0 15 10 15 * ?"

	tests := []struct {
		name    string
		ctx     context.Context
		want    *TasksConfig
		wantErr bool
	}{
		{
			name: "Empty",
			ctx: configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{
				"scheduled.tasks.my-task.name": "my-task",
			}),
			want: &TasksConfig{
				Tasks: map[string]TaskConfig{
					"mytask": {},
				},
			},
			wantErr: false,
		},
		{
			name: "FixedInterval",
			ctx: configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{
				"scheduled.tasks.my-task.fixed-interval": "15s",
			}),
			want: &TasksConfig{
				Tasks: map[string]TaskConfig{
					"mytask": {
						FixedInterval: &duration,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "InitialDelay",
			ctx: configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{
				"scheduled.tasks.my-task.initial-delay": "15s",
			}),
			want: &TasksConfig{
				Tasks: map[string]TaskConfig{
					"mytask": {
						InitialDelay: &duration,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "FixedDelay",
			ctx: configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{
				"scheduled.tasks.my-task.fixed-delay": "15s",
			}),
			want: &TasksConfig{
				Tasks: map[string]TaskConfig{
					"mytask": {
						FixedDelay: &duration,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "CronExpression",
			ctx: configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{
				"scheduled.tasks.my-task.cron-expression": cron,
			}),
			want: &TasksConfig{
				Tasks: map[string]TaskConfig{
					"mytask": {
						CronExpression: &cron,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newTasksConfig(tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("newTasksConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newTasksConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
