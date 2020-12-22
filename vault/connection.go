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
		return nil, err
	}

	if connectionConfig.Enabled == false {
		return nil, ErrDisabled
	}

	clientConfig, err := connectionConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	client, err := api.NewClient(clientConfig)
	if err != nil {
		return nil, err
	}

	var conn = &Connection{
		ConnectionApi: newTraceConnection(newStatsConnection(newConnectionImpl(connectionConfig, client))),
	}

	tokenSource, err := NewTokenSource(connectionConfig.TokenSource.Source, cfg, conn)
	if err != nil {
		return nil, err
	}

	token, err := tokenSource.GetToken(ctx)
	client.SetToken(token)

	if tokenSource.Renewable() {
		conn.renewer, err = newRenewer(client)
		if err != nil {
			return nil, err
		}
		go conn.renewer.Run(ctx)
	}

	return conn, nil
}
