// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package ops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/sanitize"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"cto-github.cisco.com/NFV-BU/go-msx/validate"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"io"
	"mime/multipart"
	"reflect"
)

var ErrMissingRequiredValue = errors.New("Missing required value")

type Injector interface {
	InjectPrimitive(value string) (err error)
	InjectArray(values []string) (err error)
	InjectObject(pojo types.Pojo) (err error)
	InjectFile(file *multipart.FileHeader) (err error)
	InjectFileArray(files []*multipart.FileHeader) (err error)
	InjectContent(content Content) (err error)
	InjectAny(value any) (err error)
}

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

func (i PortFieldInjector) injectIndirectScalar(fv reflect.Value, value string, pft PortFieldType) (err error) {
	// Create nested pointers as necessary
	for n := 0; n < pft.Indirections; n++ {
		here := reflect.New(fv.Type().Elem())
		next := here.Elem()

		fv.Set(here)
		fv = next
	}

	return i.injectScalar(fv, value, pft)
}

func (i PortFieldInjector) injectScalar(fv reflect.Value, value string, pft PortFieldType) (err error) {
	// Custom handlers
	switch pft.HandlerType {
	case ByteSliceType:
		fv.Set(reflect.ValueOf([]byte(value)).Convert(pft.Type))
		return nil
	case RuneSliceType:
		fv.Set(reflect.ValueOf([]rune(value)))
		return nil
	case TextUnmarshalerType:
		return fv.Addr().Interface().(types.TextUnmarshaler).UnmarshalText(value)
	}

	v := fv.Interface()

	// Core types
	switch vt := v.(type) {
	case int8:
		vt, err = cast.ToInt8E(value)
		fv.Set(reflect.ValueOf(vt))
	case int16:
		vt, err = cast.ToInt16E(value)
		fv.Set(reflect.ValueOf(vt))
	case int32:
		vt, err = cast.ToInt32E(value)
		fv.Set(reflect.ValueOf(vt))
	case int64:
		vt, err = cast.ToInt64E(value)
		fv.Set(reflect.ValueOf(vt))
	case int:
		vt, err = cast.ToIntE(value)
		fv.Set(reflect.ValueOf(vt))
	case uint8:
		vt, err = cast.ToUint8E(value)
		fv.Set(reflect.ValueOf(vt))
	case uint16:
		vt, err = cast.ToUint16E(value)
		fv.Set(reflect.ValueOf(vt))
	case uint32:
		vt, err = cast.ToUint32E(value)
		fv.Set(reflect.ValueOf(vt))
	case uint64:
		vt, err = cast.ToUint64E(value)
		fv.Set(reflect.ValueOf(vt))
	case uint:
		vt, err = cast.ToUintE(value)
		fv.Set(reflect.ValueOf(vt))
	case float32:
		vt, err = cast.ToFloat32E(value)
		fv.Set(reflect.ValueOf(vt))
	case float64:
		vt, err = cast.ToFloat64E(value)
		fv.Set(reflect.ValueOf(vt))
	case bool:
		vt, err = cast.ToBoolE(value)
		fv.Set(reflect.ValueOf(vt))
	case string:
		vt, err = cast.ToStringE(value)
		san, _ := i.portField.BoolOption("san")
		if san {
			if err = sanitize.Input(&vt, i.portField.SanitizeOptions()); err != nil {
				return err
			}
		}
		fv.Set(reflect.ValueOf(vt))

	default:
		return errors.Wrapf(ErrIncorrectShape,
			"Cannot apply primitive value to %T field",
			fv.Interface())
	}

	return err
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

	return i.injectScalar(fv, value, i.portField.Type)
}

func (i PortFieldInjector) InjectArray(values []string) (err error) {
	if err = i.portField.ExpectShape(FieldShapeArray); err != nil {
		return
	}

	fv := i.fieldValue()

	defer func() {
		if err == nil {
			err = validate.ValidateValue(fv)
		}
	}()

	sliceType := fv.Type()
	isPtr := sliceType.Kind() == reflect.Ptr
	if isPtr {
		sliceType = sliceType.Elem()
	}

	var sliceValue reflect.Value
	if sliceType.Kind() == reflect.Slice {
		sliceValue = reflect.MakeSlice(sliceType, len(values), len(values))
	} else {
		// TODO: verify array support
		sliceValue = reflect.New(sliceType).Elem()
	}

	for idx, queryValue := range values {
		err = i.injectIndirectScalar(sliceValue.Index(idx), queryValue, i.portField.Type.Items.PortFieldType)
		if err != nil {
			return err
		}
	}

	if isPtr {
		x := reflect.New(sliceType)
		x.Elem().Set(sliceValue)
		fv.Set(x)
	} else {
		fv.Set(sliceValue)
	}

	return nil
}

func (i PortFieldInjector) InjectObject(pojo types.Pojo) (err error) {
	if err = i.portField.ExpectShape(FieldShapeObject); err != nil {
		return
	}

	if pojo == nil && i.portField.Optional {
		return nil
	}

	fv := i.fieldValue()

	defer func() {
		if v := recover(); v != nil {
			err = errors.Errorf("Cannot marshal object %q into field %q: %s", pojo, i.portField.Name, v)
		} else if err == nil {
			err = validate.ValidateValue(fv)
		}
	}()

	objectType := fv.Type()
	isPtr := objectType.Kind() == reflect.Ptr
	if isPtr {
		objectType = objectType.Elem()
	}

	var objectRef reflect.Value
	var objectValue reflect.Value
	switch objectType.Kind() {
	case reflect.Map:
		objectValue = reflect.MakeMapWithSize(objectType, len(pojo))
		objectRef = reflect.New(objectType)
		objectRef.Elem().Set(objectValue)

		keyType := objectType.Key()
		valueType := objectType.Elem()
		for k, v := range pojo {
			entryKey := reflect.New(keyType).Elem()
			if err = i.injectIndirectScalar(entryKey, k, i.portField.Type.Keys.PortFieldType); err != nil {
				return err
			}

			entryValue := reflect.New(valueType).Elem()
			if err = i.injectIndirectScalar(entryValue, cast.ToString(v), i.portField.Type.Values.PortFieldType); err != nil {
				return err
			}

			objectValue.SetMapIndex(entryKey, entryValue)
		}
		break
	case reflect.Struct:
		objectRef = reflect.New(objectType)
		objectValue = objectRef.Elem()
		pft := i.portField.Type

		for _, pfet := range pft.Fields {
			var value string

			sf := objectType.FieldByIndex(pfet.Indices)

			name := strcase.ToLowerCamel(sf.Name)
			value, err = pojo.StringValue(name)
			if err == nil {
				entryValue := reflect.New(sf.Type).Elem()
				if err = i.injectIndirectScalar(entryValue, cast.ToString(value), pfet.PortFieldType); err != nil {
					return err
				}
				objectValue.FieldByIndex(sf.Index).Set(entryValue)
			} else {
				logger.WithError(err).Errorf("Failed to populate field %q", sf.Name)
			}
		}
		break

	}

	if isPtr {
		fv.Set(objectRef)
	} else {
		fv.Set(objectValue)
	}

	return nil
}

func (i PortFieldInjector) InjectFile(file *multipart.FileHeader) (err error) {
	if err = i.portField.ExpectShape(FieldShapeFile); err != nil {
		return
	}

	fieldValue := i.fieldValue()

	fieldValue.Set(reflect.ValueOf(file))

	return nil
}

func (i PortFieldInjector) InjectFileArray(files []*multipart.FileHeader) (err error) {
	if err = i.portField.ExpectShape(FieldShapeFileArray); err != nil {
		return
	}

	fieldValue := i.fieldValue()

	fieldValue.Set(reflect.ValueOf(files))

	return nil
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

func (i PortFieldInjector) InjectAny(value any) (err error) {
	if err = i.portField.ExpectShape(FieldShapeAny); err != nil {
		return
	}

	fieldValue := i.fieldValue()

	fieldValue.Set(reflect.ValueOf(value))

	return nil
}

func NewPortFieldInjector(portField *PortField, inputs interface{}) PortFieldInjector {
	return PortFieldInjector{
		portField:   portField,
		structValue: reflect.ValueOf(inputs).Elem(),
		inputs:      inputs,
	}
}
