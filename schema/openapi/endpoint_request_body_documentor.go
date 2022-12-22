// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package openapi

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/ops/restops"
	"github.com/swaggest/openapi-go/openapi3"
)

type EndpointRequestBodyDocumentor struct {
	Skip        bool
	RequestBody *openapi3.RequestBody
	Mutator     ops.DocumentElementMutator[openapi3.RequestBody]
}

func (d *EndpointRequestBodyDocumentor) WithSkip(skip bool) *EndpointRequestBodyDocumentor {
	d.Skip = skip
	return d
}

func (d *EndpointRequestBodyDocumentor) WithRequestBody(requestBody *openapi3.RequestBody) *EndpointRequestBodyDocumentor {
	d.RequestBody = requestBody
	return d
}

func (d *EndpointRequestBodyDocumentor) WithMutator(mutator ops.DocumentElementMutator[openapi3.RequestBody]) *EndpointRequestBodyDocumentor {
	d.Mutator = mutator
	return d
}

func (d *EndpointRequestBodyDocumentor) DocType() string {
	return DocType
}

func (d *EndpointRequestBodyDocumentor) documentFormField(formField restops.EndpointRequestBodyFormField) (schemaOrRef *openapi3.SchemaOrRef, err error) {
	var formFieldSchemaOrRef *openapi3.SchemaOrRef
	if formField.PortField != nil {
		formFieldSchemaOrRef = openApiSchemaFromPortField(formField.PortField)
	} else if formField.Payload.IsPresent() {
		formFieldSchemaOrRef, err = Reflect(formField.Payload.Value())
		if err != nil {
			return
		}
	} else if formField.Type != nil {
		formFieldSchema := SwaggerTypeToOpenApiSchema(formField.Type, formField.Format)
		if formFieldSchema != nil {
			formFieldSchemaOrRef = NewSchemaOrRefPtr(formFieldSchema)
		}
	}

	if formFieldSchemaOrRef == nil {
		return
	}

	var examplePtr *interface{}
	if formField.Example.IsPresent() {
		example := formField.Example.ValueInterface()
		examplePtr = &example
		if formFieldSchemaOrRef.Schema != nil {
			formFieldSchemaOrRef.Schema.Example = examplePtr
		}
	}


	schemaOrRef = formFieldSchemaOrRef
	return
}

func (d *EndpointRequestBodyDocumentor) documentFormFields(formFields []restops.EndpointRequestBodyFormField) (schemaOrRef *openapi3.SchemaOrRef, err error) {
	schemaOrRef = new(openapi3.SchemaOrRef)
	schema := ObjectSchema()

	for _, formField := range formFields {
		var formFieldSchemaOrRef *openapi3.SchemaOrRef
		formFieldSchemaOrRef, err = d.documentFormField(formField)
		if err != nil {
			return
		}

		schema.WithPropertiesItem(formField.Name, *formFieldSchemaOrRef)

		if formField.Required {
			schema.Required = append(schema.Required, formField.Name)
		}
	}

	schemaOrRef.Schema = schema
	return
}

func (d *EndpointRequestBodyDocumentor) Document(b *restops.EndpointRequestBody) (err error) {
	if d.Skip {
		return nil
	}

	var schemaOrRef *openapi3.SchemaOrRef
	if b.PortField != nil {
		schemaOrRef = openApiSchemaFromPortField(b.PortField)
	} else if b.Payload.IsPresent() {
		schemaOrRef, err = Reflect(b.Payload.Value())
		if err != nil {
			return
		}
	} else if len(b.FormFields) > 0 {
		schemaOrRef, err = d.documentFormFields(b.FormFields)
		if err != nil {
			return
		}
	}

	if schemaOrRef == nil {
		return nil
	}

	if d.RequestBody == nil {
		d.RequestBody = new(openapi3.RequestBody)
	}

	var examplePtr *interface{}
	if b.Example.IsPresent() {
		example := b.Example.ValueInterface()
		examplePtr = &example
	}

	d.RequestBody.
		WithRequired(b.Required).
		WithContentItem(b.Mime, openapi3.MediaType{
			Schema:  schemaOrRef,
			Example: examplePtr,
			//Encoding: b.Encoding,
		})

	if b.Description != "" {
		d.RequestBody.WithDescription(b.Description)
	}

	if d.Mutator != nil {
		d.Mutator(d.RequestBody)
	}

	return nil
}

func (d *EndpointRequestBodyDocumentor) Result() *openapi3.RequestBody {
	return d.RequestBody
}
