package vaultprovider

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
)

const configRootEncryptionConfig = "per-tenant-encryption"

type Config struct {
	Enabled          bool
	AlwaysCreateKeys bool
	KeyProperties    struct {
		Type                 string
		Exportable           *bool
		AllowPlaintextBackup *bool
	}
}

func NewEncryptionConfig(cfg *config.Config) (*Config, error) {
	var encryptionConfig Config
	if err := cfg.Populate(&encryptionConfig, configRootEncryptionConfig); err != nil {
		return nil, err
	}
	return &encryptionConfig, nil
}
