package cassandra

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"github.com/gocql/gocql"
	"github.com/opentracing/opentracing-go"
	"strings"
)

type TraceObserver struct{}

func (s *TraceObserver) ObserveQuery(ctx context.Context, query gocql.ObservedQuery) {
	ctx, span := trace.NewSpan(
		ctx,
		"cassandra.query",
		opentracing.StartTime(query.Start))
	span.SetTag(trace.FieldOperation, strings.Fields(query.Statement)[0])
	span.SetTag(trace.FieldKeyspace, query.Keyspace)
	if query.Err != nil {
		span.LogFields(trace.Error(query.Err))
	}
	span.FinishWithOptions(opentracing.FinishOptions{
		FinishTime: query.End,
	})
}

func (s *TraceObserver) ObserveBatch(ctx context.Context, batch gocql.ObservedBatch) {
	ctx, span := trace.NewSpan(
		ctx,
		"cassandra.batch",
		opentracing.StartTime(batch.Start))
	span.SetTag(trace.FieldKeyspace, batch.Keyspace)
	if batch.Err != nil {
		span.LogFields(trace.Error(batch.Err))
	}
	span.FinishWithOptions(opentracing.FinishOptions{
		FinishTime: batch.End,
	})
}
