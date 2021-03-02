package fs

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestNewOverlayFileSystem(t *testing.T) {
	top := new(MockFileSystem)
	bottom := new(MockFileSystem)
	want := &OverlayFileSystem{
		top:    top,
		bottom: bottom,
	}

	if got := NewOverlayFileSystem(top, bottom); !reflect.DeepEqual(got, want) {
		t.Errorf("NewOverlayFileSystem() = %v, want %v", got, want)
	}
}

func TestOverlayFileSystem_Open_File(t *testing.T) {
	top := http.Dir("testdata/json")
	bottom := http.Dir("testdata/json2")

	type args struct {
		name string
	}
	tests := []struct {
		name     string
		args     args
		wantStat os.FileInfo
		wantErr  bool
	}{
		{
			name: "Top",
			args: args{
				name: "file1.json",
			},
			wantStat: func() os.FileInfo {
				f, _ := top.Open("file1.json")
				s, _ := f.Stat()
				return s
			}(),
			wantErr: false,
		},
		{
			name: "Bottom",
			args: args{
				name: "file3.json",
			},
			wantStat: func() os.FileInfo {
				f, _ := bottom.Open("file3.json")
				s, _ := f.Stat()
				return s
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := NewOverlayFileSystem(top, bottom)
			gotF, err := o.Open(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (gotF != nil) != (tt.wantStat != nil) {
				t.Errorf("Open() gotF = %v, wantStat %v", gotF, tt.wantStat)
			}

			if tt.wantStat != nil && gotF != nil {
				gotStat, _ := gotF.Stat()
				if !reflect.DeepEqual(gotStat, tt.wantStat) {
				}
			}
		})
	}
}

func TestOverlayFileSystem_Open_Dir(t *testing.T) {
	top := http.Dir("testdata/overlay/top")
	bottom := http.Dir("testdata/overlay/bottom")

	type args struct {
		name string
	}
	tests := []struct {
		name          string
		args          args
		wantFileNames []string
		wantErr       bool
	}{
		{
			name: "Deep",
			args: args{
				name: "deep",
			},
			wantFileNames: []string{
				"file1.txt",
				"file2.txt",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := NewOverlayFileSystem(top, bottom)
			gotF, err := o.Open(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (gotF != nil) != (tt.wantFileNames != nil) {
				t.Errorf("Open() gotF = %v, wantStat %v", gotF, tt.wantFileNames)
			}

			if tt.wantFileNames != nil && gotF != nil {
				gotEntries, _ := gotF.Readdir(0)
				var gotFileNames []string
				for _, gotEntry := range gotEntries {
					gotFileNames = append(gotFileNames, gotEntry.Name())
				}
				if !reflect.DeepEqual(gotFileNames, tt.wantFileNames) {
					t.Errorf("Open() gotFileNames = %v, wantFileNames %v", gotFileNames, tt.wantFileNames)
				}
			}
		})
	}
}

func newOverlayDir(t *testing.T) http.File {
	top := http.Dir("testdata/overlay/top")
	bottom := http.Dir("testdata/overlay/bottom")
	o := NewOverlayFileSystem(top, bottom)

	f, err := o.Open("deep")
	assert.NoError(t, err)

	return f
}

func Test_overlayDir_Readdir(t *testing.T) {
	f := newOverlayDir(t)

	files, err := f.Readdir(1)
	assert.NoError(t, err)
	assert.Equal(t, "file1.txt", files[0].Name())

	pos, err := f.Seek(0, io.SeekStart)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), pos)

	files, err = f.Readdir(1)
	assert.NoError(t, err)
	assert.Equal(t, "file1.txt", files[0].Name())

	pos, err = f.Seek(1, io.SeekStart)
	assert.Error(t, err)
	pos, err = f.Seek(0, io.SeekEnd)
	assert.Error(t, err)

}

func Test_overlayDir_Thunks(t *testing.T) {
	f := newOverlayDir(t)

	count, err := f.Read(nil)
	assert.Equal(t, 0, count)
	assert.Error(t, err)

	err = f.Close()
	assert.NoError(t, err)

	s, err := f.Stat()
	assert.Equal(t, f, s)
	assert.NoError(t, err)

	name := s.Name()
	assert.Equal(t, "deep", name)

	size := s.Size()
	assert.Equal(t, int64(0), size)

	mode := s.Mode()
	assert.Equal(t, os.FileMode(0755|os.ModeDir), mode)

	modTime := s.ModTime()
	assert.NotEqual(t, time.Unix(0, 0), modTime)

	isDir := s.IsDir()
	assert.True(t, isDir)

	sys := s.Sys()
	assert.Nil(t, sys)
}
