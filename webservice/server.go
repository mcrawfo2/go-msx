package webservice

import (
	"context"
	"crypto/tls"
	"github.com/emicklei/go-restful"
	swagger "github.com/emicklei/go-restful-swagger12"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type WebServer struct {
	ctx          context.Context
	cfg          *WebServerConfig
	container    *restful.Container
	containerMtx sync.Mutex
	router       *restful.CurlyRouter
	services     []*restful.WebService
	server       *http.Server
}

func NewWebServer(cfg *WebServerConfig, ctx context.Context) *WebServer {
	return &WebServer{
		ctx:    ctx,
		cfg:    cfg,
		router: &restful.CurlyRouter{},
	}
}

func (s *WebServer) NewService() *restful.WebService {
	s.resetContainer()
	webService := new(restful.WebService)
	webService.Path(s.cfg.ContextPath)
	s.services = append(s.services, webService)
	return webService
}

func (s *WebServer) resetContainer() {
	s.containerMtx.Lock()
	defer s.containerMtx.Unlock()

	s.container = nil
}

func (s *WebServer) Handler() http.Handler {
	s.containerMtx.Lock()
	defer s.containerMtx.Unlock()

	if s.container != nil {
		return s.container
	}

	s.container = restful.NewContainer()
	s.container.Router(s.router)
	s.container.Filter(s.contextInjectorFilter(s.container, s.router))

	for _, service := range s.services {
		s.container.Add(service)
	}

	if s.cfg.Swagger.Enabled {
		swaggerConfig := swagger.Config{
			WebServices:     s.services,
			WebServicesUrl:  s.cfg.Swagger.WebServicesUrl,
			ApiPath:         s.cfg.Swagger.ApiPath,
			SwaggerPath:     s.cfg.Swagger.SwaggerPath,
			SwaggerFilePath: s.cfg.Swagger.Path,
		}

		swagger.RegisterSwaggerService(swaggerConfig, s.container)

		s.container.Add(newSwaggerService(s.cfg.ContextPath))
	}

	s.container.Add(newAdminService(s.cfg.ContextPath))

	return s.container
}

func (s *WebServer) Serve(ctx context.Context) error {
	s.ctx = ctx

	s.server = &http.Server{
		Addr:    ":" + strconv.Itoa(s.cfg.Port),
		Handler: s.Handler(),
	}

	// Start the server
	go func() {
		var err error
		if s.cfg.Tls {
			logger.Infof("Serving HTTPS on %s", s.server.Addr)

			tlsConfig := &tls.Config{
				ClientAuth: tls.VerifyClientCertIfGiven,
			}
			tlsConfig.BuildNameToCertificate()
			err = s.server.ListenAndServeTLS(s.cfg.CertFile, s.cfg.KeyFile)
		} else {
			logger.Infof("Serving HTTP on %s", s.server.Addr)

			err = s.server.ListenAndServe()
		}

		if err != nil {
			logger.WithError(err).Error("Web Server exited")
		}
	}()

	return nil
}

func (s *WebServer) StopServing(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func (s *WebServer) contextInjectorFilter(container *restful.Container, router restful.RouteSelector) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		// Inject the container and router
		ctx := ContextWithContainer(s.ctx, container)
		ctx = ContextWithRouter(ctx, router)

		// Inject the webservice and route
		service, route, _ := router.SelectRoute(container.RegisteredWebServices(), req.Request)
		ctx = ContextWithService(ctx, service)
		ctx = ContextWithRoute(ctx, route)
		if route != nil {
			ctx = ContextWithRouteOperation(ctx, route.Operation)
		} else {
			ctx = ContextWithRouteOperation(ctx, "unknown")
		}

		// Inject the user context filter config
		ctx = ContextWithUserContextFilterConfig(ctx, &s.cfg.JWT)

		// Inject the request variables
		ctx = ContextWithHttpRequest(ctx, req.Request)
		ctx = ContextWithPathParameters(ctx, req.PathParameters())

		// Add the context into the request
		req.Request = req.Request.WithContext(ctx)

		chain.ProcessFilter(req, resp)
	}
}
