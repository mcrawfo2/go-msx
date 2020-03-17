package security

import (
	"context"
	"github.com/pkg/errors"
)

var ErrNotRegistered = errors.New("Token Details Provider not registered")

type TokenDetailsProvider interface {
	TokenDetails(ctx context.Context) (*UserContextDetails, error)
	IsTokenActive(ctx context.Context) (bool, error)
}

var tokenDetailsProvider TokenDetailsProvider

func SetTokenDetailsProvider(provider TokenDetailsProvider) {
	if provider != nil {
		tokenDetailsProvider = provider
	}
}

func NewUserContextDetails(ctx context.Context) (userContextDetails *UserContextDetails, err error) {
	if tokenDetailsProvider == nil {
		return nil, ErrNotRegistered
	}
	return tokenDetailsProvider.TokenDetails(ctx)
}

func IsTokenActive(ctx context.Context) (bool, error) {
	if tokenDetailsProvider == nil {
		return false, ErrNotRegistered
	}
	return tokenDetailsProvider.IsTokenActive(ctx)
}
