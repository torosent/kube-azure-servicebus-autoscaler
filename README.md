# kube-azure-servicebus-autoscaler
Kubernetes pod autoscaler based on queue size in Azure Service Bus Queues. It periodically retrieves the number of messages in your queue and scales pods accordingly.

## Setting up
Setting up kube-azure-servicebus-autoscaler requires two steps:
1) Deploying it as an incluster service in your cluster
2) Adding Service Prinicipal credentials, subscription id and tenant id in Secrets so it can read the number of messages in your queues.

### Deploying kube-azure-servicebus-autoscaler
Deploying kube-azure-servicebus-autoscaler should be as simple as applying this deployment:
```yaml
---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: kube-azure-servicebus-autoscaler
  labels:
    app: kube-azure-servicebus-autoscaler
spec:
  replicas: 1
  selector:
    matchLabels:
      app: kube-azure-servicebus-autoscaler
  template:
    metadata:
      labels:
        app: kube-azure-servicebus-autoscaler
    spec:
      containers:
      - name: kube-azure-servicebus-autoscaler
        image: torosent/kube-azure-servicebus-autoscaler:1.0.0
        command:
          - /kube-azure-servicebus-autoscaler
          - --resourcegroup=queuerg  #required
          - --queuename=somequeuename  #required
          - --namespace=somenamespace  #required
          - --kubernetes-deployment=your-kubernetes-deployment-name # required
          - --kubernetes-namespace=$(POD_NAMESPACE) # optional
          - --poll-period=5s # optional
          - --scale-down-cool-down=30s # optional
          - --scale-up-cool-down=5m # optional
          - --scale-up-messages=100 # optional
          - --scale-down-messages=10 # optional
          - --max-pods=5 # optional
          - --min-pods=1 # optional
        env:
          - name: POD_NAMESPACE
            valueFrom:
              fieldRef:
                fieldPath: metadata.namespace
          - name: AZURE_CLIENT_ID
            valueFrom:
              secretKeyRef:
                name: azureserviceprincipal
                key: clientid
          - name: AZURE_CLIENT_SECRET
            valueFrom:
              secretKeyRef:
                name: azureserviceprincipal
                key: clientsecret
          - name: AZURE_SUBSCRIPTION_ID
            valueFrom:
              secretKeyRef:
                name: azureserviceprincipal
                key: subscriptionid
          - name: AZURE_TENANT_ID
            valueFrom:
              secretKeyRef:
                name: azureserviceprincipal
                key: tenantid
        resources:
          requests:
            memory: "200Mi"
            cpu: "100m"
          limits:
            memory: "200Mi"
            cpu: "100m"
        volumeMounts:
          - name: ssl-certs
            mountPath: /etc/ssl/certs/ca-certificates.crt
            readOnly: true
      volumes:
        - name: ssl-certs
          hostPath:
            path: "/etc/ssl/certs/ca-certificates.crt"
```
