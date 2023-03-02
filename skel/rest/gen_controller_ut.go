package rest

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel"
	"cto-github.cisco.com/NFV-BU/go-msx/skel/text"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"fmt"
	"github.com/mcrawfo2/go-jsonschema/pkg/codegen"
	"path"
	"strings"
)

type DomainControllerUnitTestGenerator struct {
	Domain     string
	Folder     string
	Style      string
	Tenant     string
	Actions    types.ComparableSlice[string]
	Components []string
	Variables  map[string]string

	Spec Spec

	*text.GoFile
}

func (g DomainControllerUnitTestGenerator) createFixtureSnippet() error {
	return g.AddNewText(
		"Fixture/Structure",
		"struct",
		`
			// lowerCamelSingularControllerTestFixture contains test dependencies for Controller tests.
			type lowerCamelSingularControllerTestFixture struct {
				// SUT
				Controller *lowerCamelSingularController
			
				// Mocks
				UpperCamelSingularService *MockUpperCamelSingularServiceApi
			
				// Canned API Values
				ApiData lowerCamelSingularTestApiData
			}
		`,
		[]codegen.Import{})
}

func (g DomainControllerUnitTestGenerator) createFixtureConstructorSnippet() error {
	return g.AddNewText(
		"Fixture/Constructor",
		"constructor",
		`
			// newUpperCamelSingularControllerTestFixture creates a new fixture for Controller tests
			func newUpperCamelSingularControllerTestFixture() lowerCamelSingularControllerTestFixture {
				return lowerCamelSingularControllerTestFixture{
					ApiData: newUpperCamelSingularTestApiData(),
				}
			}
		`,
		[]codegen.Import{})
}

func (g DomainControllerUnitTestGenerator) createTestCaseConstructorSnippet() error {
	return g.AddNewText(
		"TestCase/Constructor",
		"constructor",
		`
			func newUpperCamelSingularControllerTest() *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture] {
				return testhelpers.
					NewFixtureCase(new(controllertest.ControllerTest), newUpperCamelSingularControllerTestFixture()).
					WithSetup(func(c *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture]) {
						c.Fixture.UpperCamelSingularService = NewMockUpperCamelSingularServiceApi(c.T)
					}).
					WithSetup(func(c *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture]) {
						// Set token provider, controller
						c.Testable.
							WithTokenDetailsProvider(securitytest.NewMockTokenDetailsProvider()).
							WithEndpointProducerSourceFactory(func(ctx context.Context) (restops.EndpointsProducer, error) {
								// Save the controller
								controller, err := newUpperCamelSingularController(ctx)
								if err == nil {
									c.Fixture.Controller = controller.(*lowerCamelSingularController)
								}
								return controller, err
							}).
							WithContextInjector(func(ctx context.Context) context.Context {
								return contextUpperCamelSingularService().Set(ctx, c.Fixture.UpperCamelSingularService)
							})
					}).
					WithNamedSetup("permissions", func(p *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture]) {
						// Allow view and manage permissions by default
						p.Testable.Context.TokenDetails.Permissions = []string{
							permissionViewUpperCamelPlural,
							permissionManageUpperCamelPlural,
						}
					})
			}
		`,
		[]codegen.Import{
			text.ImportTestHelpers,
			text.ImportControllerTest,
			text.ImportSecurityTest,
			text.ImportRestOps,
			text.ImportContext,
		})
}

func (g DomainControllerUnitTestGenerator) cleanOperationId(operationId string) string {
	cleanOperationId := operationId
	if strings.Contains(cleanOperationId, ".") {
		lastPeriod := strings.LastIndex(cleanOperationId, ".")
		cleanOperationId = cleanOperationId[lastPeriod+1:]
	}
	return cleanOperationId
}

func (g DomainControllerUnitTestGenerator) createEndpointActionListTestSnippet(operation Operation) error {
	if g.Style == StyleV2 {
		g.Variables["list.injected.field"] = "responseObject.content"
		g.Variables["list.validation.field"] = "debug"
	} else {
		g.Variables["list.injected.field"] = "contents"
		g.Variables["list.validation.field"] = "details"
	}

	return g.AddNewText(
		"Tests/List",
		"tests",
		`
			func Test_lowerCamelSingularController_listUpperCamelPlural(t *testing.T) {
				test := newUpperCamelSingularControllerTest().
					WithSetup(func(p *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture]) {
						// Always point to the list endpoint
						p.Testable.
							WithRequestMethod(http.MethodGet).
							WithRequestPath("/api/${domain.style}/lowerCamelPlural", nil)
					}).
					WithNamedSetup("query", func(p *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture]) {
						p.Testable.
							WithRequestQueryParameter("page", "0").
							WithRequestQueryParameter("pageSize", "10")
					}).
					WithNamedSetup("service.ListUpperCamelPlural", func(p *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture]) {
						p.Fixture.UpperCamelSingularService.EXPECT().
							ListUpperCamelPlural(
								mock.MatchedBy(testhelpers.AnyContext),
								mock.AnythingOfType("${domain.style}.PagingSortingInputs"),
								lowerCamelSingularFilterQueryInputs{}).
							Return(p.Fixture.ApiData.PagingResponse, p.Fixture.ApiData.UpperCamelSingularResponses, nil)
					}).
					WithNamedSetup("response", func(p *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture]) {
						p.Testable.
							WithResponsePredicate(webservicetest.ResponseHasStatus(http.StatusOK)).
							WithResponsePredicate(webservicetest.ResponseHasBodyJsonValue(
								"${list.injected.field}.#",
								float64(len(p.Fixture.ApiData.UpperCamelSingularResponses))))
					})
			
				tests := []struct {
					name string
					test testhelpers.Testable
				}{
					{
						name: "Success",
						test: test.Clone(),
					},
					{
						name: "MissingQuery",
						test: test.Clone().
							WithoutNamedSetup("query").
							WithoutNamedSetup("service.ListUpperCamelPlural").
							WithNamedSetup("response", func(p *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture]) {
								p.Testable.
									WithResponsePredicate(webservicetest.ResponseHasStatus(http.StatusBadRequest)).
									WithResponsePredicate(webservicetest.ResponseHasBodyJsonValue(
										"${list.validation.field}.Page.\\.failures.#",
										float64(1))).
									WithResponsePredicate(webservicetest.ResponseHasBodyJsonValue(
										"${list.validation.field}.Page.\\.failures.0",
										"expected integer, but got null")).
									WithResponsePredicate(webservicetest.ResponseHasBodyJsonValue(
										"${list.validation.field}.PageSize.\\.failures.#",
										float64(1))).
									WithResponsePredicate(webservicetest.ResponseHasBodyJsonValue(
										"${list.validation.field}.PageSize.\\.failures.0",
										"expected integer, but got null"))
							}),
					},
					{
						name: "MissingPermission",
						test: test.Clone().
							WithoutNamedSetup("permissions").
							WithoutNamedSetup("service.ListUpperCamelPlural").
							WithNamedSetup("response", func(p *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture]) {
								p.Testable.
									WithResponsePredicate(webservicetest.ResponseHasStatus(http.StatusForbidden))
							}),
					},
				}
			
				for _, tt := range tests {
					t.Run(tt.name, tt.test.Test)
				}
			}
		`,
		[]codegen.Import{
			text.ImportTestHelpers,
			text.ImportControllerTest,
			text.ImportWebServiceTest,
			text.ImportHttp,
			text.ImportTestifyMock,
			text.ImportTesting,
		})
}

func (g DomainControllerUnitTestGenerator) createEndpointActionRetrieveTestSnippet(operation Operation) error {
	return g.AddNewText(
		"Tests/Retrieve",
		"tests",
		`
			func Test_lowerCamelSingularController_getUpperCamelSingular(t *testing.T) {
				test := newUpperCamelSingularControllerTest().
					WithSetup(func(p *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture]) {
						// Always point to the retrieve endpoint
						p.Testable.
							WithRequestMethod(http.MethodGet).
							WithRequestPath("/api/${domain.style}/lowerCamelPlural/{lowerCamelSingularId}", map[string]string{
								"lowerCamelSingularId": p.Fixture.ApiData.UpperCamelSingularId.String(),
							})
					}).
					WithNamedSetup("service.GetUpperCamelSingular", func(p *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture]) {
						p.Fixture.UpperCamelSingularService.EXPECT().
							GetUpperCamelSingular(mock.MatchedBy(testhelpers.AnyContext), p.Fixture.ApiData.UpperCamelSingularId).
							Return(p.Fixture.ApiData.UpperCamelSingularResponse, nil)
					}).
					WithNamedSetup("response", func(p *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture]) {
						p.Testable.
							WithResponsePredicate(webservicetest.ResponseHasStatus(http.StatusOK))
					})
			
				tests := []struct {
					name string
					test testhelpers.Testable
				}{
					{
						name: "Success",
						test: test.Clone(),
					},
					{
						name: "NotFound",
						test: test.Clone().
							WithNamedSetup("service.GetUpperCamelSingular", func(p *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture]) {
								p.Fixture.UpperCamelSingularService.EXPECT().
									GetUpperCamelSingular(mock.MatchedBy(testhelpers.AnyContext), p.Fixture.ApiData.UpperCamelSingularId).
									Return(UpperCamelSingularResponse{}, errors.Wrapf(repository.ErrNotFound, "UpperCamelSingular not found"))
							}).
							WithNamedSetup("response", func(p *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture]) {
								p.Testable.
									WithResponsePredicate(webservicetest.ResponseHasStatus(http.StatusNotFound))
							}),
					},
					{
						name: "MissingPermission",
						test: test.Clone().
							WithoutNamedSetup("permissions").
							WithoutNamedSetup("service.GetUpperCamelSingular").
							WithNamedSetup("response", func(p *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture]) {
								p.Testable.
									WithResponsePredicate(webservicetest.ResponseHasStatus(http.StatusForbidden))
							}),
					},
				}
			
				for _, tt := range tests {
					t.Run(tt.name, tt.test.Test)
				}
			}
		`,
		[]codegen.Import{
			text.ImportTestHelpers,
			text.ImportControllerTest,
			text.ImportWebServiceTest,
			text.ImportHttp,
			text.ImportTestifyMock,
			text.ImportTesting,
			text.ImportRepository,
			text.ImportErrors,
		})
}

func (g DomainControllerUnitTestGenerator) createEndpointActionCreateTestSnippet(operation Operation) error {
	return g.AddNewText(
		"Tests/Create",
		"tests",
		`
			func Test_lowerCamelSingularController_createUpperCamelSingular(t *testing.T) {
				test := newUpperCamelSingularControllerTest().
					WithSetup(func(p *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture]) {
						// Always point to the create endpoint
						p.Testable.
							WithRequestMethod(http.MethodPost).
							WithRequestPath("/api/${domain.style}/lowerCamelPlural", nil).
							WithRequestBodyJson(p.Fixture.ApiData.UpperCamelSingularCreateRequest).
							WithRequestHeader(restops.HeaderContentType, restops.ContentTypeJson)
					}).
					WithNamedSetup("service.CreateUpperCamelSingular", func(p *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture]) {
						p.Fixture.UpperCamelSingularService.EXPECT().
							CreateUpperCamelSingular(mock.MatchedBy(testhelpers.AnyContext), p.Fixture.ApiData.UpperCamelSingularCreateRequest).
							Return(p.Fixture.ApiData.UpperCamelSingularResponse, nil)
					}).
					WithNamedSetup("response", func(p *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture]) {
						p.Testable.
							WithResponsePredicate(webservicetest.ResponseHasStatus(http.StatusCreated))
					})
			
				tests := []struct {
					name string
					test testhelpers.Testable
				}{
					{
						name: "Success",
						test: test.Clone(),
					},
					{
						name: "AlreadyExists",
						test: test.Clone().
							WithNamedSetup("service.CreateUpperCamelSingular", func(p *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture]) {
								p.Fixture.UpperCamelSingularService.EXPECT().
									CreateUpperCamelSingular(mock.MatchedBy(testhelpers.AnyContext), p.Fixture.ApiData.UpperCamelSingularCreateRequest).
									Return(UpperCamelSingularResponse{}, errors.Wrapf(repository.ErrAlreadyExists, "UpperCamelSingular already exists"))
							}).
							WithNamedSetup("response", func(p *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture]) {
								p.Testable.
									WithResponsePredicate(webservicetest.ResponseHasStatus(http.StatusConflict))
							}),
					},
					{
						name: "MissingPermission",
						test: test.Clone().
							WithoutNamedSetup("permissions").
							WithoutNamedSetup("service.CreateUpperCamelSingular").
							WithNamedSetup("response", func(p *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture]) {
								p.Testable.
									WithResponsePredicate(webservicetest.ResponseHasStatus(http.StatusForbidden))
							}),
					},
				}
			
				for _, tt := range tests {
					t.Run(tt.name, tt.test.Test)
				}
			}
		`,
		[]codegen.Import{
			text.ImportTestHelpers,
			text.ImportControllerTest,
			text.ImportWebServiceTest,
			text.ImportHttp,
			text.ImportTestifyMock,
			text.ImportTesting,
			text.ImportRepository,
			text.ImportErrors,
		})
}

func (g DomainControllerUnitTestGenerator) createEndpointActionUpdateTestSnippet(operation Operation) error {
	return g.AddNewText(
		"Tests/Update",
		"tests",
		`
			func Test_lowerCamelSingularController_updateUpperCamelSingular(t *testing.T) {
				test := newUpperCamelSingularControllerTest().
					WithSetup(func(p *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture]) {
						// Always point to the update endpoint
						p.Testable.
							WithRequestMethod(http.MethodPut).
							WithRequestPath("/api/${domain.style}/lowerCamelPlural/{lowerCamelSingularId}", map[string]string{
								"lowerCamelSingularId": p.Fixture.ApiData.UpperCamelSingularId.String(),
							}).
							WithRequestBodyJson(p.Fixture.ApiData.UpperCamelSingularUpdateRequest).
							WithRequestHeader(restops.HeaderContentType, restops.ContentTypeJson)
					}).
					WithNamedSetup("service.UpdateUpperCamelSingular", func(p *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture]) {
						p.Fixture.UpperCamelSingularService.EXPECT().
							UpdateUpperCamelSingular(mock.MatchedBy(testhelpers.AnyContext), p.Fixture.ApiData.UpperCamelSingularId, p.Fixture.ApiData.UpperCamelSingularUpdateRequest).
							Return(p.Fixture.ApiData.UpperCamelSingularResponse, nil)
					}).
					WithNamedSetup("response", func(p *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture]) {
						p.Testable.
							WithResponsePredicate(webservicetest.ResponseHasStatus(http.StatusOK))
					})
			
				tests := []struct {
					name string
					test testhelpers.Testable
				}{
					{
						name: "Success",
						test: test.Clone(),
					},
					{
						name: "NotFound",
						test: test.Clone().
							WithNamedSetup("service.UpdateUpperCamelSingular", func(p *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture]) {
								p.Fixture.UpperCamelSingularService.EXPECT().
									UpdateUpperCamelSingular(mock.MatchedBy(testhelpers.AnyContext), p.Fixture.ApiData.UpperCamelSingularId, p.Fixture.ApiData.UpperCamelSingularUpdateRequest).
									Return(UpperCamelSingularResponse{}, errors.Wrapf(repository.ErrNotFound, "UpperCamelSingular not found"))
							}).
							WithNamedSetup("response", func(p *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture]) {
								p.Testable.
									WithResponsePredicate(webservicetest.ResponseHasStatus(http.StatusNotFound))
							}),
					},
					{
						name: "MissingPermission",
						test: test.Clone().
							WithoutNamedSetup("permissions").
							WithoutNamedSetup("service.UpdateUpperCamelSingular").
							WithNamedSetup("response", func(p *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture]) {
								p.Testable.
									WithResponsePredicate(webservicetest.ResponseHasStatus(http.StatusForbidden))
							}),
					},
				}
			
				for _, tt := range tests {
					t.Run(tt.name, tt.test.Test)
				}
			}
		`,
		[]codegen.Import{
			text.ImportTestHelpers,
			text.ImportControllerTest,
			text.ImportWebServiceTest,
			text.ImportHttp,
			text.ImportTestifyMock,
			text.ImportTesting,
			text.ImportRepository,
			text.ImportErrors,
		})
}

func (g DomainControllerUnitTestGenerator) createEndpointActionDeleteTestSnippet(operation Operation) error {
	if g.Style == StyleV2 {
		g.Variables["delete.success.code"] = "http.StatusOK"
	} else {
		g.Variables["delete.success.code"] = "http.StatusNoContent"
	}

	return g.AddNewText(
		"Tests/Delete",
		"tests",
		`
			func Test_lowerCamelSingularController_deleteUpperCamelSingular(t *testing.T) {
				test := newUpperCamelSingularControllerTest().
					WithSetup(func(p *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture]) {
						// Always point to the update endpoint
						p.Testable.
							WithRequestMethod(http.MethodDelete).
							WithRequestPath("/api/${domain.style}/lowerCamelPlural/{lowerCamelSingularId}", map[string]string{
								"lowerCamelSingularId": p.Fixture.ApiData.UpperCamelSingularId.String(),
							})
					}).
					WithNamedSetup("service.DeleteUpperCamelSingular", func(p *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture]) {
						p.Fixture.UpperCamelSingularService.EXPECT().
							DeleteUpperCamelSingular(mock.MatchedBy(testhelpers.AnyContext), p.Fixture.ApiData.UpperCamelSingularId).
							Return(nil)
					}).
					WithNamedSetup("response", func(p *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture]) {
						p.Testable.
							WithResponsePredicate(webservicetest.ResponseHasStatus(${delete.success.code}))
					})
			
				tests := []struct {
					name string
					test testhelpers.Testable
				}{
					{
						name: "Success",
						test: test.Clone(),
					},
					{
						name: "NotFound",
						test: test.Clone().
							WithNamedSetup("service.DeleteUpperCamelSingular", func(p *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture]) {
								p.Fixture.UpperCamelSingularService.EXPECT().
									DeleteUpperCamelSingular(mock.MatchedBy(testhelpers.AnyContext), p.Fixture.ApiData.UpperCamelSingularId).
									Return(errors.Wrapf(repository.ErrNotFound, "UpperCamelSingular not found"))
							}).
							WithNamedSetup("response", func(p *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture]) {
								p.Testable.
									WithResponsePredicate(webservicetest.ResponseHasStatus(http.StatusNotFound))
							}),
					},
					{
						name: "MissingPermission",
						test: test.Clone().
							WithoutNamedSetup("permissions").
							WithoutNamedSetup("service.DeleteUpperCamelSingular").
							WithNamedSetup("response", func(p *testhelpers.FixtureCase[*controllertest.ControllerTest, lowerCamelSingularControllerTestFixture]) {
								p.Testable.
									WithResponsePredicate(webservicetest.ResponseHasStatus(http.StatusForbidden))
							}),
					},
				}
			
				for _, tt := range tests {
					t.Run(tt.name, tt.test.Test)
				}
			}
		`,
		[]codegen.Import{
			text.ImportTestHelpers,
			text.ImportControllerTest,
			text.ImportWebServiceTest,
			text.ImportHttp,
			text.ImportTestifyMock,
			text.ImportTesting,
			text.ImportRepository,
			text.ImportErrors,
		})
}

func (g DomainControllerUnitTestGenerator) Apply(options skel.RenderOptions) skel.RenderOptions {
	options.AddVariable("domain.style", g.Style)
	options.AddVariables(g.Variables)
	return options
}

func (g DomainControllerUnitTestGenerator) Generate() error {
	errs := types.ErrorList{
		g.createFixtureSnippet(),
		g.createFixtureConstructorSnippet(),
		g.createTestCaseConstructorSnippet(),
	}

	var operationMethods []string
	for _, operation := range g.Spec.Operations {
		if !g.Actions.Contains(operation.Action) {
			continue
		}

		// For the Endpoints Producer
		operationMethods = append(operationMethods, g.cleanOperationId(*operation.Operation.ID))

		var err error
		switch operation.Action {
		case ActionList:
			err = g.createEndpointActionListTestSnippet(operation)
		case ActionRetrieve:
			err = g.createEndpointActionRetrieveTestSnippet(operation)
		case ActionCreate:
			err = g.createEndpointActionCreateTestSnippet(operation)
		case ActionUpdate:
			err = g.createEndpointActionUpdateTestSnippet(operation)
		case ActionDelete:
			err = g.createEndpointActionDeleteTestSnippet(operation)
		}

		errs = append(errs, err)
	}

	return errs.Filter()
}

func (g DomainControllerUnitTestGenerator) Filename() string {
	target := path.Join(g.Folder, fmt.Sprintf("controller_lowersingular_%s_test.go", g.Style))
	return g.GoFile.Inflector.Inflect(target)
}

func (g DomainControllerUnitTestGenerator) importStyle() codegen.Import {
	switch g.Style {
	case StyleV2:
		return text.ImportRestOpsV2
	default:
		return text.ImportRestOpsV8
	}
}

func NewDomainControllerUnitTestGenerator(spec Spec) ComponentGenerator {
	inflector := skel.NewInflector(generatorConfig.Domain)

	return DomainControllerUnitTestGenerator{
		// Configuration
		Domain:     generatorConfig.Domain,
		Folder:     generatorConfig.Folder,
		Tenant:     generatorConfig.Tenant,
		Actions:    generatorConfig.Actions,
		Components: generatorConfig.Components,
		Style:      generatorConfig.Style,
		Variables:  map[string]string{},

		// Sources
		Spec: spec,

		GoFile: &text.GoFile{
			File: &text.File[text.GoSnippet]{
				Comment: fmt.Sprintf(
					"%s API REST Controller Unit Tests for %s",
					generatorConfig.Style,
					generatorConfig.Domain),
				Inflector: inflector,
				Sections: text.NewGoSections(
					&text.Section[text.GoSnippet]{
						Name: "Fixture",
						Sections: text.NewGoSections(
							"Structure",
							"Constructor"),
					},
					&text.Section[text.GoSnippet]{
						Name: "TestCase",
						Sections: text.NewGoSections(
							"Constructor"),
					},
					&text.Section[text.GoSnippet]{
						Name: "Tests",
						Sections: text.NewGoSections(
							"List",
							"Retrieve",
							"Create",
							"Update",
							"Delete",
						),
					},
				),
			},
			Package: generatorConfig.PackageName(),
		},
	}
}
