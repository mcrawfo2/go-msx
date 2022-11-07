// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package vaultprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/retry"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/vault"
	"fmt"
	"github.com/pkg/errors"
	"github.com/thejerf/abtime"
	"reflect"
	"time"
)

const (
	configRootVaultConfigProvider = "spring.cloud.vault.generic"
	configKeyAppName              = "spring.application.name"
	VaultClockBackoffTimerId      = iota
	configRootInstallerPasswords  = "installer.passwords"
)

type ProviderConfig struct {
	Enabled          bool          `config:"default=false"`
	Backend          string        `config:"default=secret"`
	ProfileSeparator string        `config:"default=/"`
	DefaultContext   string        `config:"default=defaultapplication"`
	Delay            time.Duration `config:"default=1h"`
}

type InstallerPasswordsConfig struct {
	Enabled       bool   `config:"default=false"`
	Context       string `config:"default=deploymentpasswords/ansible_pass_file"`
	Prefix        string `config:"default=installer.passwords.data"`
	YamlDataField string `config:"default=ansible_pass_file"`
}

type Provider struct {
	name         string
	cfg          *ProviderConfig
	contextPath  string
	connection   vault.ConnectionApi
	loaded       chan map[string]string
	notify       chan struct{}
	clock        abtime.AbstractTime
	installerCfg *InstallerPasswordsConfig
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

	isInstallerPasswords := p.installerCfg != nil

	for k, v := range settings {
		if !isInstallerPasswords {
			entries = append(entries, config.NewEntry(p, k, v))

		} else if k == p.installerCfg.YamlDataField {
			// convert yaml
			ents, err := config.ParseYaml(func() ([]byte, error) { return []byte(v), nil }, p)
			if err != nil {
				return nil, err
			}
			for _, e1 := range ents {
				k1 := config.PrefixWithName(p.installerCfg.Prefix, e1.Name) // add prefix
				entries = append(entries, config.NewEntry(p, k1, e1.Value))
			}
		}
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

		newHash := config.SettingsHash(settings)
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

func newInstallerPasswordsProviderConfig(cfg *config.Config) (*InstallerPasswordsConfig, error) {
	var installerPasswordsProviderConfig = &InstallerPasswordsConfig{}
	var err = cfg.Populate(installerPasswordsProviderConfig, configRootInstallerPasswords)
	if err != nil {
		return nil, err
	}
	return installerPasswordsProviderConfig, nil
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

	providers := []config.Provider{
		NewProvider(name, providerConfig, providerConfig.DefaultContext, conn, clock),
		NewProvider(name, providerConfig, appContext, conn, clock),
	}

	installerPasswordsProviderConfig, err := newInstallerPasswordsProviderConfig(cfg)
	if err != nil {
		err = errors.Wrap(err, "Failed to load installer passwords provider config")
		return nil, err
	}

	if !installerPasswordsProviderConfig.Enabled {
		logger.Warn("installer passwords configuration source disabled")

	} else { // installer passwords enabled
		// kv v1 only for now
		installerPasswordsProvider := NewProvider(name, providerConfig, installerPasswordsProviderConfig.Context, conn, clock)
		installerPasswordsProvider.installerCfg = installerPasswordsProviderConfig
		providers = append(providers, installerPasswordsProvider)
	}

	return providers, nil
}
