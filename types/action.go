package types

import (
	"context"
	"github.com/pkg/errors"
)

type ActionFunc func(ctx context.Context) error

type ActionFuncDecorator func(action ActionFunc) ActionFunc

func RecoverErrorDecorator() ActionFuncDecorator {
	return func(action ActionFunc) ActionFunc {
		return func(ctx context.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					var e error
					if err, ok := r.(error); ok {
						e = err
					} else {
						e = errors.Errorf("Exception: %v", r)
					}

					// TODO: decorate error with backtrace
					//bt := BackTraceFromDebugStackTrace(debug.Stack())
					err = e
				}
			}()

			err = action(ctx)
			return err
		}
	}
}
