package integration

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
)

var (
	logger = log.NewLogger("msx.integration")
)

type MsxEndpointRequest struct {
	EndpointName       string
	EndpointParameters map[string]string
	Headers            http.Header
	QueryParameters    url.Values
	Body               []byte
	ExpectEnvelope     bool
	NoToken            bool
	Payload            interface{}
	ErrorPayload       interface{}
}

type ServiceType string

const (
	ServiceTypeMicroservice     ServiceType = "managedMicroservice"
	ServiceTypeResourceProvider ServiceType = "resourceProvider"
)

type MsxServiceEndpoint struct {
	Method string
	Path   string
}

type MsxService struct {
	serviceName string
	endpoints   map[string]MsxServiceEndpoint
	serviceType ServiceType
	ctx         context.Context
}

func (v *MsxService) Target(endpointName string) (Target, error) {
	endpoint, ok := v.endpoints[endpointName]
	if !ok {
		return Target{}, errors.Errorf("Endpoint %q not found for service %q", endpointName, v.serviceName)
	}

	return Target{
		ServiceName:  v.serviceName,
		ServiceType:  v.serviceType,
		EndpointName: endpointName,
		Method:       endpoint.Method,
		Path:         endpoint.Path,
	}, nil
}

func (v *MsxService) ServiceRequest(request *MsxEndpointRequest) (*MsxRequest, error) {
	target, err := v.Target(request.EndpointName)
	if err != nil {
		return nil, err
	}

	return &MsxRequest{
		Target:             target,
		EndpointParameters: request.EndpointParameters,
		Headers:            request.Headers,
		QueryParameters:    request.QueryParameters,
		Body:               request.Body,
		ExpectEnvelope:     request.ExpectEnvelope,
		NoToken:            request.NoToken,
		Payload:            request.Payload,
		ErrorPayload:       request.ErrorPayload,
	}, nil
}

func (v *MsxService) Execute(request *MsxEndpointRequest) (response *MsxResponse, err error) {
	serviceRequest, err := v.ServiceRequest(request)
	if err != nil {
		return nil, err
	}
	return serviceRequest.Execute(v.ctx)
}

func (v *MsxService) ExecuteWithContext(ctx context.Context, request *MsxEndpointRequest) (response *MsxResponse, err error) {
	serviceRequest, err := v.ServiceRequest(request)
	if err != nil {
		return nil, err
	}
	return serviceRequest.Execute(ctx)
}

func (v *MsxService) Context() context.Context {
	return v.ctx
}

func NewMsxService(ctx context.Context, serviceName string, endpoints map[string]MsxServiceEndpoint) *MsxService {
	return &MsxService{
		serviceName: serviceName,
		endpoints:   endpoints,
		ctx:         ctx,
		serviceType: ServiceTypeMicroservice,
	}
}

func NewMsxServiceResourceProvider(ctx context.Context, serviceName string, endpoints map[string]MsxServiceEndpoint) *MsxService {
	return &MsxService{
		serviceName: serviceName,
		endpoints:   endpoints,
		serviceType: ServiceTypeResourceProvider,
		ctx:         ctx,
	}
}
