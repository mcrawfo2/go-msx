// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package schema

import (
	"fmt"
	"github.com/swaggest/refl"
	"go.uber.org/atomic"
	"reflect"
	"sync"
)

const (
	VoidTypeName = "types.Void"
)

type DottedTypeNamer struct {
	Named    map[reflect.Type]string
	Suffixes map[string]string
	Anon     *atomic.Int64
}

func (n DottedTypeNamer) SetTypeName(t reflect.Type, name string) {
	dottedTypeNamerNamedMtx.Lock()
	defer dottedTypeNamerNamedMtx.Unlock()
	n.Named[t] = name
}

func (n DottedTypeNamer) getNamed(t reflect.Type) (string, bool) {
	dottedTypeNamerNamedMtx.Lock()
	defer dottedTypeNamerNamedMtx.Unlock()
	val, ok := n.Named[t]
	return val, ok
}

func (n DottedTypeNamer) SetParameterizedTypeName(t reflect.Type, wrapperType reflect.Type, wrappedType reflect.Type) {
	n.Named[t] = n.ParameterizedTypeName(wrapperType, wrappedType)
}

func (n DottedTypeNamer) ParameterizedTypeNameWithWrappedName(wrapperType reflect.Type, wrappedTypeName string) string {
	wrapperTypeName := n.TypeName(wrapperType)

	deepWrapperType := refl.DeepIndirect(wrapperType)
	if deepWrapperType.Kind() == reflect.Struct {
		indices, wrapperSuffixName, err := FindParameterizedStructField(deepWrapperType)
		if len(indices) != 0 && wrapperSuffixName != "" && err == nil {
			n.Suffixes[wrapperTypeName] = wrapperSuffixName
		}
	}

	if suffix, ok := n.Suffixes[wrapperTypeName]; ok {
		return wrappedTypeName + "." + suffix
	} else {
		return wrappedTypeName + "." + wrapperTypeName
	}
}

func (n DottedTypeNamer) ParameterizedTypeName(wrapperType, wrappedType reflect.Type) string {
	wrappedTypeName := n.TypeName(wrappedType)
	return n.ParameterizedTypeNameWithWrappedName(wrapperType, wrappedTypeName)
}

// TypeInstanceName returns the canonical schema name for the type of this instance
func (n DottedTypeNamer) TypeInstanceName(instance interface{}) string {
	if instance == nil {
		return VoidTypeName
	}

	t := reflect.TypeOf(instance)
	return n.TypeName(t)
}

// TypeName returns the canonical schema name of this type
func (n DottedTypeNamer) TypeName(instanceType reflect.Type) string {
	typeName := ""
	ok := false

	if instanceType == nil {
		return VoidTypeName
	}

	// Unwrap the pointers
	instanceType = refl.DeepIndirect(instanceType)

	// Lookup any overridden type name
	if typeName, ok = n.getNamed(instanceType); ok {
		return typeName
	}

	if typeName, ok = GetNamedTypeName(instanceType); ok {
		return typeName
	}

	switch instanceType.Kind() {
	case reflect.Array, reflect.Slice:
		typeName = n.TypeName(instanceType.Elem())
		return typeName + "." + n.Suffixes["array"]

	case reflect.Map:
		typeName = n.TypeName(instanceType.Elem())
		return typeName + "." + n.Suffixes["object"]
	}

	switch instanceType.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Bool, reflect.String, reflect.Complex64, reflect.Complex128:
		return instanceType.Name()

	case reflect.Struct:
		return fmt.Sprintf(".anonymous%d", n.Anon.Inc())

	default:
		return instanceType.Name()
	}
}

var dottedTypeNamerNamed = make(map[reflect.Type]string)
var dottedTypeNamerNamedMtx sync.Mutex

func NewDottedTypeNamer() TypeNamer {
	return &DottedTypeNamer{
		Named: dottedTypeNamerNamed,
		Suffixes: map[string]string{
			"array":  "List",
			"object": "Map",
		},
		Anon: new(atomic.Int64),
	}
}
