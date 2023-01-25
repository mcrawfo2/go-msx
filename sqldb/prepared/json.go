// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package prepared

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"database/sql/driver"
	"encoding/json"
	"github.com/pkg/errors"
)

// NullableJson represents a JSON value in the database, or NULL.
// Expected column type: JSON/JSONB/BINARY/BLOB
type NullableJson[I any] struct {
	valid bool
	value I
}

func (a NullableJson[I]) Value() (driver.Value, error) {
	if !a.valid {
		return nil, nil
	}

	return json.Marshal(a.value)
}

func (a *NullableJson[I]) Scan(value interface{}) error {
	switch vt := value.(type) {
	case []byte:
		return json.Unmarshal(vt, &a.value)
	case nil:
		var empty I
		a.valid = false
		a.value = empty
		return nil
	default:
		return errors.Errorf("Cannot convert %T to Json", value)
	}
}

func (a NullableJson[I]) Optional() types.Optional[I] {
	if !a.valid {
		return types.Optional[I]{}
	}
	return types.OptionalOf[I](a.value)
}

func (a NullableJson[I]) Null() bool {
	return !a.valid
}

func (a NullableJson[I]) Unwrap() I {
	return a.value
}

func NewNullableJsonFromOptional[I any](opt types.Optional[I]) NullableJson[I] {
	if !opt.IsPresent() {
		return NullableJson[I]{}
	}
	return NullableJson[I]{
		valid: true,
		value: opt.Value(),
	}
}

func NewNullableJson[I any](ptr *I) NullableJson[I] {
	if ptr == nil {
		return NullableJson[I]{}
	}
	return NullableJson[I]{
		valid: true,
		value: *ptr,
	}
}

// Json represents a JSON value in the database.
// Expected column type: JSON/JSONB/BINARY/BLOB
type Json[I any] struct {
	value I
}

func (a Json[I]) Value() (driver.Value, error) {
	return json.Marshal(a.value)
}

func (a *Json[I]) Scan(value interface{}) error {
	switch vt := value.(type) {
	case []byte:
		return json.Unmarshal(vt, &a.value)
	case string:
		bv := []byte(vt)
		return json.Unmarshal(bv, &a.value)
	default:
		return errors.Errorf("Cannot convert %T to Json", value)
	}
}

func (a Json[I]) Nullable() NullableJson[I] {
	return NullableJson[I]{
		valid: true,
		value: a.value,
	}
}

func (a Json[I]) Unwrap() I {
	return a.value
}

func NewJson[I any](value I) Json[I] {
	return Json[I]{
		value: value,
	}
}
