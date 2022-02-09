package js

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/iancoleman/strcase"
	"github.com/swaggest/jsonschema-go"
	"reflect"
	"strings"
)

func TypeTitleDecorator() func(*jsonschema.ReflectContext) {
	return jsonschema.InterceptType(func(value reflect.Value, schema *jsonschema.Schema) (bool, error) {
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
			r = false
		}

		jsonTag, ok := structField.Tag.Lookup("json")
		if !ok {
			continue
		}

		res := strings.Split(jsonTag, ",")
		name := res[0]

		if name == "-" {
			continue
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

func StructRequiredDecorator() func(*jsonschema.ReflectContext) {
	return jsonschema.InterceptType(func(value reflect.Value, schema *jsonschema.Schema) (bool, error) {
		required := FindRequiredJsonFields(value.Type())

		if len(required) > 0 {
			schema.WithRequired(required...)
		}

		// Continue with type interceptor chain
		return false, nil
	})
}

func CustomizeSchema() func(*jsonschema.ReflectContext) {
	return jsonschema.InterceptType(func(value reflect.Value, schema *jsonschema.Schema) (bool, error) {
		valueIface := value.Interface()

		if exampler, ok := valueIface.(Exampler); ok {
			example := exampler.Example()
			schema.Examples = []interface{}{example}
		}

		return false, nil
	})
}

func EnvelopNullability() func(rc *jsonschema.ReflectContext) {
	return func(rc *jsonschema.ReflectContext) {
		rc.EnvelopNullability = true
	}
}

type Exampler interface {
	Example() interface{}
}

func NewType(simpleType jsonschema.SimpleType) *jsonschema.Type {
	val := simpleType.Type()
	return &val
}
