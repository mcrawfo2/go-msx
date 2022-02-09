package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/schema/js"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/swaggest/jsonschema-go"
	"github.com/swaggest/openapi-go/openapi3"
	"reflect"
	"regexp"
	"strings"
)

var Reflector = openapi3.Reflector{
	Reflector: jsonschema.Reflector{
		DefaultOptions: []func(*jsonschema.ReflectContext){
			jsonschema.DefinitionsPrefix("#/components/schemas/"),
			jsonschema.RootRef,
			js.TypeTitleDecorator(),
			OverrideSchema(),
			js.StructRequiredDecorator(),
			js.EnvelopNullability(),
			js.CustomizeSchema(),
		},
	},
	Spec: &openapi3.Spec{
		Openapi: "3.0.3",
	},
}

var JsonReflector = jsonschema.Reflector{
	DefaultOptions: []func(ctx *jsonschema.ReflectContext){
		jsonschema.RootRef,
		OverrideSchema(),
		js.StructRequiredDecorator(),
	},
}

var openApiSchemaIndex = make(map[reflect.Type]openapi3.SchemaOrRef)

func GetSchemaOrRef(value interface{}) (openapi3.SchemaOrRef, bool) {
	vt := reflect.TypeOf(value)
	v, e := openApiSchemaIndex[vt]
	return v, e
}

func GetSchemaOrRefEns(value interface{}) (openapi3.SchemaOrRef, error) {
	v, e := GetSchemaOrRef(value)
	if !e {
		vp, err := Reflect(value)
		if err != nil {
			return openapi3.SchemaOrRef{}, err
		}
		v = *vp
	}
	return v, nil
}

func SetSchemaOrRef(t reflect.Type, s openapi3.SchemaOrRef) {
	openApiSchemaIndex[t] = s
}

var overrideSchema = make(map[reflect.Type]jsonschema.Schema)
var overrideSchemaById = make(map[string]reflect.Type)

func OverrideSchema() func(*jsonschema.ReflectContext) {
	return jsonschema.InterceptType(func(value reflect.Value, schema *jsonschema.Schema) (bool, error) {
		valueType := value.Type()
		if valueSchemaOverride, ok := overrideSchema[valueType]; ok {
			*schema = valueSchemaOverride
		}

		// Continue with type interceptor chain
		return false, nil
	})
}

func RegisterTypeSchema(t reflect.Type, schema jsonschema.Schema) {
	overrideSchema[t] = schema
	if schema.ID != nil {
		overrideSchemaById[*schema.ID] = t
	}

	s := ConvertToOpenApiSchema(schema)
	Reflector.SpecEns().ComponentsEns().SchemasEns().WithMapOfSchemaOrRefValuesItem(*schema.ID, s)
}

func ReflectOpenApiJson(value interface{}) (jsonschema.Schema, error) {
	s, err := Reflector.Reflect(
		value,
		jsonschema.CollectDefinitions(CollectDefinition),
	)

	if err != nil {
		return jsonschema.Schema{}, err
	}

	return s, nil
}

func Reflect(value interface{}) (*openapi3.SchemaOrRef, error) {
	schema, err := ReflectOpenApiJson(value)
	if err != nil {
		return nil, err
	}

	s := ConvertToOpenApiSchema(schema)
	s = PromoteNullableSchema(s)
	SetSchemaOrRef(reflect.TypeOf(value), s)

	if s.SchemaReference != nil {
		// Get type name
		schemaRefName := SchemaRefName(s.SchemaReference)

		// patch the schema if desired
		if directSchema, ok := LookupSchema(schemaRefName); ok {
			// TODO: referenced schemas
			if patcher, ok := value.(Patcher); ok {
				// filter the schema through the patcher
				patcher.PatchOpenApiSchema(directSchema)
			}
		}
	}

	return &s, nil
}

type Patcher interface {
	PatchOpenApiSchema(schema *openapi3.Schema)
}

func ConvertToOpenApiSchema(schema jsonschema.Schema) openapi3.SchemaOrRef {
	s := openapi3.SchemaOrRef{}
	s.FromJSONSchema(schema.ToSchemaOrBool())
	return s
}

func PromoteNullableSchema(osr openapi3.SchemaOrRef) openapi3.SchemaOrRef {
	if osr.Schema == nil {
		return osr
	}

	// Convert enveloped schema
	if len(osr.Schema.AnyOf) == 2 {
		for i := 0; i < len(osr.Schema.AnyOf); i++ {
			one := osr.Schema.AnyOf[i]
			if one.Schema == nil {
				continue
			}
			if one.Schema.Nullable != nil && *one.Schema.Nullable {
				osr.Schema.Nullable = one.Schema.Nullable
				osr.Schema.AllOf = []openapi3.SchemaOrRef{
					osr.Schema.AnyOf[2-i-1],
				}
				osr.Schema.AnyOf = nil
				break
			}
		}
	}

	for key, property := range osr.Schema.Properties {
		osr.Schema.Properties[key] = PromoteNullableSchema(property)
	}

	if osr.Schema.Items != nil {
		*osr.Schema.Items = PromoteNullableSchema(*osr.Schema.Items)
	}

	if osr.Schema.AdditionalProperties != nil && osr.Schema.AdditionalProperties.SchemaOrRef != nil {
		*osr.Schema.AdditionalProperties.SchemaOrRef = PromoteNullableSchema(*osr.Schema.AdditionalProperties.SchemaOrRef)
	}

	return osr
}

func CollectDefinition(name string, schema jsonschema.Schema) {
	if _, ok := LookupSchema(name); ok {
		// Don't override
		return
	}

	s := ConvertToOpenApiSchema(schema)
	s = PromoteNullableSchema(s)
	SetSchemaOrRef(schema.ReflectType, s)

	if s.SchemaReference != nil {
		// Get type name
		schemaRefName := SchemaRefName(s.SchemaReference)

		// patch the schema if desired
		if directSchema, ok := LookupSchema(schemaRefName); ok {
			value := reflect.New(schema.ReflectType).Interface()
			if patcher, ok := value.(Patcher); ok {
				// filter the schema through the patcher
				patcher.PatchOpenApiSchema(directSchema)
			}
		}
	}

	Reflector.SpecEns().ComponentsEns().SchemasEns().WithMapOfSchemaOrRefValuesItem(name, s)
}

func LookupSchema(schemaName string) (*openapi3.Schema, bool) {
	result, ok := Reflector.SpecEns().ComponentsEns().SchemasEns().MapOfSchemaOrRefValues[schemaName]
	if !ok {
		return nil, false
	}
	return result.Schema, ok
}

func DeepLookupSchema(schemaName string) (*openapi3.Schema, bool) {
	result, ok := Reflector.SpecEns().ComponentsEns().SchemasEns().MapOfSchemaOrRefValues[schemaName]
	if !ok {
		return nil, false
	}
	if result.SchemaReference != nil {
		refName := SchemaRefName(result.SchemaReference)
		return DeepLookupSchema(refName)
	}
	return result.Schema, ok
}

var regexFindPathParameter = regexp.MustCompile(`{([^}:]+)(:[^/]+)?(?:})`)

func CleanPath(p string) string {
	pathParametersSubmatches := regexFindPathParameter.FindAllStringSubmatch(p, -1)
	for _, submatch := range pathParametersSubmatches {
		if submatch[2] != "" { // Remove gorilla.Mux-style regexp in path
			p = strings.Replace(p, submatch[0], "{"+submatch[1]+"}", 1)
		}
	}
	return p
}

func LookupOperation(method, path string) (result openapi3.Operation, ok bool) {
	path = CleanPath(path)

	pathItem, ok := Reflector.SpecEns().Paths.MapOfPathItemValues[path]
	if !ok {
		return
	}

	result, ok = pathItem.MapOfOperationValues[method]
	return
}

func ResolveParameter(p openapi3.ParameterOrRef) openapi3.ParameterOrRef {
	if p.Parameter != nil {
		return p
	}
	if p.ParameterReference == nil {
		return p
	}

	refName := ParameterRefName(p.ParameterReference)
	if target, ok := Reflector.SpecEns().ComponentsEns().ParametersEns().MapOfParameterOrRefValues[refName]; ok {
		return ResolveParameter(target)
	}

	return p
}

func init() {
	RegisterTypeSchema(
		reflect.TypeOf(types.UUID{}),
		jsonschema.Schema{
			ID:      types.NewStringPtr("TypesUUID"),
			Type:    NewType(jsonschema.String),
			Title:   types.NewStringPtr("UUID"),
			Format:  types.NewStringPtr("uuid"),
			Pattern: types.NewStringPtr(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`),
		})

	RegisterTypeSchema(
		reflect.TypeOf(types.Time{}),
		jsonschema.Schema{
			ID:      types.NewStringPtr("TypesTime"),
			Type:    NewType(jsonschema.String),
			Title:   types.NewStringPtr("Time"),
			Format:  types.NewStringPtr("date-time"),
			Pattern: types.NewStringPtr(`^([0-9]+)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])[Tt]([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]|60)(\.[0-9]+)?(([Zz])|([\+|\-]([01][0-9]|2[0-3]):[0-5][0-9]))$`),
		})

	RegisterTypeSchema(
		reflect.TypeOf(types.Duration(0)),
		jsonschema.Schema{
			ID:      types.NewStringPtr("TypesDuration"),
			Type:    NewType(jsonschema.String),
			Title:   types.NewStringPtr("Duration"),
			Format:  types.NewStringPtr("duration"),
			Pattern: types.NewStringPtr(`^(\d+(\.\d+)?h)?(\d+(\.\d+)m)?(\d+(\.\d+)?s)?(\d+(\.\d+)?ms)?(\d+(\.\d+)?us)?(\d+ns)?$`),
		})

	RegisterTypeSchema(
		reflect.TypeOf([]byte{}),
		jsonschema.Schema{
			ID:     types.NewStringPtr("TypesBinary"),
			Type:   NewType(jsonschema.String),
			Title:  types.NewStringPtr("Binary"),
			Format: types.NewStringPtr("binary"),
		})
	Reflector.AddTypeMapping([]byte{}, types.Binary{})

	RegisterTypeSchema(
		reflect.TypeOf(types.Base64Bytes{}),
		jsonschema.Schema{
			ID:     types.NewStringPtr("TypesBytes"),
			Type:   NewType(jsonschema.String),
			Title:  types.NewStringPtr("Base64"),
			Format: types.NewStringPtr("byte"),
		})

	RegisterTypeSchema(
		multipartFileHeaderType,
		jsonschema.Schema{
			ID:     types.NewStringPtr("MultipartFileHeader"),
			Type:   NewType(jsonschema.String),
			Title:  types.NewStringPtr("File"),
			Format: types.NewStringPtr("binary"),
		})
}
