package infoprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"github.com/emicklei/go-restful"
	"time"
)


const (
	configKeyInfo = "info"
)

type InfoProvider struct{}

func (h InfoProvider) infoReport(ctx context.Context) (interface{}, error) {
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
		return nil, webservice.NewStatusError(err, 500)
	}

	i.App.Version = i.Build.Version
	buildTime, err := time.Parse("2006-01-02T15:04:05.999999999Z", i.Build.BuildDateTime)
	if err == nil {
		i.Build.Time = newEpochSeconds(buildTime)
	}

	return i, nil
}

func (h InfoProvider) Actuate(infoService *restful.WebService) error {
	infoService.Consumes(restful.MIME_JSON, restful.MIME_XML)
	infoService.Produces(restful.MIME_JSON, restful.MIME_XML)

	infoService.Path(infoService.RootPath() + "/admin/info")

	// Unsecured routes for info
	infoService.Route(infoService.GET("").
		To(webservice.RawContextController(h.infoReport)).
		Doc("Get System info").
		Do(webservice.Returns200))

	return nil
}

func RegisterInfoProvider(ctx context.Context) error {
	server := webservice.WebServerFromContext(ctx)
	if server != nil {
		server.RegisterActuator(new(InfoProvider))
	}
	return nil
}