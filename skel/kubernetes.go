package skel

func init() {
	AddTarget("kubernetes-manifest-templates", "Create production kubernetes manifest templates", GenerateKubernetesManifestTemplates)
	AddTarget("kubernetes-manifests", "Create dev kubernetes manifests", GenerateKubernetesManifests)
}

func GenerateKubernetesManifestTemplates(args []string) error {
	logger.Info("Generating kubernetes manifest templates")
	return renderTemplate("k8s/kubernetes-deployment.yml.tpl")
}

func GenerateKubernetesManifests(args []string) error {
	logger.Info("Generating kubernetes manifests")
	return nil
}
