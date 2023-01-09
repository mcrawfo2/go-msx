// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package sanitize

import (
	"github.com/kennygrant/sanitize"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

type Options struct {
	Accents  bool
	BaseName bool
	Xss      bool
	Name     bool
	Path     bool
	Secret   bool

	Inherit bool
}

func NewOptionsInherit() Options {
	return Options{Inherit: true}
}

func NewOptions(tag string) Options {
	var result Options
	for _, opt := range strings.Split(tag, ",") {
		switch opt {
		case "accents":
			result.Accents = true
		case "basename":
			result.BaseName = true
		case "xss":
			result.Xss = true
		case "name":
			result.Name = true
		case "path":
			result.Path = true
		case "-":
			result = Options{}
		default:
			// ignore
		}
	}
	return result
}

// String sanitizes a string value and returns the sanitized result.
func String(originalValue string, options Options) string {
	resultValue := originalValue

	if options.Accents {
		resultValue = sanitize.Accents(resultValue)
	}
	if options.BaseName {
		resultValue = sanitize.BaseName(resultValue)
	}
	if options.Xss {
		resultValue = sanitize.HTML(resultValue)
	}
	if options.Name {
		resultValue = sanitize.Name(resultValue)
	}
	if options.Path {
		resultValue = sanitize.Path(resultValue)
	}
	if options.Secret {
		resultValue = secretSanitizer.Secrets(resultValue)
	}

	return resultValue
}

// Input sanitizes values in-place
func Input(value interface{}, options Options) error {
	v := reflect.ValueOf(value)
	return InputValue(v, options)
}

func InputValue(v reflect.Value, options Options) error {
	vt := v.Type()
	return newSanitizer(options).walk(vt, v)
}

type optionsStack []Options

func (s *optionsStack) Push(element Options) {
	*s = append(*s, element)
}

func (s *optionsStack) Pop() {
	*s = (*s)[:len(*s)-1]
}

func (s *optionsStack) Active() Options {
	t := *s
	for i := len(t) - 1; i >= 0; i-- {
		if t[i].Inherit {
			continue
		}
		return t[i]
	}
	return Options{}
}

type sanitizer struct {
	optionsStack optionsStack
}

func (s *sanitizer) walk(vt reflect.Type, v reflect.Value) error {
	if v.Kind() == reflect.Interface {
		v = v.Elem()
		vt = v.Type()
	}

	switch vt.Kind() {
	case reflect.Ptr:
		return s.walkPtr(vt, v)
	case reflect.String:
		return s.walkString(vt, v)
	case reflect.Struct:
		return s.walkStruct(vt, v)
	case reflect.Map:
		return s.walkMap(vt, v)
	case reflect.Slice:
		return s.walkSlice(vt, v)
	default:
		return errors.New("Value must be a pointer or reference")
	}
}

func (s *sanitizer) walkMap(vt reflect.Type, v reflect.Value) error {
	if v.IsNil() {
		return nil
	}

	vte := vt.Elem()
	for _, k := range v.MapKeys() {
		ve := v.MapIndex(k)

		if ve.Kind() == reflect.Interface {
			ve = ve.Elem()
			vte = ve.Type()
		}

		pointer := false
		if !ve.CanSet() {
			ve = s.ptr(ve)
			vte = ve.Type()
			pointer = true
		}

		if err := s.walk(vte, ve); err != nil {
			return err
		}

		if pointer {
			v.SetMapIndex(k, ve.Elem())
		}

	}
	return nil
}

func (s *sanitizer) ptr(v reflect.Value) reflect.Value {
	pt := reflect.PtrTo(v.Type()) // create a *T type.
	pv := reflect.New(pt.Elem())  // create a reflect.Value of type *T.
	pv.Elem().Set(v)              // sets pv to point to underlying value of v.
	return pv
}

func (s *sanitizer) walkSlice(vt reflect.Type, v reflect.Value) error {
	if v.IsNil() {
		return nil
	}

	vte := vt.Elem()
	for i := 0; i < v.Len(); i++ {
		ve := v.Index(i)
		if err := s.walk(vte, ve); err != nil {
			return err
		}
	}

	return nil
}

func (s *sanitizer) walkPtr(_ reflect.Type, v reflect.Value) error {
	if v.IsNil() {
		// No value here
		return nil
	}

	ve := v.Elem()
	vte := ve.Type()

	return s.walk(vte, ve)
}

func (s *sanitizer) walkStruct(vt reflect.Type, v reflect.Value) error {
	for i := 0; i < vt.NumField(); i++ {
		sf := vt.Field(i)

		options := s.getStructFieldOptions(vt, sf.Name)
		s.optionsStack.Push(options)

		ve := v.FieldByName(sf.Name)
		vte := ve.Type()

		if !ve.CanSet() || !ve.IsValid() {
			return errors.Errorf("Unsettable field %q", sf.Name)
		}

		if err := s.walk(vte, ve); err != nil {
			return err
		}

		s.optionsStack.Pop()
	}

	return nil
}

var structFieldOptions = make(map[reflect.Type]map[string]Options)

func (s *sanitizer) getStructFieldOptions(vt reflect.Type, name string) Options {
	if _, ok := structFieldOptions[vt]; !ok {
		// Populate the san options for each field of the struct
		structFieldOptions[vt] = make(map[string]Options)
		for i := 0; i < vt.NumField(); i++ {
			sf := vt.Field(i)
			var options = NewOptionsInherit()
			if sanTag, ok := sf.Tag.Lookup("san"); ok {
				options = NewOptions(sanTag)
			}
			structFieldOptions[vt][sf.Name] = options
		}
	}
	// Return the options for the specified field of the struct
	return structFieldOptions[vt][name]
}

var stringType = reflect.TypeOf("")

func (s *sanitizer) walkString(vt reflect.Type, v reflect.Value) error {
	originalValue := v.Convert(stringType).Interface().(string)
	sanitizedValue := String(originalValue, s.optionsStack.Active())
	v.Set(reflect.ValueOf(sanitizedValue).Convert(vt))
	return nil
}

func newSanitizer(options Options) *sanitizer {
	return &sanitizer{
		optionsStack: optionsStack{options},
	}
}
