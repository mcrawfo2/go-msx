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
	return exec.Execute(
		"docker", "build",
		"-t", dockerImageName(),
		"-f", "docker/Dockerfile",
		"--build-arg", "BUILDER_FLAGS",
		"--build-arg", "BUILD_FLAGS",
		"--force-rm",
		"--no-cache",
		".")
}

func DockerPush(args []string) error {
	if BuildConfig.Docker.Username != "" && BuildConfig.Docker.Password != "" {
		err := exec.Execute(
			"docker", "login",
			"-u", BuildConfig.Docker.Username,
			"-p", BuildConfig.Docker.Password)
		if err != nil {
			return err
		}
	}

	return exec.Execute("docker", "push", dockerImageName())
}

func dockerImageName() string {
	return fmt.Sprintf("%s/%s:%s",
		BuildConfig.Docker.Repository,
		BuildConfig.App.Name,
		BuildConfig.FullBuildNumber())
}
