package build

func init() {
	AddTarget("install-dockerfile", "Install the build dockerfile", InstallDockerfile)
	AddTarget("docker-build", "Build the target docker image", DockerBuild)
	AddTarget("docker-push", "Push the target docker image to the upstream repository", DockerPush)
}

func InstallDockerfile(args []string) error {
	return nil
}

func DockerBuild(args []string) error {
	return nil
}

func DockerPush(args []string) error {
	return nil
}
