package asyncapi

import (
	"cto-github.cisco.com/NFV-BU/go-msx/schema/js"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/swaggest/jsonschema-go"
	"reflect"
)

type asyncApiReflector struct {
	jsonschema.Reflector
	Spec *Spec
}

func (r asyncApiReflector) SpecEns() *Spec {
	if r.Spec == nil {
		r.Spec = &Spec{}
	}
	return r.Spec
}

var Reflector = asyncApiReflector{
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
	Spec: &Spec{},
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

func LookupSchema(schemaName string) (*jsonschema.Schema, bool) {
	result, ok := Reflector.SpecEns().ComponentsEns().Schemas[schemaName]
	if !ok {
		return nil, false
	}
	return &result, ok
}

func CollectDefinition(name string, schema jsonschema.Schema) {
	if _, ok := LookupSchema(name); ok {
		// Don't override
		return
	}

	Reflector.SpecEns().ComponentsEns().WithSchemasItem(name, schema)
}

func RegisterTypeSchema(t reflect.Type, schema jsonschema.Schema) {
	overrideSchema[t] = schema
	if schema.ID != nil {
		overrideSchemaById[*schema.ID] = t
	}

	Reflector.SpecEns().ComponentsEns().WithSchemasItem(*schema.ID, schema)
}

func Reflect(value interface{}) (jsonschema.Schema, error) {
	s, err := Reflector.Reflect(
		value,
		jsonschema.CollectDefinitions(CollectDefinition),
	)

	if err != nil {
		return jsonschema.Schema{}, err
	}

	return s, nil
}

func init() {
	RegisterTypeSchema(
		reflect.TypeOf(types.UUID{}),
		jsonschema.Schema{
			ID:      types.NewStringPtr("TypesUUID"),
			Type:    js.NewType(jsonschema.String),
			Title:   types.NewStringPtr("UUID"),
			Format:  types.NewStringPtr("uuid"),
			Pattern: types.NewStringPtr(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`),
		})

	RegisterTypeSchema(
		reflect.TypeOf(types.Time{}),
		jsonschema.Schema{
			ID:      types.NewStringPtr("TypesTime"),
			Type:    js.NewType(jsonschema.String),
			Title:   types.NewStringPtr("Time"),
			Format:  types.NewStringPtr("date-time"),
			Pattern: types.NewStringPtr(`^([0-9]+)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])[Tt]([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]|60)(\.[0-9]+)?(([Zz])|([\+|\-]([01][0-9]|2[0-3]):[0-5][0-9]))$`),
		})

	RegisterTypeSchema(
		reflect.TypeOf(types.Duration(0)),
		jsonschema.Schema{
			ID:      types.NewStringPtr("TypesDuration"),
			Type:    js.NewType(jsonschema.String),
			Title:   types.NewStringPtr("Duration"),
			Format:  types.NewStringPtr("duration"),
			Pattern: types.NewStringPtr(`^(\d+(\.\d+)?h)?(\d+(\.\d+)m)?(\d+(\.\d+)?s)?(\d+(\.\d+)?ms)?(\d+(\.\d+)?us)?(\d+ns)?$`),
		})

	RegisterTypeSchema(
		reflect.TypeOf([]byte{}),
		jsonschema.Schema{
			ID:     types.NewStringPtr("TypesBinary"),
			Type:   js.NewType(jsonschema.String),
			Title:  types.NewStringPtr("Binary"),
			Format: types.NewStringPtr("binary"),
		})

	RegisterTypeSchema(
		reflect.TypeOf(types.Base64Bytes{}),
		jsonschema.Schema{
			ID:     types.NewStringPtr("TypesBytes"),
			Type:   js.NewType(jsonschema.String),
			Title:  types.NewStringPtr("Base64"),
			Format: types.NewStringPtr("byte"),
		})
}
