package transit

//go:generate mockery --inpackage --name=Provider --structname=MockProvider --filename mock_provider_test.go

import (
	"context"
	"github.com/pkg/errors"
)

var encryptionProvider Provider
var ErrNotRegistered = errors.New("Transit encryption provider not registered")

type Provider interface {
	CreateKey(ctx context.Context, keyName string) (err error)
	Encrypt(ctx context.Context, value Value) (secureValue Value, err error)
	Decrypt(ctx context.Context, secureValue Value) (value Value, err error)
	DecryptBulk(ctx context.Context, secureValues []Value) (values []Value, err error)
}

func provider() (Provider, error) {
	if encryptionProvider == nil {
		return nil, ErrNotRegistered
	}
	return encryptionProvider, nil
}

func RegisterProvider(p Provider) error {
	if p == nil {
		return ErrNotRegistered
	}
	encryptionProvider = p
	return nil
}
