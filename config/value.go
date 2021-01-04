package config

import (
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

type Values interface {
	String(key string) (string, error)
	StringOr(key, alt string) (string, error)
	Int(key string) (int, error)
	IntOr(key string, alt int) (int, error)
	Uint(key string) (uint, error)
	UintOr(key string, alt uint) (uint, error)
	Float(key string) (float64, error)
	FloatOr(key string, alt float64) (float64, error)
	Bool(key string) (bool, error)
	BoolOr(key string, alt bool) (bool, error)
	Duration(key string) (time.Duration, error)
	DurationOr(key string, alt time.Duration) (time.Duration, error)

	Value(key string) (Value, error)
	Populate(target interface{}, prefix string) error

	// Deprecated.
	Settings() map[string]string
	// TODO: Deprecate.  Replace with EachValue.
	Each(target func(string, string))
}

type Value string

func (v Value) StringPtr() *string {
	s := string(v)
	return &s
}

func (v Value) String() string {
	return string(v)
}

func (v Value) Int() (int64, error) {
	i, err := strconv.ParseInt(string(v), 10, 64)
	if err != nil {
		return 0, errors.Wrap(ErrInvalidValue, err.Error())
	}
	return i, nil
}

func (v Value) Uint() (uint64, error) {
	i, err := strconv.ParseUint(string(v), 10, 64)
	if err != nil {
		return 0, errors.Wrap(ErrInvalidValue, err.Error())
	}
	return i, nil
}

func (v Value) Float() (float64, error) {
	i, err := strconv.ParseFloat(string(v), 64)
	if err != nil {
		return 0, errors.Wrap(ErrInvalidValue, err.Error())
	}
	return i, nil
}

func (v Value) Bool() (bool, error) {
	i, err := strconv.ParseBool(string(v))
	if err != nil {
		return false, errors.Wrap(ErrInvalidValue, err.Error())
	}
	return i, nil
}

func (v Value) StringSlice(sep string) []string {
	if len(v) == 0 {
		return nil
	}
	return strings.Split(string(v), sep)
}

func (v Value) Duration() (time.Duration, error) {
	i, err := time.ParseDuration(string(v))
	if err != nil {
		return 0, errors.Wrap(ErrInvalidValue, err.Error())
	}
	return i, nil
}
