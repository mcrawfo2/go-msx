package migrate

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterMigrator(t *testing.T) {
	type args struct {
		migrator types.ActionFunc
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Nil",
			args: args{migrator: nil},
		},
		{
			name: "NotNil",
			args: args{migrator: func(ctx context.Context) error {
				return nil
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			migrators = nil
			RegisterMigrator(tt.args.migrator)
			if tt.args.migrator != nil {
				assert.Len(t, migrators, 1)
			} else {
				assert.Len(t, migrators, 0)
			}
		})
	}
}

func TestMigrate(t *testing.T) {
	var count = 0
	migratorSuccess := func(_ context.Context) error { count++; return nil }
	migratorError := func(_ context.Context) error { count++; return errors.New("error") }

	tests := []struct {
		name      string
		migrators []types.ActionFunc
		wantCount int
		wantErr   bool
	}{
		{
			name: "SingleSuccess",
			migrators: []types.ActionFunc{
				migratorSuccess,
			},
			wantCount: 1,
		},
		{
			name: "SingleError",
			migrators: []types.ActionFunc{
				migratorError,
			},
			wantCount: 1,
			wantErr:   true,
		},
		{
			name: "MultipleSuccess",
			migrators: []types.ActionFunc{
				migratorSuccess,
				migratorSuccess,
				migratorSuccess,
			},
			wantCount: 3,
		},
		{
			name: "MultipleError",
			migrators: []types.ActionFunc{
				migratorSuccess,
				migratorError,
				migratorSuccess,
			},
			wantCount: 2,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			migrators = nil
			for _, migrator := range tt.migrators {
				RegisterMigrator(migrator)
			}

			count = 0
			if err := Migrate(context.Background(), nil); (err != nil) != tt.wantErr {
				t.Errorf("Migrate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if count != tt.wantCount {
				t.Errorf("Migrate() count = %v, wantCount %v", count, tt.wantCount)
			}
		})
	}
}
