package redis

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"github.com/go-redis/redis/v7"
	"strings"
)

type traceHook struct{}

func (s *traceHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	ctx, _ = trace.NewSpan(ctx, "redis.cmd."+strings.ToLower(cmd.Name()))
	return ctx, nil
}

func (s *traceHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	span := trace.SpanFromContext(ctx)
	if cmd.Err() != nil {
		span.LogFields(trace.Error(cmd.Err()))
	}
	span.Finish()
	return nil
}

func (s *traceHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	ctx, _ = trace.NewSpan(ctx, "redis.pipeline")
	return ctx, nil
}

func (s *traceHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	span := trace.SpanFromContext(ctx)
	span.Finish()
	return nil
}
