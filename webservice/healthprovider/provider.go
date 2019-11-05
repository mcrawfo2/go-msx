package healthprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/health"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"github.com/emicklei/go-restful"
)

type HealthProvider struct{}


func (h HealthProvider) healthReport(ctx context.Context) (interface{}, error) {
	userContext := security.UserContextFromContext(ctx)
	if userContext != nil {
		return health.GenerateReport(ctx), nil
	} else {
		return health.GenerateSummary(ctx), nil
	}
}

func (h HealthProvider) Actuate(healthService *restful.WebService) error {
	healthService.Consumes(restful.MIME_JSON, restful.MIME_XML)
	healthService.Produces(restful.MIME_JSON, restful.MIME_XML)

	healthService.Path(healthService.RootPath() + "/admin/health")

	// Unsecured routes for health
	healthService.Route(healthService.GET("").
		To(webservice.RawContextController(h.healthReport)).
		Doc("Get System health").
		Do(webservice.Returns200))

	return nil
}

func RegisterHealthProvider(ctx context.Context) error {
	server := webservice.WebServerFromContext(ctx)
	if server != nil {
		server.RegisterActuator(new(HealthProvider))
	}
	return nil
}
