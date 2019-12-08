package stats

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/trace"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"time"
)

const (
	ConfigKeyStatsPush = "stats.push"
)

var (
	ErrDisabled  = errors.New("Stats collector disabled")
	globalPusher *Pusher
	logger       = log.NewLogger("msx.stats")
)

type PushConfig struct {
	Enabled   bool          `config:"default=false"`
	Url       string        `config:"default="` // no default
	JobName   string        `config:"default=go_msx"`
	Frequency time.Duration `config:"default=15s"`
}

func NewPushConfigFromConfig(cfg *config.Config) (*PushConfig, error) {
	pushConfig := &PushConfig{}
	if err := cfg.Populate(pushConfig, ConfigKeyStatsPush); err != nil {
		return nil, err
	}

	return pushConfig, nil
}

type Pusher struct {
	ctx  context.Context
	cfg  *PushConfig
	done chan struct{}
}

func (p *Pusher) Start() error {
	go p.run()
	return nil
}

func (p *Pusher) run() {
	ctx := trace.UntracedContextFromContext(p.ctx)
	ticker := time.NewTicker(p.cfg.Frequency)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if ctx.Err() != nil {
				return
			}

			if err := p.push(); err != nil {
				logger.WithContext(ctx).WithError(err).Error("Push metrics failed")
			}

		case <-p.done:
			return
		}
	}
}

func (p *Pusher) push() error {
	return push.
		New(p.cfg.Url, p.cfg.JobName).
		Gatherer(prometheus.DefaultGatherer).
		Push()
}

func (p *Pusher) Stop() error {
	close(p.done)
	return nil
}

func NewPusher(ctx context.Context, cfg *PushConfig) (*Pusher, error) {
	if !cfg.Enabled {
		return nil, ErrDisabled
	}

	return &Pusher{
		ctx:  ctx,
		cfg:  cfg,
		done: make(chan struct{}),
	}, nil
}

func NewPusherFromConfig(ctx context.Context, cfg *config.Config) (*Pusher, error) {
	pushConfig, err := NewPushConfigFromConfig(cfg)
	if err != nil {
		return nil, err
	}

	return NewPusher(ctx, pushConfig)
}

func Configure(ctx context.Context) (err error) {
	var cfg *config.Config

	if cfg = config.FromContext(ctx); cfg == nil {
		return errors.New("Failed to retrieve cfg from context")
	}

	globalPusher, err = NewPusherFromConfig(ctx, cfg)
	return err
}

func Start(context.Context) error {
	if globalPusher == nil {
		return ErrDisabled
	}
	return globalPusher.Start()
}

func Stop(context.Context) error {
	if globalPusher == nil {
		return ErrDisabled
	}

	return globalPusher.Stop()
}
