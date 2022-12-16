service:
  name: "${app.name}"
  identifier: "${app.name}"
  tags: {}
  serviceDefinition:
    spec:
      manifests:
        - manifest:
            identifier: helm
            type: K8sManifest
            spec:
              store:
                type: Github
                spec:
                  connectorRef: account.githuborg
                  gitFetchType: Branch
                  paths:
                  #if GENERATOR_BEATS
                    - "charts/beatservice-template/templates"
                  #else GENERATOR_BEATS
                    - "charts/goservice-template/templates"
                  #endif GENERATOR_BEATS
                  repoName: helm-charts
                  branch: develop
              valuesPaths:
                  #if GENERATOR_BEATS
                - "charts/beatservice-template/values_${app.name}.yaml"
                  #else GENERATOR_BEATS
                - "charts/goservice-template/values_${app.name}.yaml"
                  #endif GENERATOR_BEATS
              skipResourceVersioning: false
        - manifest:
            identifier: Ansible Generated Values File
            type: Values
            spec:
              store:
                type: CustomRemote
                spec:
                  filePath: /tmp/custom/helm_init_vars/values_ansible.yaml
                  extractionScript: |-
                    FILE=/tmp/custom/helm_init_vars/values_ansible.yaml
                    if test -f "$FILE"; then
                        echo "$FILE exists."
                    else
                       touch /tmp/custom/helm_init_vars/values_ansible.yaml
                    fi
                  delegateSelectors:
                    - <+env.name>
      variables:
        - name: slack_notification_channel
          type: String
          description: "the name of the slackchannel where all your service notifications will go prefixed with #cd-pipeline-. for example, setting the value to meraki will result in notifications going to #cd-pipeline-meraki"
          value: "${app.shortname}"
        - name: database_required
          type: String
          description: indicate if your service requires a Database. this is needed for fresh install to let the sp-init ansible playbook know if it needs to create a DB for the service
          value: "${repository.cockroach.enabled}"
        - name: database_name_if_required
          type: String
          description: the name of the database if a DB is needed for your service
          value: "${app.name}"
        - name: nso_required
          type: String
          description: set to NSO if your service pack requires initializing NSO user name and password (MD only)
          value: "false"
        - name: has_ui
          type: String
          description: set to true if your service comes with UI. this is needed to tell the deployment pipeline whether to deploy UI along with your service
          #if UI
          value: "true"
          #else UI
          value: "false"
          #endif UI
        - name: servicepack
          type: String
          description: the official name of the servicepack. most of the name it is the name of the service without 'service' postfix, except for the beats
          value: "${kubernetes.group}"
        - name: test_component
          type: String
          description: the name of the component responsible for testing the service. for example, merakigo uses meraki component tests to validate the service
          value: "${kubernetes.group}"
      artifacts:
        primary:
          sources:
            - spec:
                connectorRef: account.cisco_dockerhub_ext
                imagePath: vms-platform-dev-docker/${app.version}/latest/${app.name}
                tag: <+input>
              identifier: "${app.name}"
              type: DockerRegistry
          primaryArtifactRef: "${app.name}"
        sidecars:
          - sidecar:
              spec:
                connectorRef: account.cisco_dockerhub_ext
                imagePath: vms-platform-dev-docker/consul-vms
                tag: <+variable.consulsidecar>
              identifier: consul
              type: DockerRegistry
          - sidecar:
              spec:
                connectorRef: account.cisco_dockerhub_ext
                imagePath: vms-platform-dev-docker/vault-vms
                tag: <+variable.vaultsidecar>
              identifier: vault
              type: DockerRegistry
    type: Kubernetes