package clienttest

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

type EndpointCall func(t *testing.T, executor integration.MsxServiceExecutor) (*integration.MsxResponse, error)

type EndpointTest struct {
	Call   EndpointCall
	Checks struct {
		Request  EndpointRequestCheck
		Endpoint ServiceEndpointCheck
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
}

func (c *EndpointTest) WithCall(call EndpointCall) *EndpointTest {
	c.Call = call
	return c
}

func (c *EndpointTest) WithEndpoints(endpoints map[string]integration.MsxServiceEndpoint) *EndpointTest {
	c.Endpoints = endpoints
	return c
}

func (c *EndpointTest) WithRequestPredicate(rp EndpointRequestPredicate) *EndpointTest {
	c.Checks.Request.Validators = append(c.Checks.Request.Validators, rp)
	return c
}

func (c *EndpointTest) WithRequestPredicates(predicates ...EndpointRequestPredicate) *EndpointTest {
	for _, predicate := range predicates {
		c.WithRequestPredicate(predicate)
	}
	return c
}

func (c *EndpointTest) WithEndpointPredicate(ep ServiceEndpointPredicate) *EndpointTest {
	c.Checks.Endpoint.Validators = append(c.Checks.Endpoint.Validators, ep)
	return c
}

func (c *EndpointTest) WithEndpointPredicates(predicates ...ServiceEndpointPredicate) *EndpointTest {
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

func (c *EndpointTest) Test(t *testing.T) {
	executor := new(integration.MockMsxServiceExecutor)

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

	response, err := c.Call(t, executor)

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
