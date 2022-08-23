// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package build

import (
	"github.com/bmatcuk/doublestar"
	"io/ioutil"
	"os"
	"path/filepath"
)

func init() {
	AddTarget("copy-book-chapters", "Copy markdown files into book folder", CopyBookChapters)
}

func CopyBookChapters(_ []string) error {
	// Find all *.md files and copy them (with directories) under "book"
	for _, sourceFile := range findBookChapters() {
		if err := copyFile(sourceFile, filepath.Join("book", sourceFile)); err != nil {
			return err
		}
	}
	return nil
}

func findBookChapters() []string {
	var fileMap = make(map[string]struct{})
	_ = filepath.Walk(".", func(p string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if isVendor, _ := doublestar.Match("vendor/**/*", p); isVendor {
			return nil
		}

		if isBook, _ := doublestar.Match("book/**/*", p); isBook {
			return nil
		}

		if isSkelTemplate, _ := doublestar.Match("skel/_templates/**/*", p); isSkelTemplate {
			return nil
		}

		if isMdFile, _ := doublestar.Match("**/*.md", p); !isMdFile {
			return nil
		}

		fileMap[p] = struct{}{}
		return nil
	})

	var results []string
	for dirName := range fileMap {
		results = append(results, dirName)
	}
	return results
}

func copyFile(src, dest string) error {
	sourceBytes, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(dest, sourceBytes, 0644)
}
