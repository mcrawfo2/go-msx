// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package vault

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigurePool(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Defaults",
			args: args{
				ctx: configtest.ContextWithNewInMemoryConfig(context.Background(), nil),
			},
		},
		{
			name: "Disabled",
			args: args{
				ctx: configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{
					"spring.cloud.vault.enabled": "false",
				}),
			},
			wantErr: true,
		},
		{
			name: "ConfigError",
			args: args{
				ctx: configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{
					"spring.cloud.vault.enabled": "falsy",
				}),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool = nil
			if err := ConfigurePool(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("ConfigurePool() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				assert.NotNil(t, pool)
			} else {
				assert.Nil(t, pool)
			}
		})
	}
}

func TestConnectionPool_Connection(t *testing.T) {
	mockConnection := new(MockConnection)
	connection := &Connection{ConnectionApi: mockConnection}
	connectionPool := &ConnectionPool{conn: connection}

	assert.Equal(t, connection, connectionPool.Connection())
}

func TestConnectionPool_WithConnection(t *testing.T) {
	mockConnection := new(MockConnection)
	connection := &Connection{ConnectionApi: mockConnection}
	connectionPool := &ConnectionPool{conn: connection}

	err := connectionPool.WithConnection(func(conn *Connection) error {
		assert.Equal(t, connection, conn)
		return errors.New("err")
	})
	assert.Error(t, err)
}

func TestPool(t *testing.T) {
	pool = &ConnectionPool{}
	assert.Equal(t, pool, Pool())
}
