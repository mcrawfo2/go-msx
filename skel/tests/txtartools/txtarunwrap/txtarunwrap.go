// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.
//
// txtarwrap
// see help text for description
// this program loads everything into memory, so it's not suitable for very large files or trees

package main

import (
	"flag"
	"fmt"
	"golang.org/x/tools/txtar"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const help1 = `
%s is designed to facilitate creating testscript tests for programs that produce
trees of files as their output. 

Usage of %s:
		%s [flags] source|- destination 
where:

  * source is the source txtar file
  * destination is the directory to write the files to

note: a txtar archive contains of comments (ordinary file text) and a list of files with their contents
      The format is described here: https://godoc.org/golang.org/x/tools/txtar

Flags:
`
const (
	fromStdin = "-"
)

// Variables for command line flags
var debugOn bool   // debug output
var overwrite bool // whether to overwrite existing files

func main() {

	var err error

	flag.Usage = func() {
		_, err = fmt.Fprintf(flag.CommandLine.Output(), help1, strings.ToUpper(os.Args[0]), os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	flag.BoolVar(&overwrite, "overwrite", false, "overwrite files that already exist in dest")
	flag.BoolVar(&debugOn, "debug", false, "debug output")

	flag.Parse()

	debug("txtarunwrap: args: %s, \n", flag.Args())

	if len(flag.Args()) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	source := flag.Arg(0)
	if source == fromStdin {
		debug("reading from stdin\n")
		// read from stdin
	}

	dest := flag.Arg(1)
	if !isDir(dest) {
		_, err = fmt.Fprintf(os.Stderr, "destination must be a directory: %s\n", dest)
		os.Exit(1)

	}

	debug("reading: %s\n", source)
	_, files, err := readArchive(source)
	if err != nil {
		_, err = fmt.Fprintf(os.Stderr, "error reading archive: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		debug("no files in archive, nothing to do: %s\n", source)
		os.Exit(0)
	}

	debug("files: %d\n", len(files))

	// write everything out
	for _, f := range files {

		writeOut := false
		fullpath := filepath.Join(dest, f.Name)
		if _, err := os.Stat(fullpath); os.IsNotExist(err) {
			path := filepath.Dir(fullpath)
			err := os.MkdirAll(path, 0700)
			if err != nil {
				_, err = fmt.Fprintf(os.Stderr, "error creating directory: %v\n", err)
				os.Exit(1)
			}
			debug("created directory: %s\n", path)
			writeOut = true
		} else {
			if overwrite {
				writeOut = true
			} else {
				debug("file already exists: %s\n", fullpath)
			}
			continue
		}

		if writeOut {
			err = os.WriteFile(fullpath, f.Data, 0600)
			if err != nil {
				_, err = fmt.Fprintf(os.Stderr, "error writing file: %v\n", err)
				os.Exit(1)
			}
			debug("wrote: %s\n", fullpath)
		}

	}

	return
}

func readArchive(filename string) (comment []byte, files []txtar.File, err error) {
	var f *os.File
	if filename == "-" {
		f = os.Stdin
		err = nil
	} else {
		f, err = os.Open(filename)
	}
	if err != nil {
		return comment, files, fmt.Errorf("error reading the archive: %v\n", err)
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return comment, files, fmt.Errorf("error reading the archive: %v\n", err)
	}
	arch := txtar.Parse(data)
	_ = f.Close()
	return arch.Comment, arch.Files, err
}

func isDir(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}
	return f.IsDir()
}

func debug(format string, args ...interface{}) {
	if debugOn {
		fmt.Printf("#debug> "+format, args...)
	}
}
