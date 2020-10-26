# MSX Integration

The MSX Integration module provides the ability to interact with different types of JSON REST APIs.

## Components

MSX Integration uses various components to describe and execute REST API calls, including:

- Services – Represents a single MSX microservice or External service
- Endpoints – Represents an API exposed by a service 
- Request – Options and data for a single REST call
- Interceptor – Transform a request
- Response – Result of a single REST call
- Envelope – Body content wrapper of a response
- Error – Error result of a response

### Services

An MSX Integration Service represents a single MSX Microservice or External service.

#### MSX Microservice

The `MsxService` object is used to represent a single MSX Microservice.  The microservice can be one of the following types:

- `ServiceTypeMicroservice` - A standard MSX microservice deployed as one or more stateless instances.  Construct a new service object using `NewMsxService`.
- `ServiceTypeProbe` - An MSX monitoring probe deployed as one or more stateful instances.  Construct a new service object using `NewProbeService`.

`MsxService` objects contain a map of Endpoints which Clients can use to identify the name, path, and method of a particular operation. 

#### External Service

The `ExternalService` object is used to represent an API outside of MSX.  This could include a Controller or other Web service API.  Construct a new service object using `NewExternalService` or `NewExternalServiceFromUrl`.

### Endpoints

An Endpoint represents a single API exposed by a Service.   For example, each MSX microservices exposes a Health API endpoint at `GET {contextPath}/admin/health`.  An endpoint is composed of three parts:

- Operation Name - The endpoint name (for logging, stats, and tracing)
- Method - The HTTP method used to access the endpoint
- Path - The path (or path template) used to access the endpoint

For `MSXService`, endpoints are defined as `integration.MsxServiceEndpoint` and provided in the constructor as a map:

```go
endpoints := map[string]integration.MsxServiceEndpoint{
	"getAdminHealth": {Method: http.MethodGet, Path: "/admin/health"},
}
```

For `ExternalService`, endpoints are defined as `integration.Endpoint` and passed into `Request()`:

```go
endpoint := integration.Endpoint{
  Name: "getOrganizationDevices",
  Method: http.MethodGet,
  Path: "/organization/{organizationId}/devices",
}
```

### Request

A Request object specifies the options and data for a single REST call.  

When using `MsxService`, an `MsxEndpointRequest` object accepts the following:

- Endpoint name
- Endpoint path parameters
- Headers
- Query Parameters
- Body
- Token Expected in Request
- Envelope Expected in Response

When using `ExternalService` , the `Request()` method accepts:

- Endpoint
- Endpoint path parameters
- Headers
- Query parameters
- Body

### Interceptor

An Interceptor is a function which may transform or modify any of the Request attributes to form a new request.  By default, MSX Service applies interceptors as follows:

- `ServiceTypeMicroservice` 	
  - MSX Service will resolve the hostname (microservice name) to a single healthy microservice instance
  - MSX Service will inject the current authorization token header
- `ServiceTypeProbe` 
  - MSX Service will inject the current authorization token header

External services do not have default interceptors.  You can add interceptors to an External service using the `AddInterceptor()` method.

### Response

