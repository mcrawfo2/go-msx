package retry

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"github.com/pkg/errors"
	"time"
)

const (
	configRootRetry = "spring.retry"
)

var retryLogger = log.NewLogger("msx.retry")

type Retryable func() error

type PermanentFailure interface {
	IsPermanent() bool
}

type PermanentError struct {
	Cause error
}

func (e *PermanentError) IsPermanent() bool {
	return true
}

func (e *PermanentError) Error() string {
	return e.Cause.Error()
}

type RetryConfig struct {
	Attempts int     `config:"default=3"`
	Delay    int     `config:"default=500"`
	BackOff  float64 `config:"default=0.0"`
	Linear   bool    `config:"default=true"`
}

type Retry struct {
	Attempts int
	Delay    time.Duration
	BackOff  float64
	Linear   bool
	Context  context.Context
}

func (r Retry) Retry(retryable Retryable) (err error) {
	currentDelay := r.Delay.Nanoseconds()
	var n int
	for n < r.Attempts {
		if n > 0 {
			retryLogger.WithContext(r.Context).WithError(err).Errorf("Attempt %d failed, retrying after delay", n)
			currentDelay = r.delay(currentDelay, n)
		}

		if err = retryable(); err == nil {
			break
		} else if perm, ok := err.(PermanentFailure); ok && perm.IsPermanent() {
			break
		}

		n++
	}

	if err != nil {
		if perm, ok := err.(PermanentFailure); ok && perm.IsPermanent() {
			retryLogger.WithContext(r.Context).WithError(err).Errorf("Attempt %d failed with permanent failure", n)
		} else {
			retryLogger.WithContext(r.Context).WithError(err).Errorf("Attempt %d failed, no more attempts", n)
		}
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

	time.Sleep(time.Duration(currentDelay))

	return currentDelay
}

func NewRetry(ctx context.Context, cfg RetryConfig) Retry {
	return Retry{
		Attempts: cfg.Attempts,
		Delay:    time.Duration(cfg.Delay) * time.Millisecond,
		BackOff:  cfg.BackOff,
		Linear:   cfg.Linear,
		Context:  ctx,
	}
}

func NewRetryFromConfig(ctx context.Context, cfg *config.Config) (*Retry, error) {
	var retryConfig RetryConfig
	if err := cfg.Populate(&retryConfig, configRootRetry); err != nil {
		return nil, errors.Wrap(err, "Failed to populate default retry configuration")
	}

	retry := NewRetry(ctx, retryConfig)
	return &retry, nil
}

func NewRetryFromContext(ctx context.Context) (*Retry, error) {
	retry, err := NewRetryFromConfig(ctx, config.FromContext(ctx))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create Retry wrapper")
	}

	return retry, nil
}
