// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package retry

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
	"github.com/thejerf/abtime"
	"math"
	"math/rand"
	"time"
)

const (
	configRootRetry = "spring.retry"
	DelaySleepId    = iota
)

var logger = log.NewPackageLogger()

type Retryable func() error

type failure interface {
	IsPermanent() bool
}

type Retry struct {
	Attempts int
	Delay    time.Duration
	BackOff  float64
	Linear   bool
	Jitter   time.Duration
	Context  context.Context
	clock    abtime.AbstractTime
}

// NewRetry returns a new Retry instance using the specified RetryConfig
func NewRetry(ctx context.Context, cfg RetryConfig) Retry {
	return Retry{
		Attempts: cfg.Attempts,
		Delay:    time.Duration(cfg.Delay) * time.Millisecond,
		BackOff:  cfg.BackOff,
		Linear:   cfg.Linear,
		Jitter:   time.Duration(cfg.Jitter) * time.Millisecond,
		Context:  ctx,
		clock:    types.NewClock(ctx),
	}
}

func NewRetryConfigFromConfig(cfg *config.Config, root string) (*RetryConfig, error) {
	var retryConfig RetryConfig
	if err := cfg.Populate(&retryConfig, root); err != nil {
		return nil, errors.Wrap(err, "Failed to populate default retry configuration")
	}

	return &retryConfig, nil
}

// NewRetryFromConfig returns a new Retry instance configured from the default RetryConfig in the specified *config.Config
func NewRetryFromConfig(ctx context.Context, cfg *config.Config) (*Retry, error) {
	retryConfig, err := NewRetryConfigFromConfig(cfg, configRootRetry)
	if err != nil {
		return nil, err
	}

	retry := NewRetry(ctx, *retryConfig)
	return &retry, nil
}

// NewRetryFromContext returns a new Retry instance configured from the default RetryConfig in the context's *config.Config
func NewRetryFromContext(ctx context.Context) (*Retry, error) {
	retry, err := NewRetryFromConfig(ctx, config.FromContext(ctx))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create Retry wrapper")
	}

	return retry, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Retry executes the Retryable up to Attempts times, exiting early if there is a success (no error returned) or a PermanentError occurs
func (r Retry) Retry(retryable Retryable) (err error) {
	var n int
	for n < r.Attempts {
		if n > 0 {
			logger.WithContext(r.Context).WithError(err).Errorf("Attempt %d failed, retrying after delay", n)
			r.delay(n)
		}

		if err = retryable(); err == nil {
			break
		} else if perm, ok := err.(failure); ok && perm.IsPermanent() {
			break
		}

		n++
	}

	if err != nil {
		bt := types.BackTraceFromError(err)
		if perm, ok := err.(failure); ok && perm.IsPermanent() {
			logger.
				WithContext(r.Context).
				WithError(err).
				WithFields(bt.LogFields()).
				Errorf("Attempt %d failed with permanent failure", n)
		} else {
			logger.
				WithContext(r.Context).
				WithError(err).
				WithFields(bt.LogFields()).
				Errorf("Attempt %d failed, no more attempts", n)
		}
		log.Stack(logger, r.Context, bt)
	}

	return
}

func (r Retry) delay(n int) int64 {
	currentDelay := r.GetCurrentDelay(n)

	logger.Debug("sleep / backoff ", n, currentDelay/int64(time.Millisecond))
	r.clock.Sleep(time.Duration(currentDelay), DelaySleepId)

	return currentDelay
}

func (r Retry) GetCurrentDelay(n int) (currentDelay int64) {
	if n > 1 {
		if r.Jitter == 0 {
			if r.Linear {
				currentDelay = r.Delay.Nanoseconds() + int64(n-1)*int64(float64(r.Delay.Nanoseconds())*r.BackOff)
			} else {
				currentDelay = int64(float64(r.Delay.Nanoseconds()) * math.Pow(r.BackOff, float64(n-1)))
			}

		} else {
			jitterMin := rand.Int63n(r.Jitter.Nanoseconds())

			if r.Linear { // Jitter Type Linear Backoff
				linearDelay := r.Delay.Nanoseconds() + int64(n-1)*int64(float64(r.Delay.Nanoseconds())*r.BackOff)
				currentDelay = linearDelay + jitterMin
			} else { // Jitter Type Exponential Backoff
				exponentialDelay := int64(float64(r.Delay.Nanoseconds()) * math.Pow(r.BackOff, float64(n-1)))
				currentDelay = exponentialDelay + jitterMin
			}
		}

	} else if n == 1 {
		currentDelay = r.Delay.Nanoseconds()
	}

	if r.clock == nil {
		r.clock = types.NewClock(r.Context)
	}

	return currentDelay
}

func Decorator(cfg RetryConfig) types.ActionFuncDecorator {
	return func(action types.ActionFunc) types.ActionFunc {
		return func(ctx context.Context) error {
			retry := NewRetry(ctx, cfg)
			return retry.Retry(func() error { return action(ctx) })
		}
	}
}
