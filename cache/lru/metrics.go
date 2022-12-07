// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package lru

import (
	"cto-github.cisco.com/NFV-BU/go-msx/stats"
	"github.com/prometheus/client_golang/prometheus"
)

type metricsObserver interface {
	OnEntriesInc()
	OnEntriesResize(newSize int)
	OnHit()
	OnMiss()
	OnSet()
	OnGC(howMany int)
	OnEvict(howMany int)
	OnDeAge(ttlRemainingMs int64)
}

// Metric identifiers
const (
	subsystemName = "cache"
	entriesName   = "entries"
	hitsName      = "hits"
	missesName    = "misses"
	setsName      = "sets"
	evictionsName = "evictions"
	gcRunsName    = "gc_runs"
	gcSizesName   = "gc_sizes"
	deAgedName    = "deaged"
)

// metrics holds the metrics for a cache
type prometheusMetricsObserver struct {
	entries   prometheus.Gauge
	hits      prometheus.Counter
	misses    prometheus.Counter
	sets      prometheus.Counter
	evictions prometheus.Counter
	gcRuns    prometheus.Counter
	gcSizes   prometheus.Histogram
	deAgedAt  prometheus.Histogram
}

func newPrometheusMetrics(name string) *prometheusMetricsObserver {
	return &prometheusMetricsObserver{
		entries:   stats.NewGauge(name, entriesName),
		hits:      stats.NewCounter(name, hitsName),
		misses:    stats.NewCounter(name, missesName),
		sets:      stats.NewCounter(name, setsName),
		evictions: stats.NewCounter(name, evictionsName),
		gcRuns:    stats.NewCounter(name, gcRunsName),
		gcSizes:   stats.NewHistogram(name, gcSizesName, prometheus.ExponentialBuckets(1, 2, 10)),
		deAgedAt:  stats.NewHistogram(name, deAgedName, prometheus.ExponentialBuckets(10, 10, 6)),
	}
}

func (o *prometheusMetricsObserver) OnEntriesInc() {
	o.entries.Inc()
}

func (o *prometheusMetricsObserver) OnEntriesResize(newSize int) {
	o.entries.Set(float64(newSize))
}

func (o *prometheusMetricsObserver) OnHit() {
	o.hits.Inc()
}

func (o *prometheusMetricsObserver) OnMiss() {
	o.misses.Inc()
}

func (o *prometheusMetricsObserver) OnSet() {
	o.sets.Inc()
}

func (o *prometheusMetricsObserver) OnGC(howMany int) {
	o.gcRuns.Inc()
	o.gcSizes.Observe(float64(howMany))
}

func (o *prometheusMetricsObserver) OnEvict(howMany int) {
	o.evictions.Add(float64(howMany))
	o.entries.Add(float64(-howMany))
}

// OnDeAge expects the TTL remaining in milliseconds
func (o *prometheusMetricsObserver) OnDeAge(ttlRemainingMs int64) {
	o.deAgedAt.Observe(float64(ttlRemainingMs))
}

// Noop metrics observer
type nullMetricsObserver struct {
}

func (o *nullMetricsObserver) OnEntriesInc() {
}

func (o *nullMetricsObserver) OnEntriesResize(_ int) {
}

func (o *nullMetricsObserver) OnHit() {
}

func (o *nullMetricsObserver) OnMiss() {
}

func (o *nullMetricsObserver) OnSet() {
}

func (o *nullMetricsObserver) OnGC(_ int) {
}

func (o *nullMetricsObserver) OnEvict(_ int) {
}

func (o *nullMetricsObserver) OnDeAge(_ int64) {
}
