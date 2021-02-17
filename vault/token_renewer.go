package vault

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/vault/api"
	"time"
)

type renewer struct {
	client  *api.Client
	renewer *api.Renewer
}

func newRenewer(client *api.Client) (*renewer, error) {
	secret, err := client.Auth().Token().LookupSelf()
	if err != nil {
		return nil, err
	}
	var renewable bool
	if v, ok := secret.Data["renewable"]; ok {
		renewable, _ = v.(bool)
	}
	var increment int64
	if v, ok := secret.Data["ttl"]; ok {
		if n, ok := v.(json.Number); ok {
			increment, _ = n.Int64()
		}
	}
	r, err := client.NewRenewer(&api.RenewerInput{
		Secret: &api.Secret{
			Auth: &api.SecretAuth{
				ClientToken: client.Token(),
				Renewable:   renewable,
			},
		},
		Increment: int(increment),
	})

	return &renewer{
		client:  client,
		renewer: r,
	}, nil
}

func (r renewer) Run(ctx context.Context) {
	go r.renewer.Renew()
	defer r.renewer.Stop()
	for {
		select {
		case <-ctx.Done():
			return

		case err := <-r.renewer.DoneCh():
			if err != nil {
				logger.WithContext(ctx).Error(err)
				return
			}

			// Renewal is now over
		case renewal := <-r.renewer.RenewCh():
			logger.Infof("Successfully renewed token at %s", renewal.RenewedAt.Format(time.RFC3339Nano))
		}
	}
}
