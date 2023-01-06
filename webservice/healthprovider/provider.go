// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package healthprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/health"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/adminprovider"
	"encoding/json"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
	"net/http"
)

const endpointName = "health"

const showDetailsAlways = "ALWAYS"
const showDetailsNever = "NEVER"
const showDetailsWhenAuthorized = "WHEN_AUTHORIZED"

type HealthProvider struct {
	cfg *webservice.ManagementEndpointConfig
}

func (h HealthProvider) healthReport(req *restful.Request) (interface{}, error) {
	ctx := req.Request.Context()
	userContext := security.UserContextFromContext(ctx)

	isReport := (h.cfg.ShowDetails == showDetailsAlways) ||
		(h.cfg.ShowDetails == showDetailsWhenAuthorized && (userContext != nil && userContext.Token != ""))

	if isReport {
		return health.GenerateReport(ctx), nil
	} else {
		return health.GenerateSummary(ctx), nil
	}
}

func (h HealthProvider) healthComponentReport(req *restful.Request) (interface{}, error) {
	component := req.PathParameter("component")
	report := health.GenerateReport(req.Request.Context())
	if details, ok := report.Details[component]; !ok {
		return nil, webservice.NewStatusError(
			errors.New("Component not found"),
			http.StatusNotFound)
	} else {
		return details, nil
	}
}

func (h HealthProvider) EndpointName() string {
	return endpointName
}

func HealthAdminController(fn webservice.ControllerFunction) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
		body, err := fn(req)
		if err != nil {
			webservice.RawResponse(req, resp, nil, err)
			return
		}

		resp.Header().Set("Expires", "0")
		resp.Header().Set("Pragma", "no-cache")
		resp.Header().Set("Content-Type", "application/vnd.spring-boot.actuator.v2+json")
		resp.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")

		if body != nil {
			var bodyBytes []byte
			if report, ok := body.(*health.Report); ok {
				if report.Status != health.StatusUp {
					resp.WriteHeader(503)
				} else {
					resp.WriteHeader(200)
				}
			} else {
				webservice.RawResponse(req, resp, body, errors.New("unable to assert to health report"))
				return
			}
			bodyBytes, _ = json.Marshal(body)
			_, _ = resp.Write(bodyBytes)
		} else {
			resp.WriteHeader(204)
		}
	}
}

func (h HealthProvider) Actuate(healthService *restful.WebService) error {
	healthService.Consumes(restful.MIME_JSON, restful.MIME_XML)
	healthService.Produces(restful.MIME_JSON, restful.MIME_XML)

	healthService.Path(healthService.RootPath() + "/admin/health")

	healthService.Route(healthService.GET("").
		Operation("admin.health").
		To(HealthAdminController(h.healthReport)).
		Doc("Get System health").
		Do(webservice.Returns200))

	healthService.Route(healthService.GET("/{component}").
		Operation("admin.health-component").
		To(adminprovider.RawAdminController(h.healthComponentReport)).
		Param(healthService.PathParameter("component", "Name of component to probe")).
		Doc("Get component health").
		Do(webservice.Returns(200, 404)))

	return nil
}

func RegisterProvider(ctx context.Context) error {
	server := webservice.WebServerFromContext(ctx)
	if server != nil {
		cfg, err := webservice.NewManagementEndpointConfig(ctx, "health")
		if err != nil {
			return err
		}

		provider := new(HealthProvider)
		provider.cfg = cfg

		server.RegisterActuator(provider)
		adminprovider.RegisterLink("health", "health", false)
		adminprovider.RegisterLink("health-component", "health/{component}", true)
	}
	return nil
}
