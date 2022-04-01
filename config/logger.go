// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package config

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
)

var logger = log.NewLogger("msx.config")

type ChangeLogger struct {
	cfg *Config
}

func (c *ChangeLogger) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case snapshot := <-c.cfg.Notify():
			for _, change := range snapshot.Delta {
				if change.IsSet() {
					if change.OldEntry.ProviderEntry.Source == nil {
						logger.WithContext(ctx).Infof("Value %q added", change.NewEntry.NormalizedName)
					} else {
						logger.WithContext(ctx).Infof("Value %q updated", change.NewEntry.NormalizedName)
					}
				} else {
					logger.WithContext(ctx).Infof("Value %q removed", change.OldEntry.NormalizedName)
				}
			}
		}
	}
}

func NewChangeLogger(cfg *Config) *ChangeLogger {
	return &ChangeLogger{cfg: cfg}
}
