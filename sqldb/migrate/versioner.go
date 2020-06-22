package migrate

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/sqldb"
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"sort"
)

const tableMigrationName = "flyway_schema_history"

type Versioner struct {
	ctx     context.Context
	cfg     *sqldb.Config
	db      *sqlx.DB
	dialect goqu.DialectWrapper
}

func (v *Versioner) CreateVersionTables() error {
	logger.WithContext(v.ctx).Info("Creating migration history table if it does not exist")

	query := `create table if not exists flyway_schema_history
	(
		installed_rank int8 not null
			constraint "primary"
				primary key,
		version varchar(50),
		description varchar(200) not null,
		type varchar(20) not null,
		script varchar(1000) not null,
		checksum int8,
		installed_by varchar(100) not null,
		installed_on timestamp(6) default now() not null,
		execution_time int8 not null,
		success bool not null
	);
	
	create index if not exists flyway_schema_history_s_idx
		on flyway_schema_history (success);`

	if _, err := v.db.Exec(query); err != nil {
		logger.WithError(err)
		return err
	}

	return nil
}

func (v *Versioner) GetAppliedMigrations() ([]AppliedMigration, error) {
	logger.WithContext(v.ctx).Info("Retrieving migration history")

	stmt, args, err := goqu.
		From(tableMigrationName).
		ToSQL()
	if err != nil {
		return nil, err
	}

	var r []AppliedMigration
	err = v.db.SelectContext(v.ctx, &r, stmt, args...)
	if err != nil {
		return nil, err
	}

	sort.Slice(r, func(i, j int) bool {
		return r[i].InstalledRank < r[j].InstalledRank
	})

	return r, nil
}

func (v *Versioner) RecordAppliedMigration(appliedMigration AppliedMigration) error {
	logger.WithContext(v.ctx).Info("Recording new migration history entry")

	stmt, args, err := v.dialect.Insert(tableMigrationName).Rows(appliedMigration).ToSQL()
	if err != nil {
		return err
	}
	_, err = v.db.ExecContext(v.ctx, stmt, args...)
	return err
}

func NewVersioner(ctx context.Context, db *sqlx.DB) (Versioner, error) {
	cfg, err := sqldb.NewSqlConfig(ctx)
	if err != nil {
		return Versioner{}, err
	}

	return Versioner{
		ctx:     ctx,
		cfg:     cfg,
		db:      db,
		dialect: goqu.Dialect(cfg.Driver),
	}, nil
}
