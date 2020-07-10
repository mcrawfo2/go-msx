package transit

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func testModelEncrypterFactory() EncrypterFactory {
	return func(ctx context.Context, keyName types.UUID) Encrypter {
		encrypter := new(MockEncrypter)

		encrypter.
			On("Encrypt", map[string]*string{
				"key1": types.NewStringPtr("value1"),
			}).
			Return(`ABCD`, true, nil)

		encrypter.
			On("Encrypt", map[string]*string{
				"key1": types.NewStringPtr("value2"),
			}).
			Return(`EFGH`, true, nil)

		encrypter.
			On("Encrypt", map[string]*string{
				"key1": types.NewStringPtr("value1"),
				"key2": types.NewStringPtr("value2"),
			}).
			Return(`IJKL`, true, nil)

		encrypter.
			On("Encrypt", map[string]*string{
				"key2": types.NewStringPtr("value1"),
			}).
			Return(`1:22a342bf-3278-4126-9a02-f1ac0c9cf05f:p:["java.util.HashMap",{"key1":"value1"}]`, false, nil)

		encrypter.
			On("Encrypt", map[string]*string{}).
			Return(``, false, errors.New("encryption error"))

		encrypter.
			On("Encrypt", map[string]*string{
				"key3": types.NewStringPtr("value3"),
			}).
			Return(``, false, errors.New("encryption error"))

		encrypter.
			On("Decrypt", `ABCD`).
			Return(map[string]*string{
				"key1": types.NewStringPtr("value1"),
			}, nil)

		encrypter.
			On("Decrypt", `EFGH`).
			Return(map[string]*string{
				"key1": types.NewStringPtr("value2"),
			}, nil)

		encrypter.
			On("Decrypt", `MNOP`).
			Return(nil, errors.New("decryption error"))

		return encrypter
	}
}

func TestSecureData_UnmarshalCQL(t *testing.T) {
	ctx := context.Background()
	keyId := types.MustParseUUID(`22a342bf-3278-4126-9a02-f1ac0c9cf05f`)
	tests := []struct {
		name    string
		value   SecureData
		data    []byte
		want    SecureData
		wantErr bool
	}{
		{
			name:  "Empty",
			value: SecureData{ctx: ctx, keyId: types.EmptyUUID()},
			data:  []byte(`1:22a342bf-3278-4126-9a02-f1ac0c9cf05f::`),
			want: SecureData{
				ctx:   ctx,
				dirty: false,
				keyId: keyId,
				secure: Value{
					version:   "1",
					keyId:     keyId,
					encrypted: false,
					payload:   "",
				},
			},
			wantErr: false,
		},
		{
			name:  "Plain",
			value: SecureData{ctx: ctx, keyId: types.EmptyUUID()},
			data:  []byte(`1:22a342bf-3278-4126-9a02-f1ac0c9cf05f:p:["java.util.HashMap",{"key1":"value1"}]`),
			want: SecureData{
				ctx:   ctx,
				dirty: false,
				keyId: keyId,
				secure: Value{
					version:   "1",
					keyId:     keyId,
					encrypted: false,
					payload:   `["java.util.HashMap",{"key1":"value1"}]`,
				},
			},
			wantErr: false,
		},
		{
			name:  "Encrypted",
			value: SecureData{ctx: ctx, keyId: types.EmptyUUID()},
			data:  []byte(`1:22a342bf-3278-4126-9a02-f1ac0c9cf05f:e:ABCD`),
			want: SecureData{
				ctx:   ctx,
				dirty: false,
				keyId: keyId,
				secure: Value{
					version:   "1",
					keyId:     keyId,
					encrypted: true,
					payload:   `ABCD`,
				},
			},
			wantErr: false,
		},
		{
			name:    "InvalidParts",
			value:   SecureData{ctx: ctx},
			data:    []byte(`1:22a342bf-3278-4126-9a02-f1ac0c9cf05f:`),
			wantErr: true,
		},
		{
			name:    "InvalidEncrypted",
			value:   SecureData{ctx: ctx},
			data:    []byte(`1:22a342bf-3278-4126-9a02-f1ac0c9cf05f:s:`),
			wantErr: true,
		},
		{
			name:    "InvalidVersion",
			value:   SecureData{ctx: ctx},
			data:    []byte(`2:22a342bf-3278-4126-9a02-f1ac0c9cf05f::`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s = tt.value
			var sp = &s
			if err := sp.UnmarshalCQL(nil, tt.data); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalCQL() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.want, s) {
				t.Errorf("UnmarshalCQL() got = %+v want %+v", s, tt.want)
			}
		})
	}
}

func TestSecureData_MarshalCQL(t *testing.T) {
	ctx := context.Background()
	keyId := types.MustParseUUID(`22a342bf-3278-4126-9a02-f1ac0c9cf05f`)
	encrypterFactory = testModelEncrypterFactory()

	tests := []struct {
		name    string
		value   SecureData
		want    string
		wantErr bool
	}{
		{
			name: "Plain",
			value: SecureData{
				ctx:   ctx,
				dirty: false,
				keyId: keyId,
				secure: Value{
					version:   "1",
					keyId:     keyId,
					encrypted: false,
					payload:   `["java.util.HashMap",{"key1":"value1"}]`,
				},
				payload: nil,
			},
			want:    `1:22a342bf-3278-4126-9a02-f1ac0c9cf05f:p:["java.util.HashMap",{"key1":"value1"}]`,
			wantErr: false,
		},
		{
			name: "Encrypted",
			value: SecureData{
				ctx:   ctx,
				dirty: false,
				keyId: keyId,
				secure: Value{
					version:   "1",
					keyId:     keyId,
					encrypted: true,
					payload:   `ABCD`,
				},
				payload: nil,
			},
			want: `1:22a342bf-3278-4126-9a02-f1ac0c9cf05f:e:ABCD`,
		},
		{
			name: "DirtyPlain",
			value: SecureData{
				ctx:   ctx,
				dirty: true,
				keyId: keyId,
				secure: Value{
					version:   "1",
					keyId:     keyId,
					encrypted: false,
					payload:   `["java.util.HashMap",{"key1":"value1"}]`,
				},
				payload: map[string]*string{
					"key1": types.NewStringPtr("value2"),
				},
			},
			want: `1:22a342bf-3278-4126-9a02-f1ac0c9cf05f:e:EFGH`,
		},
		{
			name: "DirtyEncrypted",
			value: SecureData{
				ctx:   ctx,
				dirty: true,
				keyId: keyId,
				secure: Value{
					version:   "1",
					keyId:     keyId,
					encrypted: true,
					payload:   `ABCD`,
				},
				payload: map[string]*string{
					"key1": types.NewStringPtr("value2"),
				},
			},
			want: `1:22a342bf-3278-4126-9a02-f1ac0c9cf05f:e:EFGH`,
		},
		{
			name: "DirtyExtendPlain",
			value: SecureData{
				ctx:   ctx,
				dirty: true,
				keyId: keyId,
				secure: Value{
					payload: `["java.util.HashMap",{"key1":"value1"}]`,
				},
				payload: map[string]*string{
					"key1": types.NewStringPtr("value2"),
				},
			},
			want: `1:22a342bf-3278-4126-9a02-f1ac0c9cf05f:e:EFGH`,
		},
		{
			name: "DirtyExtendEncrypted",
			value: SecureData{
				ctx:   ctx,
				dirty: true,
				keyId: keyId,
				secure: Value{
					version:   "1",
					keyId:     keyId,
					encrypted: true,
					payload:   `ABCD`,
				},
				payload: map[string]*string{
					"key1": types.NewStringPtr("value1"),
					"key2": types.NewStringPtr("value2"),
				},
			},
			want: `1:22a342bf-3278-4126-9a02-f1ac0c9cf05f:e:IJKL`,
		},
		{
			name: "EncryptionDisabled",
			value: SecureData{
				ctx:    ctx,
				dirty:  true,
				keyId:  keyId,
				secure: Value{},
				payload: map[string]*string{
					"key2": types.NewStringPtr("value1"),
				},
			},
			want: `1:22a342bf-3278-4126-9a02-f1ac0c9cf05f:p:["java.util.HashMap",{"key2":"value1"}]`,
		},
		{
			name:    "KeyNotSet",
			value:   SecureData{},
			wantErr: true,
		},
		{
			name: "EncryptError",
			value: SecureData{
				ctx:     ctx,
				dirty:   true,
				keyId:   keyId,
				secure:  Value{},
				payload: map[string]*string{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s = tt.value
			var sp = &s

			got, err := sp.MarshalCQL(nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalCQL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.want != string(got) {
				t.Errorf("MarshalCQL() got = %+v, want %+v", string(got), tt.want)
			}
		})
	}
}

func TestSecureData_Field(t *testing.T) {
	ctx := context.Background()
	keyId := types.MustParseUUID(`22a342bf-3278-4126-9a02-f1ac0c9cf05f`)
	encrypterFactory = testModelEncrypterFactory()

	tests := []struct {
		name      string
		value     SecureData
		fieldName string
		want      *string
		wantErr   bool
	}{
		{
			name: "Plain",
			value: SecureData{
				ctx:   ctx,
				dirty: false,
				keyId: keyId,
				secure: Value{
					version:   "1",
					keyId:     keyId,
					encrypted: false,
					payload:   `["java.util.HashMap",{"key1":"value1"}]`,
				},
				payload: nil,
			},
			fieldName: "key1",
			want:      types.NewStringPtr("value1"),
			wantErr:   false,
		},
		{
			name: "Encrypted",
			value: SecureData{
				ctx:   ctx,
				dirty: false,
				keyId: keyId,
				secure: Value{
					version:   "1",
					keyId:     keyId,
					encrypted: true,
					payload:   `ABCD`,
				},
				payload: nil,
			},
			fieldName: "key1",
			want:      types.NewStringPtr("value1"),
			wantErr:   false,
		},
		{
			name: "Missing",
			value: SecureData{
				ctx:   ctx,
				dirty: false,
				keyId: keyId,
				secure: Value{
					version:   "1",
					keyId:     keyId,
					encrypted: true,
					payload:   `ABCD`,
				},
				payload: nil,
			},
			fieldName: "key2",
			want:      nil,
			wantErr:   false,
		},
		{
			name: "InvalidEncrypted",
			value: SecureData{
				ctx:   ctx,
				dirty: false,
				keyId: keyId,
				secure: Value{
					version:   "1",
					keyId:     keyId,
					encrypted: true,
					payload:   `MNOP`,
				},
				payload: nil,
			},
			fieldName: "key1",
			want:      nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.value
			sp := &s
			got, err := sp.Field(ctx, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Field() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Field() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSecureData_SetField(t *testing.T) {
	ctx := context.Background()
	keyId := types.MustParseUUID(`22a342bf-3278-4126-9a02-f1ac0c9cf05f`)
	encrypterFactory = testModelEncrypterFactory()

	tests := []struct {
		name       string
		value      SecureData
		fieldName  string
		fieldValue *string
		want       SecureData
		wantErr    bool
	}{
		{
			name:       "Empty",
			value:      SecureData{ctx: ctx, keyId: keyId},
			fieldName:  "key1",
			fieldValue: types.NewStringPtr("value1"),
			want: SecureData{
				ctx:   ctx,
				dirty: true,
				keyId: keyId,
				payload: map[string]*string{
					"key1": types.NewStringPtr("value1"),
				},
			},
		},
		{
			name: "Plain",
			value: SecureData{
				ctx:   ctx,
				keyId: keyId,
				secure: Value{
					version:   "1",
					keyId:     keyId,
					encrypted: false,
					payload:   `["java.util.HashMap",{"key1":"value1"}]`,
				},
			},
			fieldName:  "key1",
			fieldValue: types.NewStringPtr("value2"),
			want: SecureData{
				ctx:   ctx,
				dirty: true,
				keyId: keyId,
				secure: Value{
					version:   "1",
					keyId:     keyId,
					encrypted: false,
					payload:   `["java.util.HashMap",{"key1":"value1"}]`,
				},
				payload: map[string]*string{
					"key1": types.NewStringPtr("value2"),
				},
			},
		},
		{
			name: "NoChange",
			value: SecureData{
				ctx:   ctx,
				keyId: keyId,
				secure: Value{
					version:   "1",
					keyId:     keyId,
					encrypted: false,
					payload:   `["java.util.HashMap",{"key1":"value1"}]`,
				},
			},
			fieldName:  "key1",
			fieldValue: types.NewStringPtr("value1"),
			want: SecureData{
				ctx:   ctx,
				dirty: false,
				keyId: keyId,
				secure: Value{
					version:   "1",
					keyId:     keyId,
					encrypted: false,
					payload:   `["java.util.HashMap",{"key1":"value1"}]`,
				},
				payload: map[string]*string{
					"key1": types.NewStringPtr("value1"),
				},
			},
		},
		{
			name: "Encrypted",
			value: SecureData{
				ctx:   ctx,
				keyId: keyId,
				secure: Value{
					version:   "1",
					keyId:     keyId,
					encrypted: true,
					payload:   `ABCD`,
				},
			},
			fieldName:  "key1",
			fieldValue: types.NewStringPtr("value2"),
			want: SecureData{
				ctx:   ctx,
				dirty: true,
				keyId: keyId,
				secure: Value{
					version:   "1",
					keyId:     keyId,
					encrypted: false,
					payload:   `["java.util.HashMap",{"key1":"value1"}]`,
				},
				payload: map[string]*string{
					"key1": types.NewStringPtr("value2"),
				},
			},
		},
		{
			name: "DecryptError",
			value: SecureData{
				ctx:   ctx,
				dirty: false,
				keyId: keyId,
				secure: Value{
					version:   "1",
					keyId:     keyId,
					encrypted: true,
					payload:   `MNOP`,
				},
			},
			fieldName:  "key3",
			fieldValue: types.NewStringPtr("value3"),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s = tt.value
			var sp = &s
			var err error

			sp, err = sp.SetField(ctx, tt.fieldName, tt.fieldValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetField() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.want, s) {
				t.Errorf("SetField() got = %+v want %+v", s, tt.want)
			}
		})
	}
}

func TestSecureData_SetKeyId(t *testing.T) {
	ctx := context.Background()
	keyId := types.MustParseUUID(`22a342bf-3278-4126-9a02-f1ac0c9cf05f`)
	value := SecureData{}

	value.SetKeyId(ctx, keyId)
	assert.Equal(t, value.ctx, ctx)
	assert.Equal(t, value.keyId, keyId)
}

func TestSecureData_KeyId(t *testing.T) {
	ctx := context.Background()
	keyId := types.MustParseUUID(`22a342bf-3278-4126-9a02-f1ac0c9cf05f`)
	value := SecureData{
		ctx:   ctx,
		keyId: keyId,
	}

	assert.Equal(t, keyId, value.KeyId())
}

func TestWithSecureData_SetSecureValue(t *testing.T) {
	ctx := context.Background()
	keyId := types.MustParseUUID(`22a342bf-3278-4126-9a02-f1ac0c9cf05f`)

	tests := []struct {
		name       string
		value      WithSecureData
		keyId      types.UUID
		fieldName  string
		fieldValue *string
		want       WithSecureData
		wantErr    bool
	}{
		{
			name:    "EmptyKey",
			keyId:   types.EmptyUUID(),
			wantErr: true,
		},
		{
			name:    "InvalidKey",
			keyId:   types.UUID([]byte{1, 2, 3}),
			wantErr: true,
		},
		{
			name:       "CreateSecureData",
			keyId:      keyId,
			fieldName:  "key1",
			fieldValue: types.NewStringPtr("value1"),
			want: WithSecureData{
				SecureData: &SecureData{
					ctx:    ctx,
					dirty:  true,
					keyId:  keyId,
					secure: Value{},
					payload: map[string]*string{
						"key1": types.NewStringPtr("value1"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ws = tt.value
			var wsp = &ws
			var err error

			err = wsp.SetSecureValue(ctx, tt.keyId, tt.fieldName, tt.fieldValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetSecureValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.want, ws) {
				t.Errorf("SetSecureValue() got = %+v want %+v", ws, tt.want)
			}
		})
	}
}

func TestWithSecureData_SecureValue(t *testing.T) {
	ctx := context.Background()
	keyId := types.MustParseUUID(`22a342bf-3278-4126-9a02-f1ac0c9cf05f`)

	tests := []struct {
		name      string
		value     WithSecureData
		keyId     types.UUID
		fieldName string
		want      string
		wantErr   bool
	}{
		{
			name:      "ValueExists",
			keyId:     keyId,
			fieldName: "key1",
			value: WithSecureData{
				SecureData: &SecureData{
					ctx:    ctx,
					dirty:  true,
					keyId:  keyId,
					secure: Value{},
					payload: map[string]*string{
						"key1": types.NewStringPtr("value1"),
					},
				},
			},
			want: "value1",
		},
		{
			name:      "ValueNotExists",
			keyId:     keyId,
			fieldName: "key1",
			value: WithSecureData{
				SecureData: &SecureData{
					ctx:     ctx,
					dirty:   true,
					keyId:   keyId,
					secure:  Value{},
					payload: map[string]*string{},
				},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ws = tt.value
			var wsp = &ws
			var err error
			var got string

			got, err = wsp.SecureValue(ctx, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("SecureValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("SecureValue() got = %+v want %+v", got, tt.want)
			}
		})
	}
}

func TestWithSecureData_SecureOptionalValue(t *testing.T) {
	ctx := context.Background()
	keyId := types.MustParseUUID(`22a342bf-3278-4126-9a02-f1ac0c9cf05f`)

	tests := []struct {
		name      string
		value     WithSecureData
		keyId     types.UUID
		fieldName string
		want      types.OptionalString
		wantErr   bool
	}{
		{
			name:      "ValueExists",
			keyId:     keyId,
			fieldName: "key1",
			value: WithSecureData{
				SecureData: &SecureData{
					ctx:    ctx,
					dirty:  true,
					keyId:  keyId,
					secure: Value{},
					payload: map[string]*string{
						"key1": types.NewStringPtr("value1"),
					},
				},
			},
			want: types.NewOptionalStringFromString("value1"),
		},
		{
			name:      "ValueNotExists",
			keyId:     keyId,
			fieldName: "key1",
			value: WithSecureData{
				SecureData: &SecureData{
					ctx:     ctx,
					dirty:   true,
					keyId:   keyId,
					secure:  Value{},
					payload: map[string]*string{},
				},
			},
			want: types.NewOptionalString(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ws = tt.value
			var wsp = &ws
			var err error
			var got types.OptionalString

			got, err = wsp.SecureOptionalValue(ctx, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("SecureOptionalValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SecureOptionalValue() got = %+v want %+v", got, tt.want)
			}
		})
	}
}
