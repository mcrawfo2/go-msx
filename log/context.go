// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package log

import (
	"context"
	"github.com/sirupsen/logrus"
)

type LogContext logrus.Fields

type key int

const (
	logContextKey key = 0

	FieldSpanId   = "spanId"
	FieldTraceId  = "traceId"
	FieldParentId = "parentId"
)

func NewContextWithLogContext(ctx context.Context, logCtx LogContext) context.Context {
	return context.WithValue(ctx, logContextKey, logCtx)
}

func LogContextFromContext(ctx context.Context) (LogContext, bool) {
	logCtx, ok := ctx.Value(logContextKey).(LogContext)
	return logCtx, ok
}

func ExtendContext(ctx context.Context, logCtx LogContext) context.Context {
	newCtx := make(LogContext)
	if oldCtx, ok := LogContextFromContext(ctx); ok {
		for k, v := range oldCtx {
			newCtx[k] = v
		}
	}
	for k, v := range logCtx {
		newCtx[k] = v
	}
	return NewContextWithLogContext(ctx, newCtx)
}
