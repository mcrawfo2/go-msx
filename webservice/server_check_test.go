// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package webservice

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/fs"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/logtest"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"testing"
)

type WebServerCheck struct {
	Validators []WebServerPredicate
}

func (r WebServerCheck) Check(server *WebServer) []error {
	var results []error

	for _, predicate := range r.Validators {
		if !predicate.Matches(server) {
			results = append(results, WebServerCheckError{
				Validator: predicate,
			})
		}
	}

	return results

}

type WebServerCustomizer func(s *WebServer)
type WebServerVerifier func(t *testing.T, s *WebServer)

type WebServerCheckError struct {
	Validator WebServerPredicate
}

func (c WebServerCheckError) Error() string {
	return fmt.Sprintf("Failed web server validator: %s", c.Validator.Description)
}

type WebServerPredicate struct {
	Description string
	Matches     func(s *WebServer) bool
}

func WebServerHasActuator(provider ServiceProvider) WebServerPredicate {
	return WebServerPredicate{
		Description: "s.actuators contains provider",
		Matches: func(s *WebServer) bool {
			for _, actuator := range s.actuators {
				if reflect.DeepEqual(actuator, provider) {
					return true
				}
			}
			return false
		},
	}
}

func WebServerHasInjector() WebServerPredicate {
	return WebServerPredicate{
		Description: "s.injectors length > 0",
		Matches: func(s *WebServer) bool {
			return len(*s.injectors) > 0
		},
	}
}

func WebServerHasAlias(path, file string) WebServerPredicate {
	return WebServerPredicate{
		Description: fmt.Sprintf("s.aliases contains StaticAlias[%q => %q]", path, file),
		Matches: func(s *WebServer) bool {
			for _, alias := range s.aliases {
				if alias.Path == path && alias.File == file {
					return true
				}
			}
			return false
		},
	}
}

func WebServerHasService(path string) WebServerPredicate {
	return WebServerPredicate{
		Description: fmt.Sprintf("s.services has %v", path),
		Matches: func(s *WebServer) bool {
			for _, svc := range s.services {
				if svc.RootPath() == s.ContextPath()+path {
					return true
				}
			}
			return false
		},
	}
}

func WebServerHasDocumentation(provider DocumentationProvider) WebServerPredicate {
	return WebServerPredicate{
		Description: fmt.Sprintf("s.documentationProviders contains %p", provider),
		Matches: func(server *WebServer) bool {
			for _, p := range server.documentation {
				if reflect.DeepEqual(p, provider) {
					return true
				}
			}
			return false
		},
	}
}

func WebServerHasContextPath(path string) WebServerPredicate {
	return WebServerPredicate{
		Description: fmt.Sprintf(`server.ContextPath == %q`, path),
		Matches: func(server *WebServer) bool {
			return server.ContextPath() == path
		},
	}
}

type WebServerTest struct {
	StaticConfig map[string]string
	Customizers  []WebServerCustomizer
	Verifiers    []WebServerVerifier
	Checks       struct {
		Server WebServerCheck
		Log    []logtest.Check
	}
	Errors struct {
		Server []error
		Log    []error
	}
	Recording *logtest.Recording
}

func (w *WebServerTest) WithStaticConfig(values map[string]string) *WebServerTest {
	w.StaticConfig = values
	return w
}

func (w *WebServerTest) WithWebServerPredicate(predicate WebServerPredicate) *WebServerTest {
	w.Checks.Server.Validators = append(w.Checks.Server.Validators, predicate)
	return w
}

func (w *WebServerTest) WithWebServerVerifier(verifier WebServerVerifier) *WebServerTest {
	w.Verifiers = append(w.Verifiers, verifier)
	return w
}

func (w *WebServerTest) WithLogCheck(check logtest.Check) *WebServerTest {
	w.Checks.Log = append(w.Checks.Log, check)
	return w
}

func (w *WebServerTest) WithWebServerCustomizer(customizer WebServerCustomizer) *WebServerTest {
	w.Customizers = append(w.Customizers, customizer)
	return w
}

func (w *WebServerTest) Test(t *testing.T) {
	err := fs.SetSources()
	assert.NoError(t, err)

	if w.Recording == nil {
		w.Recording = logtest.RecordLogging()
	}

	staticConfig := make(map[string]string)
	workdir, _ := os.Getwd()
	staticConfig["server.static-path"] = workdir
	for k, v := range w.StaticConfig {
		staticConfig[k] = v
	}

	cfg := configtest.NewInMemoryConfig(staticConfig)
	ctx := config.ContextWithConfig(context.Background(), cfg)

	managementSecurityConfig, err := NewManagementSecurityConfig(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, managementSecurityConfig)

	webServerConfig := new(WebServerConfig)
	err = config.MustFromContext(ctx).Populate(webServerConfig, configRootWebServer)
	assert.NoError(t, err)

	webServer, err := NewWebServer(webServerConfig, managementSecurityConfig, ctx)
	assert.NoError(t, err)
	assert.NotNil(t, webServer)

	for _, customizer := range w.Customizers {
		customizer(webServer)
	}

	for _, verifier := range w.Verifiers {
		verifier(t, webServer)
	}

	w.Errors.Server = w.Checks.Server.Check(webServer)

	// Check the logs
	for _, logCheck := range w.Checks.Log {
		errs := logCheck.Check(w.Recording)
		w.Errors.Log = append(w.Errors.Log, errs...)
	}

	// Report any errors
	testhelpers.ReportErrors(t, "Server", w.Errors.Server)
	testhelpers.ReportErrors(t, "Log", w.Errors.Log)
}
