package types

import "github.com/pkg/errors"

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
		return "", errors.Wrap(ErrValueWrongType, "Key %q does not contain a string")
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
		return nil, errors.Wrap(ErrValueWrongType, "Key %q does not contain an Object")
	}

	return vs, nil
}

type PojoArray []map[string]interface{}

func (a PojoArray) Index(i int) Pojo {
	return a[i]
}
