// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

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
