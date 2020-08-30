package skel

import (
	"fmt"
	"path"
	"strings"

	. "github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
)

func generateController(c Controller) (err error) {
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
			return errors.Wrapf(err, "Failed to generate endpoint %q", endpoint.OperationId())
		}
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
const pkgTypes = "cto-github.cisco.com/NFV-BU/go-msx/types"

func generateEndpoint(f *File, c Controller, e Endpoint) (err error) {
	f.Line()

	var bodyStatements []Code
	var params []Parameter

	var parameterFields []Code
	if len(e.Parameters) > 0 {

		for _, p := range e.Parameters {
			// params struct field
			parameterField := Id(p.Name())

			if !p.SchemaType.Required() {
				parameterField = parameterField.Op("*")
			}

			err = generateTypeWithImport(f, c.Namespace(skeletonConfig.AppPackageUrl()), parameterField, p.SchemaType)
			if err != nil {
				return errors.Wrapf(err, "Failed to generate parameter field %q", p.Name())
			}

			parameterField.Tag(map[string]string{
				"req": p.In(),
			})

			parameterFields = append(parameterFields, parameterField)

			// param variable
			parameterVarName := strcase.ToLowerCamel(fmt.Sprintf("param %s %s", p.In(), p.Name()))
			parameterTypeName := strcase.ToCamel(p.In() + "Parameter")

			parameterVariable := Var().Id(parameterVarName).Op("=").
				Qual(pkgRestful, parameterTypeName).
				Call(Lit(p.Name()), Lit(p.Description()))

			if p.SchemaType.Required() {
				parameterVariable.Dot("Required").Call(True())
			}

			bodyStatements = append(bodyStatements, parameterVariable)

		}

		bodyStatements = append(bodyStatements, Line())
		params = append(params, e.Parameters...)
	}

	if e.RequestBody.Exists {
		parameterField := Id("Body")
		if !e.RequestBody.SchemaType.Required() {
			parameterField = parameterField.Op("*")
		}

		err = generateTypeWithImport(f, c.Package(), parameterField, e.RequestBody.SchemaType)
		if err != nil {
			return errors.Wrapf(err, "Failed to generate body field")
		}

		parameterField.Tag(map[string]string{
			"req": "body",
		})

		parameterFields = append(parameterFields, parameterField)
	}

	if len(parameterFields) > 0 {
		paramsTypeStatement := Type().Id("params").Struct(parameterFields...)
		bodyStatements = append(bodyStatements, paramsTypeStatement, Line())
	}

	controllerDefinition := Id("svc").
		Dot(strings.ToUpper(e.Method)).Call(Lit(e.Path)).
		DotWrap("Operation").Call(Lit(e.OperationId())).
		DotWrap("Doc").Call(Lit(e.Summary()))

	for _, param := range params {
		controllerDefinition.DotWrap("Param").Call(Id(param.VarName()))
	}

	if e.RequestBody.Exists {
		typeExpr := new(Statement)

		err = generateTypeWithImport(f, c.Package(), typeExpr, e.RequestBody.SchemaType)
		if err != nil {
			return errors.Wrapf(err, "Failed to generate body field")
		}

		controllerDefinition.DotWrap("Reads").Call(typeExpr.Values())
	}

	if len(params) > 0 || e.RequestBody.Exists {
		controllerDefinition.DotWrap("Do").Call(
			Qual(pkgWebservice, "PopulateParams").Call(
				New(Id("params")),
			))

		validatorFunction := Func().
			Params(
				Id("req").Op("*").Qual(pkgRestful, "Request")).
			Parens(List(
				Id("err").Id("error"))).
			Block(
				// params := webservice.Params(req).(*params)
				Id("_").Op("=").Qual(pkgWebservice, "Params").Call(
					Id("req")).Op(".").Parens(Op("*").Id("params")),

				Return(Qual(pkgTypes, "ErrorMap").Values(Dict{
					// ...
				})))

		controllerDefinition.DotWrap("Do").Call(
			Qual(pkgWebservice, "ValidateParams").Call(
				validatorFunction,
			))
	}

	if len(e.RequestBody.ContentTypes) > 0 {
		controllerDefinition.DotWrap("Consumes").Call(
			stringLiterals(e.RequestBody.ContentTypes)...)
	}

	for _, code := range e.ReturnCodes() {
		controllerDefinition.DotWrap("Do").
			Call(Qual(pkgWebservice, fmt.Sprintf("Returns%s", code)))
	}

	if e.ResponseBody.Exists {
		typeExpr := new(Statement)
		err = generateTypeWithImport(f, c.Package(), typeExpr, e.ResponseBody.SchemaType)
		if err != nil {
			return errors.Wrapf(err, "Failed to generate response body type")
		}

		controllerDefinition.DotWrap("Do").Call(
			Qual(pkgWebservice, "ResponseRawPayload").
				Call(typeExpr.Values()))

		controllerDefinition.DotWrap("Produces").Call(stringLiterals(e.ResponseBody.ContentTypes)...)
	}

	//if permissions, ok := e.Operation.Extensions.GetStringSlice("x-msx-permissions"); ok {
	//	expressions := stringLiterals(permissions)
	//	for n, expression := range expressions {
	//		expressions[n] = Qual(pkgWebservice, "PermissionsFilter").Call(expression)
	//	}
	//	controllerDefinition.DotWrap("Filter").Call(expressions...)
	//}

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
		Id(e.OperationId()).
		Params(Id("svc").Op("*").Qual(pkgRestful, "WebService")).
		Op("*").Qual(pkgRestful, "RouteBuilder").
		Block(bodyStatements...)

	f.Add(funcStatement)

	return nil
}
