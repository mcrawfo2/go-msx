package migrate

import (
	"context"
	"github.com/jmoiron/sqlx"
	"testing"
	"time"
)

func TestMigrator_ValidateMigration(t *testing.T) {
	type fields struct {
		ctx       context.Context
		manifest  *Manifest
		db        *sqlx.DB
		versioner Versioner
	}
	type args struct {
		n                int
		migration        Migration
		appliedMigration AppliedMigration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "default_pass",
			fields: fields{
				ctx:       context.Background(),
				manifest:  &Manifest{
					migrations: []*Migration{},
					cfg:        nil,
				},
				db:        nil,
				versioner: Versioner{},
			},
			args: args{
				n:                0,
				migration:        Migration{
					Version:     nil,
					Description: "test1",
					Script:      "",
					Checksum:    nil,
					Type:        "",
					Func:        nil,
				},
				appliedMigration: AppliedMigration{
					Version:       "",
					Description:   "test1",
					Script:        "",
					Type:          "",
					Checksum:      nil,
					ExecutionTime: 0,
					InstalledBy:   "",
					InstalledOn:   time.Time{},
					InstalledRank: 1,
					Success:       true,
				},
			},
			wantErr: false,
		},
		{
			name: "checksum_fail",
			fields: fields{
				ctx:       context.Background(),
				manifest:  &Manifest{
					migrations: []*Migration{},
					cfg:        nil,
				},
				db:        nil,
				versioner: Versioner{},
			},
			args: args{
				n:                0,
				migration:        Migration{
					Version:     nil,
					Description: "test1",
					Script:      "",
					Checksum:    func()*int{i:=-999;return &i}(),
					Type:        "",
					Func:        nil,
				},
				appliedMigration: AppliedMigration{
					Version:       "",
					Description:   "test1",
					Script:        "",
					Type:          "",
					Checksum:      nil,
					ExecutionTime: 0,
					InstalledBy:   "",
					InstalledOn:   time.Time{},
					InstalledRank: 1,
					Success:       true,
				},
			},
			wantErr: true,
		},
		{
			name: "description_fail",
			fields: fields{
				ctx:       context.Background(),
				manifest:  &Manifest{
					migrations: []*Migration{},
					cfg:        nil,
				},
				db:        nil,
				versioner: Versioner{},
			},
			args: args{
				n:                0,
				migration:        Migration{
					Version:     nil,
					Description: "test1",
					Script:      "",
					Checksum:    func()*int{i:=-999;return &i}(),
					Type:        "",
					Func:        nil,
				},
				appliedMigration: AppliedMigration{
					Version:       "",
					Description:   "test2",
					Script:        "",
					Type:          "",
					Checksum:      func()*int{i:=-999;return &i}(),
					ExecutionTime: 0,
					InstalledBy:   "",
					InstalledOn:   time.Time{},
					InstalledRank: 1,
					Success:       true,
				},
			},
			wantErr: true,
		},
		{
			name: "skip_checksum_and_description_check",
			fields: fields{
				ctx:       context.Background(),
				manifest:  &Manifest{
					migrations: []*Migration{},
					cfg:        nil,
				},
				db:        nil,
				versioner: Versioner{},
			},
			args: args{
				n:                0,
				migration:        Migration{
					Version:     nil,
					Description: "test1",
					Script:      "",
					Checksum:    func()*int{i:=0;return &i}(),
					Type:        "",
					Func:        nil,
				},
				appliedMigration: AppliedMigration{
					Version:       "",
					Description:   "test2",
					Script:        "",
					Type:          "",
					Checksum:      func()*int{i:=-999;return &i}(),
					ExecutionTime: 0,
					InstalledBy:   "",
					InstalledOn:   time.Time{},
					InstalledRank: 1,
					Success:       true,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Migrator{
				ctx:       tt.fields.ctx,
				manifest:  tt.fields.manifest,
				db:        tt.fields.db,
				versioner: tt.fields.versioner,
			}
			if err := m.ValidateMigration(tt.args.n, tt.args.migration, tt.args.appliedMigration); (err != nil) != tt.wantErr {
				t.Errorf("ValidateMigration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}