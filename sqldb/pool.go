package sqldb

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"net/url"
	"strings"
	"sync"
)

var pool *ConnectionPool
var poolMtx sync.Mutex

type SqlxAction func(ctx context.Context, conn *sqlx.DB) error
type SqlAction func(ctx context.Context, conn *sql.DB) error

type ConnectionPool struct {
	cfg *Config
}

func (p *ConnectionPool) Config() *Config {
	return p.cfg
}

func (p *ConnectionPool) WithSqlxConnection(ctx context.Context, action SqlxAction) error {
	db, err := sqlx.ConnectContext(ctx, p.cfg.Driver, p.cfg.DataSourceName)
	if err != nil {
		return err
	}
	defer func() {
		err := db.Close()
		if err != nil {
			logger.WithContext(ctx).WithError(err).Error("Error closing database connection")
		}
	}()

	return action(ctx, db)
}

func (p *ConnectionPool) WithSqlConnection(ctx context.Context, action SqlAction) error {
	db, err := sql.Open(p.cfg.Driver, p.cfg.DataSourceName)
	if err != nil {
		return err
	}
	defer func() {
		err := db.Close()
		if err != nil {
			logger.WithContext(ctx).WithError(err).Error("Error closing database connection")
		}
	}()

	return action(ctx, db)
}

func CreateDatabaseForPool(ctx context.Context) error {
	logger.WithContext(ctx).Warn("Creating sql database")

	if pool == nil {
		logger.WithContext(ctx).Warn("Sqldb is disabled - skipping database creation")
		return nil
	}

	cfg := pool.Config()

	if cfg.Driver != "postgres" {
		logger.WithContext(ctx).Warn("Sqldb is not using postgres driver - skipping database creation")
		return nil
	}

	dbUrl, err := url.Parse(cfg.DataSourceName)
	if err != nil {
		return errors.Wrap(err, "Failed to parse datasource url")
	}

	if dbUrl.User == nil || dbUrl.User.Username() != "root" {
		logger.WithContext(ctx).Warn("Sqldb is not using privileged account - skipping database creation")
		return nil
	}

	dbName := strings.TrimPrefix(dbUrl.Path, "/")
	if dbName == "" {
		return errors.New("Database name not specified in datasource url")
	}

	dbUrl.Path = "/"

	db, err := sql.Open(cfg.Driver, dbUrl.String())
	if err != nil {
		return err
	}

	logger.WithContext(ctx).Infof("Ensuring database %s exists", dbName)

	_, err = db.ExecContext(ctx,
		fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName),
	)

	if err == nil {
		logger.WithContext(ctx).Infof("Sql database %s created", dbName)
	}

	return err
}

func ConfigurePool(ctx context.Context) error {
	poolMtx.Lock()
	defer poolMtx.Unlock()

	if pool != nil {
		return nil
	}

	sqlConfig, err := NewSqlConfig(ctx)
	if err != nil {
		return err
	}

	pool = &ConnectionPool{
		cfg: sqlConfig,
	}

	return nil
}
