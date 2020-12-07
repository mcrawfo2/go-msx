package retry

type TransientError struct {
	Cause error
}

func (e *TransientError) Unwrap() error {
	return e.Cause
}

func (e *TransientError) IsPermanent() bool {
	return false
}

func (e *TransientError) Error() string {
	return e.Cause.Error()
}

func TransientErrorInterceptor(fn Retryable) error {
	return TransientErrorDecorator(fn)()
}

func TransientErrorDecorator(fn Retryable) Retryable {
	return func() error {
		err := fn()
		if err != nil {
			return &TransientError{
				Cause: err,
			}
		}

		return nil
	}
}

type PermanentError struct {
	Cause error
}
func (e *PermanentError) Unwrap() error {
	return e.Cause
}

func (e *PermanentError) IsPermanent() bool {
	return true
}

func (e *PermanentError) Error() string {
	return e.Cause.Error()
}

func PermanentErrorInterceptor(fn Retryable) error {
	return PermanentErrorDecorator(fn)()
}

func PermanentErrorDecorator(fn Retryable) Retryable {
	return func() error {
		err := fn()
		if err != nil {
			return &PermanentError{
				Cause: err,
			}
		}

		return nil
	}
}
