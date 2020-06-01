package migrate

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/cassandra/ddl"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
	gocqlxqb "github.com/scylladb/gocqlx/qb"
	"sort"
)

var tableMigrationVersion = ddl.Table{
	Name: "cassandra_migration_version",
	Columns: []ddl.Column{
		{"version", "text"},
		{"checksum", "int"},
		{"description", "text"},
		{"execution_time", "int"},
		{"installed_by", "text"},
		{"installed_on", "timestamp"},
		{"installed_rank", "int"},
		{"script", "text"},
		{"success", "boolean"},
		{"type", "text"},
		{"version_rank", "int"},
	},
	PartitionKeys: []string{"version"},
}

type Versioner struct {
	ctx     context.Context
	session *gocql.Session
}

func (v *Versioner) CreateVersionTables() error {
	logger.WithContext(v.ctx).Info("Creating migration history table if it does not exist")

	qb := new(ddl.CreateTableQueryBuilder)

	if err := v.session.Query(qb.CreateTable(tableMigrationVersion, true)).Exec(); err != nil {
		logger.WithError(err)
		return err
	}

	return nil
}

func (v *Versioner) GetAppliedMigrations() ([]AppliedMigration, error) {
	logger.WithContext(v.ctx).Info("Retrieving migration history")

	stmt, names := gocqlxqb.
		Select(tableMigrationVersion.Name).
		Columns(tableMigrationVersion.ColumnNames()...).
		ToCql()

	var r []AppliedMigration
	err := gocqlx.
		Query(v.session.Query(stmt), names).
		SelectRelease(&r)

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

	stmt, names := gocqlxqb.
		Insert(tableMigrationVersion.Name).
		Columns(tableMigrationVersion.ColumnNames()...).
		ToCql()

	return gocqlx.
		Query(v.session.Query(stmt), names).
		Consistency(gocql.All).
		BindStruct(appliedMigration).
		ExecRelease()
}

func NewVersioner(ctx context.Context, session *gocql.Session) Versioner {
	return Versioner{
		ctx:     ctx,
		session: session,
	}
}
