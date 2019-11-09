package migrate

import "time"

type AppliedMigration struct {
	Version       string
	Description   string
	Script        string
	Type          MigrationType
	Checksum      int64
	ExecutionTime int
	InstalledBy   string
	InstalledOn   time.Time
	InstalledRank int
	Success       bool
	VersionRank   int
}
