package fs

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNewPrefixFileSystem(t *testing.T) {
	mockFileSystem := new(MockFileSystem)

	type args struct {
		fs   http.FileSystem
		root string
	}
	tests := []struct {
		name    string
		args    args
		want    http.FileSystem
		wantErr bool
	}{
		{
			name:    "Absolute",
			args:    args{
				fs:   mockFileSystem,
				root: "/root",
			},
			want:    PrefixFileSystem{
				fs:   mockFileSystem,
				root: "/root",
			},
			wantErr: false,
		},
		{
			name:    "Relative",
			args:    args{
				fs:   mockFileSystem,
				root: "root",
			},
			want:    nil,
			wantErr: true,
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPrefixFileSystem(tt.args.fs, tt.args.root)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPrefixFileSystem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPrefixFileSystem() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrefixFileSystem_Open(t1 *testing.T) {
	mockFile := new(MockFile)

	type fields struct {
		fs   http.FileSystem
		root string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    http.File
		wantErr bool
	}{
		{
			name:    "Open",
			fields:  fields{
				fs:   func() http.FileSystem {
					fs := new(MockFileSystem)
					fs.On("Open", "/root/test.json").Return(mockFile, nil)
					return fs
				}(),
				root: "/root",
			},
			args:    args{
				name: "test.json",
			},
			want:    mockFile,
			wantErr: false,
		},

		{
			name:    "Clean",
			fields:  fields{
				fs:   func() http.FileSystem {
					fs := new(MockFileSystem)
					fs.On("Open", "/test.json").Return(mockFile, nil)
					return fs
				}(),
				root: "/root",
			},
			args:    args{
				name: "../test.json",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := PrefixFileSystem{
				fs:   tt.fields.fs,
				root: tt.fields.root,
			}
			got, err := t.Open(tt.args.name)
			if (err != nil) != tt.wantErr {
				t1.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Open() got = %v, want %v", got, tt.want)
			}
		})
	}
}
