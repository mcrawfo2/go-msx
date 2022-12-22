package lowerplural

import (
	"context"
	//#if REPOSITORY_COCKROACH
	db "cto-github.cisco.com/NFV-BU/go-msx/sqldb/prepared"
	//#endif REPOSITORY_COCKROACH

	//#if TENANT_DOMAIN
	"cto-github.cisco.com/NFV-BU/go-msx/rbac"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	//#endif TENANT_DOMAIN
	"cto-github.cisco.com/NFV-BU/go-msx/repository"
	"cto-github.cisco.com/NFV-BU/go-msx/skel/_templates/code/domain/api"
	"github.com/pkg/errors"
)

var (
	lowerCamelSingularErrNotFound      = errors.Wrap(repository.ErrNotFound, "Title Singular not found")
	lowerCamelSingularErrAlreadyExists = errors.Wrap(repository.ErrAlreadyExists, "Title Singular already exists")
)

type lowerCamelSingularServiceApi interface {
	ListUpperCamelPlural(ctx context.Context,
		//#if TENANT_DOMAIN
		tenantId types.UUID,
		//#endif TENANT_DOMAIN
	) ([]api.UpperCamelSingularResponse, error)
	GetUpperCamelSingular(ctx context.Context, lowerCamelSingularId types.UUID) (api.UpperCamelSingularResponse, error)
	CreateUpperCamelSingular(ctx context.Context, request api.UpperCamelSingularCreateRequest) (api.UpperCamelSingularResponse, error)
	UpdateUpperCamelSingular(ctx context.Context, lowerCamelSingularId types.UUID, request api.UpperCamelSingularUpdateRequest) (api.UpperCamelSingularResponse, error)
	DeleteUpperCamelSingular(ctx context.Context, lowerCamelSingularId types.UUID) error
}

type lowerCamelSingularService struct {
	lowerCamelSingularRepository lowerCamelSingularRepositoryApi
	lowerCamelSingularConverter  lowerCamelSingularConverter
}

func (s *lowerCamelSingularService) ListUpperCamelPlural(ctx context.Context,
	//#if TENANT_DOMAIN
	tenantId types.UUID,
	//#endif TENANT_DOMAIN
) (body []api.UpperCamelSingularResponse, err error) {
	//#if TENANT_DOMAIN
	if err = rbac.HasTenant(ctx, tenantId); err != nil {
		return nil, err
	}
	results, err := s.lowerCamelSingularRepository.FindAllByIndexTenantId(ctx, db.ToModelUuid(tenantId))
	//#else TENANT_DOMAIN
	results, err := s.lowerCamelSingularRepository.FindAll(ctx)
	//#endif TENANT_DOMAIN
	if err == nil {
		body = s.lowerCamelSingularConverter.ToUpperCamelSingularListResponse(results)
	}
	return
}

func (s *lowerCamelSingularService) GetUpperCamelSingular(ctx context.Context, lowerCamelSingularId types.UUID) (body api.UpperCamelSingularResponse, err error) {
	optionalResult, err := s.lowerCamelSingularRepository.FindByKey(ctx, db.ToModelUuid(lowerCamelSingularId))
	if err == repository.ErrNotFound {
		err = lowerCamelSingularErrNotFound
	}
	if err == nil {
		result := *optionalResult
		//#if TENANT_DOMAIN
		if err = rbac.HasTenant(ctx, db.ToApiUuid(result.TenantId)); err != nil {
			return
		}
		//#endif TENANT_DOMAIN
		body = s.lowerCamelSingularConverter.ToUpperCamelSingularResponse(result)
	}

	return
}

func (s *lowerCamelSingularService) CreateUpperCamelSingular(ctx context.Context, request api.UpperCamelSingularCreateRequest) (body api.UpperCamelSingularResponse, err error) {
	result := s.lowerCamelSingularConverter.FromCreateRequest(request)

	//#if TENANT_DOMAIN
	if err = rbac.HasTenant(ctx, db.ToApiUuid(result.TenantId)); err != nil {
		return
	}
	//#endif TENANT_DOMAIN

	_, err = s.lowerCamelSingularRepository.FindByKey(ctx, result.UpperCamelSingularId)
	if err == nil {
		err = lowerCamelSingularErrAlreadyExists
		return
	}

	err = s.lowerCamelSingularRepository.Save(ctx, result)
	if err == nil {
		body = s.lowerCamelSingularConverter.ToUpperCamelSingularResponse(result)
	}
	return
}

func (s *lowerCamelSingularService) UpdateUpperCamelSingular(ctx context.Context, lowerCamelSingularId types.UUID, request api.UpperCamelSingularUpdateRequest) (body api.UpperCamelSingularResponse, err error) {
	a, err := s.lowerCamelSingularRepository.FindByKey(ctx, db.ToModelUuid(lowerCamelSingularId))
	if err == repository.ErrNotFound {
		err = lowerCamelSingularErrNotFound
	}
	if err != nil {
		return
	}

	//#if TENANT_DOMAIN
	if err = rbac.HasTenant(ctx, db.ToApiUuid(a.TenantId)); err != nil {
		return
	}
	//#endif TENANT_DOMAIN

	result := s.lowerCamelSingularConverter.FromUpdateRequest(*a, request)

	err = s.lowerCamelSingularRepository.Save(ctx, result)
	if err == nil {
		body = s.lowerCamelSingularConverter.ToUpperCamelSingularResponse(result)
	}
	return
}

func (s *lowerCamelSingularService) DeleteUpperCamelSingular(ctx context.Context, lowerCamelSingularId types.UUID) (err error) {
	//#if TENANT_DOMAIN
	a, err := s.lowerCamelSingularRepository.FindByKey(ctx, db.ToModelUuid(lowerCamelSingularId))
	if err == repository.ErrNotFound {
		return nil
	}
	if err != nil {
		return
	}

	if err = rbac.HasTenant(ctx, db.ToApiUuid(a.TenantId)); err != nil {
		return
	}
	//#endif TENANT_DOMAIN

	return s.lowerCamelSingularRepository.Delete(ctx, db.ToModelUuid(lowerCamelSingularId))
}

func newUpperCamelSingularService(ctx context.Context) (lowerCamelSingularServiceApi, error) {
	service := lowerCamelSingularServiceFromContext(ctx)
	if service == nil {
		lowerCamelSingularRepository, err := newUpperCamelSingularRepository(ctx)
		if err != nil {
			return nil, err
		}

		service = &lowerCamelSingularService{
			lowerCamelSingularRepository: lowerCamelSingularRepository,
			lowerCamelSingularConverter:  lowerCamelSingularConverter{},
		}
	}
	return service, nil
}
