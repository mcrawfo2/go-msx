// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package trace

import "context"

type contextKey string

const (
	contextKeySpan       = contextKey("Span")
	contextKeyParentSpan = contextKey("ParentSpan")
)

func SpanFromContext(ctx context.Context) Span {
	span, _ := ctx.Value(contextKeySpan).(Span)
	return span
}

func ContextWithSpan(ctx context.Context, span Span) context.Context {
	return context.WithValue(ctx, contextKeySpan, span)
}

func ParentContextFromContext(ctx context.Context) SpanContext {
	spanContext, _ := ctx.Value(contextKeyParentSpan).(SpanContext)
	return spanContext
}

func ContextWithParentContext(ctx context.Context, spanContext SpanContext) context.Context {
	return context.WithValue(ctx, contextKeyParentSpan, spanContext)
}
