package cassandra

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/stats"
	"github.com/gocql/gocql"
)

const (
	statsCounterCassandraConnections   = "cassandra.connections"
	statsTimerCassandraConnectTime     = "cassandra.connectTime"
	statsCounterCassandraConnectErrors = "cassandra.connectErrors"
	statsCounterCassandraQueries       = "cassandra.queries"
	statsCounterCassandraQueryErrors   = "cassandra.queryErrors"
	statsTimerCassandraQueryTime       = "cassandra.queryTime"
	statsCounterCassandraBatches       = "cassandra.batches"
	statsCounterCassandraBatchErrors   = "cassandra.batchErrors"
	statsTimerCassandraBatchTime       = "cassandra.batchTime"
)

type StatsObserver struct{}

func (s *StatsObserver) ObserveBatch(ctx context.Context, batch gocql.ObservedBatch) {
	stats.Incr(statsCounterCassandraBatches, 1)
	stats.PrecisionTiming(statsTimerCassandraBatchTime, batch.End.Sub(batch.Start))
	if batch.Err != nil {
		stats.Incr(statsCounterCassandraBatchErrors, 1)
	}
}

func (s *StatsObserver) ObserveQuery(ctx context.Context, observedQuery gocql.ObservedQuery) {
	stats.Incr(statsCounterCassandraQueries, 1)
	stats.PrecisionTiming(statsTimerCassandraQueryTime, observedQuery.End.Sub(observedQuery.Start))
	if observedQuery.Err != nil {
		stats.Incr(statsCounterCassandraQueryErrors, 1)
	}
}

func (s *StatsObserver) ObserveConnect(observedConnect gocql.ObservedConnect) {
	stats.Incr(statsCounterCassandraConnections, 1)
	stats.PrecisionTiming(statsTimerCassandraConnectTime, observedConnect.End.Sub(observedConnect.Start))
	if observedConnect.Err != nil {
		stats.Incr(statsCounterCassandraConnectErrors, 1)
	}
}
