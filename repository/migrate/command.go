// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package migrate

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/spf13/cobra"
)

var migrators []types.ActionFunc

func RegisterMigrator(migrator types.ActionFunc) {
	if migrator != nil {
		migrators = append(migrators, migrator)
	}
}

func Migrate(ctx context.Context, _ []string) error {
	for _, migrator := range migrators {
		err := migrator(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func CustomizeCommand(cmd *cobra.Command) {
	cmd.Flags().Bool("pre-upgrade", false, "Execute only the pre-upgrade migrations")
}
