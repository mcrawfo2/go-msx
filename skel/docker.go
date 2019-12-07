package skel


func init() {
	AddTarget("generate-dockerfile", "Create a dockerfile for the target executable", GenerateDockerfile)
}

func GenerateDockerfile(args []string) error {
	logger.Info("Generating Dockerfile")
	return renderTemplate("docker/Dockerfile")
}
