// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package transit

import (
	"context"
	"github.com/pkg/errors"
)

type BulkSet []*SecureData

func (s BulkSet) IsEmpty() bool {
	return len(s) == 0
}

func (s BulkSet) Valid() error {
	if s.IsEmpty() {
		return nil
	}
	keyId := s[0].KeyId()
	for i := 1; i < len(s); i++ {
		if !keyId.Equals(s[i].KeyId()) {
			return errors.Errorf("Mismatched transit key: %q", s[i].KeyId().String())
		}
	}
	return nil
}

type BulkSets interface {
	Sets() []BulkSet
}

type BulkEncrypter interface {
	DecryptSets(sets BulkSets) error
	DecryptSet(set BulkSet) error
}

type bulkEncrypter struct {
	ctx context.Context
	cfg *Config
}


func (e bulkEncrypter) DecryptSets(sets BulkSets) error {
	for _, set := range sets.Sets() {
		if err := e.DecryptSet(set); err != nil {
			return err
		}
	}
	return nil
}

func (e bulkEncrypter) DecryptSet(set BulkSet) error {
	logger.WithContext(e.ctx).Debugf("Decrypting bulk set")

	if err := set.Valid(); err != nil {
		return err
	}

	if set.IsEmpty() {
		return nil
	}

	p, err := provider()
	if err != nil {
		return err
	}

	var values []Value
	for _, entry := range set {
		values = append(values, entry.secure)
	}

	insecureValues, err := p.DecryptBulk(e.ctx, values)
	if err != nil {
		return err
	}

	// Update target in-place
	for i, insecureValue := range insecureValues {
		if err = set[i].withDecryptedValue(insecureValue); err != nil {
			return err
		}
	}

	return nil
}

type BulkEncrypterFactory func(ctx context.Context) BulkEncrypter

func (f BulkEncrypterFactory) Create(ctx context.Context) BulkEncrypter {
	return f(ctx)
}

var bulkEncrypterFactory BulkEncrypterFactory = NewDummyBulkEncrypter

func SetBulkEncrypterFactory(factory BulkEncrypterFactory) {
	if factory != nil {
		bulkEncrypterFactory = factory
	}
}

func NewProductionBulkEncrypter(ctx context.Context) BulkEncrypter {
	return &bulkEncrypter{
		ctx:   ctx,
	}
}

func NewBulkEncrypter(ctx context.Context) BulkEncrypter {
	return bulkEncrypterFactory.Create(ctx)
}
