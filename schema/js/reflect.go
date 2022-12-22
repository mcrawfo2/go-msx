// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package js

import (
	"cto-github.cisco.com/NFV-BU/go-msx/schema"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/iancoleman/strcase"
	"github.com/swaggest/jsonschema-go"
	"github.com/swaggest/refl"
	"reflect"
	"strings"
)

type ReflectContextOptionFunc func(*jsonschema.ReflectContext)

func DefNameInterceptor() ReflectContextOptionFunc {
	return func(rc *jsonschema.ReflectContext) {
		rc.DefName = func(t reflect.Type, defaultDefName string) string {
			if exposer, ok := reflect.New(t).Interface().(DefNameExposer); ok {
				return exposer.JSONSchemaDefName()
			}

			defName := schema.Namer().TypeName(t)
			if defName != "" {
				return defName
			}

			return defaultDefName
		}
	}
}

func TypeTitleInterceptor() ReflectContextOptionFunc {
	return jsonschema.InterceptType(func(value reflect.Value, schema *jsonschema.Schema) (bool, error) {
		if schema.ID != nil {
			schema.WithTitle(*schema.ID)
			return false, nil
		}

		valueType := value.Type()

		typeName, pkgPath := valueType.Name(), valueType.PkgPath()

		if pkgPath != "" {
			schema.Title = types.NewStringPtr(typeName)
		}

		// Continue with type interceptor chain
		return false, nil
	})
}

func FindRequiredJsonFields(valueType reflect.Type) []string {
	valueType = refl.DeepIndirect(valueType)

	if valueType.Kind() != reflect.Struct {
		return nil
	}

	var required []string
	for i := 0; i < valueType.NumField(); i++ {
		structField := valueType.Field(i)

		r := true
		if structField.Type.Kind() == reflect.Ptr ||
			structField.Type.Kind() == reflect.Map ||
			structField.Type.Kind() == reflect.Slice {

			r = structField.Type == reflect.TypeOf(types.UUID{})
		}

		jsonTag, ok := structField.Tag.Lookup("json")
		name := ""
		if ok {
			res := strings.Split(jsonTag, ",")
			overrideName := res[0]

			if overrideName == "-" {
				continue
			}
			name = overrideName
		}

		requiredTag, ok := structField.Tag.Lookup("required")
		if ok {
			r = requiredTag == "true"
		}

		optionalTag, ok := structField.Tag.Lookup("optional")
		if ok {
			r = optionalTag != "true"
		}

		if r {
			if name == "" {
				name = strcase.ToLowerCamel(structField.Name)
			}

			required = append(required, name)
		}
	}

	return required
}

func StructRequiredInterceptor() ReflectContextOptionFunc {
	return jsonschema.InterceptType(func(value reflect.Value, schema *jsonschema.Schema) (bool, error) {
		required := FindRequiredJsonFields(value.Type())

		if len(required) > 0 {
			schema.WithRequired(required...)
		}

		// Continue with type interceptor chain
		return false, nil
	})
}

func ExampleInterceptor() ReflectContextOptionFunc {
	return jsonschema.InterceptType(func(value reflect.Value, schema *jsonschema.Schema) (bool, error) {
		valueIface := value.Interface()

		if exampler, ok := valueIface.(Exampler); ok {
			example := exampler.Example()
			schema.Examples = []interface{}{example}
		}

		return false, nil
	})
}

func EnvelopNullability() ReflectContextOptionFunc {
	return func(rc *jsonschema.ReflectContext) {
		rc.EnvelopNullability = true
	}
}

type Exampler interface {
	Example() interface{}
}
