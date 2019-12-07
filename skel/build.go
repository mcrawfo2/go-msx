package skel

func init() {
	AddTarget("generate-build", "Create the basic build system", GenerateBuild)
}

func GenerateBuild(args []string) error {
	logger.Info("Generating build command")
	if err := generateBuildYml(); err != nil {
		return err
	}
	if err := generateBuildGo(); err != nil {
		return err
	}
	return nil
}

func generateBuildYml() error {
	logger.Info("Creating build descriptor")
	return renderTemplate("cmd/build/build.yml")
}

func generateBuildGo() error {
	logger.Info("Creating build command source")
	return renderTemplate("cmd/build/build.go")
}
