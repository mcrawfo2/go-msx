// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package ops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"reflect"
)

type PortFieldExtractor struct {
	portField    *PortField
	outputs      interface{}
	outputsValue reflect.Value
}

func (i PortFieldExtractor) ExtractRawValue() reflect.Value {
	return i.outputsValue.FieldByIndex(i.portField.Indices)
}

func (i PortFieldExtractor) extractBelow(fv reflect.Value) reflect.Value {
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

func (i PortFieldExtractor) ExtractValue() (fv reflect.Value, err error) {
	fv = i.ExtractRawValue()

	for fv.Kind() == reflect.Ptr {
		if fv.IsNil() {
			fv = reflect.Value{}
			break
		}
		fv = fv.Elem()
	}

	fv = i.extractBelow(fv)

	// No value found
	if fv.Kind() == reflect.Invalid && !i.portField.Optional {
		err = errors.Errorf("Missing required value for field %q", i.portField.Name)
	}

	return
}

func (i PortFieldExtractor) ExtractPrimitive() (value types.Optional[string], err error) {
	fv, err := i.ExtractValue()
	if err != nil {
		return
	}

	if fv.Kind() == reflect.Invalid {
		return
	}

	var done bool
	if done, value, err = i.extractPrimitiveIndirect(fv); err != nil || done {
		if value.IsPresent() || err != nil {
			return
		}

		// if the optional value is not present, check for default/const
		fv = i.extractBelow(reflect.Value{})

		if !fv.IsValid() {
			return types.OptionalEmpty[string](), nil
		}
	}

	result, err := cast.ToStringE(fv.Interface())
	if err != nil {
		err = errors.Wrapf(err, "Could not coerce %T to primitive", fv.Interface())
	} else {
		value = types.OptionalOf(result)
	}
	return
}

type optionalValue interface {
	ValuePtrInterface() interface{}
}

// extractPrimitiveIndirect converts some types that the cast module does not
func (i PortFieldExtractor) extractPrimitiveIndirect(fv reflect.Value) (bool, types.Optional[string], error) {
	fvi := fv.Interface()
	switch fvit := fvi.(type) {
	case types.TextMarshaler:
		value, err := fvit.MarshalText()
		if err != nil {
			return false, types.OptionalEmpty[string](), err
		}
		return true, types.OptionalOf(value), nil
	case types.Optional[string]:
		return true, fvit, nil
		// TODO: Other optional types
	}

	if fv.Kind() == reflect.Ptr {
		fv = fv.Elem()
		return i.extractPrimitiveIndirect(fv)
	}

	switch fvi.(type) {
	case []rune:
		result, ok := fvi.([]rune)
		if ok {
			value := string(result)
			return true, types.OptionalOf(value), nil
		}
	}

	return false, types.OptionalEmpty[string](), nil
}

func NewPortFieldExtractor(f *PortField, outputs interface{}) PortFieldExtractor {
	outputsValue := reflect.ValueOf(outputs)
	if reflect.TypeOf(outputs).Kind() == reflect.Ptr {
		outputsValue = outputsValue.Elem()
	}

	return PortFieldExtractor{
		portField:    f,
		outputs:      outputs,
		outputsValue: outputsValue,
	}
}