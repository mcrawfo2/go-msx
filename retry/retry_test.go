// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package retry

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
	"reflect"
	"testing"
	"time"
)

func TestNewRetry(t *testing.T) {
	type args struct {
		ctx context.Context
		cfg RetryConfig
	}
	tests := []struct {
		name string
		args args
		want Retry
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				cfg: RetryConfig{
					Attempts: 2,
					Delay:    250,
					BackOff:  1.0,
					Linear:   true,
				},
			},
			want: Retry{
				Attempts: 2,
				Delay:    250 * time.Millisecond,
				BackOff:  1.0,
				Linear:   true,
				Context:  context.Background(),
				clock:    types.NewClock(context.Background()),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRetry(tt.args.ctx, tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRetry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRetryFromConfig(t *testing.T) {
	cfg := configtest.NewInMemoryConfig(map[string]string{
		"spring.retry.attempts": "2",
		"spring.retry.delay":    "500",
		"spring.retry.backoff":  "2.0",
		"spring.retry.linear":   "true",
	})

	tests := []struct {
		name    string
		want    *Retry
		wantErr bool
	}{
		{
			name: "Success",
			want: &Retry{
				Attempts: 2,
				Delay:    500 * time.Millisecond,
				BackOff:  2.0,
				Linear:   true,
				Context:  context.Background(),
				clock:    types.NewClock(context.Background()),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRetryFromConfig(context.Background(), cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRetryFromConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRetryFromConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRetryFromContext(t *testing.T) {
	ctx := configtest.ContextWithNewInMemoryConfig(
		context.Background(),
		map[string]string{
			"spring.retry.attempts": "2",
			"spring.retry.delay":    "500",
			"spring.retry.backoff":  "2.0",
			"spring.retry.linear":   "true",
		})

	tests := []struct {
		name    string
		want    *Retry
		wantErr bool
	}{
		{
			name: "Success",
			want: &Retry{
				Attempts: 2,
				Delay:    500 * time.Millisecond,
				BackOff:  2.0,
				Linear:   true,
				Context:  ctx,
				clock:    types.NewClock(context.Background()),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRetryFromContext(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRetryFromConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRetryFromConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRetry_Retry(t *testing.T) {
	clock := types.NewMockClock()
	ctx := types.ContextWithClock(context.Background(), clock)

	type args struct {
		config    RetryConfig
		retryable Retryable
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Once",
			args: args{
				config: RetryConfig{
					Attempts: 1,
					Delay:    500,
					BackOff:  1.0,
					Linear:   true,
				},
				retryable: func() error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "OnceError",
			args: args{
				config: RetryConfig{
					Attempts: 1,
					Delay:    500,
					BackOff:  1.0,
					Linear:   true,
				},
				retryable: func() error {
					return errors.New("a transient error")
				},
			},
			wantErr: true,
		},
		{
			name: "Thrice",
			args: args{
				config: RetryConfig{
					Attempts: 3,
					Delay:    500,
					BackOff:  1.0,
					Linear:   true,
				},
				retryable: func() error {
					return errors.New("a transient error")
				},
			},
			wantErr: true,
		},
		{
			name: "Permanent",
			args: args{
				config: RetryConfig{
					Attempts: 3,
					Delay:    500,
					BackOff:  1.0,
					Linear:   true,
				},
				retryable: func() error {
					return &PermanentError{
						Cause: errors.New("a transient error"),
					}
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRetry(ctx, tt.args.config)

			// Move the clock ahead when there are multiple attempts
			go func() {
				for n := 1; n < tt.args.config.Attempts; n++ {
					clock.Advance(1 * time.Minute)
					clock.Trigger(DelaySleepId)
				}
			}()

			err := r.Retry(tt.args.retryable)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Retry() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestNewRetryConfigFromConfig(t *testing.T) {
	type args struct {
		cfg  *config.Config
		root string
	}
	tests := []struct {
		name    string
		args    args
		want    *RetryConfig
		wantErr bool
	}{
		{
			name: "Defaults",
			args: args{
				cfg:  configtest.NewInMemoryConfig(map[string]string{}),
				root: "some.random.spot",
			},
			want: &RetryConfig{
				Attempts: 3,
				Delay:    500,
				BackOff:  0.0,
				Linear:   true,
			},
		},
		{
			name: "Custom",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"some.random.spot.attempts": "1",
					"some.random.spot.delay":    "2",
					"some.random.spot.backoff":  "3.0",
					"some.random.spot.linear":   "false",
				}),
				root: "some.random.spot",
			},
			want: &RetryConfig{
				Attempts: 1,
				Delay:    2,
				BackOff:  3.0,
				Linear:   false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRetryConfigFromConfig(tt.args.cfg, tt.args.root)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRetryConfigFromConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRetryConfigFromConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
