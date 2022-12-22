// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package openapi

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/ops/restops"
	"github.com/pkg/errors"
	"github.com/swaggest/openapi-go/openapi3"
)

const (
	ExtensionPermissions = "x-msx-permissions"
)

type EndpointDocumentor struct {
	Skip      bool
	Operation *openapi3.Operation
	Mutator   ops.DocumentElementMutator[openapi3.Operation]
}

func (d *EndpointDocumentor) WithSkip(skip bool) *EndpointDocumentor {
	d.Skip = skip
	return d
}

func (d *EndpointDocumentor) WithOperation(op *openapi3.Operation) *EndpointDocumentor {
	d.Operation = op
	return d
}

func (d *EndpointDocumentor) WithMutator(mutator ops.DocumentElementMutator[openapi3.Operation]) *EndpointDocumentor {
	d.Mutator = mutator
	return d
}

func (d *EndpointDocumentor) DocType() string {
	return DocType
}

func (d *EndpointDocumentor) Document(e *restops.Endpoint) error {
	if d.Skip {
		return nil
	}

	if d.Operation == nil {
		d.Operation = new(openapi3.Operation)
	}

	d.Operation.
		WithID(e.OperationID).
		WithTags(e.Tags...)

	if e.Deprecated {
		d.Operation.WithDeprecated(e.Deprecated)
	}

	if e.Description != "" {
		d.Operation.WithDescription(e.Description)
	}

	if e.Summary != "" {
		d.Operation.WithSummary(e.Summary)
	}

	for _, p := range e.Request.Parameters {
		if parameterOrRef, err := d.DocumentParameter(p); err != nil {
			return err
		} else if parameterOrRef != nil {
			d.Operation.Parameters = append(d.Operation.Parameters, *parameterOrRef)
		}
	}

	if requestBody, err := d.DocumentRequestBody(e.Request.Body); err != nil {
		return err
	} else if requestBody != nil {
		d.Operation.WithRequestBody(*requestBody)
	}

	if responses, err := d.DocumentResponses(e, e.Response); err != nil {
		return err
	} else {
		d.Operation.WithResponses(*responses)
	}

	if len(e.Permissions) > 0 {
		d.Operation.WithMapOfAnythingItem(ExtensionPermissions, e.Permissions)
	}

	if d.Mutator != nil {
		d.Mutator(d.Operation)
	}

	return nil
}

func (d *EndpointDocumentor) DocumentParameter(p restops.EndpointRequestParameter) (result *openapi3.ParameterOrRef, err error) {
	doc := ops.DocumentorWithType[restops.EndpointRequestParameter](p, DocType).
		OrElse(new(EndpointParameterDocumentor))

	if err = doc.Document(&p); err != nil {
		return
	}

	if resulter, ok := doc.(ops.DocumentResult[openapi3.ParameterOrRef]); !ok {
		err = errors.Errorf("Unable to retrieve parameter %q documentation from Documentor", p.Name)
	} else {
		result = resulter.Result()
	}

	return
}

func (d *EndpointDocumentor) DocumentResponses(e *restops.Endpoint, r restops.EndpointResponse) (result *openapi3.Responses, err error) {
	doc := ops.DocumentorWithType[restops.EndpointResponse](r, DocType).
		OrElse(new(EndpointResponseDocumentor))

	doc.(*EndpointResponseDocumentor).WithEndpoint(e)

	if err = doc.Document(&r); err != nil {
		return
	}

	if resulter, ok := doc.(ops.DocumentResult[openapi3.Responses]); !ok {
		err = errors.New("Unable to retrieve responses documentation from Documentor")
	} else {
		result = resulter.Result()
	}

	return
}

func (d *EndpointDocumentor) DocumentRequestBody(b restops.EndpointRequestBody) (result *openapi3.RequestBodyOrRef, err error) {
	if b.Mime == "" {
		return
	}

	doc := ops.DocumentorWithType[restops.EndpointRequestBody](b, DocType).
		OrElse(new(EndpointRequestBodyDocumentor))

	if err = doc.Document(&b); err != nil {
		return
	}

	if resulter, ok := doc.(ops.DocumentResult[openapi3.RequestBody]); !ok {
		err = errors.New("Unable to retrieve responses documentation from Documentor")
	} else {
		result = &openapi3.RequestBodyOrRef{
			RequestBody: resulter.Result(),
		}
	}

	return
}

func (d *EndpointDocumentor) Result() *openapi3.Operation {
	return d.Operation
}

type EndpointDocumentorBuilder struct {
	Skip      bool
	Operation *openapi3.Operation
	Mutator   func(*openapi3.Operation)
}

func (b EndpointDocumentorBuilder) Build() *EndpointDocumentor {
	return &EndpointDocumentor{
		Skip:      b.Skip,
		Operation: b.Operation,
		Mutator:   b.Mutator,
	}
}
