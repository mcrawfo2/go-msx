package skel

func GenerateApp(args []string) error {
	logger.Info("Generating application")
	if err := generateGoMod(); err != nil {
		return err
	}
	if err := generateBootstrapYml(); err != nil {
		return err
	}
	if err := generateMainGo(); err != nil {
		return err
	}
	return nil
}

func generateGoMod() error {
	logger.Info("Creating go module definition")
	return renderTemplate("go.mod")
}

func generateBootstrapYml() error {
	logger.Info("Creating bootstrap configuration")
	return renderTemplate("cmd/app/bootstrap.yml")
}

func generateMainGo() error {
	logger.Info("Creating application entrypoint source")
	return renderTemplate("cmd/app/main.go")
}
