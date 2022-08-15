// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package transit

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
)

type dummyEncrypter struct {
	keyId types.UUID
}

func (d dummyEncrypter) CreateKey() (err error) {
	return nil
}

func (d dummyEncrypter) Encrypt(value map[string]*string) (secureValue string, encrypted bool, err error) {
	secureValue, err = serializePayload(value)
	if err != nil {
		return "", false, err
	}
	return secureValue, false, err
}

func (d dummyEncrypter) Decrypt(secureValue string) (value map[string]*string, err error) {
	return deserializePayload(secureValue)
}

func NewDummyEncrypter(ctx context.Context, keyId types.UUID) Encrypter {
	return dummyEncrypter{keyId}
}

type dummyBulkEncrypter struct{}

func (d dummyBulkEncrypter) DecryptSets(sets BulkSets) error {
	for _, set := range sets.Sets() {
		if err := d.DecryptSet(set); err != nil {
			return err
		}
	}
	return nil
}

func (d dummyBulkEncrypter) DecryptSet(set BulkSet) error {
	return nil
}

func NewDummyBulkEncrypter(ctx context.Context) BulkEncrypter {
	return dummyBulkEncrypter{}
}
