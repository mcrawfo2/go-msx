// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package types

import (
	"encoding/json"
	"github.com/pkg/errors"
)

var ErrNoSuchDetailsKey = errors.New("Object key does not exist")
var ErrValueWrongType = errors.New("Object key does not hold a value of the requested type")

type Pojo map[string]interface{}

func (p Pojo) StringValue(key string) (string, error) {
	v, ok := p[key]
	if !ok {
		return "", errors.Wrapf(ErrNoSuchDetailsKey, "Key %q not found", key)
	}

	var vs string
	switch vt := v.(type) {
	case string:
		vs = vt
	case json.Number:
		vs = string(vt)
	default:
		return "", errors.Wrapf(ErrValueWrongType, "Key %q does not contain a string: %T found", key, v)
	}

	return vs, nil
}

func (p Pojo) PojoValue(key string) (Pojo, error) {
	v, ok := p[key]
	if !ok {
		return nil, errors.Wrapf(ErrNoSuchDetailsKey, "Key %q not found", key)
	}

	vs, ok := v.(map[string]interface{})
	if !ok {
		return nil, errors.Wrapf(ErrValueWrongType, "Key %q does not contain an object: %T found", key, v)
	}

	return vs, nil
}

func (p Pojo) ObjectValue(key string) (Pojo, error) {
	return p.PojoValue(key)
}

func (p Pojo) BoolValue(key string) (bool, error) {

	v, ok := p[key]
	if !ok {
		return false, errors.Wrapf(ErrNoSuchDetailsKey, "Key %q not found", key)
	}

	vs, ok := v.(bool)
	if !ok {
		return false, errors.Wrapf(ErrValueWrongType, "Key %q does not contain a bool: %T found", key, v)
	}

	return vs, nil
}

func (p Pojo) FloatValue(key string) (float64, error) {
	v, ok := p[key]
	if !ok {
		return 0, errors.Wrapf(ErrNoSuchDetailsKey, "Key %q not found", key)
	}

	var vs float64
	switch vt := v.(type) {
	case float64:
		vs = vt
	case json.Number:
		var err error
		vs, err = vt.Float64()
		if err != nil {
			return 0, err
		}
	default:
		return 0, errors.Wrapf(ErrValueWrongType, "Key %q does not contain an float64: %T found", key, v)
	}

	return vs, nil
}

func (p Pojo) ArrayValue(key string) (Poja, error) {
	v, ok := p[key]
	if !ok {
		return nil, errors.Wrapf(ErrNoSuchDetailsKey, "Key %q not found", key)
	}

	vs, ok := v.([]interface{})
	if !ok {
		return nil, errors.Wrapf(ErrValueWrongType, "Key %q does not contain an array: %T found", key, v)
	}

	return vs, nil
}

func (p Pojo) AddAll(o map[string]interface{}) {
	for k, v := range o {
		p[k] = v
	}
}

func (p Pojo) Clone() Pojo {
	var result = make(Pojo)
	result.AddAll(p)
	return result
}

// Deprecated
type PojoArray []map[string]interface{}

func (a PojoArray) Index(i int) Pojo {
	return a[i]
}

type Poja []interface{}

func (p Poja) StringValue(index int) (string, error) {
	if index >= len(p) {
		return "", errors.Wrapf(ErrNoSuchDetailsKey, "Index %d out of bounds", index)
	}

	var vs string
	switch vt := p[index].(type) {
	case string:
		vs = vt
	case json.Number:
		vs = string(vt)
	default:
		return "", errors.Wrapf(ErrValueWrongType, "Index %q does not contain a string: %T found", index, p[index])
	}

	return vs, nil
}

func (p Poja) ObjectValue(index int) (Pojo, error) {
	if index >= len(p) {
		return nil, errors.Wrapf(ErrNoSuchDetailsKey, "Index %d out of bounds", index)
	}

	vs, ok := p[index].(map[string]interface{})
	if !ok {
		return nil, errors.Wrapf(ErrValueWrongType, "Index %q does not contain an object: %T found", index, p[index])
	}

	return vs, nil
}

func (p Poja) BoolValue(index int) (bool, error) {
	if index >= len(p) {
		return false, errors.Wrapf(ErrNoSuchDetailsKey, "Index %d out of bounds", index)
	}

	vs, ok := p[index].(bool)
	if !ok {
		return false, errors.Wrapf(ErrValueWrongType, "Index %q does not contain a bool: %T found", index, p[index])
	}

	return vs, nil
}

func (p Poja) FloatValue(index int) (float64, error) {
	if index >= len(p) {
		return 0, errors.Wrapf(ErrNoSuchDetailsKey, "Index %d out of bounds", index)
	}

	var vs float64
	switch vt := p[index].(type) {
	case float64:
		vs = vt
	case json.Number:
		var err error
		vs, err = vt.Float64()
		if err != nil {
			return 0, err
		}
	default:
		return 0, errors.Wrapf(ErrValueWrongType, "Index %q does not contain a string: %T found", index, p[index])
	}

	return vs, nil
}

func (p Poja) ArrayValue(index int) (Poja, error) {
	if index >= len(p) {
		return nil, errors.Wrapf(ErrNoSuchDetailsKey, "Index %d out of bounds", index)
	}

	vs, ok := p[index].([]interface{})
	if !ok {
		return nil, errors.Wrapf(ErrValueWrongType, "Index %q does not contain an array: %T found", index, p[index])
	}

	return vs, nil
}
