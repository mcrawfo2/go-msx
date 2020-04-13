# REST API Controller

MSX promotes the usage of the common Controller > Service > Repository layered architecture within microservices.

The role of the Controller is to accept REST-based API requests from callers (UI, swagger, other microservices),
and route them to the service.

## Defining the Controller structure

To define a controller, create a standard Go structure with fields for its required dependencies:

```go
type productController struct {
    productService   *productService
    productConverter *productConverter
}
```

This example shows two common dependencies:

- Service
    - The service is responsible for responding to the requests.  The controller acts as an HTTP gateway
      to the service functionality.
- Converter
    - The converter transforms data transfer objects (requests and response) to and from domain models.

## Implementing the RestController interface

For registration with the web server, the `webservice.RestController` interface defines a single required method, `Routes`.
Add a standard implementation to your controller, for example:

```go
func (c *productController) Routes(svc *restful.WebService) {
	tag := webservice.TagDefinition("Products", "Products Controller")
	webservice.Routes(svc, tag,
		c.listProducts,
		c.getProduct,
		c.createProduct,
		c.updateProduct,
		c.deleteProduct)
}
```

This implementation demonstrates:
- Adding each endpoint implementation to the supplied WebService.  
- Tagging the routes for Swagger.  This allows the Swagger UI to show the human-readable controller name and group endpoints properly.
  Note that the tag does not have to be unique, and can be declared at module level to be used across multiple controllers (eg v1, v2).
  This will show all of the endpoints from the chosen controllers in a single group.

## Implementing an Endpoint

Each endpoint on your controller should be declared inside its own method.  Here's an example implementation of a List endpoint
for the Products controller:

```go
var viewPermissionFilter   = webservice.PermissionsFilter(rbac.PermissionViewProduct)

func (c *productController) listProducts(svc *restful.WebService) *restful.RouteBuilder {
    type params struct {
        Category *string `req:"query"`
    }

	return svc.GET("").
		Operation("listProducts").
		Doc("Retrieve the list of products, optionally filtering by the specified criteria.").
		Do(webservice.StandardList).
		Do(webservice.ResponsePayload(api.ProductListResponse{})).
		Filter(viewPermissionFilter).
		To(webservice.Controller(
			func(req *restful.Request) (body interface{}, err error) {
                params = webservice.Params(req).(*params)

                products, err := c.productService.ListProducts(req.Request.Context(), params.Category)
				if err != nil {
					return nil, err
				}

				return c.productConverter.ToProductListResponse(products), nil
			}))
}
```

Here we are declaring the endpoint:
- `type params ...`
    - accepts an optional string parameter `category` as a query parameter
- `svc.GET`
    - will use the GET HTTP method
- `GET("")`
    - has the same path as the controller
- `Operation("listProducts")`
    - has the operation name `listProducts`.  This will appear in tracing, logs, and in the swagger definition.
- `Doc("...")`
    - has the supplied description in the Swagger UI
- `Do(webservice.StandardList)`
    - is an implementation of a List Collection endpoint.  Returns 200 by default.
- `Do(webservice.ResponsePayload(api.ProductListResponse{}))`
    - will return the specified response DTO, wrapped _inside_ an MsxEnvelope object.
- `Filter(viewPermissionFilter)`
    - will check that callers have the "VIEW_PRODUCT" permission, as defined in
      the viewPermissionFilter object
- `To(webservice.Controller(func ...))`
    - will execute the supplied function when this endpoint is called

These are just some of the possible route building functions available in go-msx.  You may use the go-restful routing functions, along with the [go-msx routing functions](routes.go) to define many aspects of your endpoint.

## Implementing a Constructor

To allow instantiation of your controller, you can provide a constructor:

```go
func newProductController(ctx context.Context) webservice.RestController {
	return &productController{
        productService:   newProductService(ctx),
        productConverter: productConverter{},
	}
}
```

In this case, we expect the product service to be injectable, so we use its constructor function
to create an instance of the dependency. This simplifies unit testing by allowing us to inject
a mock for the service.

In contrast, we do not expect to use a mock converter, so it is instantiated directly.

## Connecting the Controller to the Application Lifecycle

In order to instantiate your controller during application startup, you can register a simple
`init` function:

```go
func init() {
	app.OnEvent(app.EventCommand, app.CommandRoot, func(ctx context.Context) error {
		app.OnEvent(app.EventStart, app.PhaseBefore, func(ctx context.Context) error {
			controller := newProductController(ctx)
			return webservice.
				WebServerFromContext(ctx).
				RegisterRestController(pathRoot, controller)
		})
		return nil
	})
}
```

This will register your controller during normal microservice startup.  Since it
is only registering for `CommandRoot`, it will not be created during `migrate`, 
`populate` or other custom command execution.

To ensure your module is included in the built microservice, include the module from your `main.go`:

```go
import _ "cto-github.cisco.com/NFV-BU/productservice/internal/products"
```
