package payloads

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/swaggest/jsonschema-go"
	"strings"
)

type SchemaWithRefs struct {
	Schema jsonschema.Schema
	Refs   []string
}

const PrefixComponentsSchemas = "#/components/schemas/"
const PrefixDefinitions = "#/definitions/"
const PrefixPackageApi = "Api"

var typeOverrideMap = map[string]map[string]interface{}{
	"UUID": {
		"type": "types.UUID",
		"imports": []string{
			"cto-github.cisco.com/NFV-BU/go-msx/types",
		},
	},
	"Time": {
		"type": "types.Time",
		"imports": []string{
			"cto-github.cisco.com/NFV-BU/go-msx/types",
		},
	},
	"Duration": {
		"type": "types.Duration",
		"imports": []string{
			"cto-github.cisco.com/NFV-BU/go-msx/types",
		},
	},
	"v8.Error": {
		"type": "v8.Error",
		"imports": []string{
			"cto-github.cisco.com/NFV-BU/go-msx/ops/restops/v8",
		},
	},
	"v8.PagingResponse": {
		"type": "v8.PagingResponse",
		"imports": []string{
			"cto-github.cisco.com/NFV-BU/go-msx/ops/restops/v8",
		},
	},
}

type SchemaResolver func(refName string) jsonschema.Schema

func CollectSchema(schemaName string, schema jsonschema.Schema, resolver SchemaResolver) (result jsonschema.Schema, err error) {
	var collectedRefs = types.StringSet{}
	var refsToCollect []string
	var schemaWithRefs SchemaWithRefs

	if schemaWithRefs, err = ConvertRefs(schema); err != nil {
		return
	}

	collectedRefs.Add(schemaName)
	refsToCollect = schemaWithRefs.Refs
	sp := &schemaWithRefs.Schema

	// Collect all references
	for len(refsToCollect) > 0 {
		var refSchemaWithRefs SchemaWithRefs
		var ref string

		ref, refsToCollect = refsToCollect[0], refsToCollect[1:]
		if collectedRefs.Contains(ref) {
			continue
		}

		refSchema := resolver(ref)
		refSchemaWithRefs, err = ConvertRefs(refSchema)

		collectedRefs.Add(ref)
		refsToCollect = append(refsToCollect, refSchemaWithRefs.Refs...)

		sp.WithDefinitionsItem(ref, refSchemaWithRefs.Schema.ToSchemaOrBool())
	}

	// Rename all references
	definitions := map[string]jsonschema.SchemaOrBool{}
	for k, v := range sp.Definitions {
		k = strings.TrimPrefix(k, PrefixPackageApi)
		definitions[k] = v
	}
	sp.WithDefinitions(definitions)

	WalkJsonSchema(sp, VisitJsonSchemaWhen(
		JsonSchemaHasRefPrefix(PrefixDefinitions),
		func(s *jsonschema.Schema) bool {
			schemaRefName := strings.TrimPrefix(*s.Ref, PrefixDefinitions)
			schemaRefName = strings.TrimPrefix(schemaRefName, PrefixPackageApi)
			s.Ref = types.NewStringPtr(PrefixDefinitions + schemaRefName)
			return false
		}))

	result = *sp
	return
}

func ConvertRefs(schema jsonschema.Schema) (result SchemaWithRefs, err error) {
	sp := &schema
	namedRefs := types.StringSet{}

	// Move all schema references from #/components/schemas to #/definitions
	WalkJsonSchema(sp, VisitJsonSchemaWhen(
		JsonSchemaHasRefPrefix(PrefixComponentsSchemas),
		func(s *jsonschema.Schema) bool {
			schemaRefName := strings.TrimPrefix(*s.Ref, PrefixComponentsSchemas)
			s.Ref = types.NewStringPtr(PrefixDefinitions + schemaRefName)
			return false
		}))

	// Override type handling for built-in types
	WalkJsonSchema(sp, VisitJsonSchemaWhen(
		JsonSchemaHasRefPrefix(PrefixDefinitions),
		func(s *jsonschema.Schema) bool {
			schemaRefName := strings.TrimPrefix(*s.Ref, PrefixDefinitions)
			if typeOverride, ok := typeOverrideMap[schemaRefName]; ok {
				if cur, ok := s.ExtraProperties["goJSONSchema"].(map[string]interface{}); ok && cur != nil {
					for k, v := range typeOverride {
						cur[k] = v
					}
					typeOverride = cur
				}
				s.WithExtraPropertiesItem("goJSONSchema", typeOverride)
			} else {
				namedRefs.Add(schemaRefName)
			}
			return false
		}))

	result.Schema = *sp
	result.Refs = namedRefs.Values()
	return
}
