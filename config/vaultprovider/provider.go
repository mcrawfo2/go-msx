package vaultprovider

import (
	"context"
	"crypto/sha1"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/retry"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/vault"
	"fmt"
	"github.com/pkg/errors"
	"github.com/thejerf/abtime"
	"reflect"
	"sort"
	"time"
)

const (
	configRootVaultConfigProvider = "spring.cloud.vault.generic"
	configKeyAppName              = "spring.application.name"
	VaultClockBackoffTimerId      = iota
)

var logger = log.NewLogger("msx.config.vaultprovider")

type ProviderConfig struct {
	Enabled          bool          `config:"default=false"`
	Backend          string        `config:"default=secret"`
	ProfileSeparator string        `config:"default=/"`
	DefaultContext   string        `config:"default=defaultapplication"`
	Delay            time.Duration `config:"default=20s"`
}

type Provider struct {
	name        string
	cfg         *ProviderConfig
	contextPath string
	connection  vault.ConnectionApi
	loaded      chan map[string]string
	notify      chan struct{}
	clock       abtime.AbstractTime
}

func (p *Provider) Description() string {
	return fmt.Sprintf("%s: [%s]", p.name, p.ContextPath())
}

func (p *Provider) ContextPath() string {
	return fmt.Sprintf("%s%s%s", p.cfg.Backend, p.cfg.ProfileSeparator, p.contextPath)
}

func (p *Provider) Load(ctx context.Context) (entries config.ProviderEntries, err error) {
	var settings map[string]string

	select {
	case settings = <-p.loaded:
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
	// Ensure our clock is used by retry
	ctx = types.ContextWithClock(ctx, p.clock)

	err = retry.NewRetry(ctx, retry.RetryConfig{
		Attempts: 10,
		Delay:    int(p.cfg.Delay.Milliseconds()),
		BackOff:  0.0,
		Linear:   true,
	}).Retry(func() error {
		if ctx.Err() != nil {
			return &retry.PermanentError{Cause: err}
		}

		if settings, err = p.connection.ListSecrets(ctx, p.ContextPath()); err != nil {
			return errors.Wrap(err, "Failed to load configuration from vault")
		}

		return nil
	})

	return
}

func (p *Provider) backoff(ctx context.Context) {
	t := p.clock.NewTimer(p.cfg.Delay, VaultClockBackoffTimerId)
	select {
	case <-ctx.Done():
	case <-t.Channel():
	}
}

var nullSeparator = []byte{0}

func (p *Provider) settingsHash(settings map[string]string) []byte {
	var keys []string
	for k := range settings {
		keys = append(keys, k)
	}
	sort.StringSlice(keys).Sort()

	h := sha1.New()
	for _, k := range keys {
		h.Write([]byte(k))
		h.Write(nullSeparator)
		h.Write([]byte(settings[k]))
		h.Write(nullSeparator)
	}

	return h.Sum(nil)
}

func (p *Provider) Run(ctx context.Context) {
	var prefix = p.ContextPath()
	logger.WithContext(ctx).Infof("Starting config watcher for %s", p.Description())

	var lastHash []byte

	for {
		settings, err := p.connection.ListSecrets(ctx, prefix)
		if err != nil {
			if ctx.Err() != nil {
				logger.WithContext(ctx).WithError(err).Infof("Stopping vault secret watcher for %q", prefix)
				return
			}

			logger.WithContext(ctx).WithError(err).Infof("Failed to watch vault secret path %q", prefix)
			p.backoff(ctx)
			continue
		}

		newHash := p.settingsHash(settings)
		if lastHash == nil || !reflect.DeepEqual(lastHash, newHash) {
			lastHash = newHash
			p.loaded <- settings
			p.notify <- struct{}{}
		}

		p.backoff(ctx)
	}
}

func (p *Provider) Notify() <-chan struct{} {
	return p.notify
}

func NewProvider(name string, cfg *ProviderConfig, contextPath string, conn vault.ConnectionApi, clock abtime.AbstractTime) *Provider {
	return &Provider{
		name:        name,
		cfg:         cfg,
		contextPath: contextPath,
		connection:  conn,
		loaded:      make(chan map[string]string, 1),
		notify:      make(chan struct{}),
		clock:       clock,
	}
}

func NewProviderConfig(cfg *config.Config) (*ProviderConfig, error) {
	var providerConfig = &ProviderConfig{}
	var err = cfg.Populate(providerConfig, configRootVaultConfigProvider)
	if err != nil {
		return nil, err
	}
	return providerConfig, nil
}

func NewProvidersFromConfig(name string, ctx context.Context, cfg *config.Config) ([]config.Provider, error) {
	providerConfig, err := NewProviderConfig(cfg)
	if err != nil {
		err = errors.Wrap(err, "Failed to load provider config")
		return nil, err
	}

	if !providerConfig.Enabled {
		logger.Warn("Vault configuration source disabled")
		return nil, nil
	}

	var appContext string
	if appContext, err = cfg.String(configKeyAppName); err != nil {
		err = errors.Wrap(err, "Failed to retrieve app name")
		return nil, err
	}

	conn := vault.ConnectionFromContext(ctx)
	if conn == nil {
		conn, err = vault.NewConnection(ctx)
		if err == vault.ErrDisabled {
			return nil, nil
		} else if err != nil {
			err = errors.Wrap(err, "Failed to obtain vault connection")
			return nil, err
		}
	}

	clock := types.NewClock(ctx)

	return []config.Provider{
		NewProvider(name, providerConfig, providerConfig.DefaultContext, conn, clock),
		NewProvider(name, providerConfig, appContext, conn, clock),
	}, nil
}
