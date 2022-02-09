// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package validate

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"reflect"
)

type Validatable interface {
	Validate() error
}

func Validate(validatable Validatable) error {
	err := validatable.Validate()
	if err != nil {
		if filterable, ok := err.(types.Filterable); ok {
			err = filterable.Filter()
		}
	}
	return err
}

var validatableInstance Validatable
var validatableType = reflect.TypeOf(&validatableInstance).Elem()

func ValidateValue(value reflect.Value) error {
	if value.Type().Implements(validatableType) {
		return Validate(value.Interface().(Validatable))
	}

	if value.Kind() == reflect.Ptr {
		return nil
	}

	if !value.CanAddr() {
		return nil
	}

	value = value.Addr()
	if value.Type().Implements(validatableType) {
		return Validate(value.Interface().(Validatable))
	}

	return nil
}
