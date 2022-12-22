// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package ops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/sanitize"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"reflect"
)

type PortFieldExtractor struct {
	portField    *PortField
	outputsValue reflect.Value
}

// extractRootValue returns the value of the current PortField from outputs
func (i PortFieldExtractor) extractRootValue() reflect.Value {
	return i.outputsValue.FieldByIndex(i.portField.Indices)
}

// extractIndirectRootValue returns the dereferenced value of the current PortField from outputs
func (i PortFieldExtractor) extractIndirectRootValue() (PortFieldElementType, reflect.Value) {
	fv := i.extractRootValue()

	// Dereference leading pointers
	for n := 0; n < i.portField.Type.Indirections; n++ {
		if fv.IsNil() {
			fv = reflect.Value{}
			break
		}
		fv = fv.Elem()
	}

	return i.rootPortFieldElement().WithIndirections(0), fv
}

// extractBelow returns the const or default value of the current PortField
func (i PortFieldExtractor) extractBelow() reflect.Value {
	fv := reflect.Value{}

	// Const
	if fv.Kind() == reflect.Invalid {
		c := i.portField.Const()
		if c != nil {
			fv = reflect.ValueOf(c)
		}
	}

	// Default
	if fv.Kind() == reflect.Invalid {
		d := i.portField.Default()
		if d != nil {
			fv = reflect.ValueOf(d)
		}
	}

	return fv
}

// ExtractValue returns the dereferenced value of the current PortField or its default/const value
func (i PortFieldExtractor) ExtractValue() (value reflect.Value, err error) {
	_, fv := i.extractIndirectRootValue()

	if !fv.IsValid() {
		if nfv := i.extractBelow(); nfv.IsValid() {
			fv = nfv
		} else {
			var v interface{}
			fv = reflect.ValueOf(v)
		}
	}

	return fv, nil
}

func (i PortFieldExtractor) ExtractPrimitive() (value types.Optional[string], err error) {
	pfet, fv := i.extractIndirectRootValue()

	var found bool
	found, value, err = i.extractElement(pfet, fv)

	if err != nil {
		return
	}

	if found {
		if value.IsPresent() {
			fv = reflect.ValueOf(value.Value())
		} else {
			fv = reflect.Value{}
		}
	}

	// Check for port field default/const value if value is missing
	if !fv.IsValid() || fv.IsZero() {
		if nfv := i.extractBelow(); nfv.IsValid() {
			fv = nfv
		} else {
			fv = reflect.Value{}
		}
	}

	if !fv.IsValid() {
		// No value returned
		if !i.portField.Optional {
			// Field not optional
			err = errors.Wrap(ErrMissingRequiredValue, i.portField.Name)
		} else {
			// Field optional
			value = types.OptionalEmpty[string]()
		}
		return
	}

	// Return the found value
	value = types.OptionalOf(fv.Interface().(string))
	return
}

func (i PortFieldExtractor) ExtractArray() (result []string, err error) {
	pfet, fv := i.extractIndirectRootValue()

	if !fv.IsValid() || (fv.Kind() == reflect.Slice && fv.IsNil()) {
		if nfv := i.extractBelow(); nfv.IsValid() {
			fv = nfv
			result, err = cast.ToStringSliceE(fv.Interface())
			return
		} else {
			fv = reflect.Value{}
		}
	}

	if !fv.IsValid() {
		// No value returned
		if !i.portField.Optional {
			// Field not optional
			err = errors.Wrap(ErrMissingRequiredValue, i.portField.Name)
		} else {
			// Field optional
			result = []string{}
		}
		return
	}

	// Convert sequences by element
	if fv.Kind() == reflect.Slice || fv.Kind() == reflect.Array {
		var es string
		result = make([]string, fv.Len())
		for e := 0; e < fv.Len(); e++ {
			efv := fv.Index(e)
			pfe := *pfet.Items

			var ev types.Optional[string]
			var found bool
			found, ev, err = i.extractElement(pfe, efv)
			if err != nil {
				return nil, err
			} else if !found && !pfe.Optional {
				err = errors.Wrap(ErrMissingRequiredValue, i.portField.Name)
			} else if ev.IsPresent() {
				es = ev.Value()
			} else {
				continue
			}

			result[e] = es
		}

		return
	}

	// Convert non-sequences using cast
	result, err = cast.ToStringSliceE(fv.Interface())
	if err != nil {
		err = errors.Wrapf(err, "Could not coerce %T to array", fv.Interface())
	}
	return
}

func (i PortFieldExtractor) ExtractObject() (result types.Pojo, err error) {
	_, fv := i.extractIndirectRootValue()

	if !fv.IsValid() ||
		(fv.Kind() == reflect.Map && fv.IsNil()) ||
		(fv.Kind() == reflect.Pointer && fv.IsNil()) {

		if nfv := i.extractBelow(); nfv.IsValid() {
			fv = nfv
		} else {
			fv = reflect.Value{}
		}
	}

	if !fv.IsValid() {
		// No value returned
		if !i.portField.Optional {
			// Field not optional
			err = errors.Wrap(ErrMissingRequiredValue, i.portField.Name)
		} else {
			// Field optional
			result = types.Pojo{}
		}
		return
	}

	if fv.Kind() == reflect.Struct || fv.Kind() == reflect.Map {
		// TODO: better

		var data json.RawMessage

		data, err = json.Marshal(fv.Interface())
		if err != nil {
			err = errors.Wrap(err, "Could not marshal field value to JSON")
			return
		}

		err = json.Unmarshal(data, &result)
		if err != nil {
			err = errors.Wrap(err, "Could not unmarshal field value from JSON")
			return
		}
	} else {
		result, err = cast.ToStringMapE(fv.Interface())
	}

	if err != nil {
		err = errors.Wrapf(err, "Could not coerce %T to object", fv.Interface())
	}
	return
}

type optionalValue interface {
	ValuePtrInterface() interface{}
}

// extractScalarIndirect converts some types that the cast module does not
func (i PortFieldExtractor) extractElement(pfet PortFieldElementType, fv reflect.Value) (found bool, optionalValue types.Optional[string], err error) {
	if !fv.IsValid() {
		return
	}

	// Dereference leading pointers
	for n := 0; n < pfet.Indirections; n++ {
		if fv.IsNil() {
			fv = reflect.Value{}
			break
		}
		fv = fv.Elem()
	}

	if !fv.IsValid() {
		return
	}

	v := fv.Interface()

	if v == nil {
		return
	}

	// If we swapped in a default or const, the handler type might not be correct
	useHandler := (fv.Type() == pfet.HandlerType) ||
		(pfet.HandlerType == TextMarshalerType &&
			fv.Type().Implements(pfet.HandlerType))

	if useHandler {
		// Custom handlers
		switch pfet.HandlerType {
		case ByteSliceType:
			value := fv.Interface().([]byte)
			return true, types.OptionalOf(string(value)), nil
		case RuneSliceType:
			value := fv.Interface().([]rune)
			return true, types.OptionalOf(string(value)), nil
		case TextMarshalerType:
			var value string
			value, err = fv.Addr().Interface().(types.TextMarshaler).MarshalText()
			if err != nil {
				return
			}
			return true, types.OptionalOf(value), nil

		case OptionalOfStringType:
			return true, fv.Interface().(types.Optional[string]), nil
		}
	}

	// Core types
	switch v.(type) {
	case int8, int16, int32, int64, int,
		uint8, uint16, uint32, uint64, uint,
		float32, float64, bool, string,
		json.Number, []byte:

		var value string
		value, err = cast.ToStringE(v)
		if err != nil {
			return
		}

		// Sanitize inputs
		san, _ := i.portField.BoolOption("san")
		if san {
			if err = sanitize.Input(&value, i.portField.SanitizeOptions()); err != nil {
				return
			}
		}

		return true, types.OptionalOf(value), nil

	default:
		err = errors.Wrapf(ErrIncorrectShape, "Cannot retrieve primitive value from %T field", v)
	}

	return
}

func (i PortFieldExtractor) rootPortFieldElement() PortFieldElementType {
	return PortFieldElementType{
		Indices:       i.portField.Indices,
		Optional:      i.portField.Optional,
		PortFieldType: i.portField.Type,
	}
}

func NewPortFieldExtractor(f *PortField, outputs interface{}) PortFieldExtractor {
	outputsValue := reflect.ValueOf(outputs)
	if reflect.TypeOf(outputs).Kind() == reflect.Ptr {
		outputsValue = outputsValue.Elem()
	}

	return PortFieldExtractor{
		portField:    f,
		outputsValue: outputsValue,
	}
}
