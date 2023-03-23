// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package rest

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel/text"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/mcrawfo2/go-jsonschema/pkg/codegen"
	"path"
)

type RepositoryMethod struct {
	Actions []string
	Method  string
}

const (
	RepositoryMethodFindAll = "FindAllPagedFiltered"
	RepositoryMethodFindOne = "FindByKey"
	RepositoryMethodSave    = "Save"
	RepositoryMethodDelete  = "Delete"
)

type DomainRepositoryGenerator struct {
	Domain  string
	Folder  string
	Actions types.ComparableSlice[string]
	Methods []RepositoryMethod
	Spec    Spec
	*text.GoFile
}

func (g DomainRepositoryGenerator) createConstantsSnippet() error {
	return g.AddNewGenerator(
		"Constants",
		"constants",
		text.GoConstants{
			{
				Name:  "tableNameUpperCamelSingular",
				Value: "lower_snake_singular",
			},
			{
				Name:  "columnUpperCamelSingularId",
				Value: "lower_snake_singular_id",
			},
		},
		nil)
}

func (g DomainRepositoryGenerator) createApiSnippet() error {
	return g.AddNewText(
		"API",
		"interface",
		`
			//go:generate mockery --name=UpperCamelSingularRepositoryApi --testonly --case=snake --inpackage --with-expecter
			
			// UpperCamelSingularRepositoryApi declares the interface for the UpperCamelSingularRepository.  This can be used
			// to interchange implementations, such as during testing.
			type UpperCamelSingularRepositoryApi interface {
				FindAllPagedFiltered(ctx context.Context, pin paging.Request, fin lowerCamelSingularFilters) (paging.Response, []UpperCamelSingular, error)
				FindByKey(ctx context.Context, lowerCamelSingularId uuid.UUID) (UpperCamelSingular, error)
				Save(ctx context.Context, lowerCamelSingular UpperCamelSingular) error
				Delete(ctx context.Context, lowerCamelSingularId uuid.UUID) error
			}
		`,
		[]codegen.Import{
			text.ImportContext,
			text.ImportPaging,
			text.ImportUuid,
		})
}

func (g DomainRepositoryGenerator) createRepositorySnippet() error {
	return g.AddNewText(
		"Repository",
		"struct",
		`
			// lowerCamelSingularTypedRepository is the production implementation of the UpperCamelSingularRepositoryApi.
			type lowerCamelSingularTypedRepository struct {
				typedRepository sqldb.TypedRepositoryApi[UpperCamelSingular]
			}
			`,
		[]codegen.Import{
			text.ImportSqldb,
		})
}

func (g DomainRepositoryGenerator) createFindAllSnippet() error {
	return g.AddNewText(
		"Actions/List",
		"FindAllPagedFiltered",
		`
			// FindAllPagedFiltered returns a series of matching UpperCamelSingular records using the specified
			// filtering and pagination criteria
			func (r *lowerCamelSingularTypedRepository) FindAllPagedFiltered(ctx context.Context, pin paging.Request, fin lowerCamelSingularFilters) (pout paging.Response, results []UpperCamelSingular, err error) {
				logger.WithContext(ctx).WithField("paging", pin).WithField("filter", fin).Debugf("Retrieving paginated and filtered UpperCamelSingular records")
				pout, err = r.typedRepository.FindAll(ctx, &results, sqldb.Where(fin.Where()), sqldb.Paging(pin))
				return
			}
			`,
		[]codegen.Import{
			text.ImportContext,
			text.ImportPaging,
			text.ImportSqldb,
		})
}

func (g DomainRepositoryGenerator) createFindOneSnippet() error {
	return g.AddNewText(
		"Actions/Retrieve",
		"FindByKey",
		`
			// FindByKey retrieves a single UpperCamelSingular record by the specified lowerCamelSingularId
			func (r *lowerCamelSingularTypedRepository) FindByKey(ctx context.Context, lowerCamelSingularId uuid.UUID) (result UpperCamelSingular, err error) {
				logger.WithContext(ctx).Debugf("Retrieving UpperCamelSingular by key %q", lowerCamelSingularId.String())
				err = r.typedRepository.FindOne(ctx, &result, sqldb.And(map[string]any{
					columnUpperCamelSingularId: lowerCamelSingularId,
				}))
				if err == sqldb.ErrNotFound {
					err = repository.ErrNotFound
				}
				return
			}
			`,
		[]codegen.Import{
			text.ImportContext,
			text.ImportSqldb,
			text.ImportRepository,
		})
}

func (g DomainRepositoryGenerator) createSaveSnippet() error {
	return g.AddNewText(
		"Actions/Save",
		"Save",
		`
			// Save inserts or updates an existing record based on the supplied UpperCamelSingular lowerCamelSingularId
			func (r *lowerCamelSingularTypedRepository) Save(ctx context.Context, lowerCamelSingular UpperCamelSingular) (err error) {
				logger.WithContext(ctx).Debugf("Storing UpperCamelSingular with key %q", lowerCamelSingular.UpperCamelSingularId)
				return r.typedRepository.Upsert(ctx, lowerCamelSingular)
			}
			`,
		[]codegen.Import{
			text.ImportContext,
		})
}

func (g DomainRepositoryGenerator) createDeleteSnippet() error {
	return g.AddNewText(
		"Actions/Delete",
		"Delete",
		`
			// Delete removes the UpperCamelSingular record with the specified id
			func (r *lowerCamelSingularTypedRepository) Delete(ctx context.Context, lowerCamelSingularId uuid.UUID) (err error) {
				logger.WithContext(ctx).Debugf("Deleting UpperCamelSingular by key %q", lowerCamelSingularId.String())
				err = r.typedRepository.DeleteOne(ctx, map[string]any{
					columnUpperCamelSingularId: lowerCamelSingularId,
				})
				if err == sqldb.ErrNotFound {
					err = repository.ErrNotFound
				}
				return
			}
			`,
		[]codegen.Import{
			text.ImportContext,
			text.ImportUuid,
			text.ImportSqldb,
			text.ImportRepository,
		})
}

func (g DomainRepositoryGenerator) createContextSnippet() error {
	return g.AddNewText(
		"Context",
		"contextAccessor",
		`
			const contextKeyUpperCamelSingularRepository = contextKeyNamed("UpperCamelSingularRepository")

			// contextPetRepository returns a ContextKeyAccessor enabling dependency overrides
			// for UpperCamelSingularRepositoryApi.
			func contextUpperCamelSingularRepository() types.ContextKeyAccessor[UpperCamelSingularRepositoryApi] {
				return types.NewContextKeyAccessor[UpperCamelSingularRepositoryApi](contextKeyUpperCamelSingularRepository)
			}
			`,
		[]codegen.Import{
			text.ImportTypes,
		})
}

func (g DomainRepositoryGenerator) createConstructorSnippet() error {
	return g.AddNewText(
		"Constructor",
		"newRepository",
		`
			// newUpperCamelSingularRepository is an abstract factory, returning by default a production implementation
			// of the UpperCamelSingularRepositoryApi.
			func newUpperCamelSingularRepository(ctx context.Context) (UpperCamelSingularRepositoryApi, error) {
				repo := contextUpperCamelSingularRepository().Get(ctx)
				if repo == nil {
					typedRepository, err := sqldb.NewTypedRepository[UpperCamelSingular](ctx, tableNameUpperCamelSingular)
					if err != nil {
						return nil, err
					}
			
					repo = &lowerCamelSingularTypedRepository{
						typedRepository: typedRepository,
					}
				}
				return repo, nil
			}
			`,
		[]codegen.Import{
			text.ImportContext,
			text.ImportSqldb,
		})
}

func (g DomainRepositoryGenerator) Generate() error {
	errs := types.ErrorList{
		g.createConstantsSnippet(),
		g.createApiSnippet(),
		g.createRepositorySnippet(),
		g.createContextSnippet(),
		g.createConstructorSnippet(),
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

func (g DomainRepositoryGenerator) Filename() string {
	target := path.Join(g.Folder, "repository_lowersingular.go")
	return g.GoFile.Inflector.Inflect(target)
}

func NewDomainRepositoryGenerator(spec Spec) ComponentGenerator {
	return DomainRepositoryGenerator{
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
				Comment:   "Repository for " + generatorConfig.Domain,
				Inflector: text.NewInflector(generatorConfig.Domain),
				Sections: text.NewGoSections(
					"Constants",
					"API",
					"Repository",
					&text.Section[text.GoSnippet]{
						Name: "Actions",
						Sections: text.NewGoSections(
							"List",
							"Retrieve",
							"Save",
							"Delete"),
					},
					"Context",
					"Constructor",
				),
			},
			Package: generatorConfig.PackageName(),
		},
	}
}
