// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package skel

import (
	"fmt"
	"path"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/mcrawfo2/jennifer/jen"
	"github.com/pkg/errors"
)

func generateController(c Controller) (err error) {
	var ns = c.Package()
	var name = c.Name()

	f := jen.NewFile(c.Package())
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
const pkgValidation = "github.com/go-ozzo/ozzo-validation"
const pkgValidate = "cto-github.cisco.com/NFV-BU/go-msx/validate"
const pkgTypes = "cto-github.cisco.com/NFV-BU/go-msx/types"
const pkgRegexp = "regexp"

func generateParameterField(f *jen.File, c Controller, p Parameter) (*jen.Statement, error) {
	// params struct field
	parameterField := jen.Id(p.Name())

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

func generateBodyField(f *jen.File, c Controller, b Body) (*jen.Statement, error) {
	parameterField := jen.Id("Body")
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

func generateParameterVariable(p Parameter) jen.Code {
	parameterVarName := strcase.ToLowerCamel(fmt.Sprintf("param %s %s", p.In(), p.Name()))
	parameterTypeName := strcase.ToCamel(p.In() + "Parameter")

	parameterVariable := jen.Var().Id(parameterVarName).Op("=").
		Qual(pkgRestful, parameterTypeName).
		Call(jen.Lit(p.JsonName()), jen.Lit(p.Description()))

	if p.SchemaType.Required() {
		parameterVariable.Dot("Required").Call(jen.True())
	}

	return parameterVariable
}

func generateEndpointValidation(f *jen.File, e Endpoint) (jen.Dict, error) {
	result := make(jen.Dict)

	f.ImportName(pkgValidation, "validation")

	for _, p := range e.Parameters {
		result[jen.Lit(p.JsonName())] = jen.Qual(pkgValidation, "Validate").Call(
			jen.Op("&").Id("p").Dot(p.Name()),
		)
	}

	if e.RequestBody.Exists {
		f.ImportName(pkgValidate, "validate")

		validators, err := generateValidators(f, e.RequestBody.Schema)
		if err != nil {
			return nil, err
		}

		args := append([]jen.Code{
			jen.Op("&").Id("p").Dot("Body"),
		}, validators...)

		result[jen.Lit("body")] = jen.Qual(pkgValidation, "Validate").Call(args...)
	}

	return result, nil
}

func generateEndpoint(f *jen.File, c Controller, e Endpoint) (err error) {
	f.Line()

	var bodyStatements []jen.Code
	var parameterFields []jen.Code

	controllerDefinition := jen.Id("svc").
		Dot(strings.ToUpper(e.Method)).Call(jen.Lit(e.Path)).
		DotWrap("Operation").Call(jen.Lit(e.OperationId())).
		DotWrap("Doc").Call(jen.Lit(e.Summary()))

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
		controllerDefinition.DotWrap("Param").Call(jen.Id(p.VarName()))
	}

	if e.RequestBody.Exists {
		// struct field
		parameterField, err := generateBodyField(f, c, e.RequestBody)
		if err != nil {
			return errors.Wrapf(err, "Failed to generate body field")
		}
		parameterFields = append(parameterFields, parameterField)

		// controller reads
		typeExpr := new(jen.Statement)
		err = generateTypeWithImport(f, c.Package(), typeExpr, e.RequestBody.Schema)
		if err != nil {
			return errors.Wrapf(err, "Failed to generate body field")
		}
		controllerDefinition.DotWrap("Reads").Call(typeExpr.Values())
	}

	if len(bodyStatements) > 0 {
		bodyStatements = append(bodyStatements, jen.Line())
	}

	// params
	if len(parameterFields) > 0 {
		// params struct
		paramsTypeStatement := jen.Type().Id("params").Struct(parameterFields...)
		bodyStatements = append(bodyStatements, paramsTypeStatement, jen.Line())

		// populate params
		controllerDefinition.DotWrap("Do").Call(
			jen.Qual(pkgWebservice, "PopulateParams").Call(
				jen.New(jen.Id("params")),
			))

		// validate params
		validation, err := generateEndpointValidation(f, e)
		if err != nil {
			return err
		}

		paramsValidatorStatement := jen.Id("paramsValidator").Op(":=").Func().
			Params(
				jen.Id("req").Op("*").Qual(pkgRestful, "Request")).
			Parens(jen.List(
				jen.Id("err").Id("error"))).
			Block(
				// params := webservice.Params(req).(*params)
				jen.Id("p").Op(":=").Qual(pkgWebservice, "Params").Call(
					jen.Id("req")).Op(".").Parens(jen.Op("*").Id("params")),

				jen.Return(jen.Qual(pkgTypes, "ErrorMap").Values(validation)))

		bodyStatements = append(bodyStatements, paramsValidatorStatement, jen.Line())

		controllerDefinition.DotWrap("Do").Call(
			jen.Qual(pkgWebservice, "ValidateParams").Call(jen.Id("paramsValidator")))
	}

	// consumes
	if len(e.RequestBody.ContentTypes) > 0 {
		controllerDefinition.DotWrap("Consumes").Call(
			stringLiterals(e.RequestBody.ContentTypes)...)
	}

	// returns
	for _, Code := range e.ReturnCodes() {
		controllerDefinition.DotWrap("Do").
			Call(jen.Qual(pkgWebservice, fmt.Sprintf("Returns%s", Code)))
	}

	// response
	if e.ResponseBody.Exists {
		// payload
		typeExpr := new(jen.Statement)
		err = generateTypeWithImport(f, c.Package(), typeExpr, e.ResponseBody.Schema)
		if err != nil {
			return errors.Wrapf(err, "Failed to generate response body type")
		}

		controllerDefinition.DotWrap("Do").Call(
			jen.Qual(pkgWebservice, "ResponseRawPayload").
				Call(typeExpr.Values()))

		// produces
		controllerDefinition.DotWrap("Produces").Call(stringLiterals(e.ResponseBody.ContentTypes)...)
	} else {
		controllerDefinition.DotWrap("Produces").Call(jen.Lit("application/json"))
	}

	permissions := e.Permissions()
	if len(permissions) > 0 {
		expressions := stringLiterals(permissions)
		permissionsCall := jen.Qual(pkgWebservice, "Permissions").Call(expressions...)
		controllerDefinition.DotWrap("Do").Call(permissionsCall)
	}

	// controller
	controllerFunction := jen.Func().
		Params(
			jen.Id("req").Op("*").Qual(pkgRestful, "Request")).
		Parens(jen.List(
			jen.Id("body").Interface(),
			jen.Id("err").Id("error"))).
		Block( // TODO: Generate body based on archetype
			jen.Return(jen.List(jen.Nil(), jen.Nil())))

	controllerDefinition.DotWrap("To").Call(jen.Qual(pkgWebservice, "RawController").Call(controllerFunction))

	bodyStatements = append(bodyStatements, jen.Return(controllerDefinition))

	// endpoint generator func
	funcStatement := jen.Func().
		Params(jen.Id("c").Op("*").Id(c.Name())).
		Id(e.OperationId()).
		Params(jen.Id("svc").Op("*").Qual(pkgRestful, "WebService")).
		Op("*").Qual(pkgRestful, "RouteBuilder").
		Block(bodyStatements...)

	f.Add(funcStatement)

	return nil
}
