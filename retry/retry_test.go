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
	"math"
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
					Jitter:   0,
				},
			},
			want: Retry{
				Attempts: 2,
				Delay:    250 * time.Millisecond,
				BackOff:  1.0,
				Linear:   true,
				Context:  context.Background(),
				clock:    types.NewClock(context.Background()),
				Jitter:   0,
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
		"spring.retry.jitter":   "10000",
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
				Jitter:   10000 * time.Millisecond,
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
			"spring.retry.jitter":   "4000",
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
				Jitter:   4000 * time.Millisecond,
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
			name: "OnceJitter",
			args: args{
				config: RetryConfig{
					Attempts: 1,
					Delay:    500,
					BackOff:  1.0,
					Linear:   true,
					Jitter:   2000,
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
				Jitter:   0,
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
					"some.random.spot.jitter":   "5000",
				}),
				root: "some.random.spot",
			},
			want: &RetryConfig{
				Attempts: 1,
				Delay:    2,
				BackOff:  3.0,
				Linear:   false,
				Jitter:   5000,
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

func TestRetry_GetCurrentDelay_NoDelay(t *testing.T) {
	clock := types.NewMockClock()
	ctx := types.ContextWithClock(context.Background(), clock)

	// retries without delays
	r := NewRetry(ctx, RetryConfig{
		Attempts: 2,
		Delay:    0,
		BackOff:  0.0,
		Linear:   true,
	})

	for i := 1; i < r.Attempts; i++ {
		got := r.GetCurrentDelay(i)
		want := int64(0)

		if got != want {
			t.Errorf("GetCurrentDelay() got = %v, want %v", got, want)
		}
	}
}

func TestRetry_GetCurrentDelay_FixedDelay(t *testing.T) {
	clock := types.NewMockClock()
	ctx := types.ContextWithClock(context.Background(), clock)

	r := NewRetry(ctx, RetryConfig{
		Attempts: 2,
		Delay:    1000,
		BackOff:  1.0,
		Linear:   true,
	})

	for i := 1; i < r.Attempts; i++ {
		got := r.GetCurrentDelay(i)
		want := r.Delay.Nanoseconds()
		if got != want {
			t.Errorf("GetCurrentDelay() got = %v, want %v", got, want)
		}
	}
}

func TestRetry_GetCurrentDelay_LinearDelay(t *testing.T) {
	clock := types.NewMockClock()
	ctx := types.ContextWithClock(context.Background(), clock)

	r := NewRetry(ctx, RetryConfig{
		Attempts: 5,
		Delay:    1000,
		BackOff:  1.0,
		Linear:   true,
	})

	configDelay := r.Delay.Nanoseconds()
	currentDelay := configDelay
	for i := 1; i < r.Attempts; i++ {
		got := r.GetCurrentDelay(i)
		want := currentDelay + configDelay*int64(r.BackOff) // should be equal to next line but based on previous value
		// want := configDelay + int64(i -1) * int64(float64(r.Delay.Nanoseconds()) * r.BackOff)
		if i == 1 {
			want = configDelay
		}
		if got != want {
			t.Errorf("GetCurrentDelay() got = %v, want %v", got, want)
		}

		currentDelay = got
		logger.Info(i, currentDelay/int64(time.Millisecond))
	}
}

func TestRetry_GetCurrentDelay_LinearDelayPlusBackOffFactor(t *testing.T) {
	clock := types.NewMockClock()
	ctx := types.ContextWithClock(context.Background(), clock)

	r := NewRetry(ctx, RetryConfig{
		Attempts: 5,
		Delay:    1000,
		BackOff:  1.5,
		Linear:   true,
	})

	configDelay := r.Delay.Nanoseconds()
	currentDelay := configDelay
	for i := 1; i < r.Attempts; i++ {
		got := r.GetCurrentDelay(i)
		want := currentDelay + int64(float64(configDelay)*r.BackOff)
		if i == 1 {
			want = configDelay
		}
		if got != want {
			t.Errorf("GetCurrentDelay() got = %v, want %v", got, want)
		}

		currentDelay = got
		logger.Info(i, currentDelay/int64(time.Millisecond))
	}
}

func TestRetry_GetCurrentDelay_ExponentialDelay(t *testing.T) {
	clock := types.NewMockClock()
	ctx := types.ContextWithClock(context.Background(), clock)

	r := NewRetry(ctx, RetryConfig{
		Attempts: 5,
		Delay:    1000,
		BackOff:  2.0,
		Linear:   false,
	})

	configDelay := r.Delay.Nanoseconds()
	currentDelay := configDelay
	for i := 1; i < r.Attempts; i++ {
		got := r.GetCurrentDelay(i)
		want := int64(float64(r.Delay.Nanoseconds()) * math.Pow(r.BackOff, float64(i-1)))
		if i == 1 {
			want = configDelay
		}
		if got != want {
			t.Errorf("GetCurrentDelay() got = %v, want %v", got, want)
		}

		currentDelay = got
		logger.Info(i, currentDelay/int64(time.Millisecond))
	}
}

func TestRetry_GetCurrentDelay_LinearLowJitter(t *testing.T) {
	clock := types.NewMockClock()
	ctx := types.ContextWithClock(context.Background(), clock)

	r := NewRetry(ctx, RetryConfig{
		Attempts: 5,
		Delay:    1000,
		BackOff:  1.0,
		Linear:   true,
		Jitter:   1000,
	})

	configDelay := r.Delay.Nanoseconds()
	currentDelay := configDelay
	jitter := r.Jitter.Nanoseconds()
	for i := 1; i < r.Attempts; i++ {
		got := r.GetCurrentDelay(i)
		if i != 1 {
			min := configDelay + int64(i-1)*int64(float64(r.Delay.Nanoseconds())*r.BackOff)
			max := min + jitter
			if got < min || got > max {
				t.Errorf("GetCurrentDelay() got = %v, want between %v - %v", got, min, max)
			}

		} else {
			want := configDelay
			if got != want {
				t.Errorf("GetCurrentDelay() got = %v, want %v", got, want)
			}
		}

		currentDelay = got
		logger.Info(i, currentDelay/int64(time.Millisecond))
	}
}

func TestRetry_GetCurrentDelay_LinearExtremeJitter(t *testing.T) {
	clock := types.NewMockClock()
	ctx := types.ContextWithClock(context.Background(), clock)

	r := NewRetry(ctx, RetryConfig{
		Attempts: 5,
		Delay:    1000,
		BackOff:  1.0,
		Linear:   true,
		Jitter:   20000,
	})

	configDelay := r.Delay.Nanoseconds()
	currentDelay := configDelay
	jitter := r.Jitter.Nanoseconds()
	for i := 1; i < r.Attempts; i++ {
		got := r.GetCurrentDelay(i)
		if i != 1 {
			min := configDelay + int64(i-1)*int64(float64(r.Delay.Nanoseconds())*r.BackOff)
			max := min + jitter
			if got < min || got > max {
				t.Errorf("GetCurrentDelay() got = %v, want between %v - %v", got, min, max)
			}

		} else {
			want := configDelay
			if got != want {
				t.Errorf("GetCurrentDelay() got = %v, want %v", got, want)
			}
		}

		currentDelay = got
		logger.Info(i, currentDelay/int64(time.Millisecond))
	}
}

func TestRetry_GetCurrentDelay_ExponentialJitter(t *testing.T) {
	clock := types.NewMockClock()
	ctx := types.ContextWithClock(context.Background(), clock)

	r := NewRetry(ctx, RetryConfig{
		Attempts: 5,
		Delay:    1000,
		BackOff:  2.0,
		Linear:   false,
		Jitter:   1,
	})

	configDelay := r.Delay.Nanoseconds()
	currentDelay := configDelay
	jitter := r.Jitter.Nanoseconds()
	for i := 1; i < r.Attempts; i++ {
		got := r.GetCurrentDelay(i)

		if i != 1 {
			min := int64(float64(r.Delay.Nanoseconds()) * math.Pow(r.BackOff, float64(i-1)))
			max := min + jitter
			if got < min || got > max {
				t.Errorf("GetCurrentDelay() got = %v, want between %v - %v", got, min, max)
			}

		} else {
			want := configDelay
			if got != want {
				t.Errorf("GetCurrentDelay() got = %v, want %v", got, want)
			}
		}

		currentDelay = got
		logger.Info(i, currentDelay/int64(time.Millisecond))
	}
}

func TestRetry_Decorator(t *testing.T) {
	clock := types.NewMockClock()
	ctx := types.ContextWithClock(context.Background(), clock)

	types.
		NewOperation(func(ctx context.Context) error {
			return errors.New("a transient error")
		}).
		WithDecorator(Decorator(RetryConfig{
			Attempts: 1,
			Delay:    10,
			BackOff:  2.0,
			Linear:   false,
			Jitter:   1,
		})).
		Run(ctx)
}
