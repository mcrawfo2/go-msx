// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package monitor

import "context"

type contextKey int

const (
	contextKeyIntegration contextKey = iota
)

func IntegrationFromContext(ctx context.Context) Api {
	value, _ := ctx.Value(contextKeyIntegration).(Api)
	return value
}

func ContextWithIntegration(ctx context.Context, api Api) context.Context {
	return context.WithValue(ctx, contextKeyIntegration, api)
}
