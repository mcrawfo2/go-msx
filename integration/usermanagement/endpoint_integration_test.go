// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package usermanagement

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/clienttest"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"reflect"
	"testing"
)

type EndpointCall func(t *testing.T, executor integration.MsxContextServiceExecutor) (*integration.MsxResponse, error)
type MultiTenantEndpointCall func(t *testing.T, executor integration.MsxContextServiceExecutor) (*integration.MsxResponse, []types.UUID, error)

type EndpointTest struct {
	Call      EndpointCall
	MultiTenantResultCall MultiTenantEndpointCall
	Injectors types.ContextInjectors
	Checks    struct {
		Request  clienttest.EndpointRequestCheck
		Endpoint clienttest.ServiceEndpointCheck
	}
	Endpoints        map[string]integration.MsxServiceEndpoint
	Response         integration.MsxResponse
	ResponseEnvelope *integration.MsxEnvelope
	ResponsePayload  interface{}
	ResponseError    error
	Errors           struct {
		Request  []error
		Endpoint []error
	}
	Tenants []types.UUID
}

func (c *EndpointTest) WithCall(call EndpointCall) *EndpointTest {
	c.Call = call
	return c
}

func (c *EndpointTest) WithMultiTenantResultCall(call MultiTenantEndpointCall) *EndpointTest {
	c.MultiTenantResultCall = call
	return c
}

func (c *EndpointTest) WithTenants(tenants []types.UUID) *EndpointTest {
	c.Tenants = tenants
	return c
}

func (c *EndpointTest) WithInjector(injector types.ContextInjector) *EndpointTest {
	c.Injectors = append(c.Injectors, injector)
	return c
}

func (c *EndpointTest) WithEndpoints(endpoints map[string]integration.MsxServiceEndpoint) *EndpointTest {
	c.Endpoints = endpoints
	return c
}

func (c *EndpointTest) WithRequestPredicate(rp clienttest.EndpointRequestPredicate) *EndpointTest {
	c.Checks.Request.Validators = append(c.Checks.Request.Validators, rp)
	return c
}

func (c *EndpointTest) WithRequestPredicates(predicates ...clienttest.EndpointRequestPredicate) *EndpointTest {
	for _, predicate := range predicates {
		c.WithRequestPredicate(predicate)
	}
	return c
}

func (c *EndpointTest) WithEndpointPredicate(ep clienttest.ServiceEndpointPredicate) *EndpointTest {
	c.Checks.Endpoint.Validators = append(c.Checks.Endpoint.Validators, ep)
	return c
}

func (c *EndpointTest) WithEndpointPredicates(predicates ...clienttest.ServiceEndpointPredicate) *EndpointTest {
	for _, predicate := range predicates {
		c.WithEndpointPredicate(predicate)
	}
	return c
}

func (c *EndpointTest) WithResponseStatus(statusCode int) *EndpointTest {
	c.Response.StatusCode = statusCode
	c.Response.Status = http.StatusText(statusCode)
	if c.ResponseEnvelope != nil {
		c.ResponseEnvelope.HttpStatus = integration.GetSpringStatusNameForCode(c.Response.StatusCode)
		if c.Response.StatusCode < 299 {
			c.ResponseEnvelope.Success = true
		}
	}
	return c
}

func (c *EndpointTest) WithResponseEnvelope() *EndpointTest {
	c.ResponseEnvelope = new(integration.MsxEnvelope)
	c.ResponseEnvelope.Payload = c.ResponsePayload
	if c.Response.StatusCode != 0 {
		c.WithResponseStatus(c.Response.StatusCode)
	}
	c.Response.Envelope = c.ResponseEnvelope
	c.withResponseBody()
	return c
}

func (c *EndpointTest) WithResponsePayload(payload interface{}) *EndpointTest {
	c.ResponsePayload = payload
	c.Response.Payload = c.ResponsePayload
	if c.ResponseEnvelope != nil {
		c.WithResponseEnvelope()
	}
	c.withResponseBody()
	return c
}

func (c *EndpointTest) withResponseBody() {
	if c.ResponseEnvelope != nil {
		c.Response.Body, _ = json.Marshal(c.ResponseEnvelope)
	} else {
		c.Response.Body, _ = json.Marshal(c.ResponsePayload)
	}
}

func (c *EndpointTest) getExecutor() *integration.MockMsxContextServiceExecutor {
	executor := new(integration.MockMsxContextServiceExecutor)
	executor.On("Context").Return(c.Injectors.Inject(context.Background())).Maybe()

	call := executor.On("Execute", mock.AnythingOfType("*integration.MsxEndpointRequest")).
		Run(func(args mock.Arguments) {
			endpointRequest, _ := args.Get(0).(*integration.MsxEndpointRequest)
			c.Errors.Request = c.Checks.Request.Check(endpointRequest)

			endpoint, ok := c.Endpoints[endpointRequest.EndpointName]
			if !ok {
				c.Errors.Endpoint = []error{
					errors.Errorf("Endpoint %q not defined", endpointRequest.EndpointName),
				}
			} else {
				c.Errors.Endpoint = c.Checks.Endpoint.Check(endpoint)
			}
		})

	if c.ResponseError != nil {
		call.Return(&c.Response, c.ResponseError)
	} else {
		call.Return(&c.Response, nil)
	}

	return executor
}

func (c *EndpointTest) assertResponse(t *testing.T, executor *integration.MockMsxContextServiceExecutor, response *integration.MsxResponse, err error) {
	if c.ResponseError != nil {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
	}

	assert.Equal(t, c.Response, *response)

	executor.AssertExpectations(t)

	testhelpers.ReportErrors(t, "Request", c.Errors.Request)
	testhelpers.ReportErrors(t, "Endpoint", c.Errors.Endpoint)
}

func (c *EndpointTest) Test(t *testing.T) {
	executor := c.getExecutor()

	response, err := c.Call(t, executor)

	c.assertResponse(t, executor, response, err)
}

func (c *EndpointTest) TestMultiTenantResult(t *testing.T) {
	executor := c.getExecutor()

	response, tenants, err := c.MultiTenantResultCall(t, executor)

	c.assertResponse(t, executor, response, err)

	if !reflect.DeepEqual(tenants, c.Tenants) {
		t.Error("unexpected result: ")
		t.Errorf(testhelpers.Diff(c.Tenants, tenants))
	}
}

