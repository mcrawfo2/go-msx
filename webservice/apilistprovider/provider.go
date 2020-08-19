package apilistprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/adminprovider"
	"github.com/emicklei/go-restful"
	"strings"
)

const (
	endpointName = "apilist"
)

type Provider struct {
	ctx    context.Context
	report []permission
}

type permission struct {
	Name      string     `json:"name"`
	Resources []resource `json:"resources"`
}

type resource struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Endpoint    string `json:"endpoint"`
	Method      string `json:"method"`
	Consumes    string `json:"consumes,omitempty"`
	Produces    string `json:"produces,omitempty"`
}

func (h *Provider) EndpointName() string {
	return endpointName
}

func (h *Provider) Report(req *restful.Request) (body interface{}, err error) {
	return h.report, nil
}

func (h *Provider) generateEndpoints(container *restful.Container) {
	permissionMap := make(map[string][]resource)

	for _, ws := range container.RegisteredWebServices() {
		for _, route := range ws.Routes() {
			permissions, _ := route.Metadata[webservice.MetadataPermissions].([]string)
			if len(permissions) == 0 {
				permissions = []string{"NO_PERMISSION"}
			}

			for _, permission := range permissions {
				resources, _ := permissionMap[permission]

				rsc := resource{
					Name:        route.Operation,
					Description: route.Doc,
					Endpoint:    route.Path,
					Method:      route.Method,
					Consumes:    strings.Join(route.Consumes, ","),
					Produces:    strings.Join(route.Produces, ","),
				}

				resources = append(resources, rsc)
				permissionMap[permission] = resources
			}
		}
	}

	for name, resources := range permissionMap {
		h.report = append(h.report, permission{
			Name:      name,
			Resources: resources,
		})
	}
}

func (h *Provider) Actuate(container *restful.Container, webService *restful.WebService) error {
	webService.Consumes(restful.MIME_JSON)
	webService.Produces(restful.MIME_JSON)

	webService.Path(webService.RootPath() + "/admin/" + endpointName)

	webService.Route(webService.GET("").
		Operation("admin.apilist").
		To(adminprovider.RawAdminController(h.Report)).
		Do(webservice.Returns200))

	h.generateEndpoints(container)

	return nil
}

func RegisterProvider(ctx context.Context) error {
	server := webservice.WebServerFromContext(ctx)
	if server != nil {
		server.AddDocumentationProvider(&Provider{
			ctx: ctx,
		})
		adminprovider.RegisterLink(endpointName, endpointName, false)
	}
	return nil
}
