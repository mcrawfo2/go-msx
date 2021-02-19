//go:generate mockery --inpackage --name=FileSystem --structname=MockFileSystem --filename mock_FileSystem.go
//go:generate mockery --inpackage --name=File --structname=MockFile --filename mock_File.go
package fs

import "net/http"

type FileSystem interface {
	http.FileSystem
}

type File interface {
	http.File
}
