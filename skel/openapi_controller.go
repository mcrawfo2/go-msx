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
const pkgValidation = "github.com/go-ozzo/ozzo-validation"
const pkgValidate = "cto-github.cisco.com/NFV-BU/go-msx/validate"
const pkgTypes = "cto-github.cisco.com/NFV-BU/go-msx/types"
const pkgRegexp = "regexp"

func generateParameterField(f *File, c Controller, p Parameter) (Code, error) {
	// params struct field
	parameterField := Id(p.Name())

	if !p.SchemaType.Required() {
		parameterField = parameterField.Op("*")
	}

	err := generateTypeWithImport(f, c.Namespace(skeletonConfig.AppPackageUrl()), parameterField, p.SchemaType)
	if err != nil {
		return nil, err
	}

	parameterField.Tag(map[string]string{"req": p.In()})

	return parameterField, nil
}

func generateBodyField(f *File, c Controller, b Body) (Code, error) {
	parameterField := Id("Body")
	if !b.Schema.Required() {
		parameterField = parameterField.Op("*")
	}

	err := generateTypeWithImport(f, c.Package(), parameterField, b.Schema)
	if err != nil {
		return nil, err
	}

	parameterField.Tag(map[string]string{"req": "body"})

	return parameterField, nil
}

func generateParameterVariable(p Parameter) Code {
	parameterVarName := strcase.ToLowerCamel(fmt.Sprintf("param %s %s", p.In(), p.Name()))
	parameterTypeName := strcase.ToCamel(p.In() + "Parameter")

	parameterVariable := Var().Id(parameterVarName).Op("=").
		Qual(pkgRestful, parameterTypeName).
		Call(Lit(p.JsonName()), Lit(p.Description()))

	if p.SchemaType.Required() {
		parameterVariable.Dot("Required").Call(True())
	}

	return parameterVariable
}

func generateEndpointValidation(f *File, e Endpoint) (Dict, error) {
	result := make(Dict)

	f.ImportName(pkgValidation, "validation")

	for _, p := range e.Parameters {
		result[Lit(p.JsonName())] = Qual(pkgValidation, "Validate").Call(
			Op("&").Id("p").Dot(p.Name()),
			)
	}

	if e.RequestBody.Exists {
		f.ImportName(pkgValidate, "validate")

		validators, err := generateValidators(f, e.RequestBody.Schema)
		if err != nil {
			return nil, err
		}

		args := append([]Code{
			Op("&").Id("p").Dot("Body"),
		}, validators...)

		result[Lit("body")] = Qual(pkgValidation, "Validate").Call(args...)
	}

	return result, nil
}

func generateEndpoint(f *File, c Controller, e Endpoint) (err error) {
	f.Line()

	var bodyStatements []Code
	var parameterFields []Code
	var controllerDefinition *Statement

	controllerDefinition = Id("svc").
		Dot(strings.ToUpper(e.Method)).Call(Lit(e.Path)).
		DotWrap("Operation").Call(Lit(e.OperationId())).
		DotWrap("Doc").Call(Lit(e.Summary()))

	// Parameter declarations
	for _, p := range e.Parameters {
		// struct field
		parameterField, err := generateParameterField(f, c, p)
		if err != nil {
			return errors.Wrapf(err, "Failed to generate parameter field %q", p.Name())
		}
		parameterFields = append(parameterFields, parameterField)

		// param variable
		parameterVariable := generateParameterVariable(p)
		bodyStatements = append(bodyStatements, parameterVariable)

		// controller param
		controllerDefinition.DotWrap("Param").Call(Id(p.VarName()))
	}

	if e.RequestBody.Exists {
		// struct field
		parameterField, err := generateBodyField(f, c, e.RequestBody)
		if err != nil {
			return errors.Wrapf(err, "Failed to generate body field")
		}
		parameterFields = append(parameterFields, parameterField)

		// controller reads
		typeExpr := new(Statement)
		err = generateTypeWithImport(f, c.Package(), typeExpr, e.RequestBody.Schema)
		if err != nil {
			return errors.Wrapf(err, "Failed to generate body field")
		}
		controllerDefinition.DotWrap("Reads").Call(typeExpr.Values())
	}

	if len(bodyStatements) > 0 {
		bodyStatements = append(bodyStatements, Line())
	}

	// params
	if len(parameterFields) > 0 {
		// params struct
		paramsTypeStatement := Type().Id("params").Struct(parameterFields...)
		bodyStatements = append(bodyStatements, paramsTypeStatement, Line())

		// populate params
		controllerDefinition.DotWrap("Do").Call(
			Qual(pkgWebservice, "PopulateParams").Call(
				New(Id("params")),
			))

		// validate params
		validation, err := generateEndpointValidation(f, e)
		if err != nil {
			return err
		}

		paramsValidatorStatement := Id("paramsValidator").Op(":=").Func().
			Params(
				Id("req").Op("*").Qual(pkgRestful, "Request")).
			Parens(List(
				Id("err").Id("error"))).
			Block(
				// params := webservice.Params(req).(*params)
				Id("p").Op(":=").Qual(pkgWebservice, "Params").Call(
					Id("req")).Op(".").Parens(Op("*").Id("params")),

				Return(Qual(pkgTypes, "ErrorMap").Values(validation)))

		bodyStatements = append(bodyStatements, paramsValidatorStatement, Line())

		controllerDefinition.DotWrap("Do").Call(
			Qual(pkgWebservice, "ValidateParams").Call(Id("paramsValidator")))
	}

	// consumes
	if len(e.RequestBody.ContentTypes) > 0 {
		controllerDefinition.DotWrap("Consumes").Call(
			stringLiterals(e.RequestBody.ContentTypes)...)
	}

	// returns
	for _, code := range e.ReturnCodes() {
		controllerDefinition.DotWrap("Do").
			Call(Qual(pkgWebservice, fmt.Sprintf("Returns%s", code)))
	}

	// response
	if e.ResponseBody.Exists {
		// payload
		typeExpr := new(Statement)
		err = generateTypeWithImport(f, c.Package(), typeExpr, e.ResponseBody.Schema)
		if err != nil {
			return errors.Wrapf(err, "Failed to generate response body type")
		}

		controllerDefinition.DotWrap("Do").Call(
			Qual(pkgWebservice, "ResponseRawPayload").
				Call(typeExpr.Values()))

		// produces
		controllerDefinition.DotWrap("Produces").Call(stringLiterals(e.ResponseBody.ContentTypes)...)
	}

	permissions := e.Permissions()
	if len(permissions) > 0 {
		expressions := stringLiterals(permissions)
		permissionsCall := Qual(pkgWebservice, "Permissions").Call(expressions...)
		controllerDefinition.DotWrap("Do").Call(permissionsCall)
	}

	// controller
	controllerFunction := Func().
		Params(
			Id("req").Op("*").Qual(pkgRestful, "Request")).
		Parens(List(
			Id("body").Interface(),
			Id("err").Id("error"))).
		Block(	// TODO: Generate body based on archetype
			Return(List(Nil(), Nil())))

	controllerDefinition.DotWrap("To").Call(Qual(pkgWebservice, "RawController").Call(controllerFunction))

	bodyStatements = append(bodyStatements, Return(controllerDefinition))

	// endpoint generator func
	funcStatement := Func().
		Params(Id("c").Op("*").Id(c.Name())).
		Id(e.OperationId()).
		Params(Id("svc").Op("*").Qual(pkgRestful, "WebService")).
		Op("*").Qual(pkgRestful, "RouteBuilder").
		Block(bodyStatements...)

	f.Add(funcStatement)

	return nil
}
