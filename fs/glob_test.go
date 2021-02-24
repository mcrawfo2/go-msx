package fs

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNewGlobFileSystem(t *testing.T) {
	type args struct {
		source   http.FileSystem
		includes []string
		excludes []string
	}
	tests := []struct {
		name      string
		args      args
		wantFiles []string
		wantErr   bool
	}{
		{
			name: "NoOverlay",
			args: args{
				source:   http.Dir("testdata"),
				includes: []string{"**/*"},
				excludes: []string{"overlay/**/*"},
			},
			wantFiles: []string{
				"/file2.json",
				"/json/file1.json",
				"/json2/file1.json",
				"/json2/file3.json",
				"/text/file3.txt",
			},
			wantErr: false,
		},
		{
			name: "AllJson",
			args: args{
				source:   http.Dir("testdata"),
				includes: []string{"**/*.json"},
				excludes: nil,
			},
			wantFiles: []string{
				"/file2.json",
				"/json/file1.json",
				"/json2/file1.json",
				"/json2/file3.json",
			},
			wantErr: false,
		},
		{
			name: "SubFolderJson",
			args: args{
				source:   http.Dir("testdata"),
				includes: []string{"**/*.json"},
				excludes: []string{"/*.json"},
			},
			wantFiles: []string{
				"/json/file1.json",
				"/json2/file1.json",
				"/json2/file3.json",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGlobFileSystem(tt.args.source, tt.args.includes, tt.args.excludes)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGlobFileSystem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (tt.wantFiles == nil) != (got == nil) {
				t.Errorf("NewGlobFileSystem() got = %v, wanted %v", got, tt.wantFiles)
			} else if tt.wantFiles != nil {
				gotFiles, _ := ListFiles(got)
				if !reflect.DeepEqual(tt.wantFiles, gotFiles) {
					t.Errorf("NewGlobFileSystem() gotFiles = %v, wantFiles %v", gotFiles, tt.wantFiles)
				}
			}
		})
	}
}
