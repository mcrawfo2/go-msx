// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

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
					logger.
						WithContext(ctx).
						WithError(e).
						WithFields(bt.LogFields()).
						Error("Recovered from panic")
					Stack(logger, ctx, bt)
				}
			}()

			return action(ctx)
		}
	}
}

func ErrorLogDecorator(logger *Logger, actionName string) types.ActionFuncDecorator {
	return func(action types.ActionFunc) types.ActionFunc {
		return func(ctx context.Context) error {
			err := action(ctx)
			if err != nil {
				bt := types.BackTraceFromError(err)
				logger.
					WithContext(ctx).
					WithError(err).
					WithFields(bt.LogFields()).
					Errorf("Action %q returned error", actionName)
				Stack(logger, ctx, bt)
			}
			return nil
		}
	}
}
