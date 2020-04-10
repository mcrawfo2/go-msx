---
apiVersion: apps/v1
kind: Deployment
metadata:
  namespace: vms
  name: ${app.name}
spec:
  replicas: {{ deployment_mode_env[deployment_mode|lower]['replica_count']['${app.name}'] }}
  selector:
    matchLabels:
      app: ${app.name}
      group: ${kubernetes.group}
      consul-gossip: allow
  template:
    metadata:
      name: ${app.name}
      labels:
        app: ${app.name}
        group: ${kubernetes.group}
        consul-gossip: allow
      annotations:
        fluentbit.io/parser: gomsx
    spec:
      containers:
        - name: consul
          image: {{ consul_image }}:{{ consul_version }}
          command:
            - consul
            - agent
            - -bind=0.0.0.0
            - -client=0.0.0.0
            - -datacenter={{ consul_dc }}
            - -retry-join=consul.service.consul
            - -data-dir=/consul/data
            - -config-dir=/consul/config
          volumeMounts:
            - mountPath: /consul/config
              name: phi
        - name: ${app.name}
          image: {{ ${app.name}_image }}:{{ ${app.name}_version }}
          command:
            - "/usr/bin/${app.name}"
            - --profile
            - production
          livenessProbe:
            httpGet:
              path: ${server.contextPath}/admin/alive
              port: ${server.port}
            initialDelaySeconds: 300
            periodSeconds: 30
          readinessProbe:
            httpGet:
              path: ${server.contextPath}/admin/health
              port: ${server.port}
            initialDelaySeconds: 60
            periodSeconds: 30
          resources:
            requests:
              memory: "64Mi"
              cpu: "1"
            limits:
              memory: "256Mi"
              cpu: "2"
          env:
            - name: SPRING_CLOUD_CONSUL_HOST
              value: "{{ vault_scheme }}://localhost"
            - name: SPRING_CLOUD_CONSUL_PORT
              value: "8500"
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
          ports:
            - containerPort: ${server.port}
          volumeMounts:
            - mountPath: /keystore
              name: keystore
            - mountPath: /etc/ssl/certs/ca-certificates.crt
              name: certs
      volumes:
        - name: phi
          configMap:
            name: phi
        - name: keystore
          hostPath:
            path: /data/vms/keystore/
        - name: certs
          hostPath:
            path: /etc/ssl/certs/ca-bundle.crt
