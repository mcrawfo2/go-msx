// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package transit

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
)

const configRootEncryptionConfig = "per-tenant-encryption"

type Config struct {
	Enabled          bool `config:"default=false"`
}

func NewConfig(ctx context.Context) (*Config, error) {
	var encryptionConfig Config
	if err := config.FromContext(ctx).Populate(&encryptionConfig, configRootEncryptionConfig); err != nil {
		return nil, err
	}
	return &encryptionConfig, nil
}

func ConfigureEncrypterFactory(ctx context.Context) error {
	cfg, err := NewConfig(ctx)
	if err != nil {
		return err
	}

	if cfg.Enabled {
		SetEncrypterFactory(NewProductionEncrypter)
		SetBulkEncrypterFactory(NewProductionBulkEncrypter)
	} else {
		SetEncrypterFactory(NewDummyEncrypter)
		SetBulkEncrypterFactory(NewDummyBulkEncrypter)
	}

	return nil
}
