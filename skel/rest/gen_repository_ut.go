// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package rest

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel"
	"cto-github.cisco.com/NFV-BU/go-msx/skel/text"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/mcrawfo2/go-jsonschema/pkg/codegen"
	"path"
)

type DomainRepositoryUnitTestGenerator struct {
	Domain  string
	Folder  string
	Actions types.ComparableSlice[string]
	Methods []RepositoryMethod
	Spec    Spec
	*text.GoFile
}

func (g DomainRepositoryUnitTestGenerator) createFixtureStructSnippet() error {
	return g.AddNewText(
		"Fixture/Structure",
		"struct",
		`
			type lowerCamelSingularRepositoryTestFixture struct {
				// SUT
				UpperCamelSingularRepository UpperCamelSingularRepositoryApi
			
				// Mocks
				TypedRepository *sqldb.MockTypedRepositoryApi[UpperCamelSingular]
			
				// Data
				ModelData lowerCamelSingularTestModelData
			}
		`,
		[]codegen.Import{
			text.ImportSqldb,
		})

}

func (g DomainRepositoryUnitTestGenerator) createFixtureSetupSnippet() error {
	return g.AddNewText(
		"Fixture/Setup",
		"setup",
		`
			func (p lowerCamelSingularRepositoryTestFixture) Setup(t *testing.T, ctx context.Context) lowerCamelSingularRepositoryTestFixture {
				svc, err := newUpperCamelSingularRepository(ctx)
				assert.NoError(t, err)
				assert.NotNil(t, svc)
			
				p.UpperCamelSingularRepository = svc
				return p
			}
		`,
		[]codegen.Import{
			text.ImportTesting,
			text.ImportContext,
			text.ImportTestifyAssert,
		})

}

func (g DomainRepositoryUnitTestGenerator) createFixtureConstructorSnippet() error {
	return g.AddNewText(
		"Fixture/Constructor",
		"constructor",
		`
			func newUpperCamelSingularRepositoryTestFixture() lowerCamelSingularRepositoryTestFixture {
				return lowerCamelSingularRepositoryTestFixture{
					ModelData: newUpperCamelSingularTestModelData(),
				}
			}
		`,
		[]codegen.Import{})
}

func (g DomainRepositoryUnitTestGenerator) createTestCaseConstructorSnippet() error {
	return g.AddNewText(
		"TestCase/Constructor",
		"constructor",
		`
			func newUpperCamelSingularRepositoryTest() *testhelpers.ServiceFixtureCase[*sqldbtest.RepositoryTest, lowerCamelSingularRepositoryTestFixture] {
				return testhelpers.
					NewServiceFixtureCase(new(sqldbtest.RepositoryTest), newUpperCamelSingularRepositoryTestFixture()).
					WithSetup(func(s *testhelpers.ServiceFixtureCase[*sqldbtest.RepositoryTest, lowerCamelSingularRepositoryTestFixture]) {
						// Create the mock lowerCamelSingular repository
						s.Fixture.TypedRepository = sqldb.NewMockTypedRepositoryApi[UpperCamelSingular](s.T)
						// Inject it for use by the SUT
						s.Testable.WithContextInjector(func(ctx context.Context) context.Context {
							return sqldb.ContextTypedRepository[UpperCamelSingular](tableNameUpperCamelSingular).Set(ctx, s.Fixture.TypedRepository)
						})
			
						// All repository methods may return an error
						s.Testable.HasErr = true
					})
			}
		`,
		[]codegen.Import{
			text.ImportTestHelpers,
			text.ImportSqldb,
			text.ImportSqlDbTest,
			text.ImportContext,
		})
}

func (g DomainRepositoryUnitTestGenerator) createFindAllSnippet() error {
	return g.AddNewText(
		"Tests/List",
		"tests",
		`
			func Test_lowerCamelSingularTypedRepository_FindAllPagedFiltered(t *testing.T) {
				test := newUpperCamelSingularRepositoryTest().
					WithNamedSetup("typedRepository.FindAll", func(s *testhelpers.ServiceFixtureCase[*sqldbtest.RepositoryTest, lowerCamelSingularRepositoryTestFixture]) {
						s.Fixture.TypedRepository.EXPECT().
							FindAll(
								mock.MatchedBy(testhelpers.AnyContext),
								mock.AnythingOfType("*[]${domain.package.name}.UpperCamelSingular"),
								mock.MatchedBy(testhelpers.Implements[sqldb.FindAllOption]),
								mock.MatchedBy(testhelpers.Implements[sqldb.FindAllOption])).
							Run(func(ctx context.Context, dest *[]UpperCamelSingular, options ...func(*goqu.SelectDataset, paging.Request) (*goqu.SelectDataset, paging.Request)) {
								// Set our output parameter
								*dest = s.Fixture.ModelData.UpperCamelPlural
							}).
							Return(
								s.Fixture.ModelData.PagingResponse,
								nil)
					}).
					WithNamedSetup("want", func(s *testhelpers.ServiceFixtureCase[*sqldbtest.RepositoryTest, lowerCamelSingularRepositoryTestFixture]) {
						s.Testable.Want = []any{
							s.Fixture.ModelData.PagingResponse,
							s.Fixture.ModelData.UpperCamelPlural,
							nil,
						}
						s.Testable.WantErr = false
					}).
					WithFunc(func(s *testhelpers.ServiceFixtureCase[*sqldbtest.RepositoryTest, lowerCamelSingularRepositoryTestFixture], ctx context.Context) {
						s.Fixture = s.Fixture.Setup(s.T, ctx)
			
						pIn := s.Fixture.ModelData.PagingRequest
						pOut, results, err := s.Fixture.UpperCamelSingularRepository.FindAllPagedFiltered(ctx, pIn, lowerCamelSingularFilters{})
			
						s.Testable.Got = []any{
							pOut, results, err,
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
							WithNamedSetup("typedRepository.FindAll", func(s *testhelpers.ServiceFixtureCase[*sqldbtest.RepositoryTest, lowerCamelSingularRepositoryTestFixture]) {
								s.Fixture.TypedRepository.EXPECT().
									FindAll(
										mock.MatchedBy(testhelpers.AnyContext),
										mock.AnythingOfType("*[]${domain.package.name}.UpperCamelSingular"),
										mock.MatchedBy(testhelpers.Implements[sqldb.FindAllOption]),
										mock.MatchedBy(testhelpers.Implements[sqldb.FindAllOption])).
									Return(
										paging.Response{},
										sqldb.ErrNotFound)
							}).
							WithNamedSetup("want", func(s *testhelpers.ServiceFixtureCase[*sqldbtest.RepositoryTest, lowerCamelSingularRepositoryTestFixture]) {
								s.Testable.Want = []any{
									paging.Response{},
									[]UpperCamelSingular(nil),
									sqldb.ErrNotFound,
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
			text.ImportTestifyMock,
			text.ImportTesting,
			text.ImportPaging,
			text.ImportRepository,
			text.ImportSqldb,
			text.ImportSqlDbTest,
			text.ImportGoqu,
		})

}

func (g DomainRepositoryUnitTestGenerator) createFindOneSnippet() error {
	return g.AddNewText(
		"Tests/Retrieve",
		"tests",
		`
			func Test_lowerCamelSingularTypedRepository_FindByKey(t *testing.T) {
				test := newUpperCamelSingularRepositoryTest().
					WithNamedSetup("typedRepository.FindOne", func(s *testhelpers.ServiceFixtureCase[*sqldbtest.RepositoryTest, lowerCamelSingularRepositoryTestFixture]) {
						s.Fixture.TypedRepository.EXPECT().
							FindOne(
								mock.MatchedBy(testhelpers.AnyContext),
								mock.AnythingOfType("*${domain.package.name}.UpperCamelSingular"),
								mock.MatchedBy(testhelpers.Implements[sqldb.WhereOption])).
							Run(func(ctx context.Context, dest *UpperCamelSingular, where sqldb.WhereOption) {
								// Set our output parameter
								*dest = s.Fixture.ModelData.UpperCamelSingular
							}).
							Return(nil)
					}).
					WithNamedSetup("want", func(s *testhelpers.ServiceFixtureCase[*sqldbtest.RepositoryTest, lowerCamelSingularRepositoryTestFixture]) {
						s.Testable.Want = []any{
							s.Fixture.ModelData.UpperCamelSingular,
							nil,
						}
						s.Testable.WantErr = false
					}).
					WithFunc(func(s *testhelpers.ServiceFixtureCase[*sqldbtest.RepositoryTest, lowerCamelSingularRepositoryTestFixture], ctx context.Context) {
						s.Fixture = s.Fixture.Setup(s.T, ctx)
			
						result, err := s.Fixture.UpperCamelSingularRepository.FindByKey(ctx, s.Fixture.ModelData.UpperCamelSingularId)
			
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
						name: "NotFoundError",
						test: test.Clone().
							WithNamedSetup("typedRepository.FindOne", func(s *testhelpers.ServiceFixtureCase[*sqldbtest.RepositoryTest, lowerCamelSingularRepositoryTestFixture]) {
								s.Fixture.TypedRepository.EXPECT().
									FindOne(
										mock.MatchedBy(testhelpers.AnyContext),
										mock.AnythingOfType("*${domain.package.name}.UpperCamelSingular"),
										mock.MatchedBy(testhelpers.Implements[sqldb.WhereOption])).
									Return(sqldb.ErrNotFound)
							}).
							WithNamedSetup("want", func(s *testhelpers.ServiceFixtureCase[*sqldbtest.RepositoryTest, lowerCamelSingularRepositoryTestFixture]) {
								s.Testable.Want = []any{
									UpperCamelSingular{},
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
			text.ImportTestifyMock,
			text.ImportTesting,
			text.ImportRepository,
			text.ImportSqldb,
			text.ImportSqlDbTest,
			text.ImportGoqu,
		})

}

func (g DomainRepositoryUnitTestGenerator) createSaveSnippet() error {
	return g.AddNewText(
		"Tests/Save",
		"tests",
		`
			func Test_lowerCamelSingularTypedRepository_Save(t *testing.T) {
				test := newUpperCamelSingularRepositoryTest().
					WithNamedSetup("typedRepository.Upsert", func(s *testhelpers.ServiceFixtureCase[*sqldbtest.RepositoryTest, lowerCamelSingularRepositoryTestFixture]) {
						s.Fixture.TypedRepository.EXPECT().
							Upsert(
								mock.MatchedBy(testhelpers.AnyContext),
								mock.AnythingOfType("${domain.package.name}.UpperCamelSingular")).
							Return(nil)
					}).
					WithNamedSetup("want", func(s *testhelpers.ServiceFixtureCase[*sqldbtest.RepositoryTest, lowerCamelSingularRepositoryTestFixture]) {
						s.Testable.Want = []any{
							nil,
						}
						s.Testable.WantErr = false
					}).
					WithFunc(func(s *testhelpers.ServiceFixtureCase[*sqldbtest.RepositoryTest, lowerCamelSingularRepositoryTestFixture], ctx context.Context) {
						s.Fixture = s.Fixture.Setup(s.T, ctx)
			
						err := s.Fixture.UpperCamelSingularRepository.Save(ctx, s.Fixture.ModelData.UpperCamelSingular)
			
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
						name: "NotFoundError",
						test: test.Clone().
							WithNamedSetup("typedRepository.Upsert", func(s *testhelpers.ServiceFixtureCase[*sqldbtest.RepositoryTest, lowerCamelSingularRepositoryTestFixture]) {
								s.Fixture.TypedRepository.EXPECT().
									Upsert(
										mock.MatchedBy(testhelpers.AnyContext),
										mock.AnythingOfType("${domain.package.name}.UpperCamelSingular")).
									Return(sqldb.ErrNotFound)
							}).
							WithNamedSetup("want", func(s *testhelpers.ServiceFixtureCase[*sqldbtest.RepositoryTest, lowerCamelSingularRepositoryTestFixture]) {
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
			text.ImportTestifyMock,
			text.ImportTesting,
			text.ImportRepository,
			text.ImportSqldb,
			text.ImportSqlDbTest,
			text.ImportGoqu,
		})
}

func (g DomainRepositoryUnitTestGenerator) createDeleteSnippet() error {
	return g.AddNewText(
		"Tests/Delete",
		"tests",
		`
			func Test_lowerCamelSingularTypedRepository_Delete(t *testing.T) {
				test := newUpperCamelSingularRepositoryTest().
					WithNamedSetup("typedRepository.DeleteOne", func(s *testhelpers.ServiceFixtureCase[*sqldbtest.RepositoryTest, lowerCamelSingularRepositoryTestFixture]) {
						s.Fixture.TypedRepository.EXPECT().
							DeleteOne(
								mock.MatchedBy(testhelpers.AnyContext),
								mock.MatchedBy(testhelpers.Implements[sqldb.KeysOption])).
							Return(nil)
					}).
					WithNamedSetup("want", func(s *testhelpers.ServiceFixtureCase[*sqldbtest.RepositoryTest, lowerCamelSingularRepositoryTestFixture]) {
						s.Testable.Want = []any{
							nil,
						}
						s.Testable.WantErr = false
					}).
					WithFunc(func(s *testhelpers.ServiceFixtureCase[*sqldbtest.RepositoryTest, lowerCamelSingularRepositoryTestFixture], ctx context.Context) {
						s.Fixture = s.Fixture.Setup(s.T, ctx)
			
						err := s.Fixture.UpperCamelSingularRepository.Delete(ctx, s.Fixture.ModelData.UpperCamelSingularId)
			
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
						name: "NotFoundError",
						test: test.Clone().
							WithNamedSetup("typedRepository.DeleteOne", func(s *testhelpers.ServiceFixtureCase[*sqldbtest.RepositoryTest, lowerCamelSingularRepositoryTestFixture]) {
								s.Fixture.TypedRepository.EXPECT().
									DeleteOne(
										mock.MatchedBy(testhelpers.AnyContext),
										mock.MatchedBy(testhelpers.Implements[sqldb.KeysOption])).
									Return(sqldb.ErrNotFound)
							}).
							WithNamedSetup("want", func(s *testhelpers.ServiceFixtureCase[*sqldbtest.RepositoryTest, lowerCamelSingularRepositoryTestFixture]) {
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
			text.ImportTestifyMock,
			text.ImportTesting,
			text.ImportRepository,
			text.ImportSqldb,
			text.ImportSqlDbTest,
			text.ImportGoqu,
		})
}

func (g DomainRepositoryUnitTestGenerator) Generate() error {
	errs := types.ErrorList{
		g.createFixtureStructSnippet(),
		g.createFixtureSetupSnippet(),
		g.createFixtureConstructorSnippet(),
		g.createTestCaseConstructorSnippet(),
	}

	for _, method := range g.Methods {
		if !g.Actions.ContainsAny(method.Actions...) {
			continue
		}

		var err error
		switch method.Method {
		case RepositoryMethodFindAll:
			err = g.createFindAllSnippet()
		case RepositoryMethodFindOne:
			err = g.createFindOneSnippet()
		case RepositoryMethodSave:
			err = g.createSaveSnippet()
		case RepositoryMethodDelete:
			err = g.createDeleteSnippet()
		}

		errs = append(errs, err)
	}

	return errs.Filter()
}

func (g DomainRepositoryUnitTestGenerator) Apply(options skel.RenderOptions) skel.RenderOptions {
	options.AddVariable("domain.package.name", g.Package)
	return options
}

func (g DomainRepositoryUnitTestGenerator) Filename() string {
	target := path.Join(g.Folder, "repository_lowersingular_test.go")
	return g.GoFile.Inflector.Inflect(target)
}

func NewDomainRepositoryUnitTestGenerator(spec Spec) ComponentGenerator {
	return DomainRepositoryUnitTestGenerator{
		Domain:  generatorConfig.Domain,
		Folder:  generatorConfig.Folder,
		Actions: generatorConfig.Actions,
		Methods: []RepositoryMethod{
			{
				Method:  RepositoryMethodFindAll,
				Actions: []string{ActionList},
			},
			{
				Method:  RepositoryMethodFindOne,
				Actions: []string{ActionRetrieve, ActionUpdate},
			},
			{
				Method:  RepositoryMethodSave,
				Actions: []string{ActionCreate, ActionUpdate},
			},
			{
				Method:  RepositoryMethodDelete,
				Actions: []string{ActionDelete},
			},
		},
		Spec: spec,
		GoFile: &text.GoFile{
			File: &text.File[text.GoSnippet]{
				Comment:   "Repository Unit Tests for " + generatorConfig.Domain,
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
							"Save",
							"Delete"),
					},
				),
			},
			Package: generatorConfig.PackageName(),
		},
	}
}
