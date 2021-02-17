package vaultprovider

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
)

const configRootEncryptionConfig = "per-tenant-encryption"

type KeyPropertiesConfig struct {
	Type                 string `config:"default=aes256-gcm96"`
	Exportable           *bool `config:"default=false"`
	AllowPlaintextBackup *bool `config:"default=false"`
}

type Config struct {
	Enabled          bool		`config:"default=false"`
	AlwaysCreateKeys bool		`config:"default=false"`
	KeyProperties    KeyPropertiesConfig
}

func NewEncryptionConfig(cfg *config.Config) (*Config, error) {
	var encryptionConfig Config
	if err := cfg.Populate(&encryptionConfig, configRootEncryptionConfig); err != nil {
		return nil, err
	}
	return &encryptionConfig, nil
}
