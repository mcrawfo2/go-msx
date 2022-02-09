package asyncapiprovider

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"cto-github.cisco.com/NFV-BU/go-msx/webservice"
	"github.com/emicklei/go-restful"
	"github.com/pkg/errors"
	"path"
)

var (
	logger = log.NewLogger("msx.webservice.asyncapiprovider")
)

type SpecProvider interface {
	Spec() ([]byte, error)
}

type StudioProvider struct {
	ctx          context.Context
	cfg          *DocumentationConfig
	specProvider SpecProvider
}

func (p StudioProvider) GetYamlSpec(request *restful.Request, response *restful.Response) {
	specBytes, err := p.specProvider.Spec()
	if err != nil {
		webservice.WriteError(request, response, 500, errors.New("Failed to load AsyncApi specification"))
		return
	}

	response.Header().Set("Content-Type", webservice.MIME_YAML_CHARSET)
	response.WriteHeader(200)
	_, _ = response.Write(specBytes)
}

func (p StudioProvider) GetUi(req *restful.Request) (body interface{}, err error) {
	return struct {
		ReadOnly bool   `json:"readOnly"`
		Spec     string `json:"spec"`
	}{
		ReadOnly: true,
		Spec:     path.Join(p.cfg.Server.ContextPath, p.cfg.Resources.Path, p.cfg.Resources.YamlSpecPath),
	}, nil
}

func (p StudioProvider) Actuate(container *restful.Container, ws *restful.WebService) error {
	// Configure the resources webservice
	contextPath := ws.RootPath()
	ws.Path(contextPath + p.cfg.Resources.Path)

	ws.Route(ws.GET(p.cfg.Resources.YamlSpecPath).
		Operation("studio.resources.spec").
		To(p.GetYamlSpec).
		Produces(webservice.MIME_YAML).
		Do(webservice.Returns(200, 401)))

	ws.Route(ws.GET("/ui").
		Operation("studio.resources.configuration.ui").
		To(webservice.RawController(p.GetUi)).
		Produces(webservice.MIME_JSON).
		Do(webservice.Returns(200, 401)))

	logger.Infof("Serving AsyncApi on http://%s:%d%s%s/",
		p.cfg.Server.Host,
		p.cfg.Server.Port,
		contextPath,
		p.cfg.Ui.Endpoint)

	return nil
}

func NewStudioProvider(ctx context.Context) (*StudioProvider, error) {
	cfg, err := NewDocumentationConfig(ctx)
	if err != nil {
		return nil, err
	}

	if !cfg.Enabled {
		return nil, ErrDisabled
	}

	var specProvider SpecProvider
	switch cfg.Source {
	case "static-file":
		specProvider = StaticFileSpecProvider{cfg: cfg.Resources}

	case "registry":
		specProvider, err = NewRegistrySpecProvider(ctx)

	default:
		return nil, errors.Errorf("Unknown asyncapi source %q", cfg.Source)
	}

	if err != nil {
		return nil, err
	}

	return &StudioProvider{
		ctx:          ctx,
		cfg:          cfg,
		specProvider: specProvider,
	}, nil
}
