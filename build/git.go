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

	logger.Infof("Applying git build tag %q", BuildConfig.FullBuildNumber())

	return exec.ExecutePipes(
		exec.Info("Removing existing remote tag (if exists)"),
		pipe.Exec("git", "push", "origin", ":refs/tags/"+tag),

		exec.Info("Recreating local tag"),
		pipe.Exec("git", "tag", "-fa", tag, "-m", "automatic build tag"),

		exec.Info("Pushing local tag to remote repository"),
		pipe.Exec("git", "push", "origin", "refs/tags/"+tag),
	)
}
