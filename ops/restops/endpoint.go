// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/lithammer/dedent"
	"github.com/pkg/errors"
	"github.com/swaggest/refl"
	"net/http"
	"path"
	"reflect"
)

const (
	exampleTag = "example"
)

type EndpointsProducer interface {
	Endpoints() (Endpoints, error)
}

type Endpoints []*Endpoint

type EndpointFunc func(e *Endpoint) error
type EndpointsPredicate func(e Endpoints) bool

func (e Endpoints) Each(fn EndpointFunc) error {
	for _, endpoint := range e {
		if err := fn(endpoint); err != nil {
			return err
		}
	}
	return nil
}

func NewEndpoint(method string, p ...string) *Endpoint {
	return &Endpoint{
		Method: method,
		Path:   path.Join(p...),
	}
}

type Endpoint struct {
	Method         string
	Path           string
	OperationID    string
	Description    string
	Summary        string
	Tags           []string
	Deprecated     bool
	Permissions    []string
	Func           types.Optional[interface{}]
	Inputs         types.Optional[reflect.Type]
	Outputs        types.Optional[reflect.Type]
	Request        EndpointRequest
	Response       EndpointResponse
	Handler        *types.Handler
	ErrorConverter ErrorConverter
	ops.Documentors[Endpoint]
}

func (e *Endpoint) WithDocumentor(d ...ops.Documentor[Endpoint]) *Endpoint {
	e.Documentors = e.Documentors.WithDocumentor(d...)
	return e
}

func (e *Endpoint) WithMethod(method string) *Endpoint {
	e.Method = method
	return e
}

func (e *Endpoint) WithPath(parts ...string) *Endpoint {
	e.Path = path.Join(parts...)
	return e
}

func (e *Endpoint) WithOperationId(name string) *Endpoint {
	e.OperationID = name
	return e
}

func (e *Endpoint) WithDescription(description string) *Endpoint {
	e.Description = dedent.Dedent(description)
	return e
}

func (e *Endpoint) WithSummary(summary string) *Endpoint {
	e.Summary = summary
	return e
}

func (e *Endpoint) WithTags(tags ...string) *Endpoint {
	e.Tags = append(e.Tags, tags...)
	return e
}

func (e *Endpoint) WithoutTags() *Endpoint {
	e.Tags = nil
	return e
}

func (e *Endpoint) WithDeprecated(deprecated bool) *Endpoint {
	e.Deprecated = deprecated
	return e
}

func (e *Endpoint) WithResponse(response EndpointResponse) *Endpoint {
	e.Response = response
	return e
}

func (e *Endpoint) WithRequest(request EndpointRequest) *Endpoint {
	e.Request = request
	return e
}

func (e *Endpoint) WithRequestParameter(parameter EndpointRequestParameter) *Endpoint {
	e.Request = e.Request.WithParameter(parameter)
	return e
}

func (e *Endpoint) WithResponseCodes(codes EndpointResponseCodes) *Endpoint {
	e.Response = e.Response.WithResponseCodes(codes)
	return e
}

func (e *Endpoint) WithResponseSuccessHeader(name string, header EndpointResponseHeader) *Endpoint {
	e.Response = e.Response.WithSuccessHeader(name, header)
	return e
}

func (e *Endpoint) WithResponseErrorHeader(name string, header EndpointResponseHeader) *Endpoint {
	e.Response = e.Response.WithErrorHeader(name, header)
	return e
}

func (e *Endpoint) WithResponseHeader(name string, header EndpointResponseHeader) *Endpoint {
	e.Response = e.Response.WithHeader(name, header)
	return e
}

func (e *Endpoint) WithPermissionAnyOf(perms ...string) *Endpoint {
	e.Permissions = perms
	return e
}

func (e *Endpoint) WithHandler(fn interface{}) *Endpoint {
	if fn != nil {
		e.Func = types.OptionalOf(fn)
	} else {
		e.Func = types.OptionalEmpty[interface{}]()
	}

	return e
}

func (e *Endpoint) WithHttpHandler(fn http.HandlerFunc) *Endpoint {
	if fn != nil {
		e.Func = types.OptionalOf[interface{}](fn)
	} else {
		e.Func = types.OptionalEmpty[interface{}]()
	}

	return e
}

func (e *Endpoint) WithInputs(inputs interface{}) *Endpoint {
	if inputs != nil {
		e.Inputs = types.OptionalOf(reflect.TypeOf(inputs))
	} else {
		e.Inputs = types.OptionalEmpty[reflect.Type]()
	}

	return e
}

func (e *Endpoint) WithOutputs(outputs interface{}) *Endpoint {
	if outputs != nil {
		e.Outputs = types.OptionalOf(reflect.TypeOf(outputs))
	} else {
		e.Outputs = types.OptionalEmpty[reflect.Type]()
	}

	return e
}

func (e *Endpoint) Build() (*Endpoint, error) {
	if !e.Func.IsPresent() {
		return nil, errors.Errorf("No handler set for operation %q", e.OperationID)
	}

	analyzer := &EndpointHandlerTypesAnalyzer{
		endpoint:    e,
		handlerFunc: e.Func.ValueInterface(),
	}

	if e.Inputs.IsPresent() {
		e.Request = e.Request.WithInputs(e.Inputs.Value())
	} else {
		ts := analyzer.ArgsTypeSet()
		if inputsType := analyzer.getInputsType(ts); inputsType != nil {
			portStructType := refl.DeepIndirect(inputsType)
			inputs := reflect.New(portStructType).Elem().Interface()
			e.WithInputs(inputs)
			e.Request = e.Request.WithInputs(inputs)
		}
	}

	if e.Outputs.IsPresent() {
		e.Response = e.Response.WithOutputs(e.Outputs.Value())
	} else {
		ts := analyzer.ReturnsTypeSet()
		if outputsType := analyzer.getOutputsType(ts); outputsType != nil {
			portStructType := refl.DeepIndirect(outputsType)
			outputs := reflect.New(portStructType).Elem().Interface()
			e.WithOutputs(outputs)
			e.Response = e.Response.WithOutputs(outputs)
		}
	}

	argsTypeSet := analyzer.ArgsTypeSet()
	returnsTypeSet := analyzer.ReturnsTypeSet()

	handler, err := types.NewHandler(e.Func.Value(),
		types.NewHandlerValueTypeReflector(
			argsTypeSet,
			types.DefaultHandlerArgumentValueTypeSet,
		),
		types.NewHandlerValueTypeReflector(
			returnsTypeSet,
			types.DefaultHandlerResultValueTypeSet,
		))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create handler")
	}

	e.Handler = handler

	return e, nil
}

// Builder

type EndpointBuilder interface {
	Build() (*Endpoint, error)
}

type EndpointBuilders []EndpointBuilder

func (b EndpointBuilders) Endpoints() (results Endpoints, err error) {
	var endpoint *Endpoint

	for _, endpointBuilder := range b {
		endpoint, err = endpointBuilder.Build()
		if err != nil {
			return nil, err
		}

		results = append(results, endpoint)
	}

	return results, nil
}
