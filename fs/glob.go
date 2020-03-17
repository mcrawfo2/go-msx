package fs

import (
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/bmatcuk/doublestar"
	"github.com/shurcooL/httpfs/filter"
	"github.com/shurcooL/httpfs/vfsutil"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

func NewGlobFileSystem(source http.FileSystem, includes []string, excludes []string) (http.FileSystem, error) {
	var keepFiles = make(types.StringSet)
	// TODO: keepDirs
	err := vfsutil.WalkFiles(source, "/", func(path string, info os.FileInfo, rs io.ReadSeeker, err error) (err2 error) {
		included := false
		for _, inc := range includes {
			if included, err2 = doublestar.Match(inc, path); err2 != nil {
				return
			} else if included {
				break
			}
		}

		if !included {
			return nil
		}

		var excluded = false
		for _, exc := range excludes {
			if excluded, err2 = doublestar.Match(exc, path); err != nil {
				return
			} else if excluded {
				break
			}
		}

		if !excluded {
			keepFiles.Add(path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return filter.Keep(source, func(p string, fi os.FileInfo) bool {
		if keepFiles.Contains(p) {
			return true
		}
		if !fi.IsDir() {
			return false
		}
		if p == "/" {
			return true
		}
		keepPrefix := path.Clean(path.Join("/", p)) + "/"
		for keepFile := range keepFiles {
			if strings.HasPrefix(keepFile, keepPrefix) {
				return true
			}
		}
		return false
	}), nil
}
