// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package security

import "context"

type TokenProvider interface {
	UserContextFromToken(ctx context.Context, token string) (*UserContext, error)
}

var tokenProvider TokenProvider

func SetTokenProvider(provider TokenProvider) {
	if provider != nil {
		tokenProvider = provider
	}
}

func NewUserContextFromToken(ctx context.Context, token string) (userContext *UserContext, err error) {
	return tokenProvider.UserContextFromToken(ctx, token[:])
}
