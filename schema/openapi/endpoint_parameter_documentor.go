// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package openapi

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/ops/restops"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/swaggest/openapi-go/openapi3"
)

type EndpointParameterDocumentor struct {
	Skip           bool
	Parameter      *openapi3.Parameter
	ParameterOrRef *openapi3.ParameterOrRef
	Mutator        ops.DocumentElementMutator[openapi3.Parameter]
}

func (e *EndpointParameterDocumentor) WithSkip(skip bool) *EndpointParameterDocumentor {
	e.Skip = skip
	return e
}

func (e *EndpointParameterDocumentor) WithParameter(p *openapi3.Parameter) *EndpointParameterDocumentor {
	e.Parameter = p
	return e
}

func (e *EndpointParameterDocumentor) WithMutator(mutator ops.DocumentElementMutator[openapi3.Parameter]) *EndpointParameterDocumentor {
	e.Mutator = mutator
	return e
}

func (e *EndpointParameterDocumentor) DocType() string {
	return DocType
}

func (e *EndpointParameterDocumentor) Document(p *restops.EndpointRequestParameter) error {
	if e.Skip {
		return nil
	}

	if e.Parameter == nil {
		e.Parameter = new(openapi3.Parameter)
	}

	e.Parameter.
		WithName(p.Name).
		WithIn(openapi3.ParameterIn(p.In))

	if p.Description != nil {
		e.Parameter.WithDescription(*p.Description)
	}

	if p.Required != nil {
		e.Parameter.WithRequired(*p.Required)
	}

	if p.Deprecated != nil {
		e.Parameter.WithDeprecated(*p.Deprecated)
	}

	if p.AllowEmptyValue != nil {
		e.Parameter.WithAllowEmptyValue(*p.AllowEmptyValue)
	}

	if p.Style != nil {
		e.Parameter.WithStyle(*p.Style)
	}

	if p.Explode != nil {
		e.Parameter.WithExplode(*p.Explode)
	}

	if p.AllowReserved != nil {
		e.Parameter.WithAllowReserved(*p.AllowReserved)
	}

	if p.PortField != nil {
		schemaOrRef := openApiSchemaFromPortField(p.PortField)
		e.Parameter.WithSchema(*schemaOrRef)
	} else if p.Payload.IsPresent() {
		schemaOrRef, err := Reflect(p.Payload.Value())
		if err != nil {
			return err
		}
		e.Parameter.WithSchema(*schemaOrRef)
	} else if p.Type != nil {
		schema := SwaggerTypeToOpenApiSchema(p.Type, p.Format)
		if types.NewOptionalBool(p.Multi).OrElse(false) {
			schema = ArraySchema(NewSchemaOrRef(schema))
		}
		if p.Enum != nil {
			schema.Enum = p.Enum
		}
		if p.Default.IsPresent() {
			value := p.Default.Value()
			schema.Default = &value
		}
		if schema != nil {
			e.Parameter.WithSchema(NewSchemaOrRef(schema))
		}
	}

	if p.Example.IsPresent() {
		e.Parameter.WithExample(p.Example.Value())
	}

	if e.Mutator != nil {
		e.Mutator(e.Parameter)
	}

	if p.Reference != nil {
		Spec().ComponentsEns().ParametersEns().WithMapOfParameterOrRefValuesItem(
			*p.Reference,
			openapi3.ParameterOrRef{
				Parameter: e.Parameter,
			})
		e.ParameterOrRef = NewParameterRefPtr(*p.Reference)
	} else {
		e.ParameterOrRef = &openapi3.ParameterOrRef{
			Parameter: e.Parameter,
		}
	}

	return nil
}

func (e *EndpointParameterDocumentor) Result() *openapi3.ParameterOrRef {
	if e.ParameterOrRef == nil {
		return nil
	}

	return e.ParameterOrRef
}
