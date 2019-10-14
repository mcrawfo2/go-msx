package stats

import (
	"context"
	"runtime"
	"time"
)

// TODO: https://github.com/prometheus/client_golang/blob/master/prometheus/go_collector.go#L64

const (
	statsGaugeProcGoroutines       = "proc.goroutines"
	statsGaugeProcThreads          = "proc.threads"
	statsGaugeMemAllocBytes        = "mem.alloc_bytes"
	statsGaugeMemAllocBytesTotal   = "mem.alloc_bytes_total"
	statsGaugeMemSysBytes          = "mem.sys_bytes"
	statsGaugeMemLookupsTotal      = "mem.lookups_total"
	statsGaugeMemMallocsTotal      = "mem.malloc_total"
	statsGaugeMemFreesTotal        = "mem.frees_total"
	statsGaugeMemHeapAllocBytes    = "mem.heap_alloc_bytes"
	statsGaugeMemHeapSysBytes      = "mem.heap_sys_bytes"
	statsGaugeMemHeapIdleBytes     = "mem.heap_idle_bytes"
	statsGaugeMemHeapInUseBytes    = "mem.heap_inuse_bytes"
	statsGaugeMemHeapReleasedBytes = "mem.heap_released_bytes"
	statsGaugeMemHeapObjects       = "mem.heap_objects"
	statsGaugeMemStackInUseBytes   = "mem.stack_inuse_bytes"
	statsGaugeMemStackSysBytes     = "mem.stack_sys_bytes"
	statsGaugeMemMSpanInUseBytes   = "mem.mspan_inuse_bytes"
	statsGaugeMemMSpanSysBytes     = "mem.mspan_sys_bytes"
	statsGaugeMemMCacheInUseBytes  = "mem.mcache_inuse_bytes"
	statsGaugeMemMCacheSysBytes    = "mem.mcache_sys_bytes"
	statsGaugeMemBuckHashSysBytes  = "mem.buck_hash_sys_bytes"
	statsGaugeMemGcSysBytes        = "mem.gc_sys_bytes"
	statsGaugeMemOtherSysBytes     = "mem.other_sys_bytes"
	statsGaugeMemNextGcBytes       = "mem.next_gc_bytes"
	statsGaugeMemLastGcTimeSeconds = "mem.last_gc_time_seconds"
	statsGaugeMemGcCpuFraction     = "mem.gc_cpu_fraction"
)

type VmStatsCollector struct {
	cfg *VmStatsConfig
	done chan struct{}
}

func (c *VmStatsCollector) Stop() {
	logger.Info("Stopping VM stats collection")
	close(c.done)
}

func (c *VmStatsCollector) Start(ctx context.Context) {
	logger.Info("Starting VM stats collection")
	go func() {
		ticker := time.NewTicker(c.cfg.Frequency)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if ctx.Err() != nil {
					c.Stop()
				} else {
					c.collectStats()
				}

			case <-c.done:
				break
			}
		}
	}()
}

func (c *VmStatsCollector) collectStats() {
	Gauge(statsGaugeProcGoroutines, int64(runtime.NumGoroutine()))
	threadCount, _ := runtime.ThreadCreateProfile(nil)
	Gauge(statsGaugeProcThreads, int64(threadCount))

	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)

	Gauge(statsGaugeMemAllocBytes, int64(memStats.Alloc))
	Gauge(statsGaugeMemAllocBytesTotal, int64(memStats.TotalAlloc))
	Gauge(statsGaugeMemSysBytes, int64(memStats.Sys))
	Gauge(statsGaugeMemLookupsTotal, int64(memStats.Lookups))
	Gauge(statsGaugeMemMallocsTotal, int64(memStats.Mallocs))
	Gauge(statsGaugeMemFreesTotal, int64(memStats.Frees))
	Gauge(statsGaugeMemHeapAllocBytes, int64(memStats.HeapAlloc))
	Gauge(statsGaugeMemHeapSysBytes, int64(memStats.HeapSys))
	Gauge(statsGaugeMemHeapIdleBytes, int64(memStats.HeapIdle))
	Gauge(statsGaugeMemHeapInUseBytes, int64(memStats.HeapInuse))
	Gauge(statsGaugeMemHeapReleasedBytes, int64(memStats.HeapReleased))
	Gauge(statsGaugeMemHeapObjects, int64(memStats.HeapObjects))
	Gauge(statsGaugeMemStackInUseBytes, int64(memStats.StackInuse))
	Gauge(statsGaugeMemStackSysBytes, int64(memStats.StackSys))
	Gauge(statsGaugeMemMSpanInUseBytes, int64(memStats.MSpanInuse))
	Gauge(statsGaugeMemMSpanSysBytes, int64(memStats.MSpanSys))
	Gauge(statsGaugeMemMCacheInUseBytes, int64(memStats.MCacheInuse))
	Gauge(statsGaugeMemMCacheSysBytes, int64(memStats.MCacheSys))
	Gauge(statsGaugeMemBuckHashSysBytes, int64(memStats.BuckHashSys))
	Gauge(statsGaugeMemGcSysBytes, int64(memStats.GCSys))
	Gauge(statsGaugeMemOtherSysBytes, int64(memStats.OtherSys))
	Gauge(statsGaugeMemNextGcBytes, int64(memStats.NextGC))
	FGauge(statsGaugeMemLastGcTimeSeconds, float64(memStats.LastGC / 1e9))
	FGauge(statsGaugeMemGcCpuFraction, memStats.GCCPUFraction)
}

func NewVmStatsCollectorFromConfig(cfg *VmStatsConfig) *VmStatsCollector {
	if !cfg.Enabled {
		logger.Info("VM stats collection disabled")
		return nil
	}

	return &VmStatsCollector{
		cfg:  cfg,
		done: make(chan struct{}),
	}
}
