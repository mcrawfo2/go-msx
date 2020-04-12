package build

import (
	"cto-github.cisco.com/NFV-BU/go-msx/exec"
	"gopkg.in/pipe.v2"
)

func init() {
	AddTarget("git-tag", "Tag the current commit", GitTag)
}

func GitTag(args []string) error {
	tag := BuildConfig.FullBuildNumber()
	return exec.ExecutePipes(
		pipe.Exec("git", "tag", "-d", tag),
		pipe.Exec("git", "push", "--delete", "origin", tag),
		pipe.Exec("git", "tag", "-a", tag),
		pipe.Exec("git", "push", "origin", tag),
	)
}
