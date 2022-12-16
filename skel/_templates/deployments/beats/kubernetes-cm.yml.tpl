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

    elasticsearch.host: "https://opensearch:9200"
    probe.outputs.elasticsearch.tls.enabled: true
    probe.outputs.elasticsearch.tls.certificate-source: "identitycache"

    certificate.source.identity.path: "{{ vms_subdomain }}/pki/vms"
    certificate.source.identity.role: "${app.name}"

    certificate.source.identitycache.provider: cache
    certificate.source.identitycache.upstream-source: identity
    certificate.source.identitycache.key-file: "/certs/${app.name}-key.pem"
    certificate.source.identitycache.cert-file: "/certs/${app.name}.pem"
    certificate.source.identitycache.ca-cert-file: "/certs/ca-identity.crt"

    spring.cloud.stream.kafka.binder.brokers: "kafka.service.consul"

    spring.cloud.consul.config.required:
