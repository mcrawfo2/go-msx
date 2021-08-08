package log

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
	"runtime/debug"
)

func RecoverLogDecorator(logger *Logger) types.ActionFuncDecorator {
	return func(action types.ActionFunc) types.ActionFunc {
		return func(ctx context.Context) error {
			defer func() {
				if r := recover(); r != nil {
					var e error
					if err, ok := r.(error); ok {
						e = err
					} else {
						e = errors.Errorf("Exception: %v", r)
					}

					logger.WithContext(ctx).WithError(e).Error("Recovered from panic")
					bt := types.BackTraceFromDebugStackTrace(debug.Stack())
					Stack(logger, ctx, bt)
				}
			}()

			return action(ctx)
		}
	}
}
