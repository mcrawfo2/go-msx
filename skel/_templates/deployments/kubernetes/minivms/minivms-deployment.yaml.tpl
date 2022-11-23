apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: ${app.name}
  name: ${app.name}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: ${app.name}
  strategy:
    type: Recreate
  template:
    metadata:
      annotations:
        fluentbit.io/parser: logfmt
        tagprefix: logfmt
      labels:
        app: ${app.name}
    spec:
      containers:
        - env:
            - name: SERVER_PORT
              value: "${server.port}"
            - name: SPRING_CLOUD_CONSUL_DISCOVERY_IPADDRESS
              valueFrom:
                fieldRef:
                  fieldPath: status.hostIP
            - name: SPRING_CLOUD_CONSUL_HOST
              valueFrom:
                fieldRef:
                  fieldPath: status.hostIP
            - name: SPRING_CLOUD_VAULT_HOST
              valueFrom:
                fieldRef:
                  fieldPath: status.hostIP
          image: dockerhub.cisco.com/vms-platform-dev-docker/${app.name}:${app.version}
          name: ${app.name}
          ports:
            - containerPort: ${server.port}
          volumeMounts:
            - mountPath: /etc/ssl/certs
              name: ${app.name}-user-volume0
              subPath: config/certs
      enableServiceLinks: false
      restartPolicy: Always
      volumes:
        - hostPath:
            path: /home/ubuntu/vms/volumes
          name: ${app.name}-user-volume0