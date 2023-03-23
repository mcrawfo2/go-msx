// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package rest

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel"
	"cto-github.cisco.com/NFV-BU/go-msx/skel/text"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"fmt"
	"github.com/mcrawfo2/go-jsonschema/pkg/codegen"
	"path"
)

type DomainServiceUnitTestGenerator struct {
	Domain    string
	Folder    string
	Style     string
	Actions   types.ComparableSlice[string]
	Spec      Spec
	Variables map[string]string
	*text.GoFile
}

func (g DomainServiceUnitTestGenerator) createFixtureStructSnippet() error {
	return g.AddNewText(
		"Fixture/Structure",
		"struct",
		`
			type lowerCamelSingularServiceTestFixture struct {
				// SUT
				UpperCamelSingularService UpperCamelSingularServiceApi
			
				// Mocks
				UpperCamelSingularRepository *MockUpperCamelSingularRepositoryApi
			
				// Canned Values
				ApiData   lowerCamelSingularTestApiData
				ModelData lowerCamelSingularTestModelData
			}
		`,
		[]codegen.Import{})
}

func (g DomainServiceUnitTestGenerator) createFixtureSetupSnippet() error {
	return g.AddNewText(
		"Fixture/Setup",
		"setup",
		`
			func (p lowerCamelSingularServiceTestFixture) Setup(t *testing.T, ctx context.Context) lowerCamelSingularServiceTestFixture {
				svc, err := newUpperCamelSingularService(ctx)
				assert.NoError(t, err)
				assert.NotNil(t, svc)
			
				p.UpperCamelSingularService = svc
				return p
			}
		`,
		[]codegen.Import{
			text.ImportTesting,
			text.ImportContext,
			text.ImportTestifyAssert,
		})
}

func (g DomainServiceUnitTestGenerator) createFixtureConstructorSnippet() error {
	return g.AddNewText(
		"Fixture/Constructor",
		"constructor",
		`
			func newUpperCamelSingularServiceTestFixture() lowerCamelSingularServiceTestFixture {
				return lowerCamelSingularServiceTestFixture{
					ApiData:   newUpperCamelSingularTestApiData(),
					ModelData: newUpperCamelSingularTestModelData(),
				}
			}
		`,
		[]codegen.Import{})
}

func (g DomainServiceUnitTestGenerator) createTestCaseConstructorSnippet() error {
	return g.AddNewText(
		"TestCase/Constructor",
		"constructor",
		`
			func newUpperCamelSingularServiceTest() *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture] {
				return testhelpers.
					NewServiceFixtureCase(new(servicetest.ServiceTest), newUpperCamelSingularServiceTestFixture()).
					WithSetup(func(s *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture]) {
						// Create the mock lowerCamelSingular repository
						s.Fixture.UpperCamelSingularRepository = NewMockUpperCamelSingularRepositoryApi(s.T)
						// Inject it for use by the SUT
						s.Testable.WithContextInjector(func(ctx context.Context) context.Context {
							return contextUpperCamelSingularRepository().Set(ctx, s.Fixture.UpperCamelSingularRepository)
						})
			
						// Inject the mock transaction manager
						s.Testable.WithContextInjector(sqldb.InjectMockTransactionManager)
			
						// All service methods may return an error
						s.Testable.HasErr = true
					})
			}
		`,
		[]codegen.Import{
			text.ImportTestHelpers,
			text.ImportServiceTest,
			text.ImportSqldb,
			text.ImportContext,
		})
}

func (g DomainServiceUnitTestGenerator) createActionListTestSnippet() error {
	if g.Style == StyleV2 {
		g.Variables["list.paging.response"] = `
							v2.PagingResponse{
								Size: 10,
								Pageable: v2.PageableResponse{
									Size: 10,
								},
							}`
	} else {
		g.Variables["list.paging.response"] = `
							v8.PagingResponse{
								PageSize:   10,
								TotalItems: types.PtrTo[int](2),
							}`
	}

	return g.AddNewText(
		"Tests/List",
		"tests",
		`
			func Test_lowerCamelSingularService_ListUpperCamelPlural(t *testing.T) {
				test := newUpperCamelSingularServiceTest().
					WithNamedSetup("repository.FindAllPagedFiltered", func(s *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture]) {
						s.Fixture.UpperCamelSingularRepository.EXPECT().
							FindAllPagedFiltered(
								mock.MatchedBy(testhelpers.AnyContext),
								s.Fixture.ModelData.PagingRequest,
								s.Fixture.ModelData.UpperCamelSingularFilters).
							Return(
								s.Fixture.ModelData.PagingResponse,
								s.Fixture.ModelData.UpperCamelPlural,
								nil)
					}).
					WithNamedSetup("want", func(s *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture]) {
						s.Testable.Want = []any{
							s.Fixture.ApiData.PagingResponse,
							s.Fixture.ApiData.UpperCamelSingularResponses,
							nil,
						}
						s.Testable.WantErr = false
					}).
					WithFunc(func(s *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture], ctx context.Context) {
						// Instantiate the SUT
						s.Fixture = s.Fixture.Setup(s.T, ctx)
			
						pReq := s.Fixture.ApiData.PagingRequest
						fReq := lowerCamelSingularFilterQueryInputs{}
			
						pResp, results, err := s.Fixture.UpperCamelSingularService.ListUpperCamelPlural(ctx, pReq, fReq)
			
						s.Testable.Got = []any{
							pResp, results, err,
						}
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
						name: "SortingRequestError",
						test: test.Clone().
							WithFunc(func(s *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture], ctx context.Context) {
								// Instantiate the SUT
								s.Fixture = s.Fixture.Setup(s.T, ctx)
			
								// Give a faulty sort order
								pReq := ${domain.style}.PagingSortingInputs{
									PagingInputs: ${domain.style}.PagingInputs{
										PageSize: 10,
									},
									SortingInputs: ${domain.style}.SortingInputs{
										SortBy:    "ERROR",
										SortOrder: "asc",
									},
								}
			
								fReq := lowerCamelSingularFilterQueryInputs{}
			
								pResp, results, err := s.Fixture.UpperCamelSingularService.ListUpperCamelPlural(ctx, pReq, fReq)
			
								s.Testable.Got = []any{
									pResp, results, err,
								}
							}).
							WithoutNamedSetup("repository.FindAllPagedFiltered").
							WithNamedSetup("want", func(s *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture]) {
								s.Testable.Want = []any{
									${domain.style}.PagingResponse{},
									[]UpperCamelSingularResponse(nil),
									errors.WithStack(errors.New("ERROR: Unknown sort by")),
								}
								s.Testable.WantErr = true
							}),
					},
					{
						name: "RepositoryError",
						test: test.Clone().
							WithNamedSetup("repository.FindAllPagedFiltered", func(s *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture]) {
								s.Fixture.UpperCamelSingularRepository.EXPECT().
									FindAllPagedFiltered(
										mock.MatchedBy(testhelpers.AnyContext),
										s.Fixture.ModelData.PagingRequest,
										s.Fixture.ModelData.UpperCamelSingularFilters).
									Return(
										paging.Response{},
										nil,
										repository.ErrNotFound)
							}).
							WithNamedSetup("want", func(s *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture]) {
								s.Testable.Want = []any{
									${domain.style}.PagingResponse{},
									[]UpperCamelSingularResponse(nil),
									repository.ErrNotFound,
								}
								s.Testable.WantErr = true
							}),
					},
					{
						name: "SortingResponseError",
						test: test.Clone().
							WithNamedSetup("repository.FindAllPagedFiltered", func(s *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture]) {
								s.Fixture.UpperCamelSingularRepository.EXPECT().
									FindAllPagedFiltered(
										mock.MatchedBy(testhelpers.AnyContext),
										s.Fixture.ModelData.PagingRequest,
										s.Fixture.ModelData.UpperCamelSingularFilters).
									Return(
										paging.Response{
											Size:       10,
											TotalItems: types.PtrTo[uint](2),
											Sort: []paging.SortOrder{{
												Property:  "ERROR",
												Direction: "ASC",
											}},
										},
										s.Fixture.ModelData.UpperCamelPlural,
										nil)
							}).
							WithNamedSetup("want", func(s *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture]) {
								s.Testable.Want = []any{
									${list.paging.response},
									[]UpperCamelSingularResponse(nil),
									errors.WithStack(errors.New("ERROR: Unknown sort by")),
								}
								s.Testable.WantErr = true
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
			text.ImportServiceTest,
			text.ImportTestifyMock,
			text.ImportTesting,
			text.ImportTypes,
			text.ImportErrors,
			text.ImportPaging,
			text.ImportRepository,
			g.importStyle(),
		})
}

func (g DomainServiceUnitTestGenerator) createActionRetrieveTestSnippet() error {
	return g.AddNewText(
		"Tests/Retrieve",
		"tests",
		`
			func Test_lowerCamelSingularService_GetUpperCamelSingular(t *testing.T) {
				test := newUpperCamelSingularServiceTest().
					WithNamedSetup("repository.FindByKey", func(s *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture]) {
						s.Fixture.UpperCamelSingularRepository.EXPECT().
							FindByKey(
								mock.MatchedBy(testhelpers.AnyContext),
								s.Fixture.ModelData.UpperCamelSingularId).
							Return(
								s.Fixture.ModelData.UpperCamelSingular,
								nil)
					}).
					WithNamedSetup("want", func(s *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture]) {
						s.Testable.Want = []any{
							s.Fixture.ApiData.UpperCamelSingularResponse,
							nil,
						}
						s.Testable.WantErr = false
					}).
					WithFunc(func(s *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture], ctx context.Context) {
						// Instantiate the SUT
						s.Fixture = s.Fixture.Setup(s.T, ctx)
			
						result, err := s.Fixture.UpperCamelSingularService.GetUpperCamelSingular(ctx, s.Fixture.ApiData.UpperCamelSingularId)
			
						s.Testable.Got = []any{
							result, err,
						}
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
						name: "RepositoryError",
						test: test.Clone().
							WithNamedSetup("repository.FindByKey", func(s *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture]) {
								s.Fixture.UpperCamelSingularRepository.EXPECT().
									FindByKey(
										mock.MatchedBy(testhelpers.AnyContext),
										s.Fixture.ModelData.UpperCamelSingularId).
									Return(
										UpperCamelSingular{},
										repository.ErrNotFound)
							}).
							WithNamedSetup("want", func(s *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture]) {
								s.Testable.Want = []any{
									UpperCamelSingularResponse{},
									repository.ErrNotFound,
								}
								s.Testable.WantErr = true
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
			text.ImportServiceTest,
			text.ImportTestifyMock,
			text.ImportTesting,
			text.ImportTypes,
			text.ImportErrors,
			text.ImportPaging,
			text.ImportRepository,
			g.importStyle(),
		})
}

func (g DomainServiceUnitTestGenerator) createActionCreateTestSnippet() error {
	return g.AddNewText(
		"Tests/Create",
		"tests",
		`
			func Test_lowerCamelSingularService_CreateUpperCamelSingular(t *testing.T) {
				test := newUpperCamelSingularServiceTest().
					WithNamedSetup("repository.Save", func(s *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture]) {
						s.Fixture.UpperCamelSingularRepository.EXPECT().
							Save(
								mock.MatchedBy(testhelpers.AnyContext),
								mock.AnythingOfType("UpperCamelSingular")).
							Return(
								nil)
					}).
					WithNamedSetup("want", func(s *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture]) {
						s.Testable.Want = []any{
							s.Fixture.ApiData.UpperCamelSingularResponse,
							nil,
						}
						s.Testable.WantErr = false
					}).
					WithFunc(func(s *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture], ctx context.Context) {
						// Instantiate the SUT
						s.Fixture = s.Fixture.Setup(s.T, ctx)
			
						result, err := s.Fixture.UpperCamelSingularService.CreateUpperCamelSingular(ctx, s.Fixture.ApiData.UpperCamelSingularCreateRequest)
						if err == nil {
							// fix for test result diff (lowerCamelSingularid auto-generated by converter)
							result.UpperCamelSingularId = s.Fixture.ApiData.UpperCamelSingularId
						}
			
						s.Testable.Got = []any{
							result, err,
						}
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
						name: "RepositoryError",
						test: test.Clone().
							WithNamedSetup("repository.Save", func(s *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture]) {
								s.Fixture.UpperCamelSingularRepository.EXPECT().
									Save(
										mock.MatchedBy(testhelpers.AnyContext),
										mock.AnythingOfType("UpperCamelSingular")).
									Return(
										repository.ErrNotFound)
							}).
							WithNamedSetup("want", func(s *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture]) {
								s.Testable.Want = []any{
									UpperCamelSingularResponse{},
									repository.ErrNotFound,
								}
								s.Testable.WantErr = true
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
			text.ImportServiceTest,
			text.ImportTestifyMock,
			text.ImportTesting,
			text.ImportTypes,
			text.ImportErrors,
			text.ImportPaging,
			text.ImportRepository,
			g.importStyle(),
		})
}

func (g DomainServiceUnitTestGenerator) createActionUpdateTestSnippet() error {
	return g.AddNewText(
		"Tests/Update",
		"tests",
		`
			func Test_lowerCamelSingularService_UpdateUpperCamelSingular(t *testing.T) {
				test := newUpperCamelSingularServiceTest().
					WithNamedSetup("repository.FindByKey", func(s *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture]) {
						s.Fixture.UpperCamelSingularRepository.EXPECT().
							FindByKey(
								mock.MatchedBy(testhelpers.AnyContext),
								s.Fixture.ModelData.UpperCamelSingularId).
							Return(
								s.Fixture.ModelData.UpperCamelSingular,
								nil)
					}).
					WithNamedSetup("repository.Save", func(s *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture]) {
						s.Fixture.UpperCamelSingularRepository.EXPECT().
							Save(
								mock.MatchedBy(testhelpers.AnyContext),
								s.Fixture.ModelData.UpperCamelSingular).
							Return(
								nil)
					}).
					WithNamedSetup("want", func(s *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture]) {
						s.Testable.Want = []any{
							s.Fixture.ApiData.UpperCamelSingularResponse,
							nil,
						}
						s.Testable.WantErr = false
					}).
					WithFunc(func(s *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture], ctx context.Context) {
						// Instantiate the SUT
						s.Fixture = s.Fixture.Setup(s.T, ctx)
			
						result, err := s.Fixture.UpperCamelSingularService.UpdateUpperCamelSingular(ctx, s.Fixture.ApiData.UpperCamelSingularId, s.Fixture.ApiData.UpperCamelSingularUpdateRequest)
			
						s.Testable.Got = []any{
							result, err,
						}
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
						name: "RepositoryFindError",
						test: test.Clone().
							WithNamedSetup("repository.FindByKey", func(s *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture]) {
								s.Fixture.UpperCamelSingularRepository.EXPECT().
									FindByKey(
										mock.MatchedBy(testhelpers.AnyContext),
										s.Fixture.ModelData.UpperCamelSingularId).
									Return(
										UpperCamelSingular{},
										repository.ErrNotFound)
							}).
							WithoutNamedSetup("repository.Save").
							WithNamedSetup("want", func(s *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture]) {
								s.Testable.Want = []any{
									UpperCamelSingularResponse{},
									repository.ErrNotFound,
								}
								s.Testable.WantErr = true
							}),
					},
					{
						name: "RepositorySaveError",
						test: test.Clone().
							WithNamedSetup("repository.Save", func(s *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture]) {
								s.Fixture.UpperCamelSingularRepository.EXPECT().
									Save(
										mock.MatchedBy(testhelpers.AnyContext),
										s.Fixture.ModelData.UpperCamelSingular).
									Return(
										repository.ErrNotFound)
							}).
							WithNamedSetup("want", func(s *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture]) {
								s.Testable.Want = []any{
									UpperCamelSingularResponse{},
									repository.ErrNotFound,
								}
								s.Testable.WantErr = true
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
			text.ImportServiceTest,
			text.ImportTestifyMock,
			text.ImportTesting,
			text.ImportTypes,
			text.ImportErrors,
			text.ImportPaging,
			text.ImportRepository,
			g.importStyle(),
		})
}

func (g DomainServiceUnitTestGenerator) createActionDeleteTestSnippet() error {
	return g.AddNewText(
		"Tests/Update",
		"tests",
		`
			func Test_lowerCamelSingularService_DeleteUpperCamelSingular(t *testing.T) {
				test := newUpperCamelSingularServiceTest().
					WithNamedSetup("repository.Delete", func(s *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture]) {
						s.Fixture.UpperCamelSingularRepository.EXPECT().
							Delete(
								mock.MatchedBy(testhelpers.AnyContext),
								s.Fixture.ModelData.UpperCamelSingularId).
							Return(
								nil)
					}).
					WithNamedSetup("want", func(s *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture]) {
						s.Testable.Want = []any{
							nil,
						}
						s.Testable.WantErr = false
					}).
					WithFunc(func(s *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture], ctx context.Context) {
						// Instantiate the SUT
						s.Fixture = s.Fixture.Setup(s.T, ctx)
			
						err := s.Fixture.UpperCamelSingularService.DeleteUpperCamelSingular(ctx, s.Fixture.ApiData.UpperCamelSingularId)
			
						s.Testable.Got = []any{
							err,
						}
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
						name: "RepositoryError",
						test: test.Clone().
							WithNamedSetup("repository.Delete", func(s *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture]) {
								s.Fixture.UpperCamelSingularRepository.EXPECT().
									Delete(
										mock.MatchedBy(testhelpers.AnyContext),
										s.Fixture.ModelData.UpperCamelSingularId).
									Return(
										repository.ErrNotFound)
							}).
							WithNamedSetup("want", func(s *testhelpers.ServiceFixtureCase[*servicetest.ServiceTest, lowerCamelSingularServiceTestFixture]) {
								s.Testable.Want = []any{
									repository.ErrNotFound,
								}
								s.Testable.WantErr = true
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
			text.ImportServiceTest,
			text.ImportTestifyMock,
			text.ImportTesting,
			text.ImportTypes,
			text.ImportErrors,
			text.ImportPaging,
			text.ImportRepository,
			g.importStyle(),
		})
}

func (g DomainServiceUnitTestGenerator) Apply(options skel.RenderOptions) skel.RenderOptions {
	options.AddVariable("domain.style", g.Style)
	options.AddVariables(g.Variables)
	return options
}

func (g DomainServiceUnitTestGenerator) Generate() error {
	errs := types.ErrorList{
		g.createFixtureStructSnippet(),
		g.createFixtureSetupSnippet(),
		g.createFixtureConstructorSnippet(),
		g.createTestCaseConstructorSnippet(),
	}

	for _, operation := range g.Spec.Operations {
		if !g.Actions.Contains(operation.Action) {
			continue
		}

		var err error
		switch operation.Action {
		case ActionList:
			err = g.createActionListTestSnippet()
		case ActionRetrieve:
			err = g.createActionRetrieveTestSnippet()
		case ActionCreate:
			err = g.createActionCreateTestSnippet()
		case ActionUpdate:
			err = g.createActionUpdateTestSnippet()
		case ActionDelete:
			err = g.createActionDeleteTestSnippet()
		default:
			// TODO
		}

		errs = append(errs, err)
	}

	return errs.Filter()
}

func (g DomainServiceUnitTestGenerator) Filename() string {
	target := path.Join(g.Folder, fmt.Sprintf("service_lowersingular_%s_test.go", g.Style))
	return g.GoFile.Inflector.Inflect(target)
}

func (g DomainServiceUnitTestGenerator) importStyle() codegen.Import {
	switch g.Style {
	case StyleV2:
		return text.ImportRestOpsV2
	default:
		return text.ImportRestOpsV8
	}
}

func NewDomainServiceUnitTestGenerator(spec Spec) ComponentGenerator {
	return DomainServiceUnitTestGenerator{
		Domain:    generatorConfig.Domain,
		Folder:    generatorConfig.Folder,
		Actions:   generatorConfig.Actions,
		Style:     generatorConfig.Style,
		Spec:      spec,
		Variables: map[string]string{},
		GoFile: &text.GoFile{
			File: &text.File[text.GoSnippet]{
				Comment:   generatorConfig.Style + " API Service Unit Tests for " + generatorConfig.Domain,
				Inflector: text.NewInflector(generatorConfig.Domain),
				Sections: text.NewGoSections(
					&text.Section[text.GoSnippet]{
						Name: "Fixture",
						Sections: text.NewGoSections(
							"Structure",
							"Setup",
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
