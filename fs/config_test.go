// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package fs

import (
	"cto-github.cisco.com/NFV-BU/go-msx/config"
	"cto-github.cisco.com/NFV-BU/go-msx/testhelpers/configtest"
	"reflect"
	"testing"
)

func TestNewFileSystemConfig(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *FileSystemConfig
		wantErr bool
	}{
		{
			name: "StructDefaults",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"spring.application.name": "TestNewFileSystemConfig",
				}),
			},
			want: &FileSystemConfig{
				Root:      "/",
				Resources: "/var/lib/TestNewFileSystemConfig",
				Configs:   "/etc/TestNewFileSystemConfig",
				Binaries:  "/usr/bin",
				Sources:   "",
				Mode:      "detect",
			},
			wantErr: false,
		},
		{
			name: "WithSources",
			args: args{
				cfg: configtest.NewInMemoryConfig(map[string]string{
					"spring.application.name": "TestNewFileSystemConfig",
					"fs.sources":              "/home/ubuntu/go-msx",
					"fs.mode":                 "release",
				}),
			},
			want: &FileSystemConfig{
				Root:      "/",
				Resources: "/var/lib/TestNewFileSystemConfig",
				Configs:   "/etc/TestNewFileSystemConfig",
				Binaries:  "/usr/bin",
				Sources:   "/home/ubuntu/go-msx",
				Mode:      "release",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFileSystemConfig(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFileSystemConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFileSystemConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
