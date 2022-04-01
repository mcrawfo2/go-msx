// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package security

import "context"

type securityContextKey int

const (
	contextKeyUserContext securityContextKey = iota
	contextKeyTokenDetailsProvider
)

func ContextWithUserContext(ctx context.Context, userContext *UserContext) context.Context {
	return context.WithValue(ctx, contextKeyUserContext, userContext)
}

func UserContextFromContext(ctx context.Context) *UserContext {
	userContextInterface := ctx.Value(contextKeyUserContext)
	if userContextInterface == nil {
		return defaultUserContext.Clone()
	}
	return userContextInterface.(*UserContext)
}

func UserNameFromContext(ctx context.Context) string {
	return UserContextFromContext(ctx).UserName
}

func ContextWithTokenDetailsProvider(ctx context.Context, provider TokenDetailsProvider) context.Context {
	return context.WithValue(ctx, contextKeyTokenDetailsProvider, provider)
}

func TokenDetailsProviderFromContext(ctx context.Context) TokenDetailsProvider {
	tokenDetailsProviderInterface := ctx.Value(contextKeyTokenDetailsProvider)
	if tokenDetailsProviderInterface == nil {
		return nil
	}

	return tokenDetailsProviderInterface.(TokenDetailsProvider)
}
