apiVersion: skaffold/v3
kind: Config
metadata:
  name: ${app.name}

build:
  artifacts:
    - image: dockerhub.cisco.com/vms-platform-dev-docker/${app.name}
      custom:
        buildCommand: |
          export DOCKER_PUBLISHTAGS=devtest
          export DOCKER_TAGS_DEVTEST="${IMAGE}"
          make docker-debug docker-push

manifests:
  rawYaml:
    - deployments/kubernetes/minivms/${app.name}-deployment.yaml

profiles:
  - name: msxlite
    patches:
      - op: replace
        path: /manifests/rawYaml/0
        value: deployments/kubernetes/msxlite/${app.name}-deployment.yaml
    activation:
      - kubeContext: msxlite
      - env: SKAFFOLD_PROFILE=msxlite