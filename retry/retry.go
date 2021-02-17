package retry

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
	"github.com/thejerf/abtime"
	"time"
)

const (
	configRootRetry = "spring.retry"
	DelaySleepId    = iota
)

var logger = log.NewLogger("msx.retry")

type Retryable func() error

type failure interface {
	IsPermanent() bool
}

type Retry struct {
	Attempts int
	Delay    time.Duration
	BackOff  float64
	Linear   bool
	Context  context.Context
	clock    abtime.AbstractTime
}

// Retry executes the Retryable up to Attempts times, exiting early if there is a success (no error returned) or a PermanentError occurs
func (r Retry) Retry(retryable Retryable) (err error) {
	currentDelay := r.Delay.Nanoseconds()
	var n int
	for n < r.Attempts {
		if n > 0 {
			logger.WithContext(r.Context).WithError(err).Errorf("Attempt %d failed, retrying after delay", n)
			currentDelay = r.delay(currentDelay, n)
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
				WithField(log.FieldStack, bt.Stanza()).
				Errorf("Attempt %d failed with permanent failure", n)
		} else {
			logger.
				WithContext(r.Context).
				WithError(err).
				WithField(log.FieldStack, bt.Stanza()).
				Errorf("Attempt %d failed, no more attempts", n)
		}
		log.Stack(logger, r.Context, bt)
	}

	return
}

func (r Retry) delay(currentDelay int64, n int) int64 {
	if n > 1 {
		if r.Linear {
			currentDelay += int64(float64(r.Delay.Nanoseconds()) * r.BackOff)
		} else {
			currentDelay = int64(float64(currentDelay) * r.BackOff)
		}
	}

	if r.clock == nil {
		r.clock = types.NewClock(r.Context)
	}

	r.clock.Sleep(time.Duration(currentDelay), DelaySleepId)

	return currentDelay
}

// NewRetry returns a new Retry instance using the specified RetryConfig
func NewRetry(ctx context.Context, cfg RetryConfig) Retry {
	return Retry{
		Attempts: cfg.Attempts,
		Delay:    time.Duration(cfg.Delay) * time.Millisecond,
		BackOff:  cfg.BackOff,
		Linear:   cfg.Linear,
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
