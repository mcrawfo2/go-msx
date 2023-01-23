// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package rest

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops/restops"
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
	"github.com/swaggest/jsonschema-go"
	"github.com/swaggest/openapi-go/openapi3"
	"path"
	"sort"
	"strings"
)

type DomainControllerGeneratorV8 struct {
	Domain     string
	Folder     string
	Tenant     string
	Actions    types.ComparableSlice[string]
	Components []string

	Spec Spec

	OutputVariables  map[string]string
	OutputConditions map[string]bool

	*File
}

func (g DomainControllerGeneratorV8) createConstantsSnippet() error {
	constants := Constants{
		{
			Name:  "pathSuffixUpperCamelSingularId",
			Value: "{lowerCamelSingularId}",
		},
	}

	if g.Actions.ContainsAny(ActionList, ActionRetrieve) {
		constants = append(constants, &codegen.Constant{
			Name:  "permissionViewUpperCamelPlural",
			Value: g.File.Inflector.Inflect("VIEW_SCREAMING_SNAKE_PLURAL"),
		})
	}

	if g.Actions.ContainsAny(ActionCreate, ActionUpdate, ActionDelete) {
		constants = append(constants, &codegen.Constant{
			Name:  "permissionManageUpperCamelPlural",
			Value: g.File.Inflector.Inflect("MANAGE_SCREAMING_SNAKE_PLURAL"),
		})
	}

	return g.AddNewDecl(
		"Constants",
		"pathConstants",
		constants,
		[]codegen.Import{})
}

func (g DomainControllerGeneratorV8) createControllerSnippet() error {
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

func (g DomainControllerGeneratorV8) createEndpointsConstructionSnippet(methods []string) error {
	sort.Strings(methods)

	var operations = new(strings.Builder)
	for _, method := range methods {
		err := jen.Id("c").Dot(method).Call().Render(operations)
		if err != nil {
			return err
		}
		operations.WriteString(",\n")
	}

	renderOptions := skel.RenderOptions{
		Variables: map[string]string{
			"operations": operations.String(),
		},
	}

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
			importRestOps,
		})
}

func (g DomainControllerGeneratorV8) createEndpointTransformationSnippet() error {
	return g.AddNewText(
		"Endpoint Transformation",
		"endpointTransformers",
		`
			// EndpointTransformers provides a set of transformations to be applied to each Endpoint
			// created by this controller.
			func (c *lowerCamelSingularController) EndpointTransformers() restops.EndpointTransformers {
				const tagName = "UpperCamelSingular"
				const pathPrefix = "/v8/lowerplural"

				openapi.AddTag(tagName, "Title Plural")
			
				return restops.EndpointTransformers{
					restops.AddEndpointPathPrefix(pathPrefix),
					restops.AddEndpointTag(tagName),
				}
			}
			`,
		[]codegen.Import{
			importRestOps,
			importOpenapi,
		})
}

func (g DomainControllerGeneratorV8) cleanOperationId(operationId string) string {
	cleanOperationId := operationId
	if strings.Contains(cleanOperationId, ".") {
		lastPeriod := strings.LastIndex(cleanOperationId, ".")
		cleanOperationId = cleanOperationId[lastPeriod+1:]
	}
	return cleanOperationId
}

func (g DomainControllerGeneratorV8) createEndpointActionListSnippet(operation Operation) error {
	renderOptions := skel.RenderOptions{
		Variables: map[string]string{
			"snippet.outputs.content.tag": "`resp:\"body\"`",
			"operation.id":                *operation.Operation.ID,
			"operation.id.clean":          g.cleanOperationId(*operation.Operation.ID),
			"operation.summary":           types.OptionalOfPtr(operation.Operation.Summary).OrElse(""),
		},
	}

	template, err := skel.Template{
		SourceData: []byte(`
			// lowerCamelSingularFilterQueryInputs is used to declare the query string filters
			// for the ${operation.id.clean} endpoint.
			type lowerCamelSingularFilterQueryInputs struct {
			}

			// ${operation.id.clean} creates an endpoint providing a filtered, sorted, and paginated 
			// sequence of ${UpperCamelSingular} instances.
			func (c *lowerCamelSingularController) ${operation.id.clean}() restops.EndpointBuilder {
				type inputs struct {
					v8.PagingSortingInputs
					lowerCamelSingularFilterQueryInputs
				}

				type outputs struct {
					v8.PagingOutputs
					Content []UpperCamelSingularResponse ${snippet.outputs.content.tag}
				}

				return v8.
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
			importRestOps,
			importRestOpsV8,
			importOpenApi3,
			importContext,
		})
}

func (g DomainControllerGeneratorV8) createEndpointActionRetrieveSnippet(operation Operation) error {
	portStructVariables, err := g.generatePortStructVariables(operation)
	if err != nil {
		return err
	}

	renderOptions := skel.RenderOptions{
		Variables: portStructVariables,
	}

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

				return v8.
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
			importRestOps,
			importRestOpsV8,
			importOpenApi3,
			importContext,
		})
}

func (g DomainControllerGeneratorV8) createEndpointActionCreateSnippet(operation Operation) error {
	portStructVariables, err := g.generatePortStructVariables(operation)
	if err != nil {
		return err
	}

	renderOptions := skel.RenderOptions{
		Variables: portStructVariables,
	}

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

				return v8.
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
			importRestOps,
			importRestOpsV8,
			importOpenApi3,
			importContext,
		})
}

func (g DomainControllerGeneratorV8) createEndpointActionUpdateSnippet(operation Operation) error {
	portStructVariables, err := g.generatePortStructVariables(operation)
	if err != nil {
		return err
	}

	renderOptions := skel.RenderOptions{
		Variables: portStructVariables,
	}

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

				return v8.
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
			importRestOps,
			importRestOpsV8,
			importOpenApi3,
			importContext,
		})
}

func (g DomainControllerGeneratorV8) createEndpointActionDeleteSnippet(operation Operation) error {
	portStructVariables, err := g.generatePortStructVariables(operation)
	if err != nil {
		return err
	}

	renderOptions := skel.RenderOptions{
		Variables: portStructVariables,
	}

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

				return v8.
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
			importRestOps,
			importRestOpsV8,
			importOpenApi3,
			importContext,
		})
}

func (g DomainControllerGeneratorV8) generatePortStructVariables(operation Operation) (vars map[string]string, err error) {
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

func (g DomainControllerGeneratorV8) generatePortStructSnippet(portStructSchema *jsonschema.Schema) (string, error) {
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

		snippet = Decls(file.Package.Decls).Render(1)
		g.File.AddImports(file.Package.Imports)
	}

	return snippet, nil
}

func (g DomainControllerGeneratorV8) generateInputPortStructSchema(operation openapi3.Operation) (*jsonschema.Schema, error) {
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

func (g DomainControllerGeneratorV8) generateOutputPortStructSchema(operation openapi3.Operation) (*jsonschema.Schema, error) {
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
				if strings.HasSuffix(*portFieldSchema.Ref, "v8.Error") && isError && !isSuccess {
					continue
				}
			}

			payloads.AddJsonSchemaObjectProperty(portSchema, propertyName, portFieldSchema, true)
		}

	}

	return portSchema, nil
}

func (g DomainControllerGeneratorV8) createContextAccessorSnippet() error {
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
			importTypes,
			importRestOps,
		})
}

func (g DomainControllerGeneratorV8) createConstructorSnippet() error {
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
			importContext,
			importRestOps,
		})
}

func (g DomainControllerGeneratorV8) createLifecycleSnippet() error {
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
			importRestOps,
			importContext,
			importApp,
		})
}

func (g DomainControllerGeneratorV8) Generate() error {
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

func (g DomainControllerGeneratorV8) Filename() string {
	target := path.Join(g.Folder, "controller_lowersingular_v8.go")
	return g.File.Inflector.Inflect(target)
}

func (g DomainControllerGeneratorV8) Variables() map[string]string {
	return g.OutputVariables
}

func (g DomainControllerGeneratorV8) Conditions() map[string]bool {
	return g.OutputConditions
}

func NewDomainControllerGeneratorV8(spec Spec) ComponentGenerator {
	inflector := skel.NewInflector(generatorConfig.Domain)

	return DomainControllerGeneratorV8{
		// Configuration
		Domain:     generatorConfig.Domain,
		Folder:     generatorConfig.Folder,
		Tenant:     generatorConfig.Tenant,
		Actions:    generatorConfig.Actions,
		Components: generatorConfig.Components,

		// Sources
		Spec: spec,

		// Results
		OutputVariables:  make(map[string]string),
		OutputConditions: make(map[string]bool),

		File: &File{
			Comment:   "V8 API REST Controller for " + generatorConfig.Domain,
			Package:   generatorConfig.PackageName(),
			Inflector: inflector,
			Sections: NewSections(
				"Constants",
				"Components",
				"Controller",
				"Endpoint Construction",
				"Endpoint Transformation",
				&Section{
					Name: "Endpoints",
					Sections: NewSections(
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
	}
}
