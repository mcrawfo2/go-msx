// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package ops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/validate"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"io"
	"reflect"
)

var ErrMissingRequiredValue = errors.New("Missing required value")

type PortFieldInjector struct {
	portField   *PortField
	inputs      interface{}
	structValue reflect.Value
}

func (i PortFieldInjector) fieldValue() reflect.Value {
	fieldValue := i.structValue.FieldByIndex(i.portField.Indices)

	// Create nested pointers as necessary
	for n := 0; n < i.portField.Type.Indirections; n++ {
		here := reflect.New(fieldValue.Type().Elem())
		next := here.Elem()

		fieldValue.Set(here)
		fieldValue = next
	}

	// Return the addressable and settable Value
	return fieldValue
}

func (i PortFieldInjector) InjectPrimitive(value string) (err error) {
	if err = i.portField.ExpectShape(FieldShapePrimitive); err != nil {
		return
	}

	fv := i.fieldValue()

	defer func() {
		if err == nil {
			err = validate.ValidateValue(fv)
		}
	}()

	// Custom handlers
	switch i.portField.Type.HandlerType {
	case ByteSliceType:
		fv.Set(reflect.ValueOf([]byte(value)).Convert(i.portField.Type.Type))
		return nil
	case RuneSliceType:
		fv.Set(reflect.ValueOf([]rune(value)))
		return nil
	case TextUnmarshalerType:
		return fv.Addr().Interface().(types.TextUnmarshaler).UnmarshalText(value)
	}

	// Core types
	pv := fv.Addr().Interface()
	switch pvt := pv.(type) {
	case *int8:
		*pvt, err = cast.ToInt8E(value)
		return
	case *int16:
		*pvt, err = cast.ToInt16E(value)
		return
	case *int32:
		*pvt, err = cast.ToInt32E(value)
		return
	case *int64:
		*pvt, err = cast.ToInt64E(value)
		return
	case *int:
		*pvt, err = cast.ToIntE(value)
		return
	case *uint8:
		*pvt, err = cast.ToUint8E(value)
		return
	case *uint16:
		*pvt, err = cast.ToUint16E(value)
		return
	case *uint32:
		*pvt, err = cast.ToUint32E(value)
		return
	case *uint64:
		*pvt, err = cast.ToUint64E(value)
		return
	case *uint:
		*pvt, err = cast.ToUintE(value)
		return
	case *float32:
		*pvt, err = cast.ToFloat32E(value)
		return
	case *float64:
		*pvt, err = cast.ToFloat64E(value)
		return
	case *string:
		*pvt, err = cast.ToStringE(value)
		return
	case *bool:
		*pvt, err = cast.ToBoolE(value)
		return
	}

	return errors.Wrapf(ErrIncorrectShape,
		"Cannot apply primitive value to %T field",
		fv.Interface())
}

func (i PortFieldInjector) InjectContent(content Content) (err error) {
	if err = i.portField.ExpectShape(FieldShapeContent); err != nil {
		return
	}

	fv := i.fieldValue()

	defer func() {
		if err == nil {
			err = validate.ValidateValue(fv)
		}
	}()

	// Raw data
	switch i.portField.Type.HandlerType {
	case ByteSliceType:
		var data []byte
		if data, err = content.ReadBytes(); err != nil {
			return
		}
		fv.Set(reflect.ValueOf(data).Convert(i.portField.Type.Type))
		return nil
	case RuneSliceType:
		var data []byte
		if data, err = content.ReadBytes(); err != nil {
			return
		}
		var runes []rune
		runes = []rune(string(data))
		fv.Set(reflect.ValueOf(runes))
		return nil
	case ContentType:
		fv.Set(reflect.ValueOf(content))
		return nil
	case IoReadCloserType:
		var r io.ReadCloser
		if r, err = content.Reader(); err != nil {
			return err
		}
		fv.Set(reflect.ValueOf(r))
		return nil
	}

	// Marshaled Data
	pv := fv.Addr().Interface()
	return content.ReadEntity(pv)
}

func NewPortFieldInjector(portField *PortField, inputs interface{}) PortFieldInjector {
	return PortFieldInjector{
		portField:   portField,
		structValue: reflect.ValueOf(inputs).Elem(),
		inputs:      inputs,
	}
}
