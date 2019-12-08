package build

func init() {
	AddTarget("install-kubernetes-manifests", "Install the distribution kubernetes manifests", InstallKubernetesManifests)
}

func InstallKubernetesManifests(args []string) error {
	return nil
}
