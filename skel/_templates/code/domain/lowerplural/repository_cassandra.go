package lowerplural

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/cassandra"
	"cto-github.cisco.com/NFV-BU/go-msx/cassandra/ddl"
	"cto-github.cisco.com/NFV-BU/go-msx/repository"
	//#if TENANT_DOMAIN
	"github.com/gocql/gocql"
	//#endif TENANT_DOMAIN
)

const (
	columnUpperCamelSingularId  = "lower_snake_singular_id"
	tableNameUpperCamelSingular = "lower_snake_singular"
)

var tableUpperCamelSingular = ddl.Table{
	Name: tableNameUpperCamelSingular,
	Columns: []ddl.Column{
		{columnUpperCamelSingularId, ddl.DataTypeUuid},
		//#if TENANT_DOMAIN
		{"tenant_id", ddl.DataTypeUuid},
		//#endif TENANT_DOMAIN
		{"data", ddl.DataTypeText},
	},
	PartitionKeys: []string{columnUpperCamelSingularId},
}

type lowerCamelSingularRepositoryApi interface {
	FindAll(context.Context) ([]lowerCamelSingular, error)
	//#if TENANT_DOMAIN
	FindAllByIndexTenantId(ctx context.Context, id gocql.UUID) ([]lowerCamelSingular, error)
	//#endif TENANT_DOMAIN
	FindByKey(context.Context, gocql.UUID) (*lowerCamelSingular, error)
	Save(context.Context, lowerCamelSingular) error
	Delete(context.Context, gocql.UUID) error
}

type lowerCamelSingularCassandraRepository struct {
	cassandra.CrudRepositoryApi
}

func (r *lowerCamelSingularCassandraRepository) FindAll(ctx context.Context) (results []lowerCamelSingular, err error) {
	logger.WithContext(ctx).Info("Retrieving all Title Singular records")
	err = r.CrudRepositoryApi.FindAll(ctx, &results)
	return
}

//#if TENANT_DOMAIN
func (r *lowerCamelSingularCassandraRepository) FindAllByIndexTenantId(ctx context.Context, tenantId gocql.UUID) (results []lowerCamelSingular, err error) {
	logger.WithContext(ctx).Debugf("Retrieving all Title Singular records with tenantId %q", tenantId.String())
	err = r.CrudRepositoryApi.FindAllBy(ctx, map[string]interface{}{
		"tenant_id": tenantId,
	}, &results)
	return
}

//#endif TENANT_DOMAIN

func (r *lowerCamelSingularCassandraRepository) FindByKey(ctx context.Context, lowerCamelSingularId gocql.UUID) (result *lowerCamelSingular, err error) {
	logger.WithContext(ctx).Debugf("Retrieving Title Singular by key %q", lowerCamelSingularId.String())
	var res lowerCamelSingular
	err = r.CrudRepositoryApi.FindOneBy(ctx, map[string]interface{}{
		columnUpperCamelSingularId: lowerCamelSingularId,
	}, &res)
	if err == cassandra.ErrNotFound {
		err = repository.ErrNotFound
	} else if err == nil {
		result = &res
	}
	return
}

func (r *lowerCamelSingularCassandraRepository) Save(ctx context.Context, lowerCamelSingular lowerCamelSingular) (err error) {
	logger.WithContext(ctx).Debugf("Storing Title Singular with key %q", lowerCamelSingular.UpperCamelSingularId.String())
	err = r.CrudRepositoryApi.Save(ctx, lowerCamelSingular)
	return err
}

func (r *lowerCamelSingularCassandraRepository) Delete(ctx context.Context, lowerCamelSingularId gocql.UUID) (err error) {
	logger.WithContext(ctx).Debugf("Deleting Title Singular by key %q", lowerCamelSingularId.String())
	err = r.CrudRepositoryApi.DeleteBy(ctx, map[string]interface{}{
		columnUpperCamelSingularId: lowerCamelSingularId,
	})
	return
}

func newUpperCamelSingularRepository(ctx context.Context) (lowerCamelSingularRepositoryApi, error) {
	repo := lowerCamelSingularRepositoryFromContext(ctx)
	if repo == nil {
		repo = &lowerCamelSingularCassandraRepository{
			CrudRepositoryApi: cassandra.
				CrudRepositoryFactoryFromContext(ctx).
				NewCrudRepository(tableUpperCamelSingular),
		}
	}
	return repo, nil
}
