package lowerplural

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/repository"
	"cto-github.cisco.com/NFV-BU/go-msx/sqldb"

	//#if TENANT_DOMAIN
	"github.com/google/uuid"
	//#endif TENANT_DOMAIN
)

const (
	columnUpperCamelSingularId  = "lower_snake_singular_id"
	tableNameUpperCamelSingular = "lower_snake_singular"
)

type lowerCamelSingularRepositoryApi interface {
	FindAll(ctx context.Context) (results []lowerCamelSingular, err error)
	//#if TENANT_DOMAIN
	FindAllByIndexTenantId(ctx context.Context, lowerCamelSingularId uuid.UUID) (results []lowerCamelSingular, err error)
	//#endif TENANT_DOMAIN
	FindByKey(ctx context.Context, lowerCamelSingularId uuid.UUID) (optionalResult *lowerCamelSingular, err error)
	Save(ctx context.Context, lowerCamelSingular lowerCamelSingular) (err error)
	Delete(ctx context.Context, lowerCamelSingularId uuid.UUID) (err error)
}

type lowerCamelSingularSqlRepository struct {
	sqldb.CrudRepositoryApi
}

func (r *lowerCamelSingularSqlRepository) FindAll(ctx context.Context) (results []lowerCamelSingular, err error) {
	logger.WithContext(ctx).Debugf("Retrieving all Title Singular records")
	err = r.CrudRepositoryApi.FindAll(ctx, &results)
	return
}

//#if TENANT_DOMAIN
func (r *lowerCamelSingularSqlRepository) FindAllByIndexTenantId(ctx context.Context, tenantId uuid.UUID) (results []lowerCamelSingular, err error) {
	logger.WithContext(ctx).Debugf("Retrieving all Title Singular records with tenantId %q", tenantId.String())
	err = r.CrudRepositoryApi.FindAllBy(ctx, map[string]interface{}{
		"tenant_id": tenantId,
	}, &results)
	return
}

//#endif TENANT_DOMAIN

func (r *lowerCamelSingularSqlRepository) FindByKey(ctx context.Context, lowerCamelSingularId uuid.UUID) (result *lowerCamelSingular, err error) {
	logger.WithContext(ctx).Debugf("Retrieving Title Singular by key %q", lowerCamelSingularId.String())
	var res lowerCamelSingular
	err = r.CrudRepositoryApi.FindOneBy(ctx, map[string]interface{}{
		columnUpperCamelSingularId: lowerCamelSingularId,
	}, &res)
	if err == sqldb.ErrNotFound {
		err = repository.ErrNotFound
	} else if err == nil {
		result = &res
	}
	return
}

func (r *lowerCamelSingularSqlRepository) Save(ctx context.Context, lowerCamelSingular lowerCamelSingular) (err error) {
	logger.WithContext(ctx).Debugf("Storing Title Singular with key %q", lowerCamelSingular.UpperCamelSingularId.String())
	err = r.CrudRepositoryApi.Save(ctx, lowerCamelSingular)
	return err
}

func (r *lowerCamelSingularSqlRepository) Delete(ctx context.Context, lowerCamelSingularId uuid.UUID) (err error) {
	logger.WithContext(ctx).Debugf("Deleting Title Singular by key %q", lowerCamelSingularId.String())
	err = r.CrudRepositoryApi.DeleteBy(ctx, map[string]interface{}{
		columnUpperCamelSingularId: lowerCamelSingularId,
	})
	return
}

func newUpperCamelSingularRepository(ctx context.Context) (lowerCamelSingularRepositoryApi, error) {
	repo := lowerCamelSingularRepositoryFromContext(ctx)
	if repo == nil {
		repo = &lowerCamelSingularSqlRepository{
			CrudRepositoryApi: sqldb.
				CrudRepositoryFactoryFromContext(ctx).
				NewCrudRepository(tableNameUpperCamelSingular),
		}
	}
	return repo, nil
}
