package rest

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/mcrawfo2/go-jsonschema/pkg/codegen"
	"path"
)

type DomainServiceGeneratorV8 struct {
	Domain  string
	Folder  string
	Actions types.ComparableSlice[string]
	Spec    Spec
	*File
}

func (g DomainServiceGeneratorV8) createApiSnippet() error {
	return g.AddNewText(
		"API",
		"interface",
		`
			//go:generate mockery --name=UpperCamelSingularServiceApi --testonly --case=snake --inpackage --with-expecter
			
			// UpperCamelSingularServiceApi declares the interface for the UpperCamelSingularService.  This can be used
            // to interchange implementations, such as during testing. 
			type UpperCamelSingularServiceApi interface {
				ListUpperCamelPlural(ctx context.Context, pageReq v8.PagingSortingInputs, freq lowerCamelSingularFilterQueryInputs) (pageResp v8.PagingResponse, payload []UpperCamelSingularResponse, err error)
				GetUpperCamelSingular(ctx context.Context, id types.UUID) (UpperCamelSingularResponse, error)
				CreateUpperCamelSingular(ctx context.Context, request UpperCamelSingularCreateRequest) (UpperCamelSingularResponse, error)
				UpdateUpperCamelSingular(ctx context.Context, id types.UUID, request UpperCamelSingularUpdateRequest) (UpperCamelSingularResponse, error)
				DeleteUpperCamelSingular(ctx context.Context, id types.UUID) error
			}
		`,
		[]codegen.Import{
			importContext,
			importRestOpsV8,
			importTypes,
		})
}

func (g DomainServiceGeneratorV8) createServiceSnippet() error {
	return g.AddNewText(
		"Service",
		"implementation",
		`
			// lowerCamelSingularService is the production implementation of the UpperCamelSingularServiceApi.
			type lowerCamelSingularService struct {
				lowerCamelSingularRepository   UpperCamelSingularRepositoryApi
				lowerCamelSingularConverter    lowerCamelSingularConverter
				pagingConverter    v8.PagingConverter
				transactionManager sqldb.TransactionManager
			}
			`,
		[]codegen.Import{
			importRestOpsV8,
			importSqldb,
		})

}

func (g DomainServiceGeneratorV8) createActionListSnippet() error {
	return g.AddNewText(
		"Actions/List",
		"list",
		`
			// ListUpperCamelPlural returns a paginated series of UpperCamelSingular instances matching the supplied criteria. 
			func (s *lowerCamelSingularService) ListUpperCamelPlural(ctx context.Context, pageReq v8.PagingSortingInputs, freq lowerCamelSingularFilterQueryInputs) (pageResp v8.PagingResponse, payload []UpperCamelSingularResponse, err error) {
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
			importContext,
			importRestOpsV8,
		})
}

func (g DomainServiceGeneratorV8) createActionRetrieveSnippet() error {
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
			importContext,
			importTypes,
		})
}

func (g DomainServiceGeneratorV8) createActionCreateSnippet() error {
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
			importContext,
		})
}

func (g DomainServiceGeneratorV8) createActionUpdateSnippet() error {
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
			importContext,
			importTypes,
			importPrepared,
		})
}

func (g DomainServiceGeneratorV8) createActionDeleteSnippet() error {
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
			importContext,
			importTypes,
			importPrepared,
		})
}

func (g DomainServiceGeneratorV8) createContextSnippet() error {
	return g.AddNewText(
		"Context",
		"contextAccessor",
		`
			// contextUpperCamelSingularService returns a ContextKeyAccessor enabling dependency overrides
			// for UpperCamelSingularServiceApi.
			func contextUpperCamelSingularService() types.ContextKeyAccessor[UpperCamelSingularServiceApi] {
				return types.NewContextKeyAccessor[UpperCamelSingularServiceApi](contextKeyNamed("UpperCamelSingularService"))
			}
		`,
		[]codegen.Import{
			importTypes,
		})
}

func (g DomainServiceGeneratorV8) createConstructorSnippet() error {
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
						pagingConverter: v8.PagingConverter{
							SortByOptions: lowerCamelSingularSortByOptions,
						},
					}
				}
			
				return
			}
		`,
		[]codegen.Import{
			importContext,
			importSqldb,
		})
}

func (g DomainServiceGeneratorV8) Generate() error {
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

func (g DomainServiceGeneratorV8) Filename() string {
	target := path.Join(g.Folder, "service_lowersingular_v8.go")
	return g.File.Inflector.Inflect(target)
}

func (g DomainServiceGeneratorV8) Variables() map[string]string {
	return nil
}

func (g DomainServiceGeneratorV8) Conditions() map[string]bool {
	return nil
}

func NewDomainServiceGeneratorV8(spec Spec) ComponentGenerator {
	return DomainServiceGeneratorV8{
		Domain:  generatorConfig.Domain,
		Folder:  generatorConfig.Folder,
		Actions: generatorConfig.Actions,
		Spec:    spec,
		File: &File{
			Comment:   "V8 API Service for " + generatorConfig.Domain,
			Package:   generatorConfig.PackageName(),
			Inflector: skel.NewInflector(generatorConfig.Domain),
			Sections: NewSections(
				"API",
				"Service",
				&Section{
					Name: "Actions",
					Sections: NewSections(
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
	}
}
