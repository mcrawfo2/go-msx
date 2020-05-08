package vaultprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/transit"
	"cto-github.cisco.com/NFV-BU/go-msx/vault"
)

var logger = log.NewLogger("msx.transit.vaultprovider")

type Provider struct {
	cfg *Config
}

func (p Provider) CreateKey(ctx context.Context, keyName string) (err error) {
	if !p.cfg.Enabled && !p.cfg.AlwaysCreateKeys {
		logger.WithContext(ctx).Debugf("Skipping key creation for tenant %s - per-tenant encryption and automatic key creation are disabled", keyName)
		return nil
	}

	return vault.PoolFromContext(ctx).WithConnection(func(connection *vault.Connection) error {
		createRequest := vault.CreateTransitKeyRequest{
			Type:                 p.cfg.KeyProperties.Type,
			Exportable:           p.cfg.KeyProperties.Exportable,
			AllowPlaintextBackup: p.cfg.KeyProperties.AllowPlaintextBackup,
		}
		return connection.CreateTransitKey(ctx, keyName, createRequest)
	})
}

func (p Provider) Encrypt(ctx context.Context, value transit.Value) (secureValue transit.Value, err error) {
	if value.IsEmpty() || !p.cfg.Enabled {
		return value, nil
	}

	var ciphertext string
	err = vault.PoolFromContext(ctx).WithConnection(func(connection *vault.Connection) error {
		ciphertext, err = connection.TransitEncrypt(ctx, value.KeyName(), value.RawPayload())
		return err
	})
	if err != nil {
		return
	}

	return value.WithEncryptedPayload(ciphertext), nil
}

func (p Provider) Decrypt(ctx context.Context, secureValue transit.Value) (value transit.Value, err error) {
	if secureValue.IsEmpty() || !p.cfg.Enabled {
		return secureValue, nil
	}

	var plaintext string
	err = vault.PoolFromContext(ctx).WithConnection(func(connection *vault.Connection) error {
		plaintext, err = connection.TransitDecrypt(ctx, secureValue.KeyName(), secureValue.RawPayload())
		return err
	})
	if err != nil {
		return
	}

	return secureValue.WithDecryptedPayload(plaintext), nil
}

func RegisterVaultTransitProvider(ctx context.Context) error {
	cfg, err := NewEncryptionConfig(config.FromContext(ctx))
	if err != nil {
		return err
	}

	if err := config.FromContext(ctx).Populate(&cfg, configRootEncryptionConfig); err != nil {
		return err
	}
	logger.Infof("Per-Tenant Encryption Enabled: %t", cfg.Enabled)
	return transit.RegisterProvider(&Provider{
		cfg: cfg,
	})
}
