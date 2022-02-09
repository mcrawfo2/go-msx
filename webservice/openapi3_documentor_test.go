package webservice

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
	"github.com/swaggest/jsonschema-go"
	"net/http"
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

func newUpdateEndpoint() Endpoint {
	type inputs struct {
		PolicyName string              `req:"path"`
		Name       *string             `req:"query,description=Name" pattern:"^[a-z][A-Za-z0-9]{0,126}$"`
		Request    PolicyUpdateRequest `req:"body"`
	}

	type outputs struct {
		Response PolicyResponse `resp:"body"`
	}

	paramPolicyName := PathParameter("policyName", "Secret Policy Name").
		WithSchema(NewSchemaOrRef(StringSchema().WithPattern(`^[a-z][A-Za-z0-9]{0,126}$`))).
		WithExample("ciscoIosSshPassword")

	return NewEndpoint(http.MethodPut, pathPrefixPolicy, pathSuffixPolicyName).
		WithOperationId("updateSecretPolicy").
		WithDescription("Update the specified Secret Policy").
		WithPermissionAnyOf(permissionNameViewSecretPolicies, permissionNameManageSecretPolicies).
		WithRequest(NewEndpointRequest().
			WithPortStruct(inputs{}).
			WithValidator(func(p interface{}) (err error) {
				return types.ErrorMap{
					"policyName": errors.New("some error"),
					"name":       nil,
				}
			}).
			WithParameter(paramPolicyName).
			PatchParameter("name", func(p EndpointRequestParameter) EndpointRequestParameter {
				return p.WithExample("bob")
			})).
		WithResponse(NewEndpointResponse().WithPortStruct(outputs{})).
		WithController(
			func(req *restful.Request) (body interface{}, err error) {
				a := InputsFromRequest(req).(*inputs)

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

func TestNewEndpoint(t *testing.T) {
	Reflector.Spec.Info.Title = "Testing Spec"
	Reflector.Spec.Info.Version = "v2"

	endpoint := newUpdateEndpoint()
	adapter := NewEndpointOpenApiDocumentor(endpoint)
	err := Reflector.Spec.SetupOperation(endpoint.Method, endpoint.Path, adapter.DocumentOpenApiOperation)
	fmt.Println(err)

	yamlBytes, _ := Reflector.Spec.MarshalYAML()
	fmt.Println(string(yamlBytes))
}
