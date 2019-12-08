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
