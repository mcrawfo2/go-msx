package migrate

import (
	"bufio"
	"bytes"
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/resource"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"hash/crc32"
	"path"
	"regexp"
	"sort"
	"strings"
)

type MigrationType string

const (
	MigrationTypeSql      MigrationType = "SQL"
	MigrationTypeGoDriver MigrationType = "GO_DRIVER"

	configRootManifest = "migrate"
)

type Migration struct {
	Version     types.Version
	Description string
	Script      string
	Checksum    *int
	Type        MigrationType
	Func        MigrationFunc
}

type MigrationFunc func(ctx context.Context, db *sqlx.DB) error

func SqlMigration(stmts string) MigrationFunc {
	return func(ctx context.Context, db *sqlx.DB) error {
		for _, query := range strings.Split(stmts, ";") {
			query = strings.TrimSpace(query)
			if query == "" {
				continue
			}
			_, err := db.Query(query)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

type ManifestConfig struct {
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

func (m *Manifest) AddSqlStringMigration(version, description, stmts string) error {
	parsedVersion, err := types.NewVersion(version)
	if err != nil {
		return err
	}
	if len(parsedVersion) < 3 {
		return errors.Errorf("Invalid version: %s", version)
	}

	return m.addMigration(&Migration{
		Version:     parsedVersion,
		Description: description,
		Script:      "sql-inline",
		Checksum:    checksum([]byte(stmts)),
		Type:        MigrationTypeSql,
		Func:        SqlMigration(stmts),
	})
}

func (m *Manifest) AddSqlResourceMigration(version, description string, res resource.Ref) error {
	stmts, err := res.ReadAll()
	if err != nil {
		return errors.Wrap(err, "Failed to read resource")
	}

	parsedVersion, err := types.NewVersion(version)
	if err != nil {
		return err
	}
	if len(parsedVersion) < 1 {
		return errors.Errorf("Invalid version: %s", version)
	}

	return m.addMigration(&Migration{
		Version:     parsedVersion,
		Description: description,
		Script:      script(res.String()),
		Checksum:    checksum(stmts),
		Type:        MigrationTypeSql,
		Func:        SqlMigration(string(stmts)),
	})
}

func (m *Manifest) AddSqlResourceMigrations(refs ...resource.Ref) error {
	regExSuffix := regexp.MustCompile(`V([\d_]+)__(.*)\.sql$`)

	for _, ref := range refs {
		fileName := path.Base(ref.String())
		fileSuffixMatch := regExSuffix.FindStringSubmatch(fileName)
		if fileSuffixMatch[0] == "" {
			return errors.Errorf("Invalid filename format: %q", fileName)
		}

		versionParts := strings.Split(fileSuffixMatch[1], "_")
		version := strings.Join(versionParts, ".")

		description := strings.ReplaceAll(fileSuffixMatch[2], "_", " ")

		if err := m.AddSqlResourceMigration(version, description, ref); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manifest) AddGoMigration(version, description string, fn MigrationFunc) error {
	parsedVersion, err := types.NewVersion(version)
	if err != nil {
		return err
	}
	if len(parsedVersion) < 1 {
		return errors.Errorf("Invalid version: %s", version)
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

func checksum(data []byte) *int {
	checksum := crc32.NewIEEE()
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		_, _ = checksum.Write(scanner.Bytes())
	}
	result := int(int32(checksum.Sum32()))
	return &result
}

func script(filename string) string {
	return path.Base(filename)
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
