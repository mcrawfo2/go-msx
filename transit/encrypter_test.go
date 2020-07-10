package transit

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func testEncryptionProvider() Provider {
	mockEncryptionProvider := new(MockProvider)

	mockEncryptionProvider.
		On("CreateKey", context.Background(), `22a342bf-3278-4126-9a02-f1ac0c9cf05f`).
		Return(nil)

	mockEncryptionProvider.
		On("CreateKey", context.Background(), `22a342bf-3278-4126-9a02-f1ac0c9cf05e`).
		Return(errors.New("Error creating key"))

	mockEncryptionProvider.
		On("Decrypt", context.Background(), Value{
			version:   "1",
			keyId:     types.MustParseUUID(`22a342bf-3278-4126-9a02-f1ac0c9cf05f`),
			encrypted: true,
			payload:   "ABCD",
		}).
		Return(Value{
			version:   "1",
			keyId:     types.MustParseUUID(`22a342bf-3278-4126-9a02-f1ac0c9cf05f`),
			encrypted: false,
			payload:   `["java.util.HashMap",{"key1":"value1"}]`,
		}, nil)

	mockEncryptionProvider.
		On("Decrypt", context.Background(), Value{
			version:   "1",
			keyId:     types.MustParseUUID(`22a342bf-3278-4126-9a02-f1ac0c9cf05f`),
			encrypted: true,
			payload:   "MNOP",
		}).
		Return(Value{}, errors.New("Error decrypting value"))

	mockEncryptionProvider.
		On("Encrypt", context.Background(), Value{
			version:   "1",
			keyId:     types.MustParseUUID(`22a342bf-3278-4126-9a02-f1ac0c9cf05f`),
			encrypted: false,
			payload:   `["java.util.HashMap",{"key1":"value1"}]`,
		}).
		Return(Value{
			version:   "1",
			keyId:     types.MustParseUUID(`22a342bf-3278-4126-9a02-f1ac0c9cf05f`),
			encrypted: true,
			payload:   "ABCD",
		}, nil)

	return mockEncryptionProvider
}

func TestSetEncrypterFactory(t *testing.T) {
	type args struct {
		factory EncrypterFactory
	}
	tests := []struct {
		name             string
		encrypterFactory EncrypterFactory
		args             args
		wantNil          bool
	}{
		{
			name:             "NoSet",
			encrypterFactory: nil,
			args: args{
				factory: nil,
			},
			wantNil: true,
		},
		{
			name:             "Set",
			encrypterFactory: nil,
			args: args{
				factory: newEncrypter,
			},
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypterFactory = tt.encrypterFactory
			SetEncrypterFactory(tt.args.factory)
			if tt.wantNil {
				assert.Nil(t, encrypterFactory)
			} else {
				assert.NotNil(t, encrypterFactory)
			}
		})
	}
}

func Test_encrypter_CreateKey(t *testing.T) {
	ctx := context.Background()
	encryptionProvider = testEncryptionProvider()

	type fields struct {
		ctx   context.Context
		keyId types.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				ctx:   ctx,
				keyId: types.MustParseUUID(`22a342bf-3278-4126-9a02-f1ac0c9cf05f`),
			},
			wantErr: false,
		},
		{
			name: "Failure",
			fields: fields{
				ctx:   ctx,
				keyId: types.MustParseUUID(`22a342bf-3278-4126-9a02-f1ac0c9cf05e`),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := encrypter{
				ctx:   tt.fields.ctx,
				keyId: tt.fields.keyId,
			}

			if err := e.CreateKey(); (err != nil) != tt.wantErr {
				t.Errorf("CreateKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_encrypter_Decrypt(t *testing.T) {
	ctx := context.Background()
	encryptionProvider = testEncryptionProvider()

	tests := []struct {
		name    string
		keyId   types.UUID
		value   string
		want    map[string]*string
		wantErr bool
	}{
		{
			name:  "Success",
			keyId: types.MustParseUUID(`22a342bf-3278-4126-9a02-f1ac0c9cf05f`),
			value: `ABCD`,
			want: map[string]*string{
				"key1": types.NewStringPtr("value1"),
			},
			wantErr: false,
		},
		{
			name:    "DecryptError",
			keyId:   types.MustParseUUID(`22a342bf-3278-4126-9a02-f1ac0c9cf05f`),
			value:   `MNOP`,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := encrypter{
				ctx:   ctx,
				keyId: tt.keyId,
			}
			got, err := e.Decrypt(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decrypt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_encrypter_Encrypt(t *testing.T) {
	ctx := context.Background()
	encryptionProvider = testEncryptionProvider()

	tests := []struct {
		name              string
		keyId             types.UUID
		value             map[string]*string
		wantSecurePayload string
		wantEncrypted     bool
		wantErr           bool
	}{
		{
			name:              "Success",
			keyId:             types.MustParseUUID(`22a342bf-3278-4126-9a02-f1ac0c9cf05f`),
			value:             map[string]*string{"key1": types.NewStringPtr("value1")},
			wantSecurePayload: `ABCD`,
			wantEncrypted:     true,
			wantErr:           false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := encrypter{
				ctx:   ctx,
				keyId: tt.keyId,
			}
			gotSecurePayload, gotEncrypted, err := e.Encrypt(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSecurePayload != tt.wantSecurePayload {
				t.Errorf("Encrypt() gotSecurePayload = %v, want %v", gotSecurePayload, tt.wantSecurePayload)
			}
			if gotEncrypted != tt.wantEncrypted {
				t.Errorf("Encrypt() gotEncrypted = %v, want %v", gotEncrypted, tt.wantEncrypted)
			}
		})
	}
}

func Test_newEncrypter(t *testing.T) {
	e := newEncrypter(context.Background(), types.MustParseUUID(`22a342bf-3278-4126-9a02-f1ac0c9cf05f`))
	assert.NotNil(t, e)
}
