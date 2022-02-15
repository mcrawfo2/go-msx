package webservice

import (
	"context"
	"crypto/tls"
	"cto-github.cisco.com/NFV-BU/go-msx/background"
	"cto-github.cisco.com/NFV-BU/go-msx/certificate/fileprovider"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/contexttest"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/emicklei/go-restful"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
	"time"
)

func TestNewWebServer(t *testing.T) {
	new(WebServerTest).
		WithWebServerPredicate(WebServerPredicate{
			Description: "NotNil",
			Matches: func(server *WebServer) bool {
				return server != nil
			},
		}).
		Test(t)
}

func TestWebServer_AddDocumentationProvider(t *testing.T) {
	documentationProvider := new(MockDocumentationProvider)

	new(WebServerTest).
		WithWebServerCustomizer(func(s *WebServer) {
			s.AddDocumentationProvider(documentationProvider)
		}).
		WithWebServerPredicate(WebServerHasDocumentation(documentationProvider)).
		Test(t)
}

func TestWebServer_ContextPath(t *testing.T) {
	const contextPath = "/bob"
	new(WebServerTest).
		WithStaticConfig(map[string]string{
			"server.context-path": contextPath,
		}).
		WithWebServerPredicate(WebServerHasContextPath(contextPath)).
		Test(t)
}

func TestWebServer_Handler(t *testing.T) {
	new(WebServerTest).
		WithWebServerCustomizer(func(s *WebServer) {
			_, _ = s.NewService("/bob")
		}).
		WithWebServerVerifier(func(t *testing.T, s *WebServer) {
			var container http.Handler = s.generateContainer()
			handler := s.Handler()
			assert.Equal(t, container, handler)
		}).
		Test(t)

}

func TestWebServer_NewService(t *testing.T) {
	const rootPath = "/bob"
	const contextPath = "/app"

	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "MultipleRegistrations",
			test: new(WebServerTest).
				WithWebServerVerifier(func(t *testing.T, s *WebServer) {
					svc, err := s.NewService(rootPath)
					assert.NoError(t, err)
					assert.Equal(t, svc.RootPath(), contextPath+rootPath)

					otherSvc, err := s.NewService(rootPath)
					assert.NoError(t, err)
					assert.Equal(t, svc, otherSvc)
					assert.Equal(t, otherSvc.RootPath(), contextPath+rootPath)
				}),
		},
		{
			name: "NoSlash",
			test: new(WebServerTest).
				WithWebServerVerifier(func(t *testing.T, s *WebServer) {
					svc, err := s.NewService("bob")
					assert.NoError(t, err)
					assert.Equal(t, svc.RootPath(), contextPath+rootPath)
				}),
		},
		{
			name: "NoPathError",
			test: new(WebServerTest).
				WithWebServerVerifier(func(t *testing.T, s *WebServer) {
					svc, err := s.NewService("")
					assert.Error(t, err)
					assert.Nil(t, svc)
				}),
		},
		{
			name: "NoServerError",
			test: testhelpers.TestFunc(func(t *testing.T) {
				var s *WebServer = nil
				svc, err := s.NewService("")
				assert.Error(t, err)
				assert.Nil(t, svc)
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}

func TestWebServer_RegisterActuator(t *testing.T) {
	provider := new(MockServiceProvider)
	provider2 := new(MockServiceProvider)

	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "Single",
			test: new(WebServerTest).
				WithWebServerCustomizer(func(s *WebServer) {
					s.RegisterActuator(provider)
				}).
				WithWebServerPredicate(WebServerHasActuator(provider)),
		},
		{
			name: "Double",
			test: new(WebServerTest).
				WithWebServerCustomizer(func(s *WebServer) {
					s.RegisterActuator(provider)
					s.RegisterActuator(provider2)
				}).
				WithWebServerPredicate(WebServerHasActuator(provider)).
				WithWebServerPredicate(WebServerHasActuator(provider2)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}

func TestWebServer_RegisterAlias(t *testing.T) {
	const (
		path1 = "/swagger"
		path2 = "/swapper"
		file1 = "/swagger-ui.html"
		file2 = "/swapper-ui.html"
	)

	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "Single",
			test: new(WebServerTest).
				WithWebServerCustomizer(func(s *WebServer) {
					s.RegisterAlias(path1, file1)
				}).
				WithWebServerPredicate(WebServerHasAlias(path1, file1)),
		},
		{
			name: "Double",
			test: new(WebServerTest).
				WithWebServerCustomizer(func(s *WebServer) {
					s.RegisterAlias(path1, file1)
					s.RegisterAlias(path2, file2)
				}).
				WithWebServerPredicate(WebServerHasAlias(path1, file1)).
				WithWebServerPredicate(WebServerHasAlias(path2, file2)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}

func TestWebServer_RegisterInjector(t *testing.T) {
	injector := types.ContextInjector(func(ctx context.Context) context.Context { return ctx })

	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "Single",
			test: new(WebServerTest).
				WithWebServerCustomizer(func(s *WebServer) {
					s.RegisterInjector(injector)
				}).
				WithWebServerPredicate(WebServerHasInjector()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}

func TestWebServer_RegisterRestController(t *testing.T) {
	const servicePath = "/bob"

	controller := new(MockRestController)
	controller.On("Routes", mock.AnythingOfType("*restful.WebService"))

	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "Single",
			test: new(WebServerTest).
				WithWebServerVerifier(func(t *testing.T, s *WebServer) {
					err := s.RegisterRestController(servicePath, controller)
					assert.NoError(t, err)
				}).
				WithWebServerPredicate(WebServerHasService(servicePath)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}

func TestWebServer_Serve(t *testing.T) {
	new(WebServerTest).
		WithStaticConfig(map[string]string{
			"server.host": "127.0.0.10",
			"server.port": "24203",
		}).
		WithWebServerVerifier(func(t *testing.T, s *WebServer) {
			errorReporter := new(background.MockErrorReporter)
			errorChan := make(chan struct{}, 1)
			errorReporter.On("C").Return(errorChan)
			errorReporter.On("NonFatal", mock.Anything)

			rootCtx := background.ContextWithErrorReporter(context.Background(), errorReporter)

			ctx, cancel := context.WithCancel(rootCtx)
			err := s.Serve(ctx)
			assert.NoError(t, err)
			cancel()

			ctx2, cancel2 := context.WithDeadline(rootCtx, time.Now().Add(100*time.Millisecond))
			defer cancel2()
			err = s.StopServing(ctx2)
			assert.NoError(t, err)
		}).
		Test(t)

}

func TestWebServer_SetAuthenticationProvider(t *testing.T) {
	provider := new(MockAuthenticationProvider)

	new(WebServerTest).
		WithWebServerCustomizer(func(s *WebServer) {
			s.SetAuthenticationProvider(provider)
		}).
		WithWebServerVerifier(func(t *testing.T, s *WebServer) {
			assert.Equal(t, provider, s.security)
		}).
		Test(t)
}

func TestWebServer_SetHandler(t *testing.T) {
	const handlerPath = "/bob"
	handler := new(MockHttpHandler)

	new(WebServerTest).
		WithWebServerCustomizer(func(s *WebServer) {
			s.SetHandler(handlerPath, handler)
		}).
		WithWebServerVerifier(func(t *testing.T, s *WebServer) {
			h, ok := s.handlers[s.ContextPath()+handlerPath]
			assert.True(t, ok)
			assert.Equal(t, handler, h)
		}).
		Test(t)
}

func TestWebServer_StopServing(t *testing.T) {
	t.Skipped()
}

func TestWebServer_Url(t *testing.T) {
	tests := []struct {
		name string
		test testhelpers.Testable
	}{
		{
			name: "Http",
			test: new(WebServerTest).
				WithStaticConfig(map[string]string{
					"server.tls.enabled":  "false",
					"server.host":         "10.10.10.10",
					"server.port":         "9090",
					"server.context-path": "/bob",
				}).
				WithWebServerVerifier(func(t *testing.T, s *WebServer) {
					assert.Equal(t, "http://10.10.10.10:9090/bob", s.Url())
				}),
		},
		{
			name: "Https",
			test: new(WebServerTest).
				WithStaticConfig(map[string]string{
					"server.tls.enabled":  "true",
					"server.host":         "10.20.30.40",
					"server.port":         "7070",
					"server.context-path": "/bob",
				}).
				WithWebServerVerifier(func(t *testing.T, s *WebServer) {
					assert.Equal(t, "https://10.20.30.40:7070/bob", s.Url())
				}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test.Test)
	}
}

func TestWebServer_activateActuator(t *testing.T) {
	t.Skipped()
}

func TestWebServer_activateDocumentation(t *testing.T) {
	t.Skipped()
}

func TestWebServer_activateStatic(t *testing.T) {
	t.Skipped()
}

func TestWebServer_generateContainer(t *testing.T) {
	new(WebServerTest).
		WithWebServerCustomizer(func(s *WebServer) {
			_, _ = s.NewService("/api/v1/bob")
		}).
		WithWebServerCustomizer(func(s *WebServer) {
			actuator := new(MockServiceProvider)
			actuator.On("EndpointName").Return("mock")
			actuator.
				On("Actuate", mock.Anything).
				Return(nil).
				Run(func(args mock.Arguments) {
					ws := args.Get(0).(*restful.WebService)
					ws.Path(s.ContextPath() + "/admin/mock")
				})
			s.RegisterActuator(actuator)
		}).
		WithWebServerCustomizer(func(s *WebServer) {
			documentation := new(MockDocumentationProvider)
			documentation.
				On("Actuate", mock.Anything, mock.Anything).
				Return(nil).
				Run(func(args mock.Arguments) {
					ws := args.Get(1).(*restful.WebService)
					ws.Path(s.ContextPath() + "/docs")
				})
			s.AddDocumentationProvider(documentation)
		}).
		WithWebServerCustomizer(func(s *WebServer) {
			s.RegisterAlias("/bob", "/alice.html")
		}).
		WithWebServerVerifier(func(t *testing.T, s *WebServer) {
			container := s.generateContainer()
			assert.NotNil(t, container)
		}).
		Test(t)
}

func TestWebServer_generateContainer_Empty(t *testing.T) {
	new(WebServerTest).
		WithWebServerVerifier(func(t *testing.T, s *WebServer) {
			container := s.generateContainer()
			assert.NotNil(t, container)
		}).
		Test(t)
}

func TestWebServer_getTLSListener(t *testing.T) {
	t.Skipped()
}

func TestWebServer_resetContainer(t *testing.T) {
	new(WebServerTest).
		WithWebServerCustomizer(func(s *WebServer) {
			_, _ = s.NewService("/alice")
			s.generateContainer()
		}).
		WithWebServerVerifier(func(t *testing.T, s *WebServer) {
			assert.NotNil(t, s.container)
			s.resetContainer()
			assert.Nil(t, s.container)
		}).
		Test(t)
}

func Test_buildTlsConfig(t *testing.T) {
	err := fileprovider.RegisterFactory(context.Background())
	assert.NoError(t, err)

	new(WebServerTest).
		WithStaticConfig(map[string]string{
			"server.tls.enabled":                     "true",
			"server.tls.ca-file":                     "testdata/server.crt",
			"server.tls.cert-file":                   "testdata/server.crt",
			"server.tls.key-file":                    "testdata/server.key",
			"server.tls.min-version":                 "tls13",
			"server.tls.cipher-suites":               "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
			"certificate.source.server.provider":     "file",
			"certificate.source.server.ca-cert-file": "${server.tls.ca-file}",
			"certificate.source.server.cert-file":    "${server.tls.cert-file}",
			"certificate.source.server.key-file":     "${server.tls.key-file}",
		}).
		WithWebServerVerifier(func(t *testing.T, s *WebServer) {
			tlsConfig, err := s.cfg.Tls.TlsConfig(s.ctx)
			assert.NoError(t, err)
			assert.NotNil(t, tlsConfig)
			assert.Equal(t, tlsConfig.MinVersion, uint16(tls.VersionTLS13))
			assert.Equal(t, tls.VerifyClientCertIfGiven, tlsConfig.ClientAuth)
			assert.NotEmpty(t, tlsConfig.ClientCAs)
			assert.Equal(t, tlsConfig.CipherSuites, []uint16{
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			})
			assert.NotNil(t, tlsConfig.GetCertificate)
		}).
		Test(t)
}

func Test_requestContextInjectorFilter(t *testing.T) {
	const key = "testContextInjectorKey"
	const value = "testContextInjectorValue"

	ctx := configtest.ContextWithNewInMemoryConfig(context.Background(), map[string]string{})

	ws := new(restful.WebService)
	ws.Route(ws.
		GET("/bob").
		Operation("bob").
		To(func(_ *restful.Request, resp *restful.Response) {
			resp.WriteHeader(200)
		}))

	router := new(restful.CurlyRouter)

	container := new(restful.Container)
	container.ServeMux = http.DefaultServeMux
	container.Router(router)
	container.Add(ws)

	authentication := new(MockAuthenticationProvider)

	injectors := new(types.ContextInjectors)
	injectors.Register(func(ctx context.Context) context.Context {
		return context.WithValue(ctx, key, value)
	})

	new(RouteBuilderTest).
		WithRequestPath("/bob").
		WithRequestMethod(http.MethodGet).
		WithContext(ctx).
		WithRouteFilter(requestContextInjectorFilter(
			ctx,
			container,
			new(restful.CurlyRouter),
			authentication,
			injectors.Slice())).
		WithContextPredicate(contexttest.ContextGetterHasAnyValue(ContainerFromContext)).
		WithContextPredicate(contexttest.ContextGetterHasAnyValue(RouterFromContext)).
		WithContextPredicate(contexttest.ContextGetterHasAnyValue(AuthenticationProviderFromContext)).
		WithContextPredicate(contexttest.ContextGetterHasAnyValue(ServiceFromContext)).
		WithContextPredicate(contexttest.ContextGetterHasAnyValue(RouteFromContext)).
		WithContextPredicate(contexttest.ContextGetterHasAnyValue(RouteOperationFromContext)).
		WithContextPredicate(contexttest.ContextGetterHasAnyValue(func(ctx context.Context) interface{} {
			if ctx.Value(key) == nil {
				return false
			}
			value, ok := ctx.Value(key).(string)
			return ok && value != ""
		})).
		Test(t)

}
