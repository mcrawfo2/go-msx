package transit

//go:generate mockery --inpackage --name=Encrypter --structname=MockEncrypter --filename mock_encrypter_test.go

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"strings"
)

var (
	logger = log.NewLogger("msx.transit")
)

type Encrypter interface {
	CreateKey() (err error)
	Encrypt(value map[string]*string) (secureValue string, encrypted bool, err error)
	Decrypt(secureValue string) (value map[string]*string, err error)
}

type encrypter struct {
	ctx   context.Context
	cfg   *Config
	keyId types.UUID
}

func (e encrypter) keyName() string {
	return strings.ToLower(e.keyId.String())
}

func (e encrypter) CreateKey() (err error) {
	logger.WithContext(e.ctx).Debugf("Creating transit encryption key %q", e.keyId)
	p, err := provider()
	if err != nil {
		return err
	}
	return p.CreateKey(e.ctx, e.keyName())
}

func (e encrypter) Encrypt(value map[string]*string) (securePayload string, encrypted bool, err error) {
	logger.WithContext(e.ctx).Debugf("Encrypting using transit encryption key %q", e.keyId)

	p, err := provider()
	if err != nil {
		return "", false, err
	}

	insecureValue, err := NewValue(e.keyId, value)
	if err != nil {
		return "", false, err
	}

	secureValue, err := p.Encrypt(e.ctx, insecureValue)
	if err != nil {
		return "", false, err
	}

	return secureValue.payload, secureValue.encrypted, nil
}

func (e encrypter) Decrypt(value string) (map[string]*string, error) {
	logger.WithContext(e.ctx).Debugf("Decrypting using transit encryption key %q", e.keyId)

	p, err := provider()
	if err != nil {
		return nil, err
	}

	insecureValue, err := p.Decrypt(e.ctx, NewSecureValue(e.keyId, value))
	if err != nil {
		return nil, err
	}

	payload, err := insecureValue.Payload()
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func NewProductionEncrypter(ctx context.Context, keyName types.UUID) Encrypter {
	return &encrypter{
		ctx:   ctx,
		keyId: keyName,
	}
}

type EncrypterFactory func(ctx context.Context, keyName types.UUID) Encrypter

func (f EncrypterFactory) Create(ctx context.Context, keyName types.UUID) Encrypter {
	return f(ctx, keyName)
}

var encrypterFactory EncrypterFactory = NewDummyEncrypter

func SetEncrypterFactory(factory EncrypterFactory) {
	if factory != nil {
		encrypterFactory = factory
	}
}

func NewEncrypter(ctx context.Context, keyName types.UUID) Encrypter {
	return encrypterFactory.Create(ctx, keyName)
}
