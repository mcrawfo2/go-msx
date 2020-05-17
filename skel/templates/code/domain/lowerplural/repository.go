package lowerplural

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/cassandra"
	"cto-github.cisco.com/NFV-BU/go-msx/cassandra/ddl"
	"cto-github.cisco.com/NFV-BU/go-msx/repository"
)

const (
	columnUpperCamelSingularName = "name"
)

var tableUpperCamelSingular = ddl.Table{
	Name: "lower_snake_singular",
	Columns: []ddl.Column{
		{columnUpperCamelSingularName, ddl.DataTypeText},
		{"data", ddl.DataTypeText},
	},
	PartitionKeys: []string{columnUpperCamelSingularName},
}

type lowerCamelSingularRepositoryApi interface {
	FindAll(context.Context) ([]lowerCamelSingular, error)
	FindByKey(context.Context, string) (*lowerCamelSingular, error)
	Save(context.Context, lowerCamelSingular) error
	Delete(context.Context, string) error
}

type lowerCamelSingularCassandraRepository struct {
	cassandra.CrudRepositoryApi
}

func (r *lowerCamelSingularCassandraRepository) FindAll(ctx context.Context) (results []lowerCamelSingular, err error) {
	logger.WithContext(ctx).Info("Retrieving all Title Singular records")
	err = r.CrudRepositoryApi.FindAll(ctx, &results)
	return
}

func (r *lowerCamelSingularCassandraRepository) FindByKey(ctx context.Context, themeName string) (result *lowerCamelSingular, err error) {
	logger.WithContext(ctx).Infof("Retrieving Title Singular by key %q", themeName)
	err = r.CrudRepositoryApi.FindOneBy(ctx, map[string]interface{}{
		columnUpperCamelSingularName: themeName,
	}, &result)
	if err == cassandra.ErrNotFound {
		err = repository.ErrNotFound
	}
	return
}

func (r *lowerCamelSingularCassandraRepository) Save(ctx context.Context, theme lowerCamelSingular) (err error) {
	logger.WithContext(ctx).Infof("Storing Title Singular with key %q", theme.Name)
	err = r.CrudRepositoryApi.Save(ctx, theme)
	return err
}

func (r *lowerCamelSingularCassandraRepository) Delete(ctx context.Context, themeName string) (err error) {
	logger.WithContext(ctx).Infof("Deleting Title Singular by key %q", themeName)
	err = r.CrudRepositoryApi.DeleteBy(ctx, map[string]interface{}{
		columnUpperCamelSingularName: themeName,
	})
	return
}

func newUpperCamelSingularRepository(ctx context.Context) lowerCamelSingularRepositoryApi {
	repo := lowerCamelSingularRepositoryFromContext(ctx)
	if repo == nil {
		repo = &lowerCamelSingularCassandraRepository{
			CrudRepositoryApi: cassandra.
				CrudRepositoryFactoryFromContext(ctx).
				NewCrudRepository(tableUpperCamelSingular),
		}
	}
	return repo
}
