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

type DomainServiceGenerator struct {
	Domain  string
	Folder  string
	Style   string
	Actions types.ComparableSlice[string]
	Spec    Spec
	*text.GoFile
}

func (g DomainServiceGenerator) createApiSnippet() error {
	return g.AddNewText(
		"API",
		"interface",
		`
			//go:generate mockery --name=UpperCamelSingularServiceApi --testonly --case=snake --inpackage --with-expecter
			
			// UpperCamelSingularServiceApi declares the interface for the UpperCamelSingularService.  This can be used
            // to interchange implementations, such as during testing. 
			type UpperCamelSingularServiceApi interface {
				ListUpperCamelPlural(ctx context.Context, pageReq ${domain.style}.PagingSortingInputs, freq lowerCamelSingularFilterQueryInputs) (pageResp ${domain.style}.PagingResponse, payload []UpperCamelSingularResponse, err error)
				GetUpperCamelSingular(ctx context.Context, id types.UUID) (UpperCamelSingularResponse, error)
				CreateUpperCamelSingular(ctx context.Context, request UpperCamelSingularCreateRequest) (UpperCamelSingularResponse, error)
				UpdateUpperCamelSingular(ctx context.Context, id types.UUID, request UpperCamelSingularUpdateRequest) (UpperCamelSingularResponse, error)
				DeleteUpperCamelSingular(ctx context.Context, id types.UUID) error
			}
		`,
		[]codegen.Import{
			text.ImportContext,
			g.importStyle(),
			text.ImportTypes,
		})
}

func (g DomainServiceGenerator) createServiceSnippet() error {
	return g.AddNewText(
		"Service",
		"implementation",
		`
			// lowerCamelSingularService is the production implementation of the UpperCamelSingularServiceApi.
			type lowerCamelSingularService struct {
				lowerCamelSingularRepository   UpperCamelSingularRepositoryApi
				lowerCamelSingularConverter    lowerCamelSingularConverter
				pagingConverter    ${domain.style}.PagingConverter
				transactionManager sqldb.TransactionManager
			}
			`,
		[]codegen.Import{
			g.importStyle(),
			text.ImportSqldb,
		})

}

func (g DomainServiceGenerator) createActionListSnippet() error {
	return g.AddNewText(
		"Actions/List",
		"list",
		`
			// ListUpperCamelPlural returns a paginated series of UpperCamelSingular instances matching the supplied criteria. 
			func (s *lowerCamelSingularService) ListUpperCamelPlural(ctx context.Context, pageReq ${domain.style}.PagingSortingInputs, freq lowerCamelSingularFilterQueryInputs) (pageResp ${domain.style}.PagingResponse, payload []UpperCamelSingularResponse, err error) {
				pin, err := s.pagingConverter.FromPagingSortingInputs(pageReq)
				if err != nil {
					return
				}
				fin := s.lowerCamelSingularConverter.FromFilterQueryInputs(freq)
			
				pout, results, err := s.lowerCamelSingularRepository.FindAllPagedFiltered(ctx, pin, fin)
				if err != nil {
					return
				}
			
				if pageResp, err = s.pagingConverter.ToPagingResponse(pout); err != nil {
					return
				}
				payload = s.lowerCamelSingularConverter.ToUpperCamelSingularListResponse(results)
				return
			}
			`,
		[]codegen.Import{
			text.ImportContext,
			g.importStyle(),
		})
}

func (g DomainServiceGenerator) createActionRetrieveSnippet() error {
	return g.AddNewText(
		"Actions/Retrieve",
		"retrieve",
		`
			// GetUpperCamelSingular returns a single UpperCamelSingular instance matching the specified key.
			func (s *lowerCamelSingularService) GetUpperCamelSingular(ctx context.Context, id types.UUID) (UpperCamelSingularResponse, error) {
				lowerCamelSingular, err := s.lowerCamelSingularRepository.FindByKey(ctx, db.ToModelUuid(id))
				if err == nil {
					return s.lowerCamelSingularConverter.ToUpperCamelSingularResponse(lowerCamelSingular), nil
				}
				return UpperCamelSingularResponse{}, err
			}
			`,
		[]codegen.Import{
			text.ImportContext,
			text.ImportTypes,
		})
}

func (g DomainServiceGenerator) createActionCreateSnippet() error {
	return g.AddNewText(
		"Actions/Create",
		"create",
		`
			// CreateUpperCamelSingular instantiates and saves a new UpperCamelSingular instance using the specified values.
			func (s *lowerCamelSingularService) CreateUpperCamelSingular(ctx context.Context, request UpperCamelSingularCreateRequest) (response UpperCamelSingularResponse, err error) {
				lowerCamelSingular := s.lowerCamelSingularConverter.FromUpperCamelSingularCreateRequest(request)

				err = s.transactionManager.WithTransaction(ctx, func(ctx context.Context) (err error) {
					return s.lowerCamelSingularRepository.Save(ctx, lowerCamelSingular)
				})

				if err == nil {
					response = s.lowerCamelSingularConverter.ToUpperCamelSingularResponse(lowerCamelSingular)
				}

				return
			}
			`,
		[]codegen.Import{
			text.ImportContext,
		})
}

func (g DomainServiceGenerator) createActionUpdateSnippet() error {
	return g.AddNewText(
		"Actions/Update",
		"update",
		`
			// UpdateUpperCamelSingular modifies and saves an existing UpperCamelSingular instance using the specified values.
			func (s *lowerCamelSingularService) UpdateUpperCamelSingular(ctx context.Context, id types.UUID, request UpperCamelSingularUpdateRequest) (response UpperCamelSingularResponse, err error) {
				var lowerCamelSingular UpperCamelSingular
			
				err = s.transactionManager.WithTransaction(ctx, func(ctx context.Context) (err error) {
					lowerCamelSingular, err = s.lowerCamelSingularRepository.FindByKey(ctx, db.ToModelUuid(id))
					if err != nil {
						return
					}
			
					lowerCamelSingular = s.lowerCamelSingularConverter.FromUpperCamelSingularUpdateRequest(lowerCamelSingular, request)
					return s.lowerCamelSingularRepository.Save(ctx, lowerCamelSingular)
				})
			
				if err == nil {
					response = s.lowerCamelSingularConverter.ToUpperCamelSingularResponse(lowerCamelSingular)
				}
			
				return
			}
			`,
		[]codegen.Import{
			text.ImportContext,
			text.ImportTypes,
			text.ImportPrepared,
		})
}

func (g DomainServiceGenerator) createActionDeleteSnippet() error {
	return g.AddNewText(
		"Actions/Delete",
		"delete",
		`
			// DeleteUpperCamelSingular destroys an existing UpperCamelSingular instance with the specified key.
			func (s *lowerCamelSingularService) DeleteUpperCamelSingular(ctx context.Context, id types.UUID) error {
				return s.transactionManager.WithTransaction(ctx, func(ctx context.Context) error {
					return s.lowerCamelSingularRepository.Delete(ctx, db.ToModelUuid(id))
				})
			}
			`,
		[]codegen.Import{
			text.ImportContext,
			text.ImportTypes,
			text.ImportPrepared,
		})
}

func (g DomainServiceGenerator) createContextSnippet() error {
	return g.AddNewText(
		"Context",
		"contextAccessor",
		`
			const contextKeyUpperCamelSingularService = contextKeyNamed("UpperCamelSingularService")
		    
			// contextUpperCamelSingularService returns a ContextKeyAccessor enabling dependency overrides
			// for UpperCamelSingularServiceApi.
			func contextUpperCamelSingularService() types.ContextKeyAccessor[UpperCamelSingularServiceApi] {
				return types.NewContextKeyAccessor[UpperCamelSingularServiceApi](contextKeyUpperCamelSingularService)
			}
		`,
		[]codegen.Import{
			text.ImportTypes,
		})
}

func (g DomainServiceGenerator) createConstructorSnippet() error {
	return g.AddNewText(
		"Constructor",
		"constructor",
		`
			// newUpperCamelSingularService is an abstract factory, returning by default a production implementation
			// of the UpperCamelSingularServiceApi.
			func newUpperCamelSingularService(ctx context.Context) (svc UpperCamelSingularServiceApi, err error) {
				svc = contextUpperCamelSingularService().Get(ctx)
				if svc == nil {
					repo, err := newUpperCamelSingularRepository(ctx)
					if err != nil {
						return nil, err
					}
			
					transactionManager, err := sqldb.NewTransactionManager(ctx)
					if err != nil {
						return nil, err
					}
			
					svc = &lowerCamelSingularService{
						lowerCamelSingularRepository:   repo,
						transactionManager: transactionManager,
						pagingConverter: ${domain.style}.PagingConverter{
							SortByOptions: lowerCamelSingularSortByOptions,
						},
					}
				}
			
				return
			}
		`,
		[]codegen.Import{
			text.ImportContext,
			text.ImportSqldb,
			g.importStyle(),
		})
}

func (g DomainServiceGenerator) Apply(options skel.RenderOptions) skel.RenderOptions {
	options.AddVariable("domain.style", g.Style)
	return options
}

func (g DomainServiceGenerator) Generate() error {
	errs := types.ErrorList{
		g.createApiSnippet(),
		g.createServiceSnippet(),
		g.createContextSnippet(),
		g.createConstructorSnippet(),
	}

	for _, operation := range g.Spec.Operations {
		if !g.Actions.Contains(operation.Action) {
			continue
		}

		var err error
		switch operation.Action {
		case ActionList:
			err = g.createActionListSnippet( /*operation*/ )
		case ActionRetrieve:
			err = g.createActionRetrieveSnippet( /*operation*/ )
		case ActionCreate:
			err = g.createActionCreateSnippet( /*operation*/ )
		case ActionUpdate:
			err = g.createActionUpdateSnippet( /*operation*/ )
		case ActionDelete:
			err = g.createActionDeleteSnippet( /*operation*/ )
		default:
			// TODO
		}

		errs = append(errs, err)
	}

	return errs.Filter()
}

func (g DomainServiceGenerator) Filename() string {
	target := path.Join(g.Folder, fmt.Sprintf("service_lowersingular_%s.go", g.Style))
	return g.GoFile.Inflector.Inflect(target)
}

func (g DomainServiceGenerator) importStyle() codegen.Import {
	switch g.Style {
	case StyleV2:
		return text.ImportRestOpsV2
	default:
		return text.ImportRestOpsV8
	}
}

func NewDomainServiceGenerator(spec Spec) ComponentGenerator {
	return DomainServiceGenerator{
		Domain:  generatorConfig.Domain,
		Folder:  generatorConfig.Folder,
		Actions: generatorConfig.Actions,
		Style:   generatorConfig.Style,
		Spec:    spec,
		GoFile: &text.GoFile{
			File: &text.File[text.GoSnippet]{
				Comment:   generatorConfig.Style + " API Service for " + generatorConfig.Domain,
				Inflector: text.NewInflector(generatorConfig.Domain),
				Sections: text.NewGoSections(
					"API",
					"Service",
					&text.Section[text.GoSnippet]{
						Name: "Actions",
						Sections: text.NewGoSections(
							"List",
							"Retrieve",
							"Create",
							"Update",
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
