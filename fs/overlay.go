package fs

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

// `overlayDir` Based on vfsgen template output
// https://github.com/shurcooL/vfsgen

type overlayDir struct {
	name    string
	modTime time.Time
	entries []os.FileInfo
	pos     int
}

func (d *overlayDir) Read([]byte) (int, error) {
	return 0, os.ErrInvalid
}

func (d *overlayDir) Close() error { return nil }

func (d *overlayDir) Stat() (os.FileInfo, error) { return d, nil }

func (d *overlayDir) Name() string { return d.name }

func (d *overlayDir) Size() int64 { return 0 }

func (d *overlayDir) Mode() os.FileMode { return 0755 | os.ModeDir }

func (d *overlayDir) ModTime() time.Time { return d.modTime }

func (d *overlayDir) IsDir() bool { return true }

func (d *overlayDir) Sys() interface{} { return nil }

func (d *overlayDir) Seek(offset int64, whence int) (int64, error) {
	if offset == 0 && whence == io.SeekStart {
		d.pos = 0
		return 0, nil
	}
	return 0, fmt.Errorf("unsupported Seek in directory %s", d.name)
}

func (d *overlayDir) Readdir(count int) ([]os.FileInfo, error) {
	if d.pos >= len(d.entries) && count > 0 {
		return nil, io.EOF
	}
	if count <= 0 || count > len(d.entries)-d.pos {
		count = len(d.entries) - d.pos
	}
	e := d.entries[d.pos : d.pos+count]
	d.pos += count
	return e, nil
}

type stat struct {
	Exists  bool
	IsDir   bool
	ModTime time.Time
}

type OverlayFileSystem struct {
	top    http.FileSystem
	bottom http.FileSystem
}

func (o OverlayFileSystem) Open(name string) (f http.File, err error) {
	topStat := o.openStat(name, o.top)
	bottomStat := o.openStat(name, o.bottom)
	useTop, useBottom := o.which(topStat, bottomStat)

	if useTop && !useBottom {
		return o.top.Open(name)
	} else if !useTop && useBottom {
		return o.bottom.Open(name)
	} else if !useTop && !useBottom {
		return nil, &os.PathError{
			Op:   "open",
			Path: name,
			Err:  errors.Wrap(os.ErrNotExist, "Object does in either top or bottom of overlay"),
		}
	}

	// Directories are merged
	mTime := topStat.ModTime
	if mTime.Before(bottomStat.ModTime) {
		mTime = bottomStat.ModTime
	}

	tEntries, tErr := o.readdir(name, o.top)
	if tErr != nil {
		return nil, tErr
	}

	bEntries, bErr := o.readdir(name, o.bottom)
	if bErr != nil {
		return nil, bErr
	}

	return &overlayDir{
		name:    path.Base(name),
		modTime: mTime,
		entries: o.merge(tEntries, bEntries),
		pos:     0,
	}, nil
}

func (o OverlayFileSystem) which(top, bottom stat) (bool, bool) {
	// Missing bottom file overridden by top
	if !bottom.Exists {
		return true, false
	}

	// Missing top file overridden by bottom
	if !top.Exists {
		return false, true
	}

	// Top file overrides bottom directory
	// Top directory overrides bottom file
	if top.IsDir != bottom.IsDir {
		return true, false
	}

	if !top.IsDir {
		// Newest file overrides oldest file
		if top.ModTime.After(bottom.ModTime) {
			return true, false
		} else {
			return false, true
		}
	} else {
		// Combine both
		return true, true
	}
}

func (o OverlayFileSystem) merge(top, bottom []os.FileInfo) []os.FileInfo {
	var files = make(map[string]os.FileInfo)

	for _, topFi := range top {
		files[topFi.Name()] = topFi
	}

	for _, bottomFi := range bottom {
		if topFi, ok := files[bottomFi.Name()]; !ok {
			files[bottomFi.Name()] = bottomFi
		} else {
			useTop, _ := o.which(
				o.fileStat(topFi),
				o.fileStat(bottomFi))
			if !useTop {
				files[bottomFi.Name()] = bottomFi
			}
		}
	}

	rfi := make([]os.FileInfo, len(files))

	i := 0
	for _, fi := range files {
		rfi[i] = fi
		i++
	}

	return rfi
}

func (o OverlayFileSystem) openStat(name string, layer http.FileSystem) stat {
	fi, err := o.stat(name, layer)
	var exists = err == nil && fi != nil
	var dir bool
	var modTime time.Time
	if fi != nil {
		dir = fi.IsDir()
		modTime = fi.ModTime()
	}

	return stat{
		Exists:  exists,
		IsDir:   dir,
		ModTime: modTime,
	}
}

func (o OverlayFileSystem) fileStat(fi os.FileInfo) stat {
	return stat{
		Exists:  true,
		IsDir:   fi.IsDir(),
		ModTime: fi.ModTime(),
	}
}

func (o OverlayFileSystem) stat(name string, layer http.FileSystem) (os.FileInfo, error) {
	f, err := layer.Open(name)
	if err != nil {
		return nil, err
	}
	return f.Stat()
}

func (o OverlayFileSystem) readdir(name string, layer http.FileSystem) ([]os.FileInfo, error) {
	f, err := layer.Open(name)
	if err != nil {
		return nil, err
	}
	return f.Readdir(-1)
}

func NewOverlayFileSystem(top, bottom http.FileSystem) http.FileSystem {
	return &OverlayFileSystem{
		top:    top,
		bottom: bottom,
	}
}
