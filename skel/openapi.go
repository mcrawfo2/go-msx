package skel

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
	"fmt"
	. "github.com/dave/jennifer/jen"
	"github.com/gedex/inflector"
	"github.com/go-openapi/spec"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

type parameter struct {
	VarName string
	Param   spec.Parameter
}

type endpoint struct {
	Path      string
	Method    string
	Operation *spec.Operation
}

type controller struct {
	TagName        string
	TagDescription string
	RootPath       string
	Endpoints      []endpoint
}

func (c controller) Name() string {
	modelEnglish := inflector.Singularize(c.TagName)
	return strcase.ToLowerCamel(modelEnglish + " Controller")
}

func (c controller) Package() string {
	modelEnglish := inflector.Singularize(c.TagName)
	return strings.ReplaceAll(strcase.ToSnake(modelEnglish), "_", "")
}

func GenerateWebservices(args []string) error {
	if len(args) == 0 {
		return errors.New("No OpenAPI spec provided.")
	}

	bytes, err := ioutil.ReadFile(args[0])
	if err != nil {
		return errors.Wrap(err, "Failed to read OpenAPI spec")
	}

	var swagger spec.Swagger
	err = json.Unmarshal(bytes, &swagger)
	if err != nil {
		return errors.Wrap(err, "Failed to parse OpenAPI spec")
	}

	controllers, err := getControllers(swagger)
	if err != nil {
		return errors.Wrap(err, "Failed to identify controllers from spec")
	}

	for name, definition := range swagger.Definitions {
		err := generateDefinition(name, definition)
		if err != nil {
			return errors.Wrap(err, "Failed to generate definition")
		}
	}

	for _, controller := range controllers {
		err := generateController(controller, swagger.Definitions)
		if err != nil {
			return errors.Wrap(err, "Failed to generate controller")
		}
	}

	return nil
}

func getControllers(swagger spec.Swagger) ([]controller, error) {
	var controllers = make(map[string]controller)

	for _, tagDefinition := range swagger.Tags {
		controllers[tagDefinition.Name] = controller{
			TagName:        tagDefinition.Name,
			TagDescription: tagDefinition.Description,
			RootPath:       "",
			Endpoints:      nil,
		}
	}

	for pathKey, pathDefinition := range swagger.Paths.Paths {
		err := types.ErrorMap{
			http.MethodGet:    addControllerEndpoint(controllers, pathKey, http.MethodGet, pathDefinition.Get),
			http.MethodPut:    addControllerEndpoint(controllers, pathKey, http.MethodPut, pathDefinition.Put),
			http.MethodPost:   addControllerEndpoint(controllers, pathKey, http.MethodPost, pathDefinition.Post),
			http.MethodDelete: addControllerEndpoint(controllers, pathKey, http.MethodDelete, pathDefinition.Delete),
			http.MethodPatch:  addControllerEndpoint(controllers, pathKey, http.MethodPatch, pathDefinition.Patch),
		}.Filter()
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to add controller endpoint %q", pathKey)
		}
	}

	var result []controller
	for _, c := range controllers {
		// TODO: post-process controllers
		result = append(result, c)
	}
	return result, nil
}

func addControllerEndpoint(controllers map[string]controller, pathKey, method string, operation *spec.Operation) error {
	if operation == nil {
		return nil
	}

	controllerTagKey := operation.Tags[0]
	controllerEntry := controllers[controllerTagKey]

	ep := endpoint{
		Path:      pathKey,
		Method:    method,
		Operation: operation,
	}

	controllerEntry.Endpoints = append(controllerEntry.Endpoints, ep)

	controllers[controllerTagKey] = controllerEntry
	return nil
}

func generateDefinition(name string, definition spec.Schema) error {
	_, ns, name := splitNamespace(name)

	properties, imports, err := generateProperties(definition, ns)
	if err != nil {
		return errors.Wrap(err, "Failed to generate struct fields")
	}

	f := NewFile(ns)
	f.ImportNames(imports)
	f.Type().Id(name).Struct(properties...)

	targetFileName := path.Join(
		skeletonConfig.TargetDirectory(),
		"pkg",
		ns,
		strcase.ToSnake(name)+".go")

	return writeFile(targetFileName, f)
}

func generateProperties(definition spec.Schema, ns string) ([]Code, map[string]string, error) {
	var properties []Code
	var imports = make(map[string]string)

	for name, propertyDefinition := range definition.Properties {
		// Name
		property := Id(strcase.ToCamel(name))

		// Type
		typeImports, err := generateSchemaType(property, propertyDefinition, ns)
		if err != nil {
			return nil, nil, errors.Wrap(err, "Failed to generate property type")
		}

		for pkg, namespace := range typeImports {
			imports[pkg] = namespace
		}

		// Tag
		property = property.Tag(map[string]string{
			"json": strcase.ToLowerCamel(name),
		})

		properties = append(properties, property)
	}

	return properties, imports, nil
}

func generateSimpleSchemaType(s *Statement, parameterDefinition spec.SimpleSchema, ns string) (map[string]string, error) {
	propertyType := parameterDefinition.Type
	propertyItems := parameterDefinition.Items

	switch propertyType {
	case "string":
		s = s.String()
	case "boolean":
		s = s.Bool()
	case "array":
		s = s.Index()
		return generateSimpleSchemaType(s, propertyItems.SimpleSchema, ns)
	default:
		return nil, errors.Errorf("Unknown property type %q", propertyType[0])
	}

	return nil, nil
}

func generateSchemaType(s *Statement, propertyDefinition spec.Schema, ns string) (map[string]string, error) {
	if len(propertyDefinition.Type) == 0 {
		return generateRefType(s, propertyDefinition.Ref, ns)
	}

	propertyType := propertyDefinition.Type[0]
	propertyItems := propertyDefinition.Items

	switch propertyType {
	case "string":
		s = s.String()
	case "boolean":
		s = s.Bool()
	case "array":
		s = s.Index()
		return generateSchemaType(s, *propertyItems.Schema, ns)
	default:
		return nil, errors.Errorf("Unknown property type %q", propertyType[0])
	}

	return nil, nil
}

func generateRefType(s *Statement, ref spec.Ref, ns string) (map[string]string, error) {
	refParts := strings.Split(ref.String(), "/")
	refName := refParts[len(refParts)-1]
	pkg, namespace, name := splitNamespace(refName)

	imports := map[string]string{}
	if ns != namespace {
		imports[pkg] = namespace
		s = s.Qual(pkg, name)
	} else {
		s = s.Id(name)
	}
	return imports, nil
}

func splitNamespace(qualifiedName string) (string, string, string) {
	parts := strings.SplitN(qualifiedName, ".", 2)
	if len(parts) == 1 {
		parts = append([]string{"api"}, parts...)
	}
	pkg := path.Join("cto-github.cisco.com", "NFV-BU", skeletonConfig.AppName, "pkg", parts[0])
	return pkg, parts[0], parts[1]
}

func writeFile(targetFileName string, f *File) (err error) {
	err = os.MkdirAll(path.Dir(targetFileName), 0755)
	if err != nil {
		return errors.Wrap(err, "Failed to create directory")
	}

	writer, err := os.Create(targetFileName)
	if err != nil {
		return errors.Wrap(err, "Failed to create file")
	}

	err = f.Render(writer)
	if err != nil {
		return errors.Wrap(err, "Failed to write file")
	}

	return nil
}

func generateController(c controller, definitions spec.Definitions) (err error) {
	var ns = c.Package()
	var name = c.Name()

	f := NewFile(c.Package())
	f.ImportNames(map[string]string{
		"github.com/emicklei/go-restful":                "restful",
		"cto-github.cisco.com/NFV-BU/go-msx/webservice": "",
		"cto-github.cisco.com/NFV-BU/go-msx/log":        "",
	})
	f.Type().Id(name).Struct()

	for _, endpoint := range c.Endpoints {
		if err = generateEndpoint(f, c, endpoint); err != nil {
			return errors.Wrapf(err, "Failed to generate endpoint %q", endpoint.Operation)
		}
	}

	if err = generateEventHooks(f, c); err != nil {
		return errors.Wrap(err, "Failed to generate event hooks")
	}

	targetFileName := path.Join(
		skeletonConfig.TargetDirectory(),
		"internal",
		ns,
		strcase.ToSnake(name)+".go")

	return writeFile(targetFileName, f)
}

const pkgWebservice = "cto-github.cisco.com/NFV-BU/go-msx/webservice"
const pkgRestful = "github.com/emicklei/go-restful"
const pkgLog = "cto-github.cisco.com/NFV-BU/go-msx/log"

func generateEndpoint(f *File, c controller, e endpoint) (err error) {
	f.Line()

	var bodyStatements []Code
	var imports map[string]string
	var params []parameter

	if len(e.Operation.Parameters) > 0 {
		var parameterFields []Code

		for _, p := range e.Operation.Parameters {
			// params struct field
			parameterField := Id(strcase.ToCamel(p.Name))
			if p.Schema != nil {
				imports, err = generateSchemaType(parameterField, *p.Schema, c.Package())
			} else {
				imports, err = generateSimpleSchemaType(parameterField, p.SimpleSchema, c.Package())
			}
			if err != nil {
				return errors.Wrapf(err, "Failed to generate parameter field %q", p.Name)
			}
			f.ImportNames(imports)

			parameterField.Tag(map[string]string{
				"req": p.In,
			})

			parameterFields = append(parameterFields, parameterField)

			// param variable
			parameterVarName := strcase.ToLowerCamel(fmt.Sprintf("param %s %s", p.In, p.Name))
			parameterTypeName := strcase.ToCamel(p.In + "Parameter")

			parameterVariable := Var().Id(parameterVarName).Op("=").
				Qual(pkgRestful, parameterTypeName).
				Call(Lit(p.Name), Lit(p.Description))

			if p.Required {
				parameterVariable.Dot("Required").Call(True())
			}

			if p.In != "body" {
				bodyStatements = append(bodyStatements, parameterVariable)
			}

			params = append(params, parameter{
				VarName: parameterVarName,
				Param:   p,
			})
		}

		paramsTypeStatement := Type().Id("params").Struct(parameterFields...)
		bodyStatements = append(bodyStatements, Line(), paramsTypeStatement, Line())
	}

	controllerDefinition := Id("svc").
		Dot(strings.ToUpper(e.Method)).Call(Lit(e.Path)).
		DotWrap("Operation").Call(Lit(e.Operation.ID)).
		DotWrap("Doc").Call(Lit(e.Operation.Summary))

	for _, param := range params {
		if param.Param.In != "body" {
			controllerDefinition.DotWrap("Param").Call(Id(param.VarName))
		} else {
			typeExpr := new(Statement)
			_, _ = generateSchemaType(typeExpr, *param.Param.Schema, c.Package())
			controllerDefinition.DotWrap("Reads").Call(typeExpr.Values())
		}
	}

	if len(e.Operation.Consumes) > 0 {
		controllerDefinition.DotWrap("Consumes").Call(stringLiterals(e.Operation.Consumes)...)
	}

	for code := range e.Operation.Responses.StatusCodeResponses {
		controllerDefinition.DotWrap("Do").
			Call(Qual(pkgWebservice, fmt.Sprintf("Returns%d", code)))
	}

	if e.Operation.Responses.Default != nil {
		typeExpr := new(Statement)
		_, _ = generateSchemaType(typeExpr, *e.Operation.Responses.Default.Schema, c.Package())
		controllerDefinition.DotWrap("Do").Call(
			Qual(pkgWebservice, "ResponseRawPayload").
				Call(typeExpr.Values()))
	}

	if len(e.Operation.Produces) > 0 {
		controllerDefinition.DotWrap("Produces").Call(stringLiterals(e.Operation.Produces)...)
	}

	if permissions, ok := e.Operation.Extensions.GetStringSlice("x-msx-permissions"); ok {
		expressions := stringLiterals(permissions)
		for n, expression := range expressions {
			expressions[n] = Qual(pkgWebservice, "PermissionsFilter").Call(expression)
		}
		controllerDefinition.DotWrap("Filter").Call(expressions...)
	}

	controllerFunction := Func().
		Params(
			Id("req").Op("*").Qual(pkgRestful, "Request")).
		Parens(List(
			Id("body").Interface(),
			Id("err").Id("error"))).
		Block(
			Return(List(Nil(), Nil())))

	controllerDefinition.DotWrap("To").Call(Qual(pkgWebservice, "RawController").Call(controllerFunction))

	bodyStatements = append(bodyStatements, Return(controllerDefinition))

	funcStatement := Func().
		Params(Id("c").Op("*").Id(c.Name())).
		Id(e.Operation.ID).
		Params(Id("svc").Op("*").Qual(pkgRestful, "WebService")).
		Op("*").Qual(pkgRestful, "RouteBuilder").
		Block(bodyStatements...)

	f.Add(funcStatement)

	return nil
}

func stringLiterals(values []string) []Code {
	var literals []Code
	for _, value := range values {
		literals = append(literals, Lit(value))
	}
	return literals
}

func generateEventHooks(f *File, c controller) error {
	return nil
}
