# autoscaling preparation
Example yaml config for horizontal autoscaling based on a metric from prometheus

```
apiVersion: autoscaling/v2beta2
kind: HorizontalPodAutoscaler
metadata:
  name: my-app-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: my-app-deployment
  minReplicas: 1
  maxReplicas: 10
  metrics:
  - type: External
    external:
      metric:
        name: custom_metric_name
        selector:
          matchLabels:
            job: my-app
      target:
        type: Value
        value: 10
```

helm install f2s f2s/f2s --set promtail.enabled=true