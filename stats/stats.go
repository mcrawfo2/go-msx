package stats

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"fmt"
	"github.com/pkg/errors"
	"github.com/smira/go-statsd"
	"time"
)

const (
	ConfigKeyStats = "stats"
)

var (
	ErrDisabled     = errors.New("Stats collector disabled")
	globalCollector *Collector
	logger          = log.NewLogger("msx.stats")
)

type CollectorConfig struct {
	Enabled       bool          `config:"default=false"`
	Host          string        `config:"default=localhost"`
	Port          int           `config:"default=8125"`
	MaxPacketSize int           `config:"default=1400"`
	Prefix        string        `config:"default=msx."`
	FlushInterval time.Duration `config:"default=10s"`
	Vm            VmStatsConfig
}

type VmStatsConfig struct {
	Enabled   bool          `config:"default=true"`
	Frequency time.Duration `config:"default=30s"`
}

func (c *CollectorConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type Collector struct {
	config *CollectorConfig
	client *statsd.Client
	vm     *VmStatsCollector
}

func (c *Collector) Incr(stat string, count int64) {
	if c != nil && c.client != nil {
		c.client.Incr(stat, count)
	}
}

func (c *Collector) Decr(stat string, count int64) {
	if c != nil && c.client != nil {
		c.client.Decr(stat, count)
	}
}

func (c *Collector) Timing(stat string, delta int64) {
	if c != nil && c.client != nil {
		c.client.Timing(stat, delta)
	}
}

func (c *Collector) PrecisionTiming(stat string, delta time.Duration) {
	if c != nil && c.client != nil {
		c.client.PrecisionTiming(stat, delta)
	}
}

func (c *Collector) Gauge(stat string, value int64) {
	if c != nil && c.client != nil {
		c.client.Gauge(stat, value)
	}
}

func (c *Collector) GaugeDelta(stat string, value int64) {
	if c != nil && c.client != nil {
		c.client.GaugeDelta(stat, value)
	}
}

func (c *Collector) FGauge(stat string, value float64) {
	if c != nil && c.client != nil {
		c.client.FGauge(stat, value)
	}
}

func (c *Collector) FGaugeDelta(stat string, value float64) {
	if c != nil && c.client != nil {
		c.client.FGaugeDelta(stat, value)
	}
}

func (c *Collector) Close() {
	if c != nil {
		if c.client != nil {
			if err := c.client.Close(); err != nil {
				logger.Error(err)
			}
		}

		if c.vm != nil {
			c.vm.Stop()
		}
	}
}

func NewCollector(collectorConfig *CollectorConfig) (*Collector, error) {
	if !collectorConfig.Enabled {
		return nil, ErrDisabled
	}

	return &Collector{
		config: collectorConfig,
		client: statsd.NewClient(
			collectorConfig.Address(),
			statsd.MaxPacketSize(collectorConfig.MaxPacketSize),
			statsd.MetricPrefix(collectorConfig.Prefix),
			statsd.FlushInterval(collectorConfig.FlushInterval),
		),
		vm: NewVmStatsCollectorFromConfig(&collectorConfig.Vm),
	}, nil
}

func NewCollectorFromConfig(cfg *config.Config) (*Collector, error) {
	collectorConfig := &CollectorConfig{}
	if err := cfg.Populate(collectorConfig, ConfigKeyStats); err != nil {
		return nil, err
	}

	return NewCollector(collectorConfig)
}

func Configure(ctx context.Context) error {
	var err error
	var cfg *config.Config

	if cfg = config.FromContext(ctx); cfg == nil {
		return errors.New("Failed to retrieve config from context")
	}

	if globalCollector, err = NewCollectorFromConfig(cfg); err != nil {
		// no-op collector
		globalCollector = &Collector{}

		if err != ErrDisabled {
			return err
		}
	} else if globalCollector.vm != nil {
		globalCollector.vm.Start(ctx)
	}

	return nil
}

func Close() {
	globalCollector.Close()
}

func Incr(stat string, count int64) {
	globalCollector.Incr(stat, count)
}

func Decr(stat string, count int64) {
	globalCollector.Decr(stat, count)
}

func Timing(stat string, delta int64) {
	globalCollector.Timing(stat, delta)
}

func PrecisionTiming(stat string, delta time.Duration) {
	globalCollector.PrecisionTiming(stat, delta)
}

func Gauge(stat string, value int64) {
	globalCollector.Gauge(stat, value)
}

func GaugeDelta(stat string, value int64) {
	globalCollector.GaugeDelta(stat, value)
}

func FGauge(stat string, value float64) {
	globalCollector.FGauge(stat, value)
}

func FGaugeDelta(stat string, value float64) {
	globalCollector.FGaugeDelta(stat, value)
}