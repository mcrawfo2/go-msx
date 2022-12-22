// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package schema

import (
	"path"
	"reflect"
)

type TypeNamer interface {
	SetTypeName(t reflect.Type, name string)
	ParameterizedTypeName(wrapperType, wrappedType reflect.Type) string
	ParameterizedTypeNameWithWrappedName(wrapperType reflect.Type, name string) string
	SetParameterizedTypeName(derivedType, wrapperType, wrappedType reflect.Type)
	TypeInstanceName(instance interface{}) string
	TypeName(t reflect.Type) string
}

// Used everywhere except Swagger
var namer = NewDottedTypeNamer()

func Namer() TypeNamer {
	return namer
}

func GetNamedTypeName(instanceType reflect.Type) (string, bool) {
	instanceTypePackagePath := instanceType.PkgPath()
	instanceTypePackageName := ""

	if instanceTypePackagePath == "" {
		return "", false
	}

	instanceTypePackageName = path.Base(instanceTypePackagePath)
	if instanceTypePackageName != "" {
		instanceTypePackageName += "."
	}

	return instanceTypePackageName + instanceType.Name(), true
}
