// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package stream

import (
	"context"
)

type contextKey int

const (
	contextKeyService contextKey = iota
)

func PublisherServiceFromContext(ctx context.Context) PublisherService {
	value, _ := ctx.Value(contextKeyService).(PublisherService)
	return value
}

func ContextWithPublisherService(ctx context.Context, service PublisherService) context.Context {
	return context.WithValue(ctx, contextKeyService, service)
}
