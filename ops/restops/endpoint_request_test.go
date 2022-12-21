// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package restops

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestEndpointRequest_WithParameter(t *testing.T) {
	e := &EndpointRequest{}
	f := e.WithParameter(EndpointRequestParameter{
		Name: "parameter",
		In:   "path",
	})
	assert.NotEmpty(t, f.Parameters)
}

func TestEndpointRequest_WithValidator(t *testing.T) {
	e := &EndpointRequest{}
	v := func(p interface{}) (err error) {
		return nil
	}
	f := e.WithValidator(v)
	assert.NotNil(t, f.Validator)
}

func TestEndpointRequest_WithBody(t *testing.T) {
	e := &EndpointRequest{}
	b := EndpointRequestBody{
		Description: "body-description",
	}
	f := e.WithBody(b)
	assert.Equal(t, b, f.Body)
}

func TestEndpointRequest_WithInputs(t *testing.T) {
	e := &EndpointRequest{}
	i := struct{}{}
	f := e.WithInputs(i)
	assert.NotNil(t, f.Port)
}

func TestEndpointRequest_PatchParameter(t *testing.T) {
	e := &EndpointRequest{
		Parameters: []EndpointRequestParameter{
			{
				Name: "parameter",
				In:   FieldGroupHttpForm,
			},
		},
	}
	found := false
	f := e.PatchParameter("parameter", func(p EndpointRequestParameter) EndpointRequestParameter {
		found = true
		p.In = FieldGroupHttpPath
		return p
	})
	assert.Equal(t, true, found)
	assert.Len(t, f.Parameters, 1)
	assert.Equal(t, FieldGroupHttpPath, e.Parameters[0].In)
}

func TestEndpointRequest_Consumes(t *testing.T) {
	e := &EndpointRequest{
		Body: EndpointRequestBody{
			Mime: MediaTypeJson,
		},
	}

	assert.Equal(t, []string{MediaTypeJson}, e.Consumes())
}

func TestEndpointRequest_HasBody(t *testing.T) {
	assert.False(t, EndpointRequest{}.HasBody())
}

func TestEndpointRequest_withPortField(t *testing.T) {
	e := EndpointRequest{}
	pf := &ops.PortField{
		Name:     "FormField",
		Peer:     "form-field",
		Group:    FieldGroupHttpForm,
		Optional: false,
		PortType: "req",
		Options: map[string]string{
			"type":   "string",
			"format": "uuid",
		},
	}
	f := e.withPortField(pf)
	assert.NotEmpty(t, f.Body.FormFields)
}

func TestNewEndpointRequest(t *testing.T) {
	assert.Equal(t, EndpointRequest{}, NewEndpointRequest())
}

func TestEndpointRequestBodyFormField_Parameter(t *testing.T) {
	f := EndpointRequestBodyFormField{
		Name:      "form-field",
		PortField: new(ops.PortField),
	}
	p := f.Parameter()
	assert.Equal(t, f.Name, p.Name)
	assert.Equal(t, f.PortField, p.PortField)
}

func TestEndpointRequestBodyFormFieldFromPortField(t *testing.T) {
	pf := ops.PortField{
		Name:     "FormField",
		Peer:     "form-field",
		Group:    FieldGroupHttpForm,
		Optional: false,
		PortType: "req",
		Options: map[string]string{
			"type":   "string",
			"format": "uuid",
		},
	}

	r := EndpointRequestBodyFormField{
		Name:      "form-field",
		Required:  true,
		Type:      types.NewStringPtr("string"),
		Format:    types.NewStringPtr("uuid"),
		PortField: &pf,
	}

	f := EndpointRequestBodyFormFieldFromPortField(&pf)
	assert.True(t,
		reflect.DeepEqual(r, f),
		testhelpers.Diff(r, f))
}

func TestEndpointRequestBody_WithDocumentor(t *testing.T) {
	e := EndpointRequestBody{}
	f := e.WithDocumentor(TestDocumentor[EndpointRequestBody]{})
	assert.NotNil(t, f)
}

func TestEndpointRequestBody_WithFormField(t *testing.T) {
	e := EndpointRequestBody{}
	f := EndpointRequestBodyFormField{
		Name:      "form-field",
		PortField: new(ops.PortField),
	}
	g := e.WithFormField(f)
	assert.NotEmpty(t, g.FormFields)
}

func TestEndpointRequestBody_HasFormField(t *testing.T) {
	e := EndpointRequestBody{}
	f := EndpointRequestBodyFormField{
		Name:      "form-field",
		PortField: new(ops.PortField),
	}
	g := e.WithFormField(f)
	assert.True(t, g.HasFormField())
}

func TestEndpointRequestBodyFromPortField(t *testing.T) {
	pf := &ops.PortField{
		Name:  "Body",
		Peer:  "body",
		Group: FieldGroupHttpBody,
		Type: ops.PortFieldType{
			Type: reflect.TypeOf(""),
		},
		Optional: false,
		PortType: "req",
		Options: map[string]string{
			"description": "The Payload",
		},
	}

	bf := EndpointRequestBody{
		Description: "The Payload",
		Required:    true,
		Mime:        MediaTypeJson,
		Payload:     types.OptionalOf[interface{}](types.NewStringPtr("")),
		PortField:   pf,
	}

	e := EndpointRequestBodyFromPortField(pf)
	assert.True(t,
		reflect.DeepEqual(bf, e),
		testhelpers.Diff(bf, e))
}

func TestEndpointRequestParameter_WithDocumentor(t *testing.T) {
	e := EndpointRequestParameter{}
	f := e.WithDocumentor(TestDocumentor[EndpointRequestParameter]{})
	assert.NotNil(t, f)
}

func TestEndpointRequestParameter_WithDescription(t *testing.T) {
	e := EndpointRequestParameter{}
	v := "description"
	f := e.WithDescription(v)
	assert.Equal(t, v, *f.Description)
}

func TestEndpointRequestParameter_WithRequired(t *testing.T) {
	e := EndpointRequestParameter{}
	v := true
	f := e.WithRequired(v)
	assert.Equal(t, v, *f.Required)
}

func TestEndpointRequestParameter_WithDeprecated(t *testing.T) {
	e := EndpointRequestParameter{}
	v := true
	f := e.WithDeprecated(v)
	assert.Equal(t, v, *f.Deprecated)
}

func TestEndpointRequestParameter_WithStyle(t *testing.T) {
	e := EndpointRequestParameter{}
	v := "csv"
	f := e.WithStyle(v)
	assert.Equal(t, v, *f.Style)
}

func TestEndpointRequestParameter_WithAllowEmptyValue(t *testing.T) {
	e := EndpointRequestParameter{}
	v := true
	f := e.WithAllowEmptyValue(v)
	assert.Equal(t, v, *f.AllowEmptyValue)
}

func TestEndpointRequestParameter_WithExplode(t *testing.T) {
	e := EndpointRequestParameter{}
	v := true
	f := e.WithExplode(v)
	assert.Equal(t, v, *f.Explode)
}

func TestEndpointRequestParameter_WithAllowReserved(t *testing.T) {
	e := EndpointRequestParameter{}
	v := true
	f := e.WithAllowReserved(v)
	assert.Equal(t, v, *f.AllowReserved)
}

func TestEndpointRequestParameter_WithReference(t *testing.T) {
	e := EndpointRequestParameter{}
	v := "csv"
	f := e.WithReference(v)
	assert.Equal(t, v, *f.Reference)
}

func TestEndpointRequestParameter_WithType(t *testing.T) {
	e := EndpointRequestParameter{}
	v := "csv"
	f := e.WithType(v)
	assert.Equal(t, v, *f.Type)
}

func TestEndpointRequestParameter_WithFormat(t *testing.T) {
	e := EndpointRequestParameter{}
	v := "csv"
	f := e.WithFormat(v)
	assert.Equal(t, v, *f.Format)
}

func TestEndpointRequestParameter_WithExample(t *testing.T) {
	e := EndpointRequestParameter{}
	v := "example"
	f := e.WithExample(v)
	assert.Equal(t, types.OptionalOf[interface{}](v), f.Example)
}

func TestEndpointRequestParameter_WithPayload(t *testing.T) {
	e := EndpointRequestParameter{}
	v := "example"
	f := e.WithPayload(v)
	assert.Equal(t, types.OptionalOf[interface{}](v), f.Payload)
}

func TestNewEndpointRequestParameter(t *testing.T) {
	e := NewEndpointRequestParameter("parameter", FieldGroupHttpCookie)
	p := EndpointRequestParameter{
		Name: "parameter",
		In:   "cookie",
	}

	assert.True(t,
		reflect.DeepEqual(p, e),
		testhelpers.Diff(p, e))
}

func TestEndpointRequestParameterFromPortField(t *testing.T) {
	pf := &ops.PortField{
		Name:  "MyCookie",
		Peer:  "my-cookie",
		Group: FieldGroupHttpCookie,
		Type: ops.PortFieldType{
			Type: reflect.TypeOf(""),
		},
		Optional: false,
		PortType: "req",
		Options: map[string]string{
			"description": "The Cookie",
		},
	}

	bf := EndpointRequestParameter{
		Name:        "my-cookie",
		In:          FieldGroupHttpCookie,
		Description: types.NewStringPtr("The Cookie"),
		Required:    types.NewBoolPtr(!pf.Optional),
		PortField:   pf,
		Style:       types.NewStringPtr("form"),
		Explode:     types.NewBoolPtr(false),
	}

	e := EndpointRequestParameterFromPortField(pf)

	assert.True(t,
		reflect.DeepEqual(bf, e),
		testhelpers.Diff(bf, e))
}

func TestPathParameter(t *testing.T) {
	pp := PathParameter("deviceId", "Device Id")
	p := EndpointRequestParameter{
		Name:        "deviceId",
		In:          FieldGroupHttpPath,
		Required:    types.NewBoolPtr(true),
		Description: types.NewStringPtr("Device Id"),
	}
	assert.True(t,
		reflect.DeepEqual(p, pp),
		testhelpers.Diff(p, pp))
}

func TestQueryParameter(t *testing.T) {
	pp := QueryParameter("deviceId", "Device Id")
	p := EndpointRequestParameter{
		Name:        "deviceId",
		In:          FieldGroupHttpQuery,
		Description: types.NewStringPtr("Device Id"),
	}
	assert.True(t,
		reflect.DeepEqual(p, pp),
		testhelpers.Diff(p, pp))
}

func TestHeaderParameter(t *testing.T) {
	pp := HeaderParameter("deviceId", "Device Id")
	p := EndpointRequestParameter{
		Name:        "deviceId",
		In:          FieldGroupHttpHeader,
		Description: types.NewStringPtr("Device Id"),
	}
	assert.True(t,
		reflect.DeepEqual(p, pp),
		testhelpers.Diff(p, pp))
}

func TestCookieParameter(t *testing.T) {
	pp := CookieParameter("deviceId", "Device Id")
	p := EndpointRequestParameter{
		Name:        "deviceId",
		In:          FieldGroupHttpCookie,
		Description: types.NewStringPtr("Device Id"),
	}
	assert.True(t,
		reflect.DeepEqual(p, pp),
		testhelpers.Diff(p, pp))
}
