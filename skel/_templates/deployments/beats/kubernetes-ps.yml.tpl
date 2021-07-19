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
  replicas: {{ deployment_mode_env[deployment_mode|lower]['replica_count']['${app.name}_count'] }}
  selector:
    matchLabels:
      name: ${app.name}
      app: ${app.name}
      group: dataplatform
  updateStrategy:
    type: RollingUpdate
  template:
    metadata:
      labels:
        name: ${app.name}
        app: ${app.name}
        group: dataplatform
      annotations:
        tagprefix: logfmt
    spec:
{% if (deployment_mode | lower) != 'lite' %}
      tolerations:
      - key: "dedicated"
        operator: "Equal"
        value: "sa-probe"
{% endif %}
      affinity:
        nodeAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 100
            preference:
              matchExpressions:
              - key: vms/nodetype
                operator: In
                values:
                - sa
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
            - name: SPRING_CLOUD_CONSUL_DISCOVERY_DEFAULTQUERYTAG
              value: "{{ deployer_version }}"
            - name: SPRING_CLOUD_CONSUL_DISCOVERY_TAGS
              value: "{{ deployer_version }}"
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
            - name: SPRING_CLOUD_VAULT_TOKEN
{% if dual_dc and  dc == "passive" %}
              value: "{{ dc1_beats_client_token_auth_client_token }}"
{% else %}
              value: "{{ beats_general_token.stdout }}"
{% endif %}
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
          initialDelaySeconds: 60
          periodSeconds: 30
          successThreshold: 1
          timeoutSeconds: 10
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

