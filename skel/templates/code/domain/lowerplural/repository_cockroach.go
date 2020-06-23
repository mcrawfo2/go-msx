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
	columnUpperCamelSingularName = "name"
	tableUpperCamelSingular      = "lower_snake_singular"
)

type lowerCamelSingularRepositoryApi interface {
	FindAll(context.Context) ([]lowerCamelSingular, error)
	//#if TENANT_DOMAIN
	FindAllByIndexTenantId(ctx context.Context, id uuid.UUID) ([]lowerCamelSingular, error)
	//#endif TENANT_DOMAIN
	FindByKey(context.Context, string) (*lowerCamelSingular, error)
	Save(context.Context, lowerCamelSingular) error
	Delete(context.Context, string) error
}

type lowerCamelSingularSqlRepository struct {
	sqldb.CrudRepositoryApi
}

func (r *lowerCamelSingularSqlRepository) FindAll(ctx context.Context) (results []lowerCamelSingular, err error) {
	logger.WithContext(ctx).Info("Retrieving all Title Singular records")
	err = r.CrudRepositoryApi.FindAll(ctx, &results)
	return
}

//#if TENANT_DOMAIN
func (r *lowerCamelSingularSqlRepository) FindAllByIndexTenantId(ctx context.Context, tenantId uuid.UUID) (results []lowerCamelSingular, err error) {
	logger.WithContext(ctx).Info("Retrieving all Title Singular records with tenantId %q", tenantId.String())
	err = r.CrudRepositoryApi.FindAllBy(ctx, map[string]interface{}{
		"tenant_id": tenantId,
	}, &results)
	return
}

//#endif TENANT_DOMAIN

func (r *lowerCamelSingularSqlRepository) FindByKey(ctx context.Context, name string) (result *lowerCamelSingular, err error) {
	logger.WithContext(ctx).Infof("Retrieving Title Singular by key %q", name)
	var res lowerCamelSingular
	err = r.CrudRepositoryApi.FindOneBy(ctx, map[string]interface{}{
		columnUpperCamelSingularName: name,
	}, &res)
	if err == sqldb.ErrNotFound {
		err = repository.ErrNotFound
	} else if err == nil {
		result = &res
	}
	return
}

func (r *lowerCamelSingularSqlRepository) Save(ctx context.Context, lowerCamelSingular lowerCamelSingular) (err error) {
	logger.WithContext(ctx).Infof("Storing Title Singular with key %q", lowerCamelSingular.Name)
	err = r.CrudRepositoryApi.Save(ctx, lowerCamelSingular)
	return err
}

func (r *lowerCamelSingularSqlRepository) Delete(ctx context.Context, name string) (err error) {
	logger.WithContext(ctx).Infof("Deleting Title Singular by key %q", name)
	err = r.CrudRepositoryApi.DeleteBy(ctx, map[string]interface{}{
		columnUpperCamelSingularName: name,
	})
	return
}

func newUpperCamelSingularRepository(ctx context.Context) lowerCamelSingularRepositoryApi {
	repo := lowerCamelSingularRepositoryFromContext(ctx)
	if repo == nil {
		repo = &lowerCamelSingularSqlRepository{
			CrudRepositoryApi: sqldb.
				CrudRepositoryFactoryFromContext(ctx).
				NewCrudRepository(tableUpperCamelSingular),
		}
	}
	return repo
}
