package openapi

import "github.com/swaggest/openapi-go/openapi3"

const ExtraPropertiesMsxInjectedProperty = "x-msx-injected-property"

func InjectedPropertySchema(s *openapi3.Schema) *openapi3.SchemaOrRef {
	for _, oneOfAllOf := range s.AllOf {
		if oneOfAllOf.Schema == nil {
			continue
		}

		propNameAny, ok := oneOfAllOf.Schema.MapOfAnything[ExtraPropertiesMsxInjectedProperty]
		if !ok {
			continue
		}

		prop := oneOfAllOf.Schema.Properties[propNameAny.(string)]
		if prop.SchemaReference == nil {
			continue
		}

		return &prop
	}

	return nil
}
