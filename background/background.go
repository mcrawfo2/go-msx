// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

//go:generate mockery --inpackage --name=ErrorReporter --structname=MockErrorReporter --filename mock_ErrorReporter.go

package background

import "context"

type backgroundContextKey int

const (
	contextKeyErrorReporter backgroundContextKey = iota
)

type ErrorReporter interface {
	Fatal(err error)
	NonFatal(err error)
	C() <-chan struct{}
}

func ContextWithErrorReporter(ctx context.Context, reporter ErrorReporter) context.Context {
	return context.WithValue(ctx, contextKeyErrorReporter, reporter)
}

func ErrorReporterFromContext(ctx context.Context) ErrorReporter {
	i := ctx.Value(contextKeyErrorReporter)
	if i == nil {
		return nil
	}
	return i.(ErrorReporter)
}

var _ ErrorReporter = new(MockErrorReporter)
