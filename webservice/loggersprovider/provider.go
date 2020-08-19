package loggersprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice/adminprovider"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
)

const (
	endpointName = "loggers"
)

type Logger struct {
	ConfiguredLevel string `json:"configuredLevel"`
	EffectiveLevel  string `json:"effectiveLevel"`
}

type Report struct {
	Levels  []string          `json:"levels"`
	Loggers map[string]Logger `json:"loggers"`
}

type Provider struct{}

func (h Provider) EndpointName() string {
	return endpointName
}

func (h Provider) Report(req *restful.Request) (interface{}, error) {
	var loggers = make(map[string]Logger)
	for name, level := range log.GetLoggerLevels() {
		levelName := log.LoggerLevel(level).Name()
		loggers[name] = Logger{
			ConfiguredLevel: levelName,
			EffectiveLevel:  levelName,
		}
	}

	return Report{
		Levels:  log.AllLevelNames,
		Loggers: loggers,
	}, nil
}

func (h Provider) Configure(req *restful.Request) (interface{}, error) {
	var logger Logger
	err := req.ReadEntity(&logger)
	if err != nil {
		return nil, err
	}

	loggerName := req.PathParameter("loggerName")
	if loggerName == "" {
		return nil, errors.New("Logger name must not be empty")
	}

	log.SetLoggerLevel(loggerName, log.LevelFromName(logger.ConfiguredLevel))

	return nil, nil
}

func (h Provider) Actuate(webService *restful.WebService) error {
	webService.Consumes(restful.MIME_JSON)
	webService.Produces(restful.MIME_JSON)

	webService.Path(webService.RootPath() + "/admin/" + endpointName)

	// Unsecured routes for info
	webService.Route(webService.GET("").
		Operation("admin.loggers").
		To(adminprovider.RawAdminController(h.Report)).
		Do(webservice.Returns200))

	webService.Route(webService.POST("{loggerName}").
		Operation("admin.loggers.configure").
		To(adminprovider.RawAdminController(h.Configure)).
		Do(webservice.Returns204))

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
