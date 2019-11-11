package cassandra

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/stats"
	"github.com/gocql/gocql"
	"time"
)

const (
	statsSubsystemCassandra              = "cassandra"
	statsTimerCassandraConnectTime       = "connect_time"
	statsCounterCassandraConnectAttempts = "connections"
	statsCounterCassandraConnectErrors   = "connect_errors"
	statsCounterCassandraQueries         = "queries"
	statsCounterCassandraQueryErrors     = "query_errors"
	statsTimerCassandraQueryTime         = "query_time"
	statsCounterCassandraBatches         = "batches"
	statsCounterCassandraBatchErrors     = "batch_errors"
	statsTimerCassandraBatchTime         = "batch_time"
)

var (
	histCassandraConnectTime    = stats.NewHistogram(statsSubsystemCassandra, statsTimerCassandraConnectTime, nil)
	countCassandraConnections   = stats.NewCounter(statsSubsystemCassandra, statsCounterCassandraConnectAttempts)
	countCassandraConnectErrors = stats.NewCounter(statsSubsystemCassandra, statsCounterCassandraConnectErrors)
	countCassandraQueries       = stats.NewCounter(statsSubsystemCassandra, statsCounterCassandraQueries)
	countCassandraQueryErrors   = stats.NewCounter(statsSubsystemCassandra, statsCounterCassandraQueryErrors)
	histCassandraQueryTime      = stats.NewHistogram(statsSubsystemCassandra, statsTimerCassandraQueryTime, nil)
	countCassandraBatches       = stats.NewCounter(statsSubsystemCassandra, statsCounterCassandraBatches)
	countCassandraBatchErrors   = stats.NewCounter(statsSubsystemCassandra, statsCounterCassandraBatchErrors)
	histCassandraBatchTime      = stats.NewHistogram(statsSubsystemCassandra, statsTimerCassandraBatchTime, nil)
)

type StatsObserver struct{}

func (s *StatsObserver) ObserveBatch(ctx context.Context, batch gocql.ObservedBatch) {
	countCassandraBatches.Inc()
	histCassandraBatchTime.Observe(float64(batch.End.Sub(batch.Start)) / float64(time.Millisecond))
	if batch.Err != nil {
		countCassandraBatchErrors.Inc()
	}
}

func (s *StatsObserver) ObserveQuery(ctx context.Context, query gocql.ObservedQuery) {
	countCassandraQueries.Inc()
	histCassandraQueryTime.Observe(float64(query.End.Sub(query.Start)) / float64(time.Millisecond))
	if query.Err != nil {
		countCassandraQueryErrors.Inc()
	}
}

func (s *StatsObserver) ObserveConnect(connect gocql.ObservedConnect) {
	countCassandraConnections.Inc()
	histCassandraConnectTime.Observe(float64(connect.End.Sub(connect.Start)) / float64(time.Millisecond))
	if connect.Err != nil {
		countCassandraConnectErrors.Inc()
	}
}
