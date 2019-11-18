package build

func init() {
	AddTarget("generate-dockerfile", "Create a dockerfile for the target executable", GenerateDockerfile)
	AddTarget("docker-build", "Build the target docker image", DockerBuild)
	AddTarget("docker-login", "Login to the upstream docker repository", DockerLogin)
	AddTarget("docker-push", "Push the target docker image to the upstream repository", DockerPush)
}

func GenerateDockerfile(args []string) error {
	return nil
}

func DockerBuild(args []string) error {
	return nil
}


func DockerLogin(args []string) error {
	return nil
}

func DockerPush(args []string) error {
	return nil
}
