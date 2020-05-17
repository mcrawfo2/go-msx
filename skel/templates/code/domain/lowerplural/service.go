package lowerplural

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/rbac"
	"cto-github.cisco.com/NFV-BU/go-msx/repository"
	"cto-github.cisco.com/NFV-BU/go-msx/skel/templates/code/domain/api"
	"github.com/pkg/errors"
)

var (
	lowerCamelSingularErrNotFound      = errors.Wrap(repository.ErrNotFound, "Title Singular not found")
	lowerCamelSingularErrAlreadyExists = errors.Wrap(repository.ErrAlreadyExists, "Title Singular already exists")
)

type lowerCamelSingularServiceApi interface {
	ListUpperCamelPlural(ctx context.Context) ([]lowerCamelSingular, error)
	GetUpperCamelSingular(ctx context.Context, name string) (lowerCamelSingular, error)
	CreateUpperCamelSingular(ctx context.Context, request api.UpperCamelSingularCreateRequest) (lowerCamelSingular, error)
	UpdateUpperCamelSingular(ctx context.Context, name string, request api.UpperCamelSingularUpdateRequest) (lowerCamelSingular, error)
	DeleteUpperCamelSingular(ctx context.Context, name string) error
}

type lowerCamelSingularService struct {
	lowerCamelSingularRepository lowerCamelSingularRepositoryApi
	lowerCamelSingularConverter  lowerCamelSingularConverter
}

func (s *lowerCamelSingularService) ListUpperCamelPlural(ctx context.Context) ([]lowerCamelSingular, error) {
	return s.lowerCamelSingularRepository.FindAll(ctx)
}

func (s *lowerCamelSingularService) GetUpperCamelSingular(ctx context.Context, name string) (result lowerCamelSingular, err error) {
	optionalResult, err := s.lowerCamelSingularRepository.FindByKey(ctx, name)
	if err == repository.ErrNotFound {
		err = lowerCamelSingularErrNotFound
	}
	if err == nil {
		result = *optionalResult
		//#if TENANT_DOMAIN
		if err = rbac.HasTenant(ctx, result.TenantId.Bytes()); err != nil {
			return
		}
		//#endif TENANT_DOMAIN
	}

	return result, err
}

func (s *lowerCamelSingularService) CreateUpperCamelSingular(ctx context.Context, request api.UpperCamelSingularCreateRequest) (result lowerCamelSingular, err error) {
	result = s.lowerCamelSingularConverter.FromCreateRequest(request)

	//#if TENANT_DOMAIN
	if err = rbac.HasTenant(ctx, result.TenantId.Bytes()); err != nil {
		return
	}
	//#endif TENANT_DOMAIN

	_, err = s.lowerCamelSingularRepository.FindByKey(ctx, result.Name)
	if err == nil {
		err = lowerCamelSingularErrAlreadyExists
		return
	}

	err = s.lowerCamelSingularRepository.Save(ctx, result)
	return
}

func (s *lowerCamelSingularService) UpdateUpperCamelSingular(ctx context.Context, name string, request api.UpperCamelSingularUpdateRequest) (result lowerCamelSingular, err error) {
	a, err := s.lowerCamelSingularRepository.FindByKey(ctx, name)
	if err == repository.ErrNotFound {
		err = lowerCamelSingularErrNotFound
	}
	if err != nil {
		return
	}

	//#if TENANT_DOMAIN
	if err = rbac.HasTenant(ctx, a.TenantId.Bytes()); err != nil {
		return
	}
	//#endif TENANT_DOMAIN

	result = s.lowerCamelSingularConverter.FromUpdateRequest(*a, request)

	err = s.lowerCamelSingularRepository.Save(ctx, result)
	return
}

func (s *lowerCamelSingularService) DeleteUpperCamelSingular(ctx context.Context, name string) (err error) {
	//#if TENANT_DOMAIN
	a, err := s.lowerCamelSingularRepository.FindByKey(ctx, name)
	if err == repository.ErrNotFound {
		return nil
	}
	if err != nil {
		return
	}

	if err = rbac.HasTenant(ctx, a.TenantId.Bytes()); err != nil {
		return
	}
	//#endif TENANT_DOMAIN

	return s.lowerCamelSingularRepository.Delete(ctx, name)
}

func newUpperCamelSingularService(ctx context.Context) lowerCamelSingularServiceApi {
	service := lowerCamelSingularServiceFromContext(ctx)
	if service == nil {
		service = &lowerCamelSingularService{
			lowerCamelSingularRepository: newUpperCamelSingularRepository(ctx),
			lowerCamelSingularConverter:  lowerCamelSingularConverter{},
		}
	}
	return service
}