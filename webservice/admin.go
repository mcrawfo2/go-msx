package webservice

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/health"
	"cto-github.cisco.com/NFV-BU/go-msx/security"
	"fmt"
	"github.com/emicklei/go-restful"
	"time"
)

const (
	configKeyInfo = "info"
)

type epochSeconds float64

func (e epochSeconds) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%.9f", float64(e))), nil
}

func newEpochSeconds(when time.Time) epochSeconds {
	return epochSeconds(float64(when.Unix()) + (float64(when.Nanosecond()) * 1e-9))
}

func healthReport(ctx context.Context) (interface{}, error) {
	userContext := security.UserContextFromContext(ctx)
	if userContext != nil {
		return health.GenerateReport(ctx), nil
	} else {
		return health.GenerateSummary(ctx), nil
	}
}

func infoReport(ctx context.Context) (interface{}, error) {
	type Info struct {
		App struct {
			Name        string `json:"name"`
			Version     string `json:"version"`
			Description string `json:"description"`
			Attributes  struct {
				DisplayName string `json:"displayName"`
				Parent      string `json:"parent"`
				Type        string `json:"type"`
			} `json:"attributes"`
		} `json:"app"`
		Build struct {
			Version       string  `json:"version"`
			BuildNumber   string  `json:"number" config:"buildNumber"`
			BuildDateTime string  `json:"-"`
			Artifact      string  `json:"artifact"`
			Name          string  `json:"name"`
			Time          epochSeconds `json:"time" config:"default=0"`
			Group         string  `json:"group"`
		} `json:"build"`
	}

	i := Info{}
	if err := config.MustFromContext(ctx).Populate(&i, configKeyInfo); err != nil {
		return nil, NewStatusError(err, 500)
	}

	i.App.Version = i.Build.Version
	buildTime, err := time.Parse("2006-01-02T15:04:05.999999999Z", i.Build.BuildDateTime)
	if err == nil {
		i.Build.Time = newEpochSeconds(buildTime)
	}

	return i, nil
}

func newAdminService(contextPath string) *restful.WebService {
	var adminService = new(restful.WebService)
	adminService.Path(contextPath + "/admin")
	adminService.Consumes(restful.MIME_JSON, restful.MIME_XML)
	adminService.Produces(restful.MIME_JSON, restful.MIME_XML)

	// Unsecured Routes for health and info
	adminService.Route(adminService.GET("/health").
		To(RawContextController(healthReport)).
		Doc("Get System health").
		Do(Returns200))

	adminService.Route(adminService.GET("/info").
		To(RawContextController(infoReport)).
		Doc("Get System info").
		Do(Returns200))

	return adminService
}
