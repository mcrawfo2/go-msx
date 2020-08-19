package prometheusprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"github.com/emicklei/go-restful"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Provider struct{}

const endpointName = "prometheus"

func (h Provider) EndpointName() string {
	return endpointName
}

func (h Provider) Actuate(webService *restful.WebService) error {
	webService.Consumes(restful.MIME_JSON)
	webService.Produces(restful.MIME_JSON)

	webService.Path(webService.RootPath() + "/admin/" + endpointName)

	// Unsecured routes for admin
	webService.Route(webService.GET("").
		Operation("prometheus.metrics").
		To(webservice.HttpHandlerController(promhttp.Handler().ServeHTTP)).
		Doc("Retrieve metric data for prometheus").
		Do(webservice.Returns200))

	return nil
}

func RegisterProvider(ctx context.Context) error {
	server := webservice.WebServerFromContext(ctx)
	if server != nil {
		server.RegisterActuator(new(Provider))
	}
	return nil
}
