package build

import (
	"cto-github.cisco.com/NFV-BU/go-msx/exec"
	"fmt"
	"path"
)

func init() {
	AddTarget("docker-build", "Build the target docker image", DockerBuild)
	AddTarget("docker-push", "Push the target docker image to the upstream repository", DockerPush)
}

func DockerBuild(args []string) error {
	logger.WithExtendedField("target", "docker-build").
		Infof("BASE_IMAGE=%s", dockerBaseImage())

	return exec.ExecutePipes(exec.ExecSimple(
		"docker", "build",
		"-t", dockerImageName(),
		"-f", BuildConfig.Docker.Dockerfile,
		"--build-arg", "BUILDER_FLAGS",
		"--build-arg", "BUILD_FLAGS",
		"--build-arg", "BASE_IMAGE="+dockerBaseImage(),
		"--force-rm",
		"--no-cache",
		"."))
}

func DockerPush(args []string) error {
	logger.WithExtendedField("target", "docker-push").
		Infof("IMAGE=%s", dockerImageName())

	if BuildConfig.Docker.Username != "" && BuildConfig.Docker.Password != "" {
		err := exec.ExecutePipes(exec.ExecSimple(
			"docker", "login",
			"-u", BuildConfig.Docker.Username,
			"-p", BuildConfig.Docker.Password,
			path.Dir(BuildConfig.Docker.Repository)))
		if err != nil {
			return err
		}
	}

	return exec.ExecutePipes(exec.ExecSimple(
		"docker", "push", dockerImageName()))
}

func dockerImageName() string {
	return fmt.Sprintf("%s/%s:%s",
		BuildConfig.Docker.Repository,
		BuildConfig.App.Name,
		BuildConfig.FullBuildNumber())
}

func dockerBaseImage() string {
	return fmt.Sprintf("%s/%s",
		BuildConfig.Docker.Repository,
		BuildConfig.Docker.BaseImage)
}
