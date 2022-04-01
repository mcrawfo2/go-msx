// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

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
