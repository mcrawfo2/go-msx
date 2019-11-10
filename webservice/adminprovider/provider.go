package adminprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"fmt"
	"github.com/emicklei/go-restful"
	"net/http"
	"strings"
)

type Link struct {
	Href      string `json:"href"`
	Templated bool   `json:"templated"`
}

var links = make(map[string]Link)

type Report struct {
	Links map[string]Link `json:"_links"`
}

type AdminProvider struct {}

func (h AdminProvider) Actuate(infoService *restful.WebService) error {
	infoService.Consumes(restful.MIME_JSON, restful.MIME_XML)
	infoService.Produces(restful.MIME_JSON, restful.MIME_XML)

	infoService.Path(infoService.RootPath() + "/admin")
	infoService.Filter(webservice.RequireAuthenticatedFilter)

	// Unsecured routes for admin
	infoService.Route(infoService.GET("").
		To(webservice.RawController(h.adminReport)).
		Doc("Get System info").
		Do(webservice.Returns200))

	infoService.Route(infoService.GET("/alive").
		To(webservice.RawController(h.emptyReport)).
		Do(webservice.Returns200))

	return nil
}

func (h AdminProvider) adminReport(req *restful.Request) (body interface{}, err error) {
	baseUrl := fmt.Sprintf("http://%s%s", req.Request.Host, req.Request.URL.String())
	var reportLinks = map[string]Link{
		"self": {baseUrl, false},
		"alive": { baseUrl + "/alive", false},
	}
	for k, v := range links {
		reportLinks[k] = Link{baseUrl + "/" + v.Href, v.Templated}
	}
	return Report{reportLinks}, nil
}

func (h AdminProvider) emptyReport(req *restful.Request) (body interface{}, err error) {
	return struct{}{}, nil
}

func (h AdminProvider) optionsFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	var container = webservice.ContainerFromContext(req.Request.Context())
	var router = webservice.RouterFromContext(req.Request.Context())
	var newHttpRequest = *req.Request
	var allowedMethods = make(map[string]struct{})
	for _, method := range []string{"PATCH", "POST", "GET", "PUT", "DELETE", "HEAD"} {
		newHttpRequest.Method = method
		_, route, err := router.SelectRoute(container.RegisteredWebServices(), &newHttpRequest)
		if err != nil || route == nil {
			continue
		}
		allowedMethods[route.Method] = struct{}{}
	}

	if len(allowedMethods) == 0 {
		http.NotFound(resp, req.Request)
		return
	}

	allowedMethods["OPTIONS"] = struct{}{}
	allowedMethods["HEAD"] = struct{}{}
	var allowMethods []string
	for k := range allowedMethods {
		allowMethods = append(allowMethods, k)
	}
	allowMethodsHeaderValue := strings.Join(allowMethods, ",")

	resp.AddHeader("Allow", allowMethodsHeaderValue)
}

func RegisterProvider(ctx context.Context) error {
	server := webservice.WebServerFromContext(ctx)
	if server != nil {
		server.RegisterActuator(new(AdminProvider))
	}
	return nil
}

func RegisterLink(name, href string, templated bool) {
	links[name] = Link{
		Href:      href,
		Templated: templated,
	}
}
