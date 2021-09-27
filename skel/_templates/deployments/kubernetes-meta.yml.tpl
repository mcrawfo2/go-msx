---
apiVersion: v1
kind: Pod
metadata:
  namespace: "{{ kubernetes_namespace }}"
  name: ${app.name}meta
  labels:
    app: ${app.name}meta
    group: ${kubernetes.group}
  annotations:
    tagprefix: logfmt
spec:
  restartPolicy: Never
  containers:
    - name: ${app.name}
      image: {{ ${app.name}_image }}:{{ ${app.name}_version }}
      command:
        - /usr/bin/${app.name}
        - --profile
        - production
        - populate
      env:
        - name: MODE
          value: {{ schema_mode }}
        - name: POPULATE
          value: all
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
      volumeMounts:
        - mountPath: /etc/ssl/certs/ca-certificates.crt
          name: rootcert
        - mountPath: /certs/zscalerservice-key.pem
          name: cockroach-client-sslkey
        - mountPath: /certs/zscalerservice.pem
          name: cockroach-client-sslcert
        - mountPath: /keystore
          name: keystore
  volumes:
    - hostPath:
        path: /etc/ssl/certs/ca-bundle.crt
      name: rootcert
    - hostPath:
        path: /etc/ssl/vms-certs/cloudsecurityservice-key.pem
      name: cockroach-client-sslkey
    - hostPath:
        path: /etc/ssl/vms-certs/cloudsecurityservice.pem
      name: cockroach-client-sslcert
    - name: keystore
      hostPath:
        path: /data/vms/keystore/