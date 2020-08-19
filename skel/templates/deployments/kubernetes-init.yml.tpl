---
apiVersion: v1
kind: Pod
metadata:
  namespace: "{{ kubernetes_namespace }}"
  name: ${app.name}
  labels:
    app: ${app.name}
    group: ${kubernetes.group}
spec:
  restartPolicy: Never
  containers:
    - name: ${app.name}
      image: {{ ${app.name}_image }}:{{ ${app.name}_version }}
      command:
        - /usr/bin/${app.name}
        - --profile
        - production
        - migrate
      resources:
        requests:
          cpu: "{{ 1*deployment_mode_env[deployment_mode|lower]['replica_count']['${app.name}'] }}"
      env:
        - name: SPRING_CLOUD_CONSUL_HOST
          value: "consul.service.consul"
        - name: SPRING_CLOUD_CONSUL_PORT
          value: "8500"
        - name: SPRING_CLOUD_CONSUL_SCHEME
          value: "{{ vault_scheme }}"
        - name: SPRING_CLOUD_CONSUL_CONFIG_ACLTOKEN
          valueFrom:
            secretKeyRef:
              name: msxconsul
              key: token
        - name: SPRING_CLOUD_VAULT_HOST
          value: "vault.service.consul"
        - name: SPRING_CLOUD_VAULT_PORT
          value: "8200"
        - name: SPRING_CLOUD_VAULT_SCHEME
          value: "{{ vault_scheme }}"
        - name: SPRING_CLOUD_VAULT_TOKEN
          valueFrom:
            secretKeyRef:
              name: msxvault
              key: token
