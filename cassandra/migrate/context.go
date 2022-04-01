// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package migrate

import "context"

type migrateContextKey int

const contextKeyMigrationManifest migrateContextKey = iota

func ContextWithManifest(ctx context.Context, manifest *Manifest) context.Context {
	return context.WithValue(ctx, contextKeyMigrationManifest, manifest)
}

func ManifestFromContext(ctx context.Context) *Manifest {
	manifest := ctx.Value(contextKeyMigrationManifest)
	if manifest == nil {
		return nil
	}
	return manifest.(*Manifest)
}
