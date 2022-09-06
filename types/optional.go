// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package types

import "encoding/json"

type OptionalString struct {
	Value *string
}

func (s OptionalString) IsPresent() bool {
	return s.Value != nil
}

func (s OptionalString) OrOptional(value *string) OptionalString {
	if s.Value == nil {
		return OptionalString{value}
	}
	return s
}

func (s OptionalString) OrElse(value string) string {
	if s.Value != nil {
		return *s.Value
	}
	return value
}

func (s OptionalString) OrEmpty() string {
	return s.OrElse("")
}

func (s OptionalString) String() string {
	return s.OrElse("<nil>")
}

func (s OptionalString) Ptr() *string {
	return s.Value
}

func (s OptionalString) NilIfEmpty() OptionalString {
	if s.Value != nil && *s.Value == "" {
		return OptionalString{}
	}
	return s
}

func (s OptionalString) Equals(other OptionalString) bool {
	if s.Value == other.Value {
		return true
	} else if s.Value == nil || other.Value == nil {
		return false
	} else {
		return *s.Value == *other.Value
	}
}

func NewOptionalString(value *string) OptionalString {
	return OptionalString{Value: value}
}

func NewOptionalStringFromString(value string) OptionalString {
	return OptionalString{Value: &value}
}

type OptionalBool struct {
	IsValid bool
	Value   bool
}

func (b OptionalBool) OrElse(value bool) bool {
	if b.IsValid {
		return b.Value
	}
	return value
}

func NewOptionalBool(value *bool) OptionalBool {
	if value != nil {
		return OptionalBool{
			IsValid: true,
			Value:   *value,
		}
	}
	return OptionalBool{}
}

func NewOptionalBoolFromBool(value bool) OptionalBool {
	return OptionalBool{
		IsValid: true,
		Value:   value,
	}
}

type OptionalInt struct {
	IsValid bool
	Value   int
}

func (b OptionalInt) OrElse(value int) int {
	if b.IsValid {
		return b.Value
	}
	return value
}

func NewOptionalInt(value *int) OptionalInt {
	if value != nil {
		return OptionalInt{
			IsValid: true,
			Value:   *value,
		}
	}
	return OptionalInt{}
}

func NewOptionalIntFromInt(value int) OptionalInt {
	return OptionalInt{
		IsValid: true,
		Value:   value,
	}
}

type Optional[I any] struct {
	valid bool
	value I
}

func (o Optional[I]) IsPresent() bool {
	return o.valid
}

func (o Optional[I]) OrElse(v I) I {
	if o.valid {
		return o.value
	}
	return v
}

func (o Optional[I]) Project(v func(I) Optional[I]) Optional[I] {
	if o.valid {
		return v(o.value)
	}
	return o
}

func (o Optional[I]) Value() I {
	var def I
	if o.valid {
		def = o.value
	}
	return def
}

func (o Optional[I]) ValuePtrInterface() interface{} {
	if !o.valid {
		return nil
	}
	return &o.value
}

func (o Optional[I]) IfPresent(fn func(v I)) {
	if o.valid {
		fn(o.value)
	}
}

func (o Optional[I]) IfPresentE(fn func(v I) error) error {
	if o.valid {
		return fn(o.value)
	}
	return nil
}

func (o *Optional[I]) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &o.value)
	o.valid = err == nil
	return err
}

func (o Optional[I]) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.value)
}

func OptionalOf[I any](v I) Optional[I] {
	return Optional[I]{
		valid: true,
		value: v,
	}
}

func OptionalEmpty[I any]() Optional[I] {
	return Optional[I]{}
}

func InvalidateZero[I comparable](v I) Optional[I] {
	var zero I
	if v == zero {
		return OptionalEmpty[I]()
	}

	return OptionalOf(v)
}

func NewRuneSlicePtr(value []rune) *[]rune {
	return &value
}

func NewByteSlicePtr(value []byte) *[]byte {
	return &value
}

func NewStringPtr(value string) *string {
	return &value
}

func NewIntPtr(value int) *int {
	return &value
}

func NewUintPtr(value uint) *uint {
	return &value
}

func NewBoolPtr(value bool) *bool {
	return &value
}

func NewFloatPtr(value float64) *float64 {
	return &value
}

func NewTimePtr(value Time) *Time {
	return &value
}
