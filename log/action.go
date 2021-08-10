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

					bt := types.BackTraceFromDebugStackTrace(debug.Stack())
					logger.WithContext(ctx).WithError(e).WithField(FieldStack, bt.Stanza()).Error("Recovered from panic")
					Stack(logger, ctx, bt)
				}
			}()

			return action(ctx)
		}
	}
}
