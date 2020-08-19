package aliveprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/adminprovider"
	"github.com/emicklei/go-restful"
)

const (
	endpointName = "alive"
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
		Operation("admin.alive").
		To(adminprovider.RawAdminController(h.emptyReport)).
		Doc("Liveness check").
		Do(webservice.Returns200))

	return nil
}

func (h Provider) emptyReport(req *restful.Request) (body interface{}, err error) {
	return struct{}{}, nil
}

func RegisterProvider(ctx context.Context) error {
	server := webservice.WebServerFromContext(ctx)
	if server != nil {
		server.RegisterActuator(new(Provider))
		adminprovider.RegisterLink(endpointName, endpointName, false)
	}
	return nil
}
