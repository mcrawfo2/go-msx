// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package cassandra

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"github.com/gocql/gocql"
	"strings"
)

type TraceObserver struct{}

func (s *TraceObserver) ObserveQuery(ctx context.Context, query gocql.ObservedQuery) {
	ctx, span := trace.NewSpan(ctx, "cassandra.query",
		trace.StartWithStartTime(query.Start),
		trace.StartWithTag(trace.FieldOperation, strings.Fields(query.Statement)[0]),
		trace.StartWithTag(trace.FieldKeyspace, query.Keyspace),
		trace.StartWithTag(trace.FieldSpanType, "db"))
	if query.Err != nil {
		span.LogFields(trace.Error(query.Err))
	}
	span.Finish(trace.FinishWithFinishTime(query.End))
}

func (s *TraceObserver) ObserveBatch(ctx context.Context, batch gocql.ObservedBatch) {
	ctx, span := trace.NewSpan(ctx, "cassandra.batch",
		trace.StartWithStartTime(batch.Start),
		trace.StartWithTag(trace.FieldKeyspace, batch.Keyspace),
		trace.StartWithTag(trace.FieldSpanType, "db"))
	if batch.Err != nil {
		span.LogFields(trace.Error(batch.Err))
	}
	span.Finish(trace.FinishWithFinishTime(batch.End))
}
