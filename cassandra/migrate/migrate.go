package migrate

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/cassandra"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	"time"
)

const (
	configKeyCassandraUsername = "spring.data.cassandra.username"
)

var (
	logger = log.NewLogger("msx.cassandra.migrate")
)

type Migrator struct {
	ctx       context.Context
	manifest  *Manifest
	session   *gocql.Session
	versioner Versioner
}

func (m *Migrator) ValidateManifest(appliedMigrations []AppliedMigration) error {
	logger.WithContext(m.ctx).Info("Validating previously applied migrations")

	n := 0
	var migration *Migration
	for n, migration = range m.manifest.Migrations() {
		if n == len(appliedMigrations) {
			break
		}

		// Ensure manifest migration matches applied migration
		appliedMigration := appliedMigrations[n]
		if appliedMigration.InstalledRank != n+1 {
			return errors.Errorf("Incorrect installed rank: %+v", appliedMigration)
		}
		if appliedMigration.VersionRank != n+1 {
			return errors.Errorf("Incorrect version rank: %+v", appliedMigration)
		}
		if !appliedMigration.Success {
			return errors.Errorf("Failed migration recorded: %+v", appliedMigration)
		}

		if appliedMigration.Description != migration.Description {
			return errors.Errorf("Mismatched description: %+v", appliedMigration)
		}

		logger.WithContext(m.ctx).
			WithField("version", migration.Version).
			Infof("Validated %s migration %s: %s",
				migration.Type,
				migration.Version,
				migration.Description)

	}

	return nil
}

func (m *Migrator) ApplyMigrations(lastAppliedMigration int, userName string) error {
	logger.WithContext(m.ctx).Info("Applying new migrations")

	migrations := m.manifest.Migrations()
	for n := lastAppliedMigration; n < len(migrations); n++ {
		migration := migrations[n]
		appliedMigration := AppliedMigration{
			Version:       migration.Version,
			Description:   migration.Description,
			Script:        migration.Script,
			Type:          migration.Type,
			InstalledOn:   time.Now(),
			InstalledBy:   userName,
			VersionRank:   n + 1,
			InstalledRank: n + 1,
		}

		logger.WithContext(m.ctx).
			WithField("version", migration.Version).
			Infof("Applying %s migration %s: %s",
				migration.Type,
				migration.Version,
				migration.Description)

		err := migration.Func(m.session)
		if err != nil {
			return errors.Wrap(err, "Migration failed")
		}

		appliedMigration.Success = true
		appliedMigration.ExecutionTime = int(time.Since(appliedMigration.InstalledOn) / time.Millisecond)

		if err = m.versioner.RecordAppliedMigration(appliedMigration); err != nil {
			return errors.Wrap(err, "Failed to record applied migration")
		}
	}
	return nil
}

func (m *Migrator) Migrate() error {
	if err := m.versioner.CreateVersionTables(); err != nil {
		return err
	}

	appliedMigrations, err := m.versioner.GetAppliedMigrations()
	if err != nil {
		return err
	}

	if err = m.ValidateManifest(appliedMigrations); err != nil {
		return err
	}

	userName, err := config.FromContext(m.ctx).StringOr(configKeyCassandraUsername, "cassandra")
	if err != nil {
		return err
	}

	if err = m.ApplyMigrations(len(appliedMigrations), userName); err != nil {
		return err
	}

	logger.WithContext(m.ctx).Info("Database migration completed successfully.")

	return nil
}

func NewMigrator(ctx context.Context, session *gocql.Session) *Migrator {
	return &Migrator{
		ctx:       ctx,
		manifest:  ManifestFromContext(ctx),
		session:   session,
		versioner: NewVersioner(ctx, session),
	}
}

func Migrate(ctx context.Context) error {
	cassandraPool, err := cassandra.PoolFromContext(ctx)
	if err != nil {
		return err
	}

	return cassandraPool.WithSession(func(session *gocql.Session) error {
		return NewMigrator(ctx, session).Migrate()
	})
}
