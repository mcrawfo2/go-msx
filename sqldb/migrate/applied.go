package migrate

import "time"

type AppliedMigration struct {
	Version       string        `db:"version"`
	Description   string        `db:"description"`
	Script        string        `db:"script"`
	Type          MigrationType `db:"type"`
	Checksum      *int          `db:"checksum"`
	ExecutionTime int           `db:"execution_time"`
	InstalledBy   string        `db:"installed_by"`
	InstalledOn   time.Time     `db:"installed_on"`
	InstalledRank int           `db:"installed_rank"`
	Success       bool          `db:"success"`
}
