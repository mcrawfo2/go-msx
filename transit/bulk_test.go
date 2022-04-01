// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package transit

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"reflect"
	"testing"
)

func Test_encrypter_DecryptSet(t *testing.T) {
	ctx := context.Background()
	encryptionProvider = testEncryptionProvider()
	keyId1 := types.MustParseUUID(`22a342bf-3278-4126-9a02-f1ac0c9cf05f`)
	keyId2 := types.MustParseUUID(`22a342bf-3278-4dd6-9a02-f1ac0c9cf05f`)

	tests := []struct {
		name    string
		set     BulkSet
		want    BulkSet
		wantErr bool
	}{
		{
			name: "Success",
			set: BulkSet{
				{
					ctx:   nil,
					dirty: false,
					keyId: keyId1,
					secure: Value{
						version:   valueVersion1,
						keyId:     keyId1,
						encrypted: true,
						payload:   "ABCD",
					},
					payload: nil,
				},
			},
			want: BulkSet{
				{
					ctx:   nil,
					dirty: false,
					keyId: keyId1,
					secure: Value{
						version:   valueVersion1,
						keyId:     keyId1,
						encrypted: false,
						payload:   `["java.util.HashMap",{"key1":"value1"}]`,
					},
					payload: map[string]*string{
						"key1": types.NewStringPtr("value1"),
					},
				},
			},
			wantErr: false,
		}, {
			name:    "EmptySuccess",
			set:     BulkSet{},
			want:    BulkSet{},
			wantErr: false,
		},
		{
			name: "KeyMismatchFailure",
			set: BulkSet{
				{
					ctx:   nil,
					dirty: false,
					keyId: keyId1,
					secure: Value{
						version:   valueVersion1,
						keyId:     keyId1,
						encrypted: true,
						payload:   "ABCD",
					},
					payload: nil,
				},
				{
					ctx:   nil,
					dirty: false,
					keyId: keyId2,
					secure: Value{
						version:   valueVersion1,
						keyId:     keyId2,
						encrypted: true,
						payload:   "ABCD",
					},
					payload: nil,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ProviderFailure",
			set: BulkSet{
				{
					ctx:   nil,
					dirty: false,
					keyId: keyId1,
					secure: Value{
						version:   valueVersion1,
						keyId:     keyId1,
						encrypted: true,
						payload:   "MNOP",
					},
					payload: nil,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := bulkEncrypter{
				ctx:   ctx,
			}
			err := e.DecryptSet(tt.set)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(tt.set, tt.want) {
				t.Errorf("Unexpected result\n%s", testhelpers.Diff(tt.want, tt.set))
			}
		})
	}
}

