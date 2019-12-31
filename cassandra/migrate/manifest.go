package migrate

import (
	"cto-github.cisco.com/NFV-BU/go-msx/cassandra/ddl"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
)

type MigrationType string

const (
	MigrationTypeCql      MigrationType = "CQL"
	MigrationTypeGoDriver MigrationType = "GO_DRIVER"

	configRootManifest = "migrate"
)

type Migration struct {
	Version     string
	Description string
	Script      string
	Type        MigrationType
	Func        MigrationFunc
}

type MigrationFunc func(session *gocql.Session) error

func CqlMigration(cql string) MigrationFunc {
	return func(session *gocql.Session) error {
		return session.Query(cql).Consistency(gocql.All).Exec()
	}
}

type ManifestConfig struct {
	CqlFilePath string `config:"default=migrate"`
}

type Manifest struct {
	migrations []*Migration
	cfg        *ManifestConfig
}

func (m *Manifest) Migrations() []*Migration {
	return m.migrations[:]
}

func (m *Manifest) AddCqlStringMigration(version, description, cql string) error {
	parsedVersion := types.NewVersion(version)
	if len(parsedVersion) < 3 {
		return errors.Errorf("Invalid version: %s", version)
	}

	migration := &Migration{
		Version:     version,
		Description: description,
		Script:      "cql-inline",
		Type:        MigrationTypeCql,
		Func:        CqlMigration(cql),
	}

	m.migrations = append(m.migrations, migration)

	return nil
}

func (m *Manifest) AddCqlFileMigration(version, description, filename string) error {
	cqlFilePath, err := filepath.Abs(path.Join(m.cfg.CqlFilePath, filename))
	if err != nil {
		return err
	}

	cql, err := ioutil.ReadFile(cqlFilePath)
	if err != nil {
		return err
	}

	parsedVersion := types.NewVersion(version)
	if len(parsedVersion) < 3 {
		return errors.Errorf("Invalid version: %s", version)
	}

	migration := &Migration{
		Version:     version,
		Description: description,
		Script:      script(filename),
		Type:        MigrationTypeCql,
		Func:        CqlMigration(string(cql)),
	}

	m.migrations = append(m.migrations, migration)

	return nil
}

func (m *Manifest) AddGoMigration(version, description string, fn MigrationFunc) error {
	parsedVersion := types.NewVersion(version)
	if len(parsedVersion) < 3 {
		return errors.Errorf("Invalid version: %s", version)
	}

	migration := &Migration{
		Version:     version,
		Description: description,
		Script:      types.FullFunctionName(fn),
		Type:        MigrationTypeGoDriver,
		Func:        fn,
	}

	m.migrations = append(m.migrations, migration)

	return nil
}

func (m *Manifest) AddCreateTableMigration(version string, table ddl.Table, ifNotExists bool) error {
	description := fmt.Sprintf("Create %s table", table.Name)
	stmt := new(ddl.CreateTableQueryBuilder).CreateTable(table, ifNotExists)
	return m.AddCqlStringMigration(version, description, stmt)
}

func (m *Manifest) AddCreateIndexMigration(version string, index ddl.Index, ifNotExists bool) error {
	description := fmt.Sprintf("Create %s index on %s", index.Name, index.Table)
	stmt := new(ddl.CreateIndexQueryBuilder).CreateIndex(index, ifNotExists)
	return m.AddCqlStringMigration(version, description, stmt)
}

func script(filename string) string {
	return strings.ToUpper(path.Base(filename))
}

func NewManifest(cfg *config.Config) (*Manifest, error) {
	var manifestConfig ManifestConfig
	if err := cfg.Populate(&manifestConfig, configRootManifest); err != nil {
		return nil, err
	}

	return &Manifest{
		cfg: &manifestConfig,
	}, nil
}
