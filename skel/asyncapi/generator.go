// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package asyncapi

import (
	"cto-github.cisco.com/NFV-BU/go-msx/schema/asyncapi"
	"cto-github.cisco.com/NFV-BU/go-msx/schema/js"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/mcrawfo2/go-jsonschema/pkg/generator"
	"github.com/pkg/errors"
	"github.com/swaggest/jsonschema-go"
	"os"
	"path/filepath"
	"strings"
)

const PrefixComponentsMessages = "#/components/messages/"
const PrefixComponentsSchemas = "#/components/schemas/"
const PrefixDefinitions = "#/definitions/"
const PrefixPackageApi = "Api"

const OperationSubscribe = "subscribe"
const OperationPublish = "publish"

type Generator struct {
	spec     asyncapi.Spec
	packages map[string]Package
}

func (g *Generator) setPackage(prefix string, packageName string, folder string) {
	g.packages[prefix] = Package{
		Name:   packageName,
		Folder: folder,
	}
}

func (g *Generator) GenerateChannels(channels ...string) (err error) {
	if len(channels) == 0 {
		for k := range g.spec.Channels {
			channels = append(channels, k)
		}
	}

	for _, k := range channels {
		if err = g.GenerateChannel(k); err != nil {
			return
		}
	}

	return
}

func (g *Generator) messages(choices asyncapi.MessageChoices) (results []*asyncapi.Message) {
	if choices.Reference != nil {
		ref := strings.TrimPrefix(choices.Reference.Ref, PrefixComponentsMessages)
		choices = g.spec.ComponentsEns().Messages[ref]
		return g.messages(choices)
	} else if choices.Message != nil {
		results = append(results, choices.Message)
	} else if choices.MessageOptions != nil {
		for _, option := range choices.MessageOptions.OneOf {
			results = append(results, g.messages(option)...)
		}
	}
	return
}

func (g *Generator) GenerateChannel(k string) (err error) {
	streamFolderName := g.generateFolderName(k)
	g.setPackage("", streamFolderName, "internal/stream/"+streamFolderName)
	g.setPackage(PrefixPackageApi, streamFolderName, "internal/stream/"+streamFolderName)

	gen, _ := generator.New(generator.Config{
		Capitalizations: nil,
		Warner: func(s string) {
			logger.Warn(s)
		},
	})

	channel := g.spec.Channels[k]

	logger.Infof("Generating components for channel %s", k)

	// TODO: generate channel
	logger.Warnf("Channel generation not implemented: %s", k)

	if channel.Publish != nil {
		// TODO: generate channel publisher
		logger.Warnf("Channel publisher generation not implemented: %s", k)

		messages := g.messages(*channel.Publish.MessageEns())
		for n, message := range messages {
			if message.ID == nil {
				message.ID = g.generateMessageId(k, OperationPublish, n)
			}

			logger.Infof("Generating publisher message %s", *message.ID)

			// generate payload dto
			if err = g.GeneratePayload(message, gen); err != nil {
				logger.WithError(err).Errorf("Failed to generate message payload DTO: %s", message.ID)
			}

			// generate message publisher output port
			if err = g.GeneratePortStruct(message, OperationSubscribe, gen); err != nil {
				logger.WithError(err).Errorf("Failed to generate message port struct: %s", *message.ID)
			}

			// TODO: generate message publisher
		}
	}

	if channel.Subscribe != nil {
		messages := g.messages(*channel.Subscribe.MessageEns())
		for n, message := range messages {
			if message.ID == nil {
				message.ID = g.generateMessageId(k, OperationSubscribe, n)
			}

			logger.Infof("Generating subscriber message %s", *message.ID)

			// generate payload dto
			if err = g.GeneratePayload(message, gen); err != nil {
				logger.WithError(err).Errorf("Failed to generate message payload DTO: %s", *message.ID)
			}

			// generate message subscriber input port
			if err = g.GeneratePortStruct(message, OperationSubscribe, gen); err != nil {
				logger.WithError(err).Errorf("Failed to generate message port struct: %s", *message.ID)
			}

			// TODO: generate message subscriber
			logger.Warnf("Message subscriber generation not implemented: %s", *message.ID)
		}

		// TODO: generate channel subscriber
		logger.Warnf("Channel subscriber generation not implemented: %s", k)
	}

	return nil
}

func (g *Generator) GeneratePayload(message *asyncapi.Message, gen *generator.Generator) error {
	if message.Payload == nil {
		return nil
	}

	logger.Infof("Generating message %s payload", message.ID)

	rawSchema := (*message.Payload).(map[string]interface{})
	schema, err := g.resolveSchema(rawSchema)
	if err != nil {
		return err
	}

	if schema.ID == nil {
		schema.ID = message.ID
	}

	return g.GenerateType(schema, gen)
}

func (g *Generator) GeneratePortStruct(message *asyncapi.Message, opType string, gen *generator.Generator) error {
	portSchema := js.ObjectSchema()
	directionTag := ""
	if opType == OperationSubscribe {
		portSchema.ID = types.NewStringPtr(*message.ID + "Input")
		directionTag = "in"
	} else {
		portSchema.ID = types.NewStringPtr(*message.ID + "Output")
		directionTag = "out"
	}

	headerSchema := message.HeadersEns().Schema
	if headerSchema != nil {
		portSchema.WithProperties(headerSchema.Properties)
		for k, v := range portSchema.Properties {
			s := v.TypeObject
			s.WithExtraPropertiesItem("goJSONSchema", map[string]interface{}{
				"tags": map[string]string{
					"json":       "",
					directionTag: "header=" + k,
				},
			})
		}
	}

	rawSchema := (*message.Payload).(map[string]interface{})
	payloadSchema, err := g.convertSchema(rawSchema)
	if err != nil {
		return err
	}

	payloadSchema.WithExtraPropertiesItem("goJSONSchema", map[string]interface{}{
		"tags": map[string]string{
			"json":       "",
			directionTag: "body",
		},
	})

	portSchema.WithPropertiesItem("payload", payloadSchema.ToSchemaOrBool())
	portSchema.Required = append([]string{}, headerSchema.Required...)
	portSchema.Required = append(portSchema.Required, "payload")

	return g.GenerateType(*portSchema, gen)
}

func (g *Generator) GenerateType(schema jsonschema.Schema, gen *generator.Generator) error {
	schemaName := *schema.ID
	packageName := ""
	packageFolder := ""

	for prefix, pkg := range g.packages {
		if prefix == "" {
			continue
		}

		if strings.HasPrefix(schemaName, prefix) {
			schemaName = strings.TrimPrefix(schemaName, prefix)
			packageName = pkg.Name
			packageFolder = pkg.Folder
		}
	}

	if packageName == "" {
		pkg, ok := g.packages[""]
		if ok {
			packageName = pkg.Name
			packageFolder = pkg.Folder
		} else {
			return errors.Errorf("Failed to locate output package for schema %s", schemaName)
		}
	}

	logger.Infof("Generating type for schema %s.%s", packageName, schemaName)

	schema, err := g.collectSchema(schemaName, schema)
	if err != nil {
		return err
	}

	// convert to JSON
	schemaBytes, err := json.Marshal(schema)
	if err != nil {
		return err
	}

	// add the spec to the generator
	err = gen.AddSource(generator.Source{
		PackageName: packageName,
		RootType:    schemaName,
		Folder:      packageFolder,
		Data:        schemaBytes,
	})
	if err != nil {
		return err
	}

	for fileName, source := range gen.Outputs() {
		logger.Infof("Generated file %s", fileName)
		_ = os.MkdirAll(filepath.Dir(fileName), 0755)
		_ = os.WriteFile(fileName, source, 0644)
	}

	return nil
}

func (g *Generator) convertSchema(schema map[string]interface{}) (result jsonschema.Schema, err error) {
	// serialize and deserialize json into jsonschema.Schema
	schemaBytes, err := json.Marshal(schema)
	if err != nil {
		return
	}

	err = json.Unmarshal(schemaBytes, &result)
	if err != nil {
		return
	}

	return
}

func (g *Generator) resolveSchema(schema map[string]interface{}) (result jsonschema.Schema, err error) {
	if ref, ok := schema["$ref"]; ok {
		refName := strings.TrimPrefix(ref.(string), PrefixComponentsSchemas)
		result = g.spec.ComponentsEns().Schemas[refName]
		if result.ID == nil {
			result.ID = &refName
		}

		return
	}

	return g.convertSchema(schema)
}

func (g *Generator) collectSchema(schemaName string, schema jsonschema.Schema) (result jsonschema.Schema, err error) {
	var collectedRefs = types.StringSet{}
	var refsToCollect []string
	var schemaWithRefs SchemaWithRefs

	if schemaWithRefs, err = g.convertRefs(schemaName, schema); err != nil {
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

		refSchema := g.spec.ComponentsEns().Schemas[ref]
		refSchemaWithRefs, err = g.convertRefs(ref, refSchema)

		collectedRefs.Add(ref)
		refsToCollect = append(refsToCollect, refSchemaWithRefs.Refs...)

		sp.WithDefinitionsItem(ref, refSchemaWithRefs.Schema.ToSchemaOrBool())
	}

	// Rename all references
	definitions := map[string]jsonschema.SchemaOrBool{}
	for k, v := range sp.Definitions {
		k = strings.TrimPrefix(k, "Api")
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
}

func (g *Generator) convertRefs(refName string, schema jsonschema.Schema) (result SchemaWithRefs, err error) {
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

func (g *Generator) generateMessageId(channelName string, operation string, index int) *string {
	// clean channel name
	channelName = strcase.ToCamel(strings.ToLower(strings.TrimSuffix(channelName, "_TOPIC")))
	operationType := "Request"
	if operation == OperationSubscribe {
		operationType = "Response"
	}
	suffix := ""
	if index > 0 {
		suffix = fmt.Sprintf("%d", index)
	}
	return types.NewStringPtr(fmt.Sprintf("%s%s%s", channelName, operationType, suffix))
}

func (g *Generator) generateFolderName(channelName string) string {
	channelName = strings.TrimSuffix(channelName, "_TOPIC")
	return strings.ReplaceAll(strcase.ToSnake(channelName), "_", "")
}

func NewGenerator(spec asyncapi.Spec) *Generator {
	return &Generator{
		spec:     spec,
		packages: map[string]Package{},
	}
}

type SchemaWithRefs struct {
	Schema jsonschema.Schema
	Refs   []string
}

type Package struct {
	Name   string
	Folder string
}
