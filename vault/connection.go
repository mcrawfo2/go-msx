// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

//go:generate mockery --inpackage --name=ConnectionApi --structname=MockConnection --filename mock_connection.go

package vault

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
)

var (
	ErrDisabled = errors.New("Vault connection disabled")
	logger      = log.NewLogger("msx.vault")
)

type Connection struct {
	ConnectionApi
	renewer *renewer
}

func NewConnection(ctx context.Context) (*Connection, error) {
	cfg := config.FromContext(ctx)
	if cfg == nil {
		return nil, errors.New("Config not found in context")
	}

	connectionConfig, err := newConnectionConfig(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create connection config")
	}

	if connectionConfig.Enabled == false {
		return nil, ErrDisabled
	}

	clientConfig, err := connectionConfig.ClientConfig()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create ClientConfig")
	}

	client, err := api.NewClient(clientConfig)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create Client")
	}

	if connectionConfig.Disconnected {
		return &Connection{
			ConnectionApi: new(DisConnection),
		}, nil
	}

	var conn = &Connection{
		ConnectionApi: newTraceConnection(newStatsConnection(newConnectionImpl(connectionConfig, client))),
	}

	logger.WithContext(ctx).Infof("Using vault token source %q", connectionConfig.TokenSource.Source)

	tokenSource, err := NewTokenSource(connectionConfig.TokenSource.Source, cfg, conn)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create token source %q", connectionConfig.TokenSource.Source)
	}

	token, err := tokenSource.GetToken(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to obtain token from token source %q", connectionConfig.TokenSource.Source)
	}

	client.SetToken(token)

	if tokenSource.Renewable() {
		conn.renewer, err = newRenewer(client)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to create token renewer")
		}
		go conn.renewer.Run(ctx)
	}

	return conn, nil
}
