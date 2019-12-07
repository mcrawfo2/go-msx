---
apiVersion: policy/v1beta1
kind: PodDisruptionBudget
metadata:
  labels:
    app: ${app.name}
    group: ${kubernetes.group}
  name: ${app.name}-pdb
  namespace: {{ kubernetes_namespace }}
spec:
  maxUnavailable: 1
  selector:
    matchLabels:
      app: ${app.name}
