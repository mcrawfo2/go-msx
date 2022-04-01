// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package transit

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test_NewDummyEncrypter(t *testing.T) {
	e := NewDummyEncrypter(context.Background(), types.MustParseUUID(`22a342bf-3278-4126-9a02-f1ac0c9cf05f`))
	assert.NotNil(t, e)
}

func Test_dummyEncrypter_CreateKey(t *testing.T) {
	d := dummyEncrypter{
		keyId: types.MustParseUUID(`22a342bf-3278-4126-9a02-f1ac0c9cf05f`),
	}
	if err := d.CreateKey(); err != nil {
		t.Errorf("CreateKey() error = %v, wantErr %v", err, false)
	}
}

func Test_dummyEncrypter_Encrypt(t *testing.T) {
	keyId := types.MustParseUUID(`22a342bf-3278-4126-9a02-f1ac0c9cf05f`)
	tests := []struct {
		name            string
		value           map[string]*string
		wantSecureValue string
		wantEncrypted   bool
		wantErr         bool
	}{
		{
			name:            "Success",
			value:           map[string]*string{"key1": types.NewStringPtr("value1")},
			wantSecureValue: `["java.util.HashMap",{"key1":"value1"}]`,
			wantEncrypted:   false,
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := dummyEncrypter{
				keyId: keyId,
			}
			gotSecureValue, gotEncrypted, err := d.Encrypt(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSecureValue != tt.wantSecureValue {
				t.Errorf("Encrypt() gotSecureValue = %v, want %v", gotSecureValue, tt.wantSecureValue)
			}
			if gotEncrypted != tt.wantEncrypted {
				t.Errorf("Encrypt() gotEncrypted = %v, want %v", gotEncrypted, tt.wantEncrypted)
			}
		})
	}
}

func Test_dummyEncrypter_Decrypt(t *testing.T) {
	keyId := types.MustParseUUID(`22a342bf-3278-4126-9a02-f1ac0c9cf05f`)
	tests := []struct {
		name        string
		secureValue string
		wantPayload map[string]*string
		wantErr         bool
	}{
		{
			name:          "Success",
			secureValue:   `["java.util.HashMap",{"key1":"value1"}]`,
			wantPayload:   map[string]*string{"key1": types.NewStringPtr("value1")},
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := dummyEncrypter{
				keyId: keyId,
			}
			gotValue, err := d.Decrypt(tt.secureValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotValue, tt.wantPayload) {
				t.Errorf("Decrypt() gotValue = %v, want %v", gotValue, tt.wantPayload)
			}
		})
	}
}
