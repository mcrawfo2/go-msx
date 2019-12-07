package skel

func GenerateSkeleton(args []string) error {
	if err := ConfigureInteractive(nil); err != nil {
		return err
	}
	if err := GenerateBuild(nil); err != nil {
		return err
	}
	if err := GenerateApp(nil); err != nil {
		return err
	}
	if err := GenerateDockerfile(nil); err != nil {
		return err
	}
	if err := GenerateKubernetesManifestTemplates(nil); err != nil {
		return err
	}
	if err := GenerateKubernetesManifests(nil); err != nil {
		return err
	}
	return nil
}
