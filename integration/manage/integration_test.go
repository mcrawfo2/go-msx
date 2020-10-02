package manage

import (
	"cto-github.cisco.com/NFV-BU/go-msx/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestIntegration_GetAdminHealth(t *testing.T) {
	executor := new(integration.MockMsxServiceExecutor)
	executor.On("Execute", mock.Anything).Run(func(args mock.Arguments) {
		request, _ := args.Get(0).(*integration.MsxEndpointRequest)
		assert.Equal(t, endpointNameGetAdminHealth, request.EndpointName)
	}).Return(&integration.MsxResponse{
		StatusCode: 200,
		Status:     "200 OK",
		Headers:    nil,
		Envelope:   nil,
		Payload: &integration.HealthDTO{
			Status: "Up",
		},
		Body:       []byte("{}"),
		BodyString: "{}",
	}, nil)

	api := NewIntegrationWithExecutor(executor)

	healthResult, err := api.GetAdminHealth()
	assert.NoError(t, err)
	assert.Equal(t, "Up", healthResult.Payload.Status)

	executor.AssertExpectations(t)
}
