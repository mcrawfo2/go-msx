#jinja2: lstrip_blocks: "True", trim_blocks: "True"
---
apiVersion: "apps/v1"
kind: StatefulSet
metadata:
  labels:
    name: ${app.name}
    version: "{{ ${app.name}_version }}"
    app: ${app.name}
    group: dataplatform
  name: ${app.name}
  namespace: {{ kubernetes_namespace }}
spec:
  serviceName: ${app.name}
  replicas: {{ ${deployment.group}.replica_count[deployment_mode|lower] }}
  selector:
    matchLabels:
      name: ${app.name}
      app: ${app.name}
      group: dataplatform
      consul-gossip: allow
  updateStrategy:
    type: RollingUpdate
  template:
    metadata:
      labels:
        name: ${app.name}
        app: ${app.name}
        group: dataplatform
        consul-gossip: allow
      annotations:
        tagprefix: logfmt
    spec:
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
{% if cloud == 'aws' %}
          - weight: 100
            podAffinityTerm:
              labelSelector:
                matchExpressions:
                - key: app
                  operator: In
                  values:
                  - ${app.name}
              topologyKey: topology.kubernetes.io/zone
{% endif %}
      containers:
      - name: ${app.name}
        image: {{ ${app.name}_image }}:{{ ${app.name}_version }}
        env:
{% if no_proxy is defined and no_proxy %}
            - name: NO_PROXY
              value: "{{ kubernetes_overlay_network }},{{ no_proxy }}"
{% endif %}
{% if proxy is defined and proxy %}
            - name: HTTP_PROXY
              value: "{{ proxy }}"
            - name: HTTPS_PROXY
              value: "{{ proxy }}"
{% endif %}
            - name: SPRING_CLOUD_CONSUL_HOST
              value: "localhost"
            - name: SPRING_CLOUD_CONSUL_SCHEME
              value: "https"
            - name: SPRING_CLOUD_CONSUL_PORT
              value: "8500"
            - name: SPRING_CLOUD_VAULT_HOST
              value: "vault.service.consul"
            - name: SPRING_CLOUD_VAULT_PORT
              value: "8200"
            - name: SPRING_CLOUD_VAULT_SCHEME
              value: "https"
            - name: SPRING_CLOUD_VAULT_TOKEN-SOURCE_SOURCE
              value: "kubernetes"
            - name: SPRING_CLOUD_VAULT_TOKEN-SOURCE_KUBERNETES_ROLE
              value: "{{.Values.name }}"
        livenessProbe:
          failureThreshold: 3
          httpGet:
            path: /admin/alive
            port: {{ ${app.name}_port }}
            scheme: HTTP
          initialDelaySeconds: 300
          periodSeconds: 30
          successThreshold: 1
          timeoutSeconds: 10
        readinessProbe:
          failureThreshold: 3
          httpGet:
            path: /admin/health
            port: {{ ${app.name}_port }}
            scheme: HTTP
          initialDelaySeconds: 0
          periodSeconds: 30
          successThreshold: 1
          timeoutSeconds: 10
        startupProbe:
          httpGet:
            path: /admin/health
            port: {{ ${app.name}_port }}
          failureThreshold: 45
          timeoutSeconds: 5
{% if deployment_mode == 'production' %}
        resources:
         limits:
           cpu: "1000m"
           memory: 500Mi
         requests:
           cpu: "1000m"
           memory: 500Mi
{% endif %}

        ports:
         - containerPort: {{ ${app.name}_port }}
           name: http
        volumeMounts:
        - name: config
          mountPath: /etc/${app.name}/${app.name}.yml
          subPath: ${app.name}.yml
        - mountPath: /keystore
          name: keystore
        - mountPath: /etc/ssl/certs/ca-certificates.crt
          name: certs
        - mountPath: /certs
          name: cachecerts
      - name: consul
        image: {{ consul_image }}:{{ consul_version }}
        command:
          - consul
          - agent
          - -bind=0.0.0.0
          - -client=0.0.0.0
          - -datacenter={{ consul_dc }}
          - -encrypt={{ consul_encrypt_key }}
          - -retry-join=consul.service.consul
          - -data-dir=/consul/data
          - -config-dir=/consul/config
        volumeMounts:
            - mountPath: /consul/config
              name: beatconsul
      volumes:
        - name: config
          configMap:
            name: ${app.name}-config
        - name: keystore
          hostPath:
            path: /data/vms/keystore/
        - name: beatconsul
          configMap:
            name: beatconsul
        - name: certs
          hostPath:
            path: /etc/ssl/certs/ca-bundle.crt
        - name: cachecerts
          emptyDir: {}
