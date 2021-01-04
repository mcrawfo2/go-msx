package consulprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/consul"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/retry"
	"fmt"
	"github.com/pkg/errors"
	"time"
)

const (
	configRootConsulConfigProvider = "spring.cloud.consul.config"
	configKeyAppName               = "spring.application.name"
)

var logger = log.NewLogger("msx.config.consulprovider")

type ProviderConfig struct {
	Enabled        bool   `config:"default=false"`
	Prefix         string `config:"default=userviceconfiguration"`
	DefaultContext string `config:"default=defaultapplication"`
	Pool           bool   `config:"default=false"`
	Delay          int    `config:"default=3"`
}

type Provider struct {
	name        string
	cfg         *ProviderConfig
	contextPath string
	connection  *consul.Connection
	loaded      chan map[string]string
	notify      chan struct{}
}

func (p *Provider) Description() string {
	return fmt.Sprintf("%s: [%s]", p.name, p.ContextPath())
}

func (p *Provider) ContextPath() string {
	return fmt.Sprintf("%s/%s", p.cfg.Prefix, p.contextPath)
}

func (p *Provider) Load(ctx context.Context) (entries config.ProviderEntries, err error) {
	var settings map[string]string

	select {
	case settings = <- p.loaded:
		// receive from async loop change notification
	default:
		settings, err = p.loadSettings(ctx)
	}

	if err != nil {
		return nil, err
	}

	for k, v := range settings {
		entries = append(entries, config.NewEntry(p, k, v))
	}

	return entries, nil
}

func (p *Provider) loadSettings(ctx context.Context) (settings map[string]string, err error) {
	var waitTime = new(time.Duration)
	*waitTime = time.Nanosecond

	err = retry.NewRetry(ctx, retry.RetryConfig{
		Attempts: 10,
		Delay:    1000 * p.cfg.Delay,
		BackOff:  0.0,
		Linear:   true,
	}).Retry(func() error {
		if ctx.Err() != nil {
			return &retry.PermanentError{Cause: err}
		}
		if _, settings, err = p.connection.WatchKeyValuePairs(ctx, p.ContextPath(), nil, waitTime); err != nil {
			return errors.Wrap(err, "Failed to load configuration from consul")
		}
		return nil
	})

	if err != nil {
		settings = nil
	}

	return settings, err
}

func (p *Provider) Run(ctx context.Context) {
	logger.WithContext(ctx).Infof("Starting config watcher for %s", p.Description())
	var lastIndex *uint64

	for {
		foundIndex, settings, err := p.connection.WatchKeyValuePairs(ctx, p.ContextPath(), lastIndex, nil)
		if err != nil {
			if ctx.Err() != nil {
				logger.WithContext(ctx).WithError(err).Infof("Stopping config watcher for %s", p.Description())
				return
			}

			logger.WithContext(ctx).WithError(err).Infof("Failed to watch config %s", p.Description())
			p.backoff(ctx)
			continue
		}

		if lastIndex == nil || foundIndex != *lastIndex {
			if lastIndex == nil {
				lastIndex = new(uint64)
			}
			*lastIndex = foundIndex
			p.loaded <- settings
			p.notify <- struct{}{}
		}
	}
}

func (p *Provider) backoff(ctx context.Context) {
	t := time.NewTimer(time.Duration(p.cfg.Delay) * time.Second)
	select {
	case <-ctx.Done():
	case <-t.C:
	}
}

func (p *Provider) Notify() <-chan struct{} {
	return p.notify
}

func NewProvidersFromConfig(name string, cfg *config.Config) ([]config.Provider, error) {
	var providerConfig = &ProviderConfig{}
	var err = cfg.Populate(providerConfig, configRootConsulConfigProvider)
	if err != nil {
		return nil, err
	}

	if !providerConfig.Enabled {
		logger.Warn("Consul configuration source disabled")
		return nil, nil
	}

	var appContext string
	if appContext, err = cfg.String(configKeyAppName); err != nil {
		return nil, err
	}

	var conn *consul.Connection
	if providerConfig.Pool {
		if err = consul.ConfigurePool(cfg); err != nil {
			return nil, err
		}
		conn = consul.Pool().Connection()
	} else if conn, err = consul.NewConnectionFromConfig(cfg); err != nil {
		return nil, err
	}

	return []config.Provider{
		&Provider{
			name:        name,
			cfg:         providerConfig,
			contextPath: providerConfig.DefaultContext,
			connection:  conn,
			loaded:      make(chan map[string]string, 1),
			notify:      make(chan struct{}),
		},
		&Provider{
			name:        name,
			cfg:         providerConfig,
			contextPath: appContext,
			connection:  conn,
			loaded:      make(chan map[string]string, 1),
			notify:      make(chan struct{}),
		},
	}, nil
}
