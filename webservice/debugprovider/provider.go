package debugprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/adminprovider"
	"github.com/emicklei/go-restful"
	"net/http/pprof"
)

const (
	endpointName = "debug"
)

type Provider struct{}

func (h Provider) EndpointName() string {
	return endpointName
}

func (h Provider) Actuate(webService *restful.WebService) error {
	webService.Consumes(restful.MIME_JSON)
	webService.Produces(restful.MIME_JSON)

	webService.Path(webService.RootPath() + "/admin/" + endpointName)

	webService.Route(webService.GET("").
		Operation("admin.debug").
		To(adminprovider.HttpHandlerController(pprof.Index)).
		Do(webservice.Returns200))

	webService.Route(webService.GET("{subPath:*}").
		Operation("admin.debug.index").
		To(func(request *restful.Request, response *restful.Response) {
			subPath := request.PathParameter("subPath")
			handler := pprof.Handler(subPath)
			handler.ServeHTTP(response, request.Request)
		}).
		Do(webservice.Returns200))

	return nil
}

func RegisterProvider(ctx context.Context) error {
	server := webservice.WebServerFromContext(ctx)
	if server != nil {
		server.RegisterActuator(new(Provider))
		adminprovider.RegisterLink(endpointName, endpointName, false)
	}
	return nil
}
