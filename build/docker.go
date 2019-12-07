package build

import (
	"cto-github.cisco.com/NFV-BU/go-msx/exec"
	"fmt"
)

func init() {
	AddTarget("docker-build", "Build the target docker image", DockerBuild)
	AddTarget("docker-push", "Push the target docker image to the upstream repository", DockerPush)
}

func DockerBuild(args []string) error {
	return exec.MustExecute("docker", "build",
		"-t", dockerImageName(),
		"-f", "docker/Dockerfile",
		"--force-rm",
		"--no-cache",
		".")
}

func DockerPush(args []string) error {
	if BuildConfig.Docker.Username != "" && BuildConfig.Docker.Password != "" {
		err := exec.MustExecute("docker", "login",
			"-u", BuildConfig.Docker.Username,
			"-p", BuildConfig.Docker.Password)
		if err != nil {
			return err
		}
	}

	return exec.MustExecute("docker", "push", dockerImageName())
}

func dockerImageName() string {
	return fmt.Sprintf("%s/%s:%s",
		BuildConfig.Docker.Repository,
		BuildConfig.App.Name,
		BuildConfig.FullBuildNumber())
}