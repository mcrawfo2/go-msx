package auditlog

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"errors"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestAction(t *testing.T) {
	type args struct {
		resourceName string
		action       string
	}
	tests := []struct {
		name string
		args args
		want []log.EntryPredicate
	}{
		{
			name: "DefaultUser",
			args: args{
				resourceName: "SERVICE_INSTANCE",
				action:       "DEPLOY",
			},
			want: []log.EntryPredicate{
				log.HasFieldValue("resource", "SERVICE_INSTANCE"),
				log.HasFieldValue("action", "DEPLOY"),
				log.HasFieldValue("state", StateAction),
				log.HasFieldValue("audit", "true"),
				log.HasFieldValue("user", "anonymous"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recording.Reset()

			entry := Action(logger, context.Background(), tt.args.resourceName, tt.args.action)
			for _, predicate := range tt.want {
				if !predicate.Matches(*entry) {
					t.Errorf("Entry did not match predicate %q", predicate.Description)
				}
			}
		})
	}
}

func TestAudit(t *testing.T) {
	type args struct {
		resourceName string
		action       string
		fn           func() error
	}
	tests := []struct {
		name   string
		args   args
		checks []log.Check
	}{
		{
			name: "NoError",
			args: args{
				resourceName: "SERVICE_INSTANCE",
				action:       "DEPLOY",
				fn: func() error {
					return nil
				},
			},
			checks: []log.Check{
				{
					Filters: []log.EntryPredicate{
						log.HasLevel(logrus.InfoLevel),
						log.HasFieldValue("state", StateInit),
					},
					Validators: []log.EntryPredicate{
						log.HasFieldValue("resource", "SERVICE_INSTANCE"),
						log.HasFieldValue("action", "DEPLOY"),
						log.HasFieldValue("audit", "true"),
						log.HasFieldValue("user", "anonymous"),
					},
				},
				{
					Filters: []log.EntryPredicate{
						log.HasLevel(logrus.InfoLevel),
						log.HasFieldValue("state", StateSuccess),
					},
					Validators: []log.EntryPredicate{
						log.HasFieldValue("resource", "SERVICE_INSTANCE"),
						log.HasFieldValue("action", "DEPLOY"),
						log.HasFieldValue("audit", "true"),
						log.HasFieldValue("user", "anonymous"),
					},
				},
			},
		},
		{
			name: "Error",
			args: args{
				resourceName: "SERVICE_INSTANCE",
				action:       "DEPLOY",
				fn: func() error {
					return errors.New("some error")
				},
			},
			checks: []log.Check{
				{
					Filters: []log.EntryPredicate{
						log.HasLevel(logrus.InfoLevel),
						log.HasFieldValue("state", StateInit),
					},
					Validators: []log.EntryPredicate{
						log.HasFieldValue("resource", "SERVICE_INSTANCE"),
						log.HasFieldValue("action", "DEPLOY"),
						log.HasFieldValue("audit", "true"),
						log.HasFieldValue("user", "anonymous"),
					},
				},
				{
					Filters: []log.EntryPredicate{
						log.HasLevel(logrus.InfoLevel),
						log.HasFieldValue("state", StateFail),
					},
					Validators: []log.EntryPredicate{
						log.HasFieldValue("resource", "SERVICE_INSTANCE"),
						log.HasFieldValue("action", "DEPLOY"),
						log.HasFieldValue("audit", "true"),
						log.HasFieldValue("user", "anonymous"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recording.Reset()

			Audit(logger, context.Background(), tt.args.resourceName, tt.args.action, tt.args.fn)
			for _, check := range tt.checks {
				errs := check.Check(recording)
				for _, err := range errs {
					t.Error(err.Error())
				}
			}
		})
	}
}

func TestEntry(t *testing.T) {
	ctx := context.Background()
	ctx2 := security.ContextWithUserContext(ctx, &security.UserContext{UserName: "mach"})
	ctx3 := ContextWithRequestDetails(ctx2, &RequestDetails{
		Source:   "10.10.10.10",
		Protocol: "https",
		Host:     "192.168.2.1",
		Port:     "8080",
	})

	type args struct {
		ctx          context.Context
		resourceName string
		action       string
	}
	tests := []struct {
		name string
		args args
		want []log.EntryPredicate
	}{
		{
			name: "DefaultUser",
			args: args{
				ctx:          ctx,
				resourceName: "SERVICE_INSTANCE",
				action:       "DEPLOY",
			},
			want: []log.EntryPredicate{
				log.HasFieldValue("resource", "SERVICE_INSTANCE"),
				log.HasFieldValue("action", "DEPLOY"),
				log.HasFieldValue("state", StateAction),
				log.HasFieldValue("audit", "true"),
				log.HasFieldValue("user", "anonymous"),
				log.NoFieldValue("source"),
				log.NoFieldValue("protocol"),
				log.NoFieldValue("host"),
				log.NoFieldValue("port"),
			},
		},
		{
			name: "WithUser",
			args: args{
				ctx:          ctx2,
				resourceName: "SERVICE_INSTANCE",
				action:       "DEPLOY",
			},
			want: []log.EntryPredicate{
				log.HasFieldValue("resource", "SERVICE_INSTANCE"),
				log.HasFieldValue("action", "DEPLOY"),
				log.HasFieldValue("state", StateAction),
				log.HasFieldValue("audit", "true"),
				log.HasFieldValue("user", "mach"),
				log.NoFieldValue("source"),
				log.NoFieldValue("protocol"),
				log.NoFieldValue("host"),
				log.NoFieldValue("port"),
			},
		},
		{
			name: "WithAudit",
			args: args{
				ctx:          ctx3,
				resourceName: "SERVICE_INSTANCE",
				action:       "DEPLOY",
			},
			want: []log.EntryPredicate{
				log.HasFieldValue("resource", "SERVICE_INSTANCE"),
				log.HasFieldValue("action", "DEPLOY"),
				log.HasFieldValue("state", StateAction),
				log.HasFieldValue("audit", "true"),
				log.HasFieldValue("user", "mach"),
				log.HasFieldValue("source", "10.10.10.10"),
				log.HasFieldValue("protocol", "https"),
				log.HasFieldValue("host", "192.168.2.1"),
				log.HasFieldValue("port", "8080"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recording.Reset()

			entry := Action(logger, tt.args.ctx, tt.args.resourceName, tt.args.action)
			for _, predicate := range tt.want {
				if !predicate.Matches(*entry) {
					t.Errorf("Entry did not match predicate %q", predicate.Description)
				}
			}
		})
	}
}

func TestError(t *testing.T) {
	e := errors.New("some error")

	type args struct {
		resourceName string
		action       string
	}
	tests := []struct {
		name string
		args args
		want []log.EntryPredicate
	}{
		{
			name: "DefaultUser",
			args: args{
				resourceName: "SERVICE_INSTANCE",
				action:       "DEPLOY",
			},
			want: []log.EntryPredicate{
				log.HasFieldValue("resource", "SERVICE_INSTANCE"),
				log.HasFieldValue("action", "DEPLOY"),
				log.HasFieldValue("state", StateFail),
				log.HasFieldValue("audit", "true"),
				log.HasFieldValue("user", "anonymous"),
				log.HasFieldValue("error", e),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recording.Reset()

			entry := Error(logger, context.Background(), tt.args.resourceName, tt.args.action, e)
			for _, predicate := range tt.want {
				if !predicate.Matches(*entry) {
					t.Errorf("Entry did not match predicate %q", predicate.Description)
				}
			}
		})
	}
}

func TestFailure(t *testing.T) {
	type args struct {
		resourceName string
		action       string
	}
	tests := []struct {
		name string
		args args
		want []log.EntryPredicate
	}{
		{
			name: "DefaultUser",
			args: args{
				resourceName: "SERVICE_INSTANCE",
				action:       "DEPLOY",
			},
			want: []log.EntryPredicate{
				log.HasFieldValue("resource", "SERVICE_INSTANCE"),
				log.HasFieldValue("action", "DEPLOY"),
				log.HasFieldValue("state", StateFail),
				log.HasFieldValue("audit", "true"),
				log.HasFieldValue("user", "anonymous"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recording.Reset()

			entry := Failure(logger, context.Background(), tt.args.resourceName, tt.args.action)
			for _, predicate := range tt.want {
				if !predicate.Matches(*entry) {
					t.Errorf("Entry did not match predicate %q", predicate.Description)
				}
			}
		})
	}
}

func TestInit(t *testing.T) {
	type args struct {
		resourceName string
		action       string
	}
	tests := []struct {
		name string
		args args
		want []log.EntryPredicate
	}{
		{
			name: "DefaultUser",
			args: args{
				resourceName: "SERVICE_INSTANCE",
				action:       "DEPLOY",
			},
			want: []log.EntryPredicate{
				log.HasFieldValue("resource", "SERVICE_INSTANCE"),
				log.HasFieldValue("action", "DEPLOY"),
				log.HasFieldValue("state", StateInit),
				log.HasFieldValue("audit", "true"),
				log.HasFieldValue("user", "anonymous"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recording.Reset()

			entry := Init(logger, context.Background(), tt.args.resourceName, tt.args.action)
			for _, predicate := range tt.want {
				if !predicate.Matches(*entry) {
					t.Errorf("Entry did not match predicate %q", predicate.Description)
				}
			}
		})
	}
}

func TestResult(t *testing.T) {
	type args struct {
		resourceName string
		action       string
		err          error
	}
	tests := []struct {
		name string
		args args
		want []log.EntryPredicate
	}{
		{
			name: "NoError",
			args: args{
				resourceName: "SERVICE_INSTANCE",
				action:       "DEPLOY",
				err:          nil,
			},
			want: []log.EntryPredicate{
				log.HasFieldValue("resource", "SERVICE_INSTANCE"),
				log.HasFieldValue("action", "DEPLOY"),
				log.HasFieldValue("state", StateSuccess),
				log.HasFieldValue("audit", "true"),
				log.HasFieldValue("user", "anonymous"),
			},
		},
		{
			name: "Error",
			args: args{
				resourceName: "SERVICE_INSTANCE",
				action:       "DEPLOY",
				err:          errors.New("another error"),
			},
			want: []log.EntryPredicate{
				log.HasFieldValue("resource", "SERVICE_INSTANCE"),
				log.HasFieldValue("action", "DEPLOY"),
				log.HasFieldValue("state", StateFail),
				log.HasFieldValue("audit", "true"),
				log.HasFieldValue("user", "anonymous"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recording.Reset()

			entry := Result(logger, context.Background(), tt.args.resourceName, tt.args.action, tt.args.err)
			for _, predicate := range tt.want {
				if !predicate.Matches(*entry) {
					t.Errorf("Entry did not match predicate %q", predicate.Description)
				}
			}
		})
	}
}

func TestResultOf(t *testing.T) {

	type args struct {
		resourceName string
		action       string
		fn           func() error
	}
	tests := []struct {
		name string
		args args
		want []log.EntryPredicate
	}{
		{
			name: "NoError",
			args: args{
				resourceName: "SERVICE_INSTANCE",
				action:       "DEPLOY",
				fn: func() error {
					return nil
				},
			},
			want: []log.EntryPredicate{
				log.HasFieldValue("resource", "SERVICE_INSTANCE"),
				log.HasFieldValue("action", "DEPLOY"),
				log.HasFieldValue("state", StateSuccess),
				log.HasFieldValue("audit", "true"),
				log.HasFieldValue("user", "anonymous"),
			},
		},
		{
			name: "Error",
			args: args{
				resourceName: "SERVICE_INSTANCE",
				action:       "DEPLOY",
				fn: func() error {
					return errors.New("another error")
				},
			},
			want: []log.EntryPredicate{
				log.HasFieldValue("resource", "SERVICE_INSTANCE"),
				log.HasFieldValue("action", "DEPLOY"),
				log.HasFieldValue("state", StateFail),
				log.HasFieldValue("audit", "true"),
				log.HasFieldValue("user", "anonymous"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recording.Reset()

			entry := ResultOf(logger, context.Background(), tt.args.resourceName, tt.args.action, tt.args.fn)
			for _, predicate := range tt.want {
				if !predicate.Matches(*entry) {
					t.Errorf("Entry did not match predicate %q", predicate.Description)
				}
			}
		})
	}}

func TestSuccess(t *testing.T) {
	type args struct {
		resourceName string
		action       string
	}
	tests := []struct {
		name string
		args args
		want []log.EntryPredicate
	}{
		{
			name: "DefaultUser",
			args: args{
				resourceName: "SERVICE_INSTANCE",
				action:       "DEPLOY",
			},
			want: []log.EntryPredicate{
				log.HasFieldValue("resource", "SERVICE_INSTANCE"),
				log.HasFieldValue("action", "DEPLOY"),
				log.HasFieldValue("state", StateSuccess),
				log.HasFieldValue("audit", "true"),
				log.HasFieldValue("user", "anonymous"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recording.Reset()

			entry := Success(logger, context.Background(), tt.args.resourceName, tt.args.action)
			for _, predicate := range tt.want {
				if !predicate.Matches(*entry) {
					t.Errorf("Entry did not match predicate %q", predicate.Description)
				}
			}
		})
	}
}
