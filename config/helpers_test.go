// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package config

import "context"

func NewInMemoryConfig(values map[string]string) *Config {
	provider := NewInMemoryProvider("testdata", values)
	cfg := NewConfig(provider)
	_ = cfg.Load(context.Background())
	return cfg
}
