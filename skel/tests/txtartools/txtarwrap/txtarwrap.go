// Copyright © 2023, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.
//
// txtarwrap
// see help text for description
// this program loads everything into memory, so it's not suitable for very large files or trees

package main

import (
	"cto-github.cisco.com/NFV-BU/go-msx/skel/tests/txtartools/clienv"
	"flag"
	"fmt"
	"github.com/bmatcuk/doublestar"
	"golang.org/x/tools/txtar"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const help1 = `
%s is designed to facilitate creating testscript tests for programs that produce
trees of files as their output. 

Usage of %s:
		%s [flags] directory|file|- [directory|file|-]... 
where:

  * [directory] is the root of a directory tree, whose files will be read and included in the 
      (stdout) output in txtar format
  * [file] will be read for enclosed files and comments (txtar style) ad combined into the output
    A named file with enclosed files will be considered a txtar archive
    A named file without enclosed files will be enclosed as a file
    Text sent to stdin (-) without enclosed files will be considered comment

note: a txtar archive contains of comments (ordinary file text) and a list of files with their contents
      The format is described here: https://godoc.org/golang.org/x/tools/txtar

Flags:
`
const help2 = `
'contents' and 'only' default to **, use -contents="" or -only="" or set their env vars to 'none' to suppress

Doublestar globs are described here: https://github.com/bmatcuk/doublestar

Output will be sent to <stdout> in the order of arguments, within the comment/file categories:
  [comment...concatenated] [files in txtar-archive format]

For testscript syntax see: https://pkg.go.dev/github.com/rogpeppe/go-internal/testscript
`

// Environment variables read by this program
const (
	EnvOnly     = "TXTAR_ONLY"
	EnvIgnore   = "TXTAR_IGNORE"
	EnvContents = "TXTAR_CONTENTS"
)

const (
	ignoreIx clienv.Ix = iota
	contentsIx
	onlyIx
)

const (
	fromStdin = "-"
	sizeUnit  = 1024
)

// Variables for command line flags
var ignoreRaw string
var contentsRaw string
var onlyRaw string
var hidden bool
var maxContents int
var prefix string
var stripPrefix string
var debugOn bool // debug output
var ignore, contents, only []string
var filterArchives bool // whether to apply filters to files from txtars as well as directories

var flag2env = clienv.Flagvars{
	ignoreIx: {Name: "ignore", Env: EnvIgnore, Default: "", RawVar: &ignoreRaw, FinVar: &ignore,
		Use: "space-separated list of globs to ignore, overrides include"},
	contentsIx: {Name: "contents", Env: EnvContents, Default: "**", RawVar: &contentsRaw, FinVar: &contents,
		Use: "space-separated list of globs to include the contents of"},
	onlyIx: {Name: "only", Env: EnvOnly, Default: "**", RawVar: &onlyRaw, FinVar: &only,
		Use: "space-separated list of globs to include"},
}

type filefilter func(fileName string) (includeIt bool)

func main() {

	var err error

	flag.Usage = func() {
		_, err = fmt.Fprintf(flag.CommandLine.Output(), help1, strings.ToUpper(os.Args[0]), os.Args[0], os.Args[0])
		flag.PrintDefaults()
		_, err = fmt.Fprintf(flag.CommandLine.Output(), help2)
	}

	flag.BoolVar(&hidden, "hidden", true, "include hidden files")
	flag.IntVar(&maxContents, "maxcontents", 1024, "max size, kB, of files to include")
	flag.BoolVar(&debugOn, "debug", false, "debug output")
	flag.StringVar(&prefix, "prefix", "", "prefix to add to output files name")
	flag.StringVar(&stripPrefix, "strip", "", "prefix to strip, if present, from output file names")
	flag.BoolVar(&filterArchives, "filter-archives", false, "apply filters to files from txtars as well as directories")

	flag2env.Register()
	flag.Parse()

	debug("txtarwrap: args: %s, \n", flag.Args())
	debug("prefix: %s, strip: %s\n", prefix, stripPrefix)

	flag2env.Fallback()
	debug(flag2env.String())

	if len(prefix) > 0 && !strings.HasSuffix(prefix, string(filepath.Separator)) {
		prefix += string(filepath.Separator)
	}

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	mainfilter := func(fileName string) (includeIt bool) {
		if !hidden && strings.HasPrefix(fileName, ".") {
			debug("%s hidden, ignored\n")
			return false
		}
		if len(ignore) > 0 && matchGlobs(fileName, ignore) {
			debug("%s ignored\n", fileName)
			return false
		}
		if len(only) > 0 && !matchGlobs(fileName, only) {
			debug("%s ignored, not in 'only'\n")
			return false
		}
		return true
	}

	var allComments []byte
	var allFiles []txtar.File

	for _, arg := range flag.Args() {

		debug("reading: %s\n", arg)

		if isDir(arg) { // read as a directory
			debug("reading directory: %s\n", arg)
			newFiles, err := scanTree(arg, mainfilter)
			if err != nil {
				_, err = fmt.Fprintf(os.Stderr, "error reading directory: %v\n", err)
				os.Exit(1)
			}
			allFiles = append(allFiles, newFiles...)
			continue
		}

		// read as a file

		named := !(arg == fromStdin)

		debug("reading regular txtar/file: %s\n", arg)
		comment, files, err := readArchive(arg)
		if err != nil {
			_, err = fmt.Fprintf(os.Stderr, "error reading archive: %v\n", err)
			os.Exit(1)
		}

		if len(files) > 0 {
			if filterArchives {
				debug("filtering archive files\n")
				for i := range files {
					if mainfilter(files[i].Name) {
						allFiles = append(allFiles, files[i])
					}
				}
			} else {
				debug("not filtering archive files\n")
				allFiles = append(allFiles, files...)
			}
		} else {
			if named {
				debug("no files in archive: %s converting to single file archive\n", arg)
				file := txtar.File{Name: arg, Data: comment}
				allFiles = append(allFiles, file)
				comment = []byte{}
			} else {
				debug("converting stdin to comment\n")
			}
		}

		if len(comment) > 0 {
			info := fmt.Sprintf("\n# ========== comments from %s ==========\n", arg)
			allComments = append(allComments, []byte(info)...)
			allComments = append(allComments, comment...)
		}

	}

	var finFiles []txtar.File
	for _, file := range allFiles {

		if len(stripPrefix) > 0 && strings.HasPrefix(file.Name, stripPrefix) {
			file.Name = file.Name[len(stripPrefix):]
		}

		if len(prefix) > 0 {
			file.Name = prefix + file.Name
		}

		debug("file: %s, size: %d\n", file.Name, len(file.Data))
		finFiles = append(finFiles, file)
	}

	// write everything out
	all := txtar.Archive{Comment: allComments, Files: finFiles}
	_, err = os.Stdout.Write(txtar.Format(&all))
	if err != nil {
		_, err = fmt.Fprintf(os.Stderr, "error writing archive: %v\n", err)
		os.Exit(1)
	}

	return
}

func matchGlobs(fileName string, globs []string) bool {
	for _, m := range globs {
		j, err := doublestar.Match(m, fileName)
		if err != nil {
			_, err = fmt.Fprintf(os.Stderr, "error matching glob: %v\n", err)
			os.Exit(1)
		}
		if j {
			return true
		}
	}
	return false
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

// Read an entire tree of files into a txtar archive
func scanTree(root string, filter filefilter) (files []txtar.File, err error) {

	files = []txtar.File{}

	var walkit fs.WalkDirFunc = func(path string, d fs.DirEntry, err error) error {

		if !hidden && strings.HasPrefix(d.Name(), ".") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !d.IsDir() {

			if !hidden && strings.HasPrefix(d.Name(), ".") {
				debug("%s hidden, ignored\n", path)
				return nil
			}

			addIt := filter(path)
			if !addIt {
				return nil
			}

			loadContents := matchGlobs(path, contents)

			var data []byte
			if loadContents {
				f, err2 := os.Open(path)
				if err2 != nil {
					return err
				}
				st, err3 := f.Stat()
				if err3 != nil {
					return err2
				}
				if st.Size() <= int64(maxContents*sizeUnit) {
					data, err = io.ReadAll(f)
					if err != nil {
						return err
					}
				} else {
					debug("%s too big, ignored\n", path)
					return nil
				}
				_ = f.Close()
			}

			files = append(files, txtar.File{Name: path, Data: data})

		}

		return nil
	}

	err = filepath.WalkDir(root, walkit)
	return files, err

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
