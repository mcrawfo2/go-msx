package tokensource

import (
	"encoding/json"
	"github.com/hashicorp/vault/api"
)

func initRenewer(client *api.Client) (*api.Renewer, error) {
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
	return r, nil
}

func startRenewer(r *api.Renewer) {
	go r.Renew()
	defer r.Stop()
	for {
		select {
		case err := <-r.DoneCh():
			if err != nil {
				logger.Error(err)
				return
			}

			// Renewal is now over
		case renewal := <-r.RenewCh():

			logger.Infof("Successfully renewed token at %s",
				renewal.RenewedAt.Format("2006-01-02T15:04:05.999999-07:00"))
		}
	}
}
