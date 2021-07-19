apiVersion: v1
kind: ConfigMap
metadata:
  name: ${app.name}-config
  namespace: {{ kubernetes_namespace }}
data:
  ${app.name}.yml: |
    ############################# ${app.name} ######################################

    ${app.name}:
      period: {{ ${app.name}_period }}
      timeout: 10s

    server:
      port: {{ ${app.name}_port }}

    spring.cloud.vault.ssl:
      ca-cert: ""
      client-cert: ""
      client-key: ""
      insecure: true

    elasticsearch.host: "es-logs.service.consul:9200"

    spring.cloud.stream.kafka.binder.brokers: "kafka.service.consul"
