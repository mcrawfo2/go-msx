package migrate

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/cassandra/ddl"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/resource"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	"io/ioutil"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type MigrationType string

const (
	MigrationTypeCql      MigrationType = "CQL"
	MigrationTypeGoDriver MigrationType = "GO_DRIVER"

	configRootManifest = "migrate"
)

type Migration struct {
	Version     types.Version
	Description string
	Script      string
	Type        MigrationType
	Func        MigrationContextFunc
}

type MigrationFunc func(session *gocql.Session) error
type MigrationContextFunc func(ctx context.Context, session *gocql.Session) error

func dropContext(fn MigrationFunc) MigrationContextFunc {
	return func(_ context.Context, session *gocql.Session) error {
		return fn(session)
	}
}

func CqlMigration(cql string) MigrationContextFunc {
	return func(ctx context.Context, session *gocql.Session) error {
		for _, query := range strings.Split(cql, ";") {
			query = strings.TrimSpace(query)
			if query == "" {
				continue
			}
			err := session.Query(query).
				WithContext(ctx).
				Consistency(gocql.All).
				Exec()
			if err != nil {
				return err
			}
		}

		return nil
	}
}

type ManifestConfig struct {
	CqlFilePath string `config:"default=migrate"`
	PostUpgrade string `config:"default="`
}

type Manifest struct {
	migrations []*Migration
	cfg        *ManifestConfig
}

func (m *Manifest) PostUpgradeVersion() (types.Version, error) {
	if m.cfg.PostUpgrade == "" {
		return nil, nil
	}

	return types.NewVersion(m.cfg.PostUpgrade)
}

func (m *Manifest) Migrations() []*Migration {
	return m.migrations[:]
}

func validateVersion(version types.Version) error {
	if len(version) < 3 {
		return errors.Errorf("Invalid version: %s", version)
	}
	return nil
}

func (m *Manifest) AddCqlStringMigration(version, description, cql string) error {
	parsedVersion, err := types.NewVersion(version)
	if err != nil {
		return err
	}
	if err = validateVersion(parsedVersion); err != nil {
		return err
	}

	return m.addMigration(&Migration{
		Version:     parsedVersion,
		Description: description,
		Script:      "cql-inline",
		Type:        MigrationTypeCql,
		Func:        CqlMigration(cql),
	})
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

	parsedVersion, err := types.NewVersion(version)
	if err != nil {
		return err
	}
	if err = validateVersion(parsedVersion); err != nil {
		return err
	}

	return m.addMigration(&Migration{
		Version:     parsedVersion,
		Description: description,
		Script:      script(filename),
		Type:        MigrationTypeCql,
		Func:        CqlMigration(string(cql)),
	})
}

func (m *Manifest) AddCqlResourceMigration(version, description string, res resource.Ref) error {
	cql, err := res.ReadAll()
	if err != nil {
		return err
	}

	parsedVersion, err := types.NewVersion(version)
	if err != nil {
		return err
	}
	if err = validateVersion(parsedVersion); err != nil {
		return err
	}

	return m.addMigration(&Migration{
		Version:     parsedVersion,
		Description: description,
		Script:      script(res.String()),
		Type:        MigrationTypeCql,
		Func:        CqlMigration(string(cql)),
	})
}

// Deprecated
func (m *Manifest) AddGoMigration(version, description string, fn MigrationFunc) error {
	parsedVersion, err := types.NewVersion(version)
	if err != nil {
		return err
	}
	if err = validateVersion(parsedVersion); err != nil {
		return err
	}

	return m.addMigration(&Migration{
		Version:     parsedVersion,
		Description: description,
		Script:      types.FullFunctionName(fn),
		Type:        MigrationTypeGoDriver,
		Func:        dropContext(fn),
	})
}

func (m *Manifest) AddGoContextMigration(version, description string, fn MigrationContextFunc) error {
	parsedVersion, err := types.NewVersion(version)
	if err != nil {
		return err
	}
	if err = validateVersion(parsedVersion); err != nil {
		return err
	}

	return m.addMigration(&Migration{
		Version:     parsedVersion,
		Description: description,
		Script:      types.FullFunctionName(fn),
		Type:        MigrationTypeGoDriver,
		Func:        fn,
	})
}

func (m *Manifest) addMigration(migration *Migration) error {
	for _, existingMigration := range m.migrations {
		if existingMigration.Version.Equals(migration.Version) {
			return errors.Errorf("Migration version %q already defined", migration.Version.String())
		}
	}

	m.migrations = append(m.migrations, migration)
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Version.Lt(m.migrations[j].Version)
	})

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

func (m *Manifest) AddCqlResourceMigrations(refs ...resource.Ref) error {
	regExFilename := regexp.MustCompile(`^V(\d+)_(\d+)_(\d+)_(\d+)__(.*)\.cql$`)

	for _, ref := range refs {
		fileName := path.Base(ref.String())
		matches := regExFilename.FindStringSubmatch(fileName)
		if len(matches) != 6 {
			return errors.Errorf("Invalid filename format: %q", fileName)
		}

		version := fmt.Sprintf("%s.%s.%s.%s", matches[1], matches[2], matches[3], matches[4])
		description := strings.Title(strings.ReplaceAll(matches[5], "_", " "))

		if err := m.AddCqlResourceMigration(version, description, ref); err != nil {
			return err
		}
	}

	return nil
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
