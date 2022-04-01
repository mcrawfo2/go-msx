// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package migrate

import (
	"context"
	"time"

	"cto-github.cisco.com/NFV-BU/go-msx/cassandra"
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
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

func (m *Migrator) ValidateManifest(appliedMigrations []AppliedMigration, migrations []*Migration, preUpgrade bool) error {
	logger.WithContext(m.ctx).Info("Validating previously applied migrations")

	postUpgradeVersion, err := m.manifest.PostUpgradeVersion()
	if err != nil {
		return errors.Wrap(err, "Failed to parse Post-Upgrade Version")
	}

	if postUpgradeVersion == nil {
		// No post-upgrade version set, run all migrations
		preUpgrade = false
	}

	for n, migration := range migrations {
		if n == len(appliedMigrations) {
			break
		}
		appliedMigration := appliedMigrations[n]

		if preUpgrade && !migration.Version.Lt(postUpgradeVersion) {
			logger.WithContext(m.ctx).
				WithField("version", migration.Version).
				Infof("Skipping verification of post-upgrade %s migration %s: %s",
					migration.Type,
					migration.Version,
					migration.Description)
			continue
		}

		// Ensure manifest migration matches applied migration
		if appliedMigration.InstalledRank != n+1 {
			return errors.Errorf("Incorrect installed rank: %+v", appliedMigration)
		}
		if appliedMigration.VersionRank != n+1 {
			return errors.Errorf("Incorrect version rank: %+v", appliedMigration)
		}
		if !appliedMigration.Success {
			return errors.Errorf("Failed migration recorded: %+v", appliedMigration)
		}

		// Baseline migrations excluded from the description check - they are expected to differ
		if appliedMigration.Type != MigrationTypeBaseline && appliedMigration.Description != migration.Description {
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

func (m *Migrator) ApplyMigrations(migrations []*Migration, lastAppliedMigration int, userName string, preUpgrade bool) (err error) {
	logger.WithContext(m.ctx).Info("Applying new migrations")

	postUpgradeVersion, err := m.manifest.PostUpgradeVersion()
	if err != nil {
		return errors.Wrap(err, "Failed to parse Post-Upgrade Version")
	}

	if postUpgradeVersion == nil {
		// No post-upgrade version set, run all migrations
		preUpgrade = false
	}

	for n := lastAppliedMigration; n < len(migrations); n++ {
		migration := migrations[n]

		if preUpgrade && !migration.Version.Lt(postUpgradeVersion) {
			logger.WithContext(m.ctx).
				WithField("version", migration.Version).
				Infof("Skipping post-upgrade %s migration %s: %s",
					migration.Type,
					migration.Version,
					migration.Description)
			continue
		}

		appliedMigration := AppliedMigration{
			Version:       migration.Version.String(),
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

		err = migration.Func(m.ctx, m.session)
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

func (m *Migrator) Migrate(preUpgrade bool) error {
	if err := m.versioner.CreateVersionTables(); err != nil {
		return err
	}

	appliedMigrations, err := m.versioner.GetAppliedMigrations()
	if err != nil {
		return err
	}

	applicableMigrations, err := m.getApplicableMigrations(appliedMigrations)

	if len(appliedMigrations) > 0 {
		if err = m.ValidateManifest(appliedMigrations, applicableMigrations, preUpgrade); err != nil {
			return err
		}
	}

	userName, err := config.FromContext(m.ctx).StringOr(configKeyCassandraUsername, "cassandra")
	if err != nil {
		return err
	}

	if err = m.ApplyMigrations(applicableMigrations, len(appliedMigrations), userName, preUpgrade); err != nil {
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
	logger.WithContext(ctx).Info("Executing Cassandra migration")

	preUpgrade, _ := config.FromContext(ctx).BoolOr("cli.flag.preupgrade", false)

	cassandraPool, err := cassandra.PoolFromContext(ctx)
	if err == cassandra.ErrDisabled {
		logger.WithContext(ctx).WithError(err).Warn("Skipping Cassandra migration.")
		return nil
	} else if err != nil {
		return err
	}

	return cassandraPool.WithSession(func(session *gocql.Session) error {
		return NewMigrator(ctx, session).Migrate(preUpgrade)
	})
}

// Returns the migrations that should be used for this environment.  This will only differ from the complete set
// of migrations defined by the service when the first applied migration is a "baseline" migration.
//
// A baseline migration will exist when the schema was created before the custom cassandra migration tool was adopted
// (i.e. before 3.6.0 for platform, later for at least some SPs).  A single baseline migration entry was created for
// to represent all migrations steps before the version using the cassandra migration tool.
func (m *Migrator) getApplicableMigrations(appliedMigrations []AppliedMigration) ([]*Migration, error) {
	allMigrations := m.manifest.Migrations()

	if len(appliedMigrations) == 0 || appliedMigrations[0].Type != MigrationTypeBaseline {
		return allMigrations, nil
	}

	baselineVersion := appliedMigrations[0].Version
	for n, migration := range allMigrations {
		if migration.Version.String() == baselineVersion {
			return allMigrations[n:], nil
		}
	}

	return nil, errors.Errorf("Baseline migration version doesn't match any defined migration")
}
