package types

import (
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

	vs, ok := v.(string)
	if !ok {
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

	vs, ok := v.(float64)
	if !ok {
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

	vs, ok := p[index].(string)
	if !ok {
		return "", errors.Wrap(ErrValueWrongType, "Index %q does not contain a string")
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

	vs, ok := p[index].(float64)
	if !ok {
		return 0, errors.Wrapf(ErrValueWrongType, "Index %q does not contain a float: %T found", index, p[index])
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
