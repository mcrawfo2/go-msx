// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package asyncapi

import (
	"cto-github.cisco.com/NFV-BU/go-msx/schema/asyncapi"
	"cto-github.cisco.com/NFV-BU/go-msx/schema/js"
	"cto-github.cisco.com/NFV-BU/go-msx/skel"
	"cto-github.cisco.com/NFV-BU/go-msx/skel/payloads"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/mcrawfo2/go-jsonschema/pkg/codegen"
	"github.com/mcrawfo2/go-jsonschema/pkg/generator"
	"github.com/mcrawfo2/jennifer/jen"
	"github.com/pkg/errors"
	"github.com/swaggest/jsonschema-go"
	"os"
	"path/filepath"
	"strings"
)

const PrefixComponentsMessages = "#/components/messages/"
const PrefixComponentsSchemas = "#/components/schemas/"
const PrefixPackageApi = "Api"

const OperationSubscribe = "subscribe"
const OperationPublish = "publish"
const OperationNone = ""

type Generator struct {
	spec     asyncapi.Spec
	opMap    map[string]string
	packages map[string]Package
	dirs     types.StringSet
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
	skeletonConfig := skel.Config()

	streamFolderName := "internal/stream/" + g.generateFolderName(k)
	g.setPackage("", skeletonConfig.AppPackageUrl()+streamFolderName, streamFolderName)

	apiFolderName := streamFolderName + "/api"
	g.setPackage(PrefixPackageApi, skeletonConfig.AppPackageUrl()+"/"+apiFolderName, apiFolderName)

	gen, _ := generator.New(generator.Config{
		Capitalizations: nil,
		Warner: func(s string) {
			logger.Warn(s)
		},
	})

	channel := g.spec.Channels[k]

	logger.Infof("Generating components for channel %s", k)

	// generate channel
	if err = g.generateChannelImpl(k); err != nil {
		logger.WithError(err).Errorf("Failed to generate channel %s", k)
	}

	if channel.Publish != nil {
		op := g.opMap[OperationPublish]

		// generate channel publisher
		if err = g.generateChannelOperationImpl(k, op, false); err != nil {
			logger.WithError(err).Errorf("Failed to generate channel %s", k)
		}

		messages := g.messages(*channel.Publish.MessageEns())
		for n, message := range messages {
			if message.ID == nil {
				message.ID = g.generateMessageId(k, op, n)
			}

			logger.Infof("Generating subscriber message %s", *message.ID)

			// generate payload dto
			if err = g.GeneratePayload(message, gen); err != nil {
				logger.WithError(err).Errorf("Failed to generate message payload DTO: %s", *message.ID)
			}

			// generate message subscriber output port
			var vars map[string]string
			if vars, err = g.GeneratePortStructSnippets(message, op, gen); err != nil {
				logger.WithError(err).Errorf("Failed to generate message port struct: %s", *message.ID)
			}

			// generate message subscriber
			if err = g.generateMessageOperationImpl(k, op, len(messages) > 1, *message.ID, vars); err != nil {
				logger.WithError(err).Errorf("Failed to generate message subscriber: %s", *message.ID)
			}
		}
	}

	if channel.Subscribe != nil {
		op := g.opMap[OperationSubscribe]

		messages := g.messages(*channel.Subscribe.MessageEns())

		// generate channel subscriber
		if err = g.generateChannelOperationImpl(k, op, len(messages) > 1); err != nil {
			logger.WithError(err).Errorf("Failed to generate channel %s", k)
		}

		for n, message := range messages {
			if message.ID == nil {
				message.ID = g.generateMessageId(k, op, n)
			}

			logger.Infof("Generating publisher message %s", *message.ID)

			// generate payload dto
			if err = g.GeneratePayload(message, gen); err != nil {
				logger.WithError(err).Errorf("Failed to generate message payload DTO: %s", *message.ID)
			}

			// generate message publisher input port
			var vars map[string]string
			if vars, err = g.GeneratePortStructSnippets(message, op, gen); err != nil {
				logger.WithError(err).Errorf("Failed to generate message port struct: %s", *message.ID)
			}

			// generate message publisher
			if err = g.generateMessageOperationImpl(k, op, len(messages) > 1, *message.ID, vars); err != nil {
				logger.WithError(err).Errorf("Failed to generate message publisher: %s", *message.ID)
			}
		}
	}

	return GoGenerate(g.dirs.Values())
}

func (g *Generator) generateChannelImpl(channelName string) error {
	cfg := GeneratorConfig{
		ChannelName: channelName,
		Deep:        false,
	}

	gen := newTemplatedGenerator(cfg)
	if err := gen.RenderChannel(); err != nil {
		return err
	}

	g.dirs.AddAll(gen.Dirs()...)
	return nil
}

func (g *Generator) generateChannelOperationImpl(channelName string, op string, multi bool) error {
	cfg := GeneratorConfig{
		ChannelName: channelName,
		Deep:        false,
		Multi:       multi,
	}

	gen := newTemplatedGenerator(cfg)
	if err := gen.RenderChannelOperation(op); err != nil {
		return err
	}

	g.dirs.AddAll(gen.Dirs()...)
	return nil
}

func (g *Generator) generateMessageOperationImpl(channelName string, op string, multi bool, messageId string, vars map[string]string) error {
	cfg := GeneratorConfig{
		ChannelName: channelName,
		Multi:       multi,
		Deep:        false,
		Domain:      "Unknown",
	}

	gen := newTemplatedGenerator(cfg)
	if err := gen.RenderMessageOperation(op, messageId, vars); err != nil {
		return err
	}

	g.dirs.AddAll(gen.Dirs()...)
	return nil
}

func (g *Generator) GeneratePayload(message *asyncapi.Message, gen *generator.Generator) error {
	if message.Payload == nil {
		return nil
	}

	logger.Infof("Generating message %s payload", *message.ID)

	rawSchema := (*message.Payload).(map[string]interface{})
	schema, err := g.resolveSchema(rawSchema)
	if err != nil {
		return err
	}

	if schema.ID == nil {
		schema.ID = message.ID
	}

	files, err := g.GenerateType(schema, gen)
	if err != nil {
		return err
	}

	return g.WriteFiles(files)
}

type EmitterWriter struct {
	*codegen.Emitter
}

func (e *EmitterWriter) Write(p []byte) (n int, err error) {
	e.Emitter.Print("%s", string(p))
	return len(p), nil
}

func (g *Generator) GeneratePortStructSnippets(message *asyncapi.Message, op string, gen *generator.Generator) (map[string]string, error) {
	portSchema, err := g.GeneratePortStructSchema(message, op)
	if err != nil {
		return nil, err
	}

	files, err := g.GenerateType(*portSchema, gen)
	if err != nil {
		return nil, err
	}
	portStructFile := files[0]

	var portStruct *codegen.StructType
	for _, decl := range portStructFile.Package.Decls {
		if typeDecl, ok := decl.(*codegen.TypeDecl); ok {
			if typeDecl.GetName() == *portSchema.ID {
				portStruct = typeDecl.Type.(*codegen.StructType)
				break
			}
		}
	}

	var implEmitter = codegen.NewEmitter(80)
	var depEmitter = codegen.NewEmitter(80)
	var handlerEmitter = codegen.NewEmitter(80)

	imports := g.GeneratePortStructCode(implEmitter, portStructFile.Package)

	switch op {
	case OperationPublish:
		if err = g.GeneratePublisherInterfaceCode(depEmitter, *message.ID, portStruct); err != nil {
			return nil, err
		}
		if err = g.GeneratePublisherImplementationCode(implEmitter, *portSchema.ID, *message.ID, portStruct); err != nil {
			return nil, err
		}

	case OperationSubscribe:
		if err = g.GenerateSubscriberHandlerInterfaceCode(depEmitter, *message.ID, portStruct); err != nil {
			return nil, err
		}
		if err = g.GenerateSubscriberHandlerCode(handlerEmitter, *portSchema.ID, *message.ID, portStruct); err != nil {
			return nil, err
		}
	}

	return map[string]string{
		"imports":        strings.Join(imports, "\n"),
		"dependencies":   depEmitter.String(),
		"implementation": implEmitter.String(),
		"handler":        strings.TrimPrefix(strings.TrimSpace(handlerEmitter.String()), "var _ = "),
	}, nil
}

func (g *Generator) GeneratePortStructSchema(message *asyncapi.Message, op string) (*jsonschema.Schema, error) {
	directionTag := ""
	structTypeName := ""
	switch op {
	case OperationSubscribe:
		structTypeName = strcase.ToLowerCamel(*message.ID) + "Input"
		directionTag = "in"
	case OperationPublish:
		structTypeName = strcase.ToLowerCamel(*message.ID) + "Output"
		directionTag = "out"
	}

	portSchema := js.ObjectSchema()
	portSchema.ID = types.NewStringPtr(structTypeName)

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
		return nil, err
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

	return portSchema, nil
}

func (g *Generator) GeneratePortStructCode(emitter *codegen.Emitter, pkg codegen.Package) []string {
	// inject imports into message operation template
	var imports []string
	for _, imp := range pkg.Imports {
		imports = append(imports, fmt.Sprintf(`%s %q`, imp.Name, imp.QualifiedName))
	}

	// generate port struct concrete code
	for _, decl := range pkg.Decls {
		decl.Generate(emitter)
		emitter.Println("")
	}

	return imports
}

func (g *Generator) GenerateSubscriberHandlerInterfaceCode(emitter *codegen.Emitter, messageId string, portStruct *codegen.StructType) (err error) {
	interfaceName := strcase.ToCamel(messageId) + "Handler"

	params := []jen.Code{
		jen.Id("ctx").Qual("context", "Context"),
	}

	for _, structField := range portStruct.Fields {
		if !structField.Tags.HasTag("in") {
			continue
		}
		if structField.Tags.HasTag("const") {
			continue
		}

		typeEmitter := codegen.NewEmitter(80)
		structField.Type.Generate(typeEmitter)

		params = append(params, jen.Id(strcase.ToLowerCamel(structField.Name)).Op(typeEmitter.String()))
	}

	structName := "drop" + messageId + "Handler"

	stmt := jen.Statement{
		jen.
			Line().Type().Id(interfaceName).
			Interface(
				jen.Id("On" + messageId).Call(params...).Error(),
			).
			Line(),
		jen.Line().Type().Id(structName).
			Struct().
			Line(),
		jen.Line().Func().
			Params(jen.Id("_").Id(structName)).
			Id("On"+messageId).Params(params...).Error().Block(
			jen.Id("logger").Dot("Error").Call(
				jen.Lit(fmt.Sprintf(
					"No handler assigned to %s message subscription.  Dropping message.",
					messageId))),
			jen.Return(jen.Nil()),
		).
			Line(),
	}

	return stmt.Render(&EmitterWriter{Emitter: emitter})
}

func (g *Generator) GenerateSubscriberHandlerCode(emitter *codegen.Emitter, portStructName string, messageId string, portStruct *codegen.StructType) (err error) {
	// Generate handler implementation
	// func(ctx context.Context, in *${async.msgtype}Input) error {
	//				return handler.On${async.upmsgtype}(ctx, in.Payload)
	//			}

	params := []jen.Code{
		jen.Id("ctx"),
	}

	for _, structField := range portStruct.Fields {
		if !structField.Tags.HasTag("in") {
			continue
		}
		if structField.Tags.HasTag("const") {
			continue
		}

		params = append(params, jen.Id("in").Dot(structField.Name))
	}

	stmt := jen.Var().Id("_").Op("=").
		Func().
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("in").Op("*").Id(portStructName),
		).
		Error().
		Block(
			jen.Return(
				jen.Id("handler").Dot("On" + messageId).Call(params...),
			),
		)

	return stmt.Render(&EmitterWriter{Emitter: emitter})
}

func (g *Generator) GeneratePublisherInterfaceCode(emitter *codegen.Emitter, messageId string, portStruct *codegen.StructType) (err error) {
	interfaceName := strcase.ToCamel(messageId) + "Publisher"

	params := []jen.Code{
		jen.Id("ctx").Qual("context", "Context"),
	}

	for _, structField := range portStruct.Fields {
		if !structField.Tags.HasTag("out") {
			continue
		}
		if structField.Tags.HasTag("const") {
			continue
		}

		typeEmitter := codegen.NewEmitter(80)
		structField.Type.Generate(typeEmitter)

		params = append(params, jen.Id(strcase.ToLowerCamel(structField.Name)).Op(typeEmitter.String()))
	}

	// Generate publisher interface
	stmt := jen.
		Line().Type().Id(interfaceName).
		Interface(
			jen.Id("Publish" + messageId).Call(params...).Error(),
		).
		Line()

	return stmt.Render(&EmitterWriter{Emitter: emitter})
}

func (g *Generator) GeneratePublisherImplementationCode(emitter *codegen.Emitter, portStructName, messageId string, portStruct *codegen.StructType) (err error) {
	implementationName := strcase.ToLowerCamel(messageId) + "Publisher"

	params := []jen.Code{
		jen.Id("ctx").Qual("context", "Context"),
	}

	fields := jen.Dict{}

	for _, structField := range portStruct.Fields {
		if !structField.Tags.HasTag("out") {
			continue
		}
		if structField.Tags.HasTag("const") {
			continue
		}

		typeEmitter := codegen.NewEmitter(80)
		structField.Type.Generate(typeEmitter)

		varName := strcase.ToLowerCamel(structField.Name)

		params = append(params, jen.Id(varName).Op(typeEmitter.String()))
		fields[jen.Id(structField.Name)] = jen.Id(varName)
	}

	stmt := jen.Line().
		Func().
		Params(jen.Id("p").Id(implementationName)).
		Id("Publish" + messageId).
		Params(params...).
		Error().
		Block(
			jen.Return(jen.Id("p").Dot("messagePublisher").Dot("Publish").Call(
				jen.Id("ctx"),
				jen.Id(portStructName).Values(fields),
			)),
		).
		Line()
	return stmt.Render(&EmitterWriter{Emitter: emitter})
}

func (g *Generator) WriteFiles(files []codegen.File) (err error) {
	for _, file := range files {
		fileName := file.FileName

		emitter := codegen.NewEmitter(80)
		file.Generate(emitter)
		source := emitter.Bytes()

		logger.Infof("Generated file %s", fileName)
		if err = os.MkdirAll(filepath.Dir(fileName), 0755); err != nil {
			return
		}
		if err = os.WriteFile(fileName, source, 0644); err != nil {
			return
		}
	}
	return
}

func (g *Generator) GenerateType(schema jsonschema.Schema, gen *generator.Generator) ([]codegen.File, error) {
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
			break
		}
	}

	if packageName == "" {
		pkg, ok := g.packages[""]
		if ok {
			packageName = pkg.Name
			packageFolder = pkg.Folder
		} else {
			return nil, errors.Errorf("Failed to locate output package for schema %s", schemaName)
		}
	}

	logger.Infof("Generating type for schema %s.%s", packageName, schemaName)

	schema, err := g.collectSchema(schemaName, schema)
	if err != nil {
		return nil, err
	}

	// convert to JSON
	schemaBytes, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}

	// add the spec to the generator
	err = gen.AddSource(generator.Source{
		PackageName: packageName,
		RootType:    schemaName,
		Folder:      packageFolder,
		Data:        schemaBytes,
	})
	if err != nil {
		return nil, err
	}

	return gen.Files(), nil
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
		result = g.resolveSchemaRefByName(refName)
		if result.ID == nil {
			result.ID = &refName
		}

		return
	}

	return g.convertSchema(schema)
}

func (g *Generator) resolveSchemaRefByName(ref string) jsonschema.Schema {
	return g.spec.ComponentsEns().Schemas[ref]
}

func (g *Generator) collectSchema(schemaName string, schema jsonschema.Schema) (result jsonschema.Schema, err error) {
	return payloads.CollectSchema(schemaName, schema, g.resolveSchemaRefByName)
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

func (g *Generator) WithInvert(invert bool) *Generator {
	if invert {
		g.opMap[OperationSubscribe] = OperationPublish
		g.opMap[OperationPublish] = OperationSubscribe
	}
	return g
}

func NewGenerator(spec asyncapi.Spec) *Generator {
	return &Generator{
		spec: spec,
		opMap: map[string]string{
			OperationSubscribe: OperationSubscribe,
			OperationPublish:   OperationPublish,
			OperationNone:      OperationNone,
		},
		packages: map[string]Package{},
		dirs:     make(types.StringSet),
	}
}

type Package struct {
	Name   string
	Folder string
}
