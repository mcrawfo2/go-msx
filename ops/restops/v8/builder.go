// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package v8

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops/restops"
	"cto-github.cisco.com/NFV-BU/go-msx/schema/openapi"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
	"github.com/swaggest/openapi-go/openapi3"
	"net/http"
	"path"
)

type EndpointArchetype struct {
	Verb  string
	Codes restops.EndpointResponseCodes
}

const (
	EndpointArchetypeList         = "list"
	EndpointArchetypeRetrieve     = "retrieve"
	EndpointArchetypeCreate       = "create"
	EndpointArchetypeUpdate       = "update"
	EndpointArchetypeDelete       = "delete"
	EndpointArchetypeCommand      = "command"
	EndpointArchetypeAsyncCommand = "asyncCommand"
)

var archetypes = map[string]EndpointArchetype{
	EndpointArchetypeList:         {Verb: http.MethodGet, Codes: restops.ListResponseCodes},
	EndpointArchetypeRetrieve:     {Verb: http.MethodGet, Codes: restops.GetResponseCodes},
	EndpointArchetypeCreate:       {Verb: http.MethodPost, Codes: restops.CreateResponseCodes},
	EndpointArchetypeUpdate:       {Verb: http.MethodPut, Codes: restops.UpdateResponseCodes},
	EndpointArchetypeDelete:       {Verb: http.MethodDelete, Codes: restops.NoContentResponseCodes},
	EndpointArchetypeCommand:      {Verb: http.MethodPost, Codes: restops.GetResponseCodes},
	EndpointArchetypeAsyncCommand: {Verb: http.MethodPost, Codes: restops.AcceptResponseCodes},
}

type EndpointBuilder struct {
	Archetype     string
	Id            string
	Path          string
	Documentation openapi.EndpointDocumentorBuilder
	Permissions   []string

	Inputs  types.Optional[interface{}]
	Outputs types.Optional[interface{}]

	Handler interface{}
}

func (b *EndpointBuilder) WithId(operationId string) *EndpointBuilder {
	b.Id = operationId
	return b
}

func (b *EndpointBuilder) WithDoc(doc *openapi3.Operation) *EndpointBuilder {
	b.Documentation.Operation = doc
	return b
}

func (b *EndpointBuilder) WithDocMutator(mutator func(op *openapi3.Operation)) *EndpointBuilder {
	oldMutator := b.Documentation.Mutator
	if oldMutator != nil {
		b.Documentation.Mutator = func(op *openapi3.Operation) {
			mutator(op)
			oldMutator(op)
		}
	} else {
		b.Documentation.Mutator = mutator
	}
	return b
}

func (b *EndpointBuilder) WithDocSkip(skip bool) *EndpointBuilder {
	b.Documentation.Skip = skip
	return b
}

func (b *EndpointBuilder) WithPermissions(anyOf ...string) *EndpointBuilder {
	b.Permissions = anyOf
	return b
}

func (b *EndpointBuilder) WithInputs(inputs interface{}) *EndpointBuilder {
	if inputs != nil {
		b.Inputs = types.OptionalOf(inputs)
	} else {
		b.Inputs = types.OptionalEmpty[interface{}]()
	}
	return b
}

func (b *EndpointBuilder) WithOutputs(outputs interface{}) *EndpointBuilder {
	if outputs != nil {
		b.Outputs = types.OptionalOf(outputs)
	} else {
		b.Outputs = types.OptionalEmpty[interface{}]()
	}
	return b
}

func (b *EndpointBuilder) WithHandler(handler interface{}) *EndpointBuilder {
	b.Handler = handler
	return b
}

func (b *EndpointBuilder) Build() (*restops.Endpoint, error) {
	arch, ok := archetypes[b.Archetype]
	if !ok {
		return nil, errors.Errorf("EndpointBuilder could not find verb %q", arch)
	}

	e := restops.NewEndpoint(arch.Verb, b.Path).
		WithOperationId(b.Id).
		WithDocumentor(b.Documentation.Build()).
		WithHandler(b.Handler).
		WithPermissionAnyOf(b.Permissions...).
		WithResponseCodes(arch.Codes)

	if b.Inputs.IsPresent() {
		e.WithInputs(b.Inputs.ValueInterface())
	}

	if b.Outputs.IsPresent() {
		e.WithOutputs(b.Outputs.ValueInterface())
	}

	return e.Build()
}

func NewEndpointBuilder(archetype string, pathParts ...string) *EndpointBuilder {
	return &EndpointBuilder{
		Archetype: archetype,
		Path:      path.Join(pathParts...),
	}
}

func NewListEndpointBuilder(path string) *EndpointBuilder {
	return NewEndpointBuilder(EndpointArchetypeList, path)
}

func NewRetrieveEndpointBuilder(path string) *EndpointBuilder {
	return NewEndpointBuilder(EndpointArchetypeRetrieve, path)

}

func NewCreateEndpointBuilder(path string) *EndpointBuilder {
	return NewEndpointBuilder(EndpointArchetypeCreate, path)
}

func NewUpdateEndpointBuilder(path string) *EndpointBuilder {
	return NewEndpointBuilder(EndpointArchetypeUpdate, path)
}

func NewDeleteEndpointBuilder(path string) *EndpointBuilder {
	return NewEndpointBuilder(EndpointArchetypeDelete, path)
}

func NewCommandEndpointBuilder(path string) *EndpointBuilder {
	return NewEndpointBuilder(EndpointArchetypeCommand, path)
}

func NewAsyncCommandEndpointBuilder(path string) *EndpointBuilder {
	return NewEndpointBuilder(EndpointArchetypeAsyncCommand, path)
}
