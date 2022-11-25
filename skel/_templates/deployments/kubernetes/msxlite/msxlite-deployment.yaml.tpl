---
  apiVersion: apps/v1
  kind: Deployment
  metadata:
    labels:
      app: ${app.name}
      group: ${kubernetes.group}
    namespace: vms
    name: ${app.name}-${app.version}
  spec:
    replicas: 1
    selector:
      matchLabels:
        app: ${app.name}
        group: ${kubernetes.group}
    template:
      metadata:
        name: ${app.name}-${app.version}
        labels:
          app: ${app.name}
          group: ${kubernetes.group}
          consul-gossip: allow
        annotations:
          tagprefix: logfmt
      spec:
        serviceAccountName: ${app.name}
        affinity:
          podAntiAffinity:
            preferredDuringSchedulingIgnoredDuringExecution:
              - weight: 100
                podAffinityTerm:
                  labelSelector:
                    matchExpressions:
                      - key: app
                        operator: In
                        values:
                          - ${app.name}
                  topologyKey: kubernetes.io/hostname
        imagePullSecrets:
          - name: ciscovms-registry
        containers:
          - name: consul
            image: registry.service.consul:5000/consul-vms:5.0.0-1.13.2-353
            command:
              - consul
              - agent
              - -bind=0.0.0.0
              - -client=0.0.0.0
              - -datacenter=vms1
              - -retry-join=consul.service.consul
              - -data-dir=/consul/data
              - -config-dir=/consul/config
            volumeMounts:
              - mountPath: /consul/config
                name: phi
          - name: ${app.name}
            image: registry.service.consul:5000/${app.name}
            command:
              - "/usr/bin/${app.name}"
              - --profile
              - production
            resources:
              requests:
                memory: "64Mi"
                cpu: "500m"
              limits:
                memory: "256Mi"
                cpu: "2000m"
            env:
              # BEGIN ANSIBLE MANAGED BLOCK
              - name: SPRING_PROFILES_ACTIVE
                value: "indepauthservice"
              - name: SPRING_CLOUD_STREAM_DEFAULT_CONSUMER_AUTOSTARTUP
                value: "true"
              - name: SPRING_CLOUD_CONSUL_DISCOVERY_DEFAULTQUERYTAG
                value: "${app.version}"
              - name: SPRING_CLOUD_CONSUL_DISCOVERY_TAGS
                value: "${app.version}"
              # END ANSIBLE MANAGED BLOCK
              - name: PROFILE
                value: production
              - name: SPRING_CLOUD_CONSUL_HOST
                value: "localhost"
              - name: SPRING_CLOUD_CONSUL_SCHEME
                value: "https"
              - name: SPRING_CLOUD_CONSUL_PORT
                value: "8500"
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
                value: "https"
              - name: SPRING_CLOUD_VAULT_TOKEN-SOURCE_SOURCE
                value: "kubernetes"
              - name: SPRING_CLOUD_VAULT_TOKEN-SOURCE_KUBERNETES_ROLE
                value: "${app.name}"
              - name: SPRING_REDIS_SENTINEL_ENABLE
                value: "true"
            ports:
              - containerPort: 7858
            volumeMounts:
              - mountPath: /etc/ssl/certs/ca-certificates.crt
                name: rootcert
              - mountPath: /certs/${app.name}-key.pem
                name: cockroach-client-sslkey
              - mountPath: /certs/${app.name}.pem
                name: cockroach-client-sslcert
              - mountPath: /keystore
                name: keystore
        volumes:
          - hostPath:
              path: /etc/ssl/certs/ca-bundle.crt
            name: rootcert
          - hostPath:
              path: /etc/ssl/vms-certs/${app.name}-key.pem
            name: cockroach-client-sslkey
          - hostPath:
              path: /etc/ssl/vms-certs/${app.name}.pem
            name: cockroach-client-sslcert
          - hostPath:
              path: /data/vms/keystore/
            name: keystore
          - configMap:
              name: phi
            name: phi