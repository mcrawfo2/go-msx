// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package vault

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
)

// TokenSource enables retrieval and management of Vault Tokens
type TokenSource interface {
	GetToken(ctx context.Context) (token string, err error)
	Renewable() bool
}

// NewTokenSource will return a TokenSource implementation based provided config
// Currently Config based source and Kubernetes Auth are implemented
func NewTokenSource(source string, cfg *config.Config, conn ConnectionApi) (tokenSource TokenSource, err error) {
	switch source {
	case "approle":
		tokenSource, err = NewAppRoleSource(cfg, conn)
	case "kubernetes":
		tokenSource, err = NewKubernetesSource(cfg, conn)
	default:
		tokenSource = NewConfigSource(cfg)
	}

	return
}
