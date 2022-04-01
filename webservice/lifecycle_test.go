// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package webservice

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/fs"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigureWebServer(t *testing.T) {
	tests := []struct {
		name        string
		cfg         *config.Config
		wantErr     bool
		wantService bool
	}{
		{
			name: "Disabled",
			cfg: configtest.NewInMemoryConfig(map[string]string{
				"server.enabled": "false",
			}),
			wantErr:     true,
			wantService: false,
		},
		{
			name: "Enabled",
			cfg: configtest.NewInMemoryConfig(map[string]string{
				"server.enabled": "true",
			}),
			wantErr:     false,
			wantService: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service = nil
			// TODO: Fix this
			ctx := config.ContextWithConfig(context.Background(), tt.cfg)
			// TODO: Fix this
			_ = fs.SetSources()
			if err := ConfigureWebServer(tt.cfg, ctx); (err != nil) != tt.wantErr {
				t.Errorf("ConfigureWebServer() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.wantService, service != nil)
		})
	}
}

func TestNewWebServerFromConfig(t *testing.T) {
	t.Skipped()
}

func TestStart(t *testing.T) {
	t.Skipped()
}

func TestStop(t *testing.T) {
	t.Skipped()
}
