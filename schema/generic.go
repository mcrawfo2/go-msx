// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package schema

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
	"github.com/swaggest/refl"
	"reflect"
	"strings"
)

// NewParameterizedStruct creates a runtime structure from an envelope (structType) and
// a payload.  This should only be used by swagger 2.0.
func NewParameterizedStruct(structType reflect.Type, payload interface{}) reflect.Type {
	payloadFieldIndex, payloadField, _ := FindParameterizedStructField(structType)
	if len(payloadFieldIndex) == 0 {
		return structType
	}

	var structFields []reflect.StructField
	for i := 0; i < structType.NumField(); i++ {
		structField := structType.Field(i)
		if structField.Tag.Get("inject") == payloadField {
			if payload == nil {
				continue
			} else {
				structField.Type = reflect.TypeOf(payload)
			}
		}
		structFields = append(structFields, structField)
	}

	parameterizedStructType := reflect.StructOf(structFields)
	NewSpringTypeNamer().SetParameterizedTypeName(parameterizedStructType, structType, reflect.TypeOf(payload))
	NewDottedTypeNamer().SetParameterizedTypeName(parameterizedStructType, structType, reflect.TypeOf(payload))

	return parameterizedStructType
}

// GetJsonFieldName returns the json field name mapping from the supplied field's tag
func GetJsonFieldName(structField reflect.StructField) types.Optional[string] {
	json := structField.Tag.Get("json")
	if json == "" {
		return types.OptionalEmpty[string]()
	}

	jsonParts := strings.SplitN(json, ",", 2)
	return types.OptionalOf(jsonParts[0])
}

// FindParameterizedStructField searches the sequence of struct fields in the supplied structType
// looking for an `inject` struct field tag.
func FindParameterizedStructField(structType reflect.Type) (result []int, name string, err error) {
	structType = refl.DeepIndirect(structType)

	if structType.Kind() != reflect.Struct {
		err = errors.Errorf("Envelope type must be a struct, found %s", structType.String())
		return
	}

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if field.Anonymous && refl.DeepIndirect(field.Type).Kind() == reflect.Struct {
			var subResult []int
			subResult, name, err = FindParameterizedStructField(field.Type)
			if err != nil {
				return
			}
			if len(subResult) > 0 {
				return append([]int{i}, subResult...), name, nil
			}
		} else if !field.Anonymous {
			tag := field.Tag.Get("inject")
			if tag != "" {
				return []int{i}, tag, nil
			}
		}
	}

	return nil, "", nil
}
