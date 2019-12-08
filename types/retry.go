package types

import (
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"time"
)

var retryLogger = log.NewLogger("msx.types.retry")

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

type Retry struct {
	Attempts int           `config:"default=3"`
	Delay    time.Duration `config:"default=100ms"`
	BackOff  float64       `config:"default=2.0"`
	Linear   bool          `config:"default=false"`
}

func (r Retry) Retry(retryable Retryable) (err error) {
	currentDelay := r.Delay.Nanoseconds()
	var n int
	for n < r.Attempts {
		if n > 0 {
			retryLogger.WithError(err).Errorf("Attempt %d failed, retrying after delay", n)
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
			retryLogger.WithError(err).Errorf("Attempt %d failed with permanent failure", n)
		} else {
			retryLogger.WithError(err).Errorf("Attempt %d failed, no more attempts", n)
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
