// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package scheduled

import "context"

type contextKey int

const (
	contextKeyService contextKey = iota
)

func SchedulerServiceFromContext(ctx context.Context) SchedulerServiceApi {
	service, _ := ctx.Value(contextKeyService).(SchedulerServiceApi)
	return service
}

func ContextWithSchedulerService(ctx context.Context, service SchedulerServiceApi) context.Context {
	return context.WithValue(ctx, contextKeyService, service)
}
