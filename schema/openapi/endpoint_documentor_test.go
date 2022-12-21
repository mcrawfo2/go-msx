// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package openapi

import (
	"cto-github.cisco.com/NFV-BU/go-msx/ops"
	"cto-github.cisco.com/NFV-BU/go-msx/ops/restops"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/swaggest/jsonschema-go"
	"github.com/swaggest/openapi-go/openapi3"
	"net/http"
	"reflect"
	"testing"
	"time"
)

const (
	permissionNameManageSecretPolicies = "MANAGE_SECRET_POLICIES"
	permissionNameViewSecretPolicies   = "VIEW_SECRET_POLICIES"
	pathPrefixPolicy                   = "/v2/secrets/policy"
	pathSuffixPolicyName               = "/{policyName}"
)

type PolicyUpdateRequest struct {
	AgingRuleEnabled   bool           `json:"agingRuleEnabled"`
	HistoryCount       int            `json:"historyCount"`
	HistoryExpireAfter types.Duration `json:"historyExpireAfter"`
	StartDate          types.Time     `json:"startDate"`
	EndDate            *types.Time    `json:"endDate"`
	Identifier         types.UUID     `json:"identifier"`
	Name               string         `json:"name" pattern:"^[a-z][A-Za-z0-9]{0,126}$"`
}

func (p PolicyUpdateRequest) PrepareJSONSchema(s *jsonschema.Schema) error {
	s.
		WithExamples(PolicyUpdateRequest{
			AgingRuleEnabled:   true,
			HistoryCount:       10,
			HistoryExpireAfter: types.Duration(time.Minute*30 + time.Hour*2),
			StartDate:          types.NewTime(time.Now()),
			EndDate:            nil,
			Identifier:         types.MustNewUUID(),
		})
	return nil
}

type PolicyResponse struct {
	PolicyName         string         `json:"policyName"`
	AgingRuleEnabled   bool           `json:"agingRuleEnabled"`
	HistoryCount       int            `json:"historyCount"`
	HistoryExpireAfter types.Duration `json:"historyExpireAfter"`
	StartDate          types.Time     `json:"startDate"`
	EndDate            *types.Time    `json:"endDate"`
}

type PolicyName string

func (p PolicyName) PrepareJSONSchema(schema *jsonschema.Schema) error {
	schema.Pattern = types.NewStringPtr(`^[a-z][A-Za-z0-9]{0,126}$`)
	return nil
}

func newUpdateEndpoint() *restops.Endpoint {
	type inputs struct {
		PolicyName string              `req:"path" description:"Name" pattern:"^[a-z][A-Za-z0-9]{0,126}$"`
		Request    PolicyUpdateRequest `req:"body"`
	}

	type outputs struct {
		Response PolicyResponse `resp:"body"`
	}

	return restops.NewEndpoint(http.MethodPut, pathPrefixPolicy, pathSuffixPolicyName).
		WithOperationId("updateSecretPolicy").
		WithSummary("Update the specified Secret Policy").
		WithDescription(`
			**Detailed Description**:
                - first item
				- second item
		`).
		WithDeprecated(true).
		WithPermissionAnyOf(permissionNameViewSecretPolicies, permissionNameManageSecretPolicies).
		WithRequest(restops.NewEndpointRequest().
			WithInputs(inputs{}).
			WithParameter(restops.
				NewEndpointRequestParameter("MyCookie", restops.FieldGroupHttpCookie).
				WithDescription("Some cookie")).
			WithValidator(func(p interface{}) (err error) {
				return types.ErrorMap{
					"policyName": errors.New("some error"),
					"name":       nil,
				}
			})).
		WithResponse(restops.NewEndpointResponse().WithOutputs(outputs{})).
		WithHandler(
			func(req *restful.Request) (body interface{}, err error) {
				a := restops.InputsFromRequest(req).(*inputs)

				return PolicyResponse{
					PolicyName:         a.PolicyName,
					AgingRuleEnabled:   a.Request.AgingRuleEnabled,
					HistoryCount:       a.Request.HistoryCount,
					HistoryExpireAfter: a.Request.HistoryExpireAfter,
					StartDate:          a.Request.StartDate,
					EndDate:            a.Request.EndDate,
				}, nil
			})
}

func newDeleteEndpoint() *restops.Endpoint {
	type inputs struct {
		PolicyName string `req:"path" description:"Name" pattern:"^[a-z][A-Za-z0-9]{0,126}$"`
	}

	type outputs struct {
		Code int `resp:"code" enum:"204,400,401,403,404"`
	}

	return restops.NewEndpoint(http.MethodPut, pathPrefixPolicy, pathSuffixPolicyName).
		WithOperationId("deleteSecretPolicy").
		WithSummary("Delete the specified Secret Policy").
		WithPermissionAnyOf(permissionNameManageSecretPolicies).
		WithRequest(restops.
			NewEndpointRequest().
			WithInputs(inputs{})).
		WithResponse(restops.
			NewEndpointResponse().WithOutputs(outputs{})).
		WithHandler(
			func(req *restful.Request) (err error) {
				return nil
			})
}

func TestEndpointDocumentor_Document_Update(t *testing.T) {
	var doc ops.Documentor[restops.Endpoint] = new(EndpointDocumentor)

	endpoint := newUpdateEndpoint()
	assert.NotNil(t, endpoint)

	err := doc.Document(endpoint)
	assert.NoError(t, err)

	op := doc.(ops.DocumentResult[openapi3.Operation]).Result()
	assert.NotNil(t, op)
}

func TestEndpointDocumentor_Document_Delete(t *testing.T) {
	var doc ops.Documentor[restops.Endpoint] = new(EndpointDocumentor)

	endpoint := newDeleteEndpoint()
	assert.NotNil(t, endpoint)

	err := doc.Document(endpoint)
	assert.NoError(t, err)

	op := doc.(ops.DocumentResult[openapi3.Operation]).Result()
	assert.NotNil(t, op)
}

func TestEndpointDocumentor_Document(t *testing.T) {
	tests := []struct {
		name     string
		doc      *EndpointDocumentor
		endpoint restops.Endpoint
		want     *openapi3.Operation
		wantErr  bool
	}{
		{
			name:    "Skip",
			doc:     new(EndpointDocumentor).WithSkip(true),
			want:    nil,
			wantErr: false,
		},
		{
			name: "Operation",
			doc: new(EndpointDocumentor).WithOperation(
				new(openapi3.Operation).
					WithDescription("Description")),
			endpoint: restops.
				Endpoint{OperationID: "Operation"},
			want: &openapi3.Operation{
				ID:          types.NewStringPtr("Operation"),
				Description: types.NewStringPtr("Description"),
			},
			wantErr: false,
		},
		{
			name: "Mutator",
			doc: new(EndpointDocumentor).WithMutator(
				func(p *openapi3.Operation) {
					p.WithDescription("Description")
				}),
			endpoint: restops.
				Endpoint{OperationID: "Mutator"},
			want: &openapi3.Operation{
				ID:          types.NewStringPtr("Mutator"),
				Description: types.NewStringPtr("Description"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.doc
			if e == nil {
				e = new(EndpointDocumentor)
			}

			err := e.Document(&tt.endpoint)
			if tt.wantErr {
				assert.Error(t, err)
				return
			} else {
				assert.NoError(t, err)
			}

			got := e.Result()
			if tt.want == nil {
				assert.Nil(t, got)
				return
			}

			assert.True(t,
				reflect.DeepEqual(tt.want, got),
				testhelpers.Diff(tt.want, got))
		})
	}
}

func TestEndpointDocumentor_DocType(t *testing.T) {
	doc := new(EndpointDocumentor)
	assert.Equal(t, DocType, doc.DocType())
}

func TestEndpointDocumentorBuilder_Build(t *testing.T) {
	b := EndpointDocumentorBuilder{
		Skip:      true,
		Operation: new(openapi3.Operation),
		Mutator:   nil,
	}

	want := &EndpointDocumentor{
		Skip:      b.Skip,
		Operation: b.Operation,
		Mutator:   b.Mutator,
	}

	got := b.Build()
	assert.True(t,
		reflect.DeepEqual(want, got),
		testhelpers.Diff(want, got))
}
