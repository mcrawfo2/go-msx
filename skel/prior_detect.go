package skel

import (
	"github.com/bmatcuk/doublestar"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

// Routines to detect and enumerate prior projects and components
// in the directory tree where Skel is running

type ProjMark int
type Score float32

// Enum of markers for project existance
const (
	ProjConfig ProjMark = iota
	ProjPom
	ProjGomod
	ProjPackage
	ProjDocker
	ProjMake
	ProjGit
)

const (
	DirFalloff = 0.75  // falloff factor for each directory level
	Small      = 0.001 // small value to ignore insignificant scores
)

var msxMatch = regexp.MustCompile(`\/go-msx.*\/`)

type ProjMarker struct {
	Matches string // glob pattern to match, from github.com/bmatcuk/doublestar
	Weight  Score
}

type ProjMarkerSet map[ProjMark]ProjMarker

// Markers we will use to score directories for containing MSX projects:
// 1. skel Project Configuration files (.skel.json)
// 2. pom.xml file
// 3. go.mod file
// 4. package.json file
// 5. Dockerfile file
// 6. Makefile file
// 7. .git directory

var ProjMarkers = ProjMarkerSet{
	ProjConfig:  {Matches: "**/.skel.json", Weight: 3},
	ProjPom:     {Matches: "**/pom.xml", Weight: 1},
	ProjGomod:   {Matches: "**/go.mod", Weight: 1},
	ProjPackage: {Matches: "**/package.json", Weight: 1},
	ProjDocker:  {Matches: "**/Dockerfile", Weight: 1},
	ProjMake:    {Matches: "**/Makefile", Weight: 1},
	ProjGit:     {Matches: "**/.git", Weight: 1},
}

type DirName string
type DirProjWeights map[DirName]Score

// FindProjects finds which subdirs of the given root likely contain MSX projects by
// comparing their scores (as returned by ScoreTree) with the given threshold
func FindProjects(root string, threshold Score) (projects []DirName, err error) {
	projects = []DirName{}
	scores, err := ScoreTree(root, ProjMarkers)
	if err != nil {
		return projects, err
	}
	for dir, score := range scores {
		if score > threshold {
			projects = append(projects, dir)
		}
	}
	return projects, nil
}

// ScoreTree examines the file tree, from the given root down, and identifies
// any directories that contain the markers defined above.
// It sums the scores for each dir's subdirs, with a falloff factor for each level,
// to provide a weighted sum to identify any likely project roots
// It also reads any *.go files to look for those that reference "/go-msx/"
func ScoreTree(root string, markers ProjMarkerSet) (scores DirProjWeights, err error) {

	scores = DirProjWeights{}

	entries, err := os.ReadDir(root)
	if err != nil {
		return scores, err
	}

	for _, entry := range entries {
		if entry.IsDir() && []rune(entry.Name())[0] != '.' {
			if s := scoreDir(filepath.Join(root, entry.Name()), markers); s > Small {
				scores[DirName(entry.Name())] = s * DirFalloff
			}
		}
	}

	return scores, nil
}

// scoreDir scores the given directory, and any subdirs recursively, against the given markers
func scoreDir(root string, markers ProjMarkerSet) (score Score) {
	entries, _ := os.ReadDir(root)
	for _, entry := range entries {
		if entry.IsDir() {
			score += scoreDir(filepath.Join(root, entry.Name()), markers)
		} else {
			for _, marker := range markers { // filename matches
				matches, _ := doublestar.PathMatch(marker.Matches, entry.Name())
				if matches {
					score += marker.Weight
				}
			}
			if filepath.Ext(entry.Name()) == ".go" { // go file
				if fileContainsMSX(filepath.Join(root, entry.Name())) {
					score += 1
				}
			}
		}
	}
	return score
}

func fileContainsMSX(filename string) bool {
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()

	buf := make([]byte, 1024) // likely to be early in the file
	n, err := file.Read(buf)
	if err != nil && err != io.EOF { // horrible error
		return false
	}
	if msxMatch.Match(buf[:n]) {
		return true
	}
	if err == io.EOF {
		return false
	}
	buf2 := make([]byte, 4096) // read bigger chunks now
	for {
		n, err := file.Read(buf2)
		if err != nil && err != io.EOF { // horrible error
			break
		}
		if msxMatch.Match(buf2[:n]) {
			return true
		}
		if err == io.EOF {
			break
		}
	}
	return false
}
