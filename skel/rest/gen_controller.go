// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package rest

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops/restops"
	"cto-github.cisco.com/NFV-BU/go-msx/schema/js"
	"cto-github.cisco.com/NFV-BU/go-msx/skel"
	"cto-github.cisco.com/NFV-BU/go-msx/skel/payloads"
	"cto-github.cisco.com/NFV-BU/go-msx/skel/text"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/mcrawfo2/go-jsonschema/pkg/codegen"
	"github.com/mcrawfo2/go-jsonschema/pkg/generator"
	"github.com/mcrawfo2/jennifer/jen"
	"github.com/swaggest/jsonschema-go"
	"github.com/swaggest/openapi-go/openapi3"
	"path"
	"sort"
	"strings"
)

type DomainControllerGenerator struct {
	Domain     string
	Folder     string
	Style      string
	Tenant     string
	Actions    types.ComparableSlice[string]
	Components []string

	Spec Spec

	*text.GoFile
}

func (g DomainControllerGenerator) createConstantsSnippet() error {
	constants := text.GoConstants{
		{
			Name:  "pathSuffixUpperCamelSingularId",
			Value: "{lowerCamelSingularId}",
		},
	}

	if g.Actions.ContainsAny(ActionList, ActionRetrieve) {
		constants = append(constants, &codegen.Constant{
			Name:  "permissionViewUpperCamelPlural",
			Value: g.GoFile.Inflector.Inflect("VIEW_SCREAMING_SNAKE_PLURAL"),
		})
	}

	if g.Actions.ContainsAny(ActionCreate, ActionUpdate, ActionDelete) {
		constants = append(constants, &codegen.Constant{
			Name:  "permissionManageUpperCamelPlural",
			Value: g.GoFile.Inflector.Inflect("MANAGE_SCREAMING_SNAKE_PLURAL"),
		})
	}

	return g.AddNewGenerator(
		"Constants",
		"pathConstants",
		constants,
		[]codegen.Import{})
}

func (g DomainControllerGenerator) createControllerSnippet() error {
	return g.AddNewText(
		"Controller",
		"controller",
		`
			// lowerCamelSingularController implements the UpperCamelSingular API.
			type lowerCamelSingularController struct {
				lowerCamelSingularService UpperCamelSingularServiceApi
			}
		`,
		[]codegen.Import{})
}

func (g DomainControllerGenerator) createEndpointsConstructionSnippet(methods []string) error {
	sort.Strings(methods)

	var operations = new(strings.Builder)
	for _, method := range methods {
		err := jen.Id("c").Dot(method).Call().Render(operations)
		if err != nil {
			return err
		}
		operations.WriteString(",\n")
	}

	renderOptions := skel.NewEmptyRenderOptions()
	renderOptions.AddVariables(map[string]string{
		"operations": operations.String(),
	})

	template, err := skel.Template{
		SourceData: []byte(`
			// Endpoints provides this controller's endpoints to the EndpointRegisterer.
			func (c *lowerCamelSingularController) Endpoints() (restops.Endpoints, error) {
				builders := restops.EndpointBuilders{
					${operations}
				}
			
				return builders.Endpoints()
			}
			`),
	}.RenderContents(renderOptions)
	if err != nil {
		return err
	}

	return g.AddNewText(
		"Endpoint Construction",
		"endpoints",
		template,
		[]codegen.Import{
			text.ImportRestOps,
		})
}

func (g DomainControllerGenerator) createEndpointTransformationSnippet() error {
	return g.AddNewText(
		"Endpoint Transformation",
		"endpointTransformers",
		`
			// EndpointTransformers provides a set of transformations to be applied to each Endpoint
			// created by this controller.
			func (c *lowerCamelSingularController) EndpointTransformers() restops.EndpointTransformers {
				const tagName = "UpperCamelSingular"
				const pathPrefix = "/${domain.style}/lowerplural"

				openapi.AddTag(tagName, "Title Plural")
			
				return restops.EndpointTransformers{
					restops.AddEndpointPathPrefix(pathPrefix),
					restops.AddEndpointTag(tagName),
				}
			}
			`,
		[]codegen.Import{
			text.ImportRestOps,
			text.ImportOpenapi,
		})
}

func (g DomainControllerGenerator) cleanOperationId(operationId string) string {
	cleanOperationId := operationId
	if strings.Contains(cleanOperationId, ".") {
		lastPeriod := strings.LastIndex(cleanOperationId, ".")
		cleanOperationId = cleanOperationId[lastPeriod+1:]
	}
	return cleanOperationId
}

func (g DomainControllerGenerator) createEndpointActionListSnippet(operation Operation) error {
	renderOptions := skel.NewEmptyRenderOptions()
	renderOptions.AddVariables(map[string]string{
		"snippet.outputs.content.tag": "`resp:\"body\"`",
		"operation.id":                *operation.Operation.ID,
		"operation.id.clean":          g.cleanOperationId(*operation.Operation.ID),
		"operation.summary":           types.OptionalOfPtr(operation.Operation.Summary).OrElse(""),
	})

	template, err := skel.Template{
		SourceData: []byte(`
			// lowerCamelSingularFilterQueryInputs is used to declare the query string filters
			// for the ${operation.id.clean} endpoint.
			type lowerCamelSingularFilterQueryInputs struct {
			}

			// ${operation.id.clean} creates an endpoint providing a filtered, sorted, and paginated 
			// sequence of UpperCamelSingular instances.
			func (c *lowerCamelSingularController) ${operation.id.clean}() restops.EndpointBuilder {
				type inputs struct {
					${domain.style}.PagingSortingInputs
					lowerCamelSingularFilterQueryInputs
				}

				type outputs struct {
					${domain.style}.PagingOutputs
					Content []UpperCamelSingularResponse ${snippet.outputs.content.tag}
				}

				return ${domain.style}.
					NewListEndpointBuilder().
					WithId("${operation.id}").
					WithDoc(new(openapi3.Operation).
						WithSummary("${operation.summary}")).
					WithPermissions(permissionViewUpperCamelPlural).
					WithHandler(
						func(ctx context.Context, inp *inputs) (out outputs, err error) {
							out.PagingOutputs.Paging, out.Content, err = c.lowerCamelSingularService.ListUpperCamelPlural(
								ctx, inp.PagingSortingInputs, inp.lowerCamelSingularFilterQueryInputs)
							return
						})
			}
			`),
	}.RenderContents(renderOptions)
	if err != nil {
		return err
	}

	return g.AddNewText(
		"Endpoints/List",
		"list",
		template,
		[]codegen.Import{
			text.ImportRestOps,
			g.importStyle(),
			text.ImportOpenApi3,
			text.ImportContext,
		})
}

func (g DomainControllerGenerator) createEndpointActionRetrieveSnippet(operation Operation) error {
	portStructVariables, err := g.generatePortStructVariables(operation)
	if err != nil {
		return err
	}

	renderOptions := skel.NewEmptyRenderOptions()
	renderOptions.AddVariables(portStructVariables)
	renderOptions.AddVariables(map[string]string{
		"operation.id":       *operation.Operation.ID,
		"operation.id.clean": g.cleanOperationId(*operation.Operation.ID),
		"operation.summary":  types.OptionalOfPtr(operation.Operation.Summary).OrElse(""),
	})

	template, err := skel.Template{
		SourceData: []byte(`
			// ${operation.id.clean} creates an endpoint providing a single UpperCamelSingular instance
			// having the specified key.
			func (c *lowerCamelSingularController) ${operation.id.clean}() restops.EndpointBuilder {
				${snippet.inputs}
				${snippet.outputs}

				return ${domain.style}.
					NewRetrieveEndpointBuilder(pathSuffixUpperCamelSingularId).
					WithId("${operation.id}").
					WithDoc(new(openapi3.Operation).
						WithSummary("${operation.summary}")).
					WithPermissions(permissionViewUpperCamelPlural).
					WithHandler(
						func(ctx context.Context, inp *inputs) (out outputs, err error) {
							out.Body, err = c.lowerCamelSingularService.GetUpperCamelSingular(ctx, inp.UpperCamelSingularId)
							return
						})
			}
			`),
	}.RenderContents(renderOptions)
	if err != nil {
		return err
	}

	return g.AddNewText(
		"Endpoints/Retrieve",
		"retrieve",
		template,
		[]codegen.Import{
			text.ImportRestOps,
			g.importStyle(),
			text.ImportOpenApi3,
			text.ImportContext,
		})
}

func (g DomainControllerGenerator) createEndpointActionCreateSnippet(operation Operation) error {
	portStructVariables, err := g.generatePortStructVariables(operation)
	if err != nil {
		return err
	}

	renderOptions := skel.NewEmptyRenderOptions()
	renderOptions.AddVariables(portStructVariables)
	renderOptions.AddVariables(map[string]string{
		"operation.id":       *operation.Operation.ID,
		"operation.id.clean": g.cleanOperationId(*operation.Operation.ID),
		"operation.summary":  types.OptionalOfPtr(operation.Operation.Summary).OrElse(""),
	})

	template, err := skel.Template{
		SourceData: []byte(`
			// ${operation.id.clean} creates an endpoint instantiating a new UpperCamelSingular instance
			// using the specified values.
			func (c *lowerCamelSingularController) ${operation.id.clean}() restops.EndpointBuilder {
				${snippet.inputs}
				${snippet.outputs}

				return ${domain.style}.
					NewCreateEndpointBuilder().
					WithId("${operation.id}").
					WithDoc(new(openapi3.Operation).
						WithSummary("${operation.summary}")).
					WithPermissions(permissionManageUpperCamelPlural).
					WithHandler(
						func(ctx context.Context, inp *inputs) (out outputs, err error) {
							out.Body, err = c.lowerCamelSingularService.CreateUpperCamelSingular(ctx, inp.Body)
							return
						})
			}
			`),
	}.RenderContents(renderOptions)
	if err != nil {
		return err
	}

	return g.AddNewText(
		"Endpoints/Create",
		"create",
		template,
		[]codegen.Import{
			text.ImportRestOps,
			g.importStyle(),
			text.ImportOpenApi3,
			text.ImportContext,
		})
}

func (g DomainControllerGenerator) createEndpointActionUpdateSnippet(operation Operation) error {
	portStructVariables, err := g.generatePortStructVariables(operation)
	if err != nil {
		return err
	}

	renderOptions := skel.NewEmptyRenderOptions()
	renderOptions.AddVariables(portStructVariables)
	renderOptions.AddVariables(map[string]string{
		"operation.id":       *operation.Operation.ID,
		"operation.id.clean": g.cleanOperationId(*operation.Operation.ID),
		"operation.summary":  types.OptionalOfPtr(operation.Operation.Summary).OrElse(""),
	})

	template, err := skel.Template{
		SourceData: []byte(`
			// ${operation.id.clean} creates an endpoint updating an existing UpperCamelSingular instance
			// using the specified values.
			func (c *lowerCamelSingularController) ${operation.id.clean}() restops.EndpointBuilder {
				${snippet.inputs}
				${snippet.outputs}

				return ${domain.style}.
					NewUpdateEndpointBuilder().
					WithId("${operation.id}").
					WithDoc(new(openapi3.Operation).
						WithSummary("${operation.summary}")).
					WithPermissions(permissionManageUpperCamelPlural).
					WithHandler(
						func(ctx context.Context, inp *inputs) (out outputs, err error) {
							out.Body, err = c.lowerCamelSingularService.UpdateUpperCamelSingular(ctx, inp.UpperCamelSingularId, inp.Body)
							return
						})
			}
			`),
	}.RenderContents(renderOptions)
	if err != nil {
		return err
	}

	return g.AddNewText(
		"Endpoints/Update",
		"update",
		template,
		[]codegen.Import{
			text.ImportRestOps,
			g.importStyle(),
			text.ImportOpenApi3,
			text.ImportContext,
		})
}

func (g DomainControllerGenerator) createEndpointActionDeleteSnippet(operation Operation) error {
	portStructVariables, err := g.generatePortStructVariables(operation)
	if err != nil {
		return err
	}

	renderOptions := skel.NewEmptyRenderOptions()
	renderOptions.AddVariables(portStructVariables)
	renderOptions.AddVariables(map[string]string{
		"operation.id":       *operation.Operation.ID,
		"operation.id.clean": g.cleanOperationId(*operation.Operation.ID),
		"operation.summary":  types.OptionalOfPtr(operation.Operation.Summary).OrElse(""),
	})

	template, err := skel.Template{
		SourceData: []byte(`
			// ${operation.id.clean} creates an endpoint deleting an existing UpperCamelSingular instance
			// matching the specified key.
			func (c *lowerCamelSingularController) ${operation.id.clean}() restops.EndpointBuilder {
				${snippet.inputs}

				return ${domain.style}.
					NewDeleteEndpointBuilder().
					WithId("${operation.id}").
					WithDoc(new(openapi3.Operation).
						WithSummary("${operation.summary}")).
					WithPermissions(permissionManageUpperCamelPlural).
					WithHandler(
						func(ctx context.Context, inp *inputs) (err error) {
							err = c.lowerCamelSingularService.DeleteUpperCamelSingular(ctx, inp.UpperCamelSingularId)
							return
						})
			}
			`),
	}.RenderContents(renderOptions)
	if err != nil {
		return err
	}

	return g.AddNewText(
		"Endpoints/Delete",
		"delete",
		template,
		[]codegen.Import{
			text.ImportRestOps,
			g.importStyle(),
			text.ImportOpenApi3,
			text.ImportContext,
		})
}

func (g DomainControllerGenerator) generatePortStructVariables(operation Operation) (vars map[string]string, err error) {
	var schema *jsonschema.Schema
	var snippet string
	vars = map[string]string{}

	// Inputs
	{
		schema, err = g.generateInputPortStructSchema(operation.Operation)
		if err != nil {
			return
		}

		snippet, err = g.generatePortStructSnippet(schema)
		if err != nil {
			return
		}

		vars["snippet.inputs"] = snippet
	}

	// Outputs
	if operation.Action != ActionDelete {
		schema, err = g.generateOutputPortStructSchema(operation.Operation)
		if err != nil {
			return
		}

		snippet, err = g.generatePortStructSnippet(schema)
		if err != nil {
			return
		}

		vars["snippet.outputs"] = snippet
	}

	return
}

func (g DomainControllerGenerator) generatePortStructSnippet(portStructSchema *jsonschema.Schema) (string, error) {
	schemaName := *portStructSchema.ID

	// Grab all the transitive schemas
	collectedSchema, err := g.Spec.CollectJsonSchema(schemaName, *portStructSchema)
	if err != nil {
		return "", err
	}

	// Convert to JSON
	schemaBytes, err := json.Marshal(collectedSchema)
	if err != nil {
		return "", err
	}

	payloadGenerator := types.May(generator.New(generator.Config{
		Capitalizations: []string{
			"inputs",
			"outputs",
		},
		OutputFiler: func(definition string) string {
			if definition == schemaName {
				return definition
			}
			return ""
		},
	}))

	err = payloadGenerator.AddSource(generator.Source{
		PackageName: g.Package,
		RootType:    schemaName,
		Folder:      "",
		Data:        schemaBytes,
	})
	if err != nil {
		return "", err
	}

	var snippet string
	for _, file := range payloadGenerator.Files() {
		if file.FileName == "." {
			continue
		}

		snippet = text.Decls(file.Package.Decls).Render(1)
		g.GoFile.AddImports(file.Package.Imports)
	}

	return snippet, nil
}

func (g DomainControllerGenerator) generateInputPortStructSchema(operation openapi3.Operation) (*jsonschema.Schema, error) {
	directionTag := "req"
	structTypeName := "inputs"

	portSchema := js.ObjectSchema()
	portSchema.ID = types.NewStringPtr(structTypeName)

	for _, parameterOrRef := range operation.Parameters {
		parameterOrRef = g.Spec.ResolveParameter(parameterOrRef)
		parameter := parameterOrRef.Parameter

		if parameter == nil {
			logger.Errorf("Failed to resolve parameter %s", parameterOrRef.ParameterReference.Ref)
		}

		if parameter.Schema == nil {
			logger.Errorf("No schema defined for parameter %s", parameter.Name)
			continue
		}

		// Get the parameter schema as jsonschema
		portFieldSchema := g.Spec.GetJsonSchema(parameter.Schema)

		k := parameter.Name
		propertyName := strcase.ToLowerCamel(k)

		payloads.
			GoJsonSchemaForSchema(portFieldSchema).Tags().
			ClearTag("json").
			AddTag(directionTag, fmt.Sprintf("%s=%s", parameter.In, k))

		required := parameter.In == openapi3.ParameterInPath ||
			(parameter.Required != nil && *parameter.Required)

		payloads.AddJsonSchemaObjectProperty(portSchema, propertyName, portFieldSchema, required)
	}

	if operation.RequestBody != nil {
		requestBodyOrRef := g.Spec.ResolveRequestBody(*operation.RequestBody)
		for mimeType, mediaType := range requestBodyOrRef.RequestBody.Content {
			requestBodySchema := mediaType.Schema
			if requestBodySchema == nil {
				continue
			}
			portFieldSchema := g.Spec.GetJsonSchema(requestBodySchema)

			tags := payloads.
				GoJsonSchemaForSchema(portFieldSchema).Tags().
				ClearTag("json").
				AddTag(directionTag, "body")

			if mimeType != restops.MediaTypeJson {
				tags.AddTag("mime", mimeType)
			}

			required := types.OptionalOfPtr(requestBodyOrRef.RequestBody.Required).OrElse(false)
			propertyName := "body"

			payloads.AddJsonSchemaObjectProperty(portSchema, propertyName, portFieldSchema, required)
		}
	}

	return portSchema, nil
}

func (g DomainControllerGenerator) generateOutputPortStructSchema(operation openapi3.Operation) (*jsonschema.Schema, error) {
	directionTag := "resp"
	structTypeName := "outputs"

	portSchema := js.ObjectSchema()
	portSchema.ID = types.NewStringPtr(structTypeName)

	for code, responseOrRef := range operation.Responses.MapOfResponseOrRefValues {
		isSuccess := strings.HasPrefix(code, "2") || strings.HasPrefix(code, "3")
		isError := strings.HasPrefix(code, "4") || strings.HasPrefix(code, "5")

		responseOrRef = g.Spec.ResolveResponse(responseOrRef)
		response := responseOrRef.Response

		if response == nil {
			logger.Errorf("Failed to resolve response %s", responseOrRef.ResponseReference.Ref)
		}

		for headerName, headerOrRef := range response.Headers {
			headerOrRef = g.Spec.ResolveHeader(headerOrRef)
			header := headerOrRef.Header

			if header == nil {
				logger.Errorf("Failed to resolve header %s", headerOrRef.HeaderReference.Ref)
			}

			// Get the parameter schema as jsonschema
			portFieldSchema := g.Spec.GetJsonSchema(header.Schema)

			k := headerName
			propertyName := strcase.ToLowerCamel(k)

			payloads.
				GoJsonSchemaForSchema(portFieldSchema).Tags().
				ClearTag("json").
				AddTag(directionTag, fmt.Sprintf("header=%s", k))

			required := types.OptionalOfPtr(headerOrRef.Header.Required).OrElse(false)

			payloads.AddJsonSchemaObjectProperty(portSchema, propertyName, portFieldSchema, required)
		}

		for mimeType, mediaType := range response.Content {
			responseBodySchema := mediaType.Schema
			if responseBodySchema == nil {
				continue
			}

			portFieldSchema := g.Spec.GetJsonSchema(responseBodySchema)

			tags := payloads.
				GoJsonSchemaForSchema(portFieldSchema).Tags().
				ClearTag("json").
				AddTag(directionTag, "body")

			if mimeType != restops.MediaTypeJson {
				tags.AddTag("mime", mimeType)
			}

			if isError && !isSuccess {
				tags.AddTag("error", "true")
			}

			propertyName := "body"
			if isError && !isSuccess {
				propertyName = "errorBody"
			}

			if portFieldSchema.Ref != nil {
				if strings.HasSuffix(*portFieldSchema.Ref, g.Style+".Error") && isError && !isSuccess {
					continue
				}
			}

			payloads.AddJsonSchemaObjectProperty(portSchema, propertyName, portFieldSchema, true)
		}

	}

	return portSchema, nil
}

func (g DomainControllerGenerator) createContextAccessorSnippet() error {
	return g.AddNewText(
		"Context",
		"contextAccessor",
		`
			// contextUpperCamelSingularController returns a ContextKeyAccessor enabling dependency overrides
			// for lowerCamelSingularController.
			func contextUpperCamelSingularController() types.ContextKeyAccessor[restops.EndpointsProducer] {
			  return types.NewContextKeyAccessor[restops.EndpointsProducer](contextKeyNamed("UpperCamelSingularController"))
			}
		`,
		[]codegen.Import{
			text.ImportTypes,
			text.ImportRestOps,
		})
}

func (g DomainControllerGenerator) createConstructorSnippet() error {
	return g.AddNewText(
		"Constructor",
		"constructor",
		`
			// newUpperCamelSingularController is an abstract factory, returning by default a production implementation
			// of the restops.EndpointsProducer.
			func newUpperCamelSingularController(ctx context.Context) (restops.EndpointsProducer, error) {
				controller := contextUpperCamelSingularController().Get(ctx)
				if controller == nil {
					lowerCamelSingularService, err := newUpperCamelSingularService(ctx)
					if err != nil {
						return nil, err
					}
			
					controller = &lowerCamelSingularController{
						lowerCamelSingularService: lowerCamelSingularService,
					}
				}
				return controller, nil
			}
		`,
		[]codegen.Import{
			text.ImportContext,
			text.ImportRestOps,
		})
}

func (g DomainControllerGenerator) createLifecycleSnippet() error {
	return g.AddNewText(
		"Lifecycle",
		"init",
		`
			// init adds an event observer to register the endpoints from UpperCamelSingularController during startup. 
			func init() {
				app.OnCommandsEvent(
					[]string{app.CommandRoot, app.CommandOpenApi},
					app.EventStart,
					app.PhaseBefore,
					func(ctx context.Context) error {
						controller, err := newUpperCamelSingularController(ctx)
						if err != nil {
							return err
						}
			
						return restops.
							ContextEndpointRegisterer(ctx).
							RegisterEndpoints(controller)
					})
			}
		`,
		[]codegen.Import{
			text.ImportRestOps,
			text.ImportContext,
			text.ImportApp,
		})
}

func (g DomainControllerGenerator) Apply(options skel.RenderOptions) skel.RenderOptions {
	options.AddVariable("domain.style", g.Style)
	return options
}

func (g DomainControllerGenerator) Generate() error {
	errs := types.ErrorList{
		g.createConstantsSnippet(),
		g.createControllerSnippet(),
		g.createEndpointTransformationSnippet(),
		g.createContextAccessorSnippet(),
		g.createConstructorSnippet(),
		g.createLifecycleSnippet(),
	}

	var operationMethods []string
	for _, operation := range g.Spec.Operations {
		if !g.Actions.Contains(operation.Action) {
			continue
		}

		// For the Endpoints Producer
		operationMethods = append(operationMethods, g.cleanOperationId(*operation.Operation.ID))

		var err error
		switch operation.Action {
		case ActionList:
			err = g.createEndpointActionListSnippet(operation)
		case ActionRetrieve:
			err = g.createEndpointActionRetrieveSnippet(operation)
		case ActionCreate:
			err = g.createEndpointActionCreateSnippet(operation)
		case ActionUpdate:
			err = g.createEndpointActionUpdateSnippet(operation)
		case ActionDelete:
			err = g.createEndpointActionDeleteSnippet(operation)
		}

		errs = append(errs, err)
	}

	errs = append(errs, g.createEndpointsConstructionSnippet(operationMethods))

	return errs.Filter()
}

func (g DomainControllerGenerator) Filename() string {
	target := path.Join(g.Folder, fmt.Sprintf("controller_lowersingular_%s.go", g.Style))
	return g.GoFile.Inflector.Inflect(target)
}

func (g DomainControllerGenerator) importStyle() codegen.Import {
	switch g.Style {
	case StyleV2:
		return text.ImportRestOpsV2
	default:
		return text.ImportRestOpsV8
	}
}

func NewDomainControllerGenerator(spec Spec) ComponentGenerator {
	inflector := skel.NewInflector(generatorConfig.Domain)

	return DomainControllerGenerator{
		// Configuration
		Domain:     generatorConfig.Domain,
		Folder:     generatorConfig.Folder,
		Tenant:     generatorConfig.Tenant,
		Actions:    generatorConfig.Actions,
		Components: generatorConfig.Components,
		Style:      generatorConfig.Style,

		// Sources
		Spec: spec,

		GoFile: &text.GoFile{
			File: &text.File[text.GoSnippet]{
				Comment:   "V8 API REST Controller for " + generatorConfig.Domain,
				Inflector: inflector,
				Sections: text.NewGoSections(
					"Constants",
					"Components",
					"Controller",
					"Endpoint Construction",
					"Endpoint Transformation",
					&text.Section[text.GoSnippet]{
						Name: "Endpoints",
						Sections: text.NewGoSections(
							"List",
							"Retrieve",
							"Create",
							"Update",
							"Delete",
						),
					},
					"Context",
					"Constructor",
					"Lifecycle",
				),
			},
			Package: generatorConfig.PackageName(),
		},
	}
}
