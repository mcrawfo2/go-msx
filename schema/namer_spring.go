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

type SpringTypeNameParts struct {
	Prefix string
	Infix  string
	Suffix string
}

type SpringTypeNamer struct {
	Named    map[reflect.Type]string
	Wrappers map[string]SpringTypeNameParts
	Anon     *atomic.Int64
}

func (n SpringTypeNamer) SetTypeName(t reflect.Type, name string) {
	springTypeNamerNamedMtx.Lock()
	defer springTypeNamerNamedMtx.Unlock()
	n.Named[t] = name
}

func (n SpringTypeNamer) getNamed(t reflect.Type) (string, bool) {
	springTypeNamerNamedMtx.Lock()
	defer springTypeNamerNamedMtx.Unlock()
	val, ok := n.Named[t]
	return val, ok
}

func (n SpringTypeNamer) SetParameterizedTypeName(t reflect.Type, wrapperType reflect.Type, wrappedType reflect.Type) {
	typeName := n.ParameterizedTypeName(wrapperType, wrappedType)
	n.SetTypeName(t, typeName)
}

func (n SpringTypeNamer) ParameterizedTypeNameWithWrappedName(wrapperType reflect.Type, wrappedTypeName string) string {
	wrapperTypeName := n.TypeName(wrapperType)
	if wrapper, ok := n.Wrappers[wrapperTypeName]; ok {
		return wrapper.Prefix + wrappedTypeName + wrapper.Suffix
	} else {
		return fmt.Sprintf("%s«%s»", wrapperTypeName, wrappedTypeName)
	}
}

func (n SpringTypeNamer) ParameterizedTypeName(wrapperType, wrappedType reflect.Type) string {
	wrappedTypeName := n.TypeName(wrappedType)
	return n.ParameterizedTypeNameWithWrappedName(wrapperType, wrappedTypeName)
}

// TypeInstanceName returns the canonical schema name for the type of this instance
func (n SpringTypeNamer) TypeInstanceName(instance interface{}) string {
	if instance == nil {
		return VoidTypeName
	}

	instanceType := reflect.TypeOf(instance)
	return n.TypeName(instanceType)
}

// TypeName returns the canonical schema name of this type
func (n SpringTypeNamer) TypeName(instanceType reflect.Type) string {
	var typeName string
	var ok bool

	if instanceType == nil {
		return VoidTypeName
	}

	instanceType = refl.DeepIndirect(instanceType)

	if typeName, ok = n.getNamed(instanceType); ok {
		return typeName
	}

	if typeName, ok = GetNamedTypeName(instanceType); ok {
		return typeName
	}

	switch instanceType.Kind() {
	case reflect.Array, reflect.Slice:
		wrapperParts := n.Wrappers["array"]
		valueName := n.TypeName(instanceType.Elem())
		return wrapperParts.Prefix + valueName + wrapperParts.Suffix

	case reflect.Map:
		wrapperParts := n.Wrappers["object"]
		keyName := n.TypeName(instanceType.Key())
		valueName := n.TypeName(instanceType.Elem())
		return wrapperParts.Prefix + keyName + wrapperParts.Infix + valueName + wrapperParts.Suffix
	}

	switch instanceType.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Bool, reflect.String:
		return instanceType.Name()

	case reflect.Struct:
		return fmt.Sprintf(".anonymous%d", n.Anon.Inc())

	default:
		return instanceType.Name()
	}
}

var springTypeNamerNamed = make(map[reflect.Type]string)
var springTypeNamerNamedMtx sync.Mutex

func NewSpringTypeNamer() TypeNamer {
	return &SpringTypeNamer{
		Named: springTypeNamerNamed,
		Wrappers: map[string]SpringTypeNameParts{
			"integration.MsxEnvelope": {
				Prefix: "Envelope«",
				Suffix: "»",
			},
			"paging.PaginatedResponse": {
				Prefix: "Page«",
				Suffix: "»",
			},
			"paging.PaginatedResponseV8": {
				Prefix: "Page«",
				Suffix: "»",
			},
			"array": {
				Prefix: "List«",
				Suffix: "»",
			},
			"object": {
				Prefix: "Map«",
				Infix:  ",",
				Suffix: "»",
			},
		},
		Anon: new(atomic.Int64),
	}
}
