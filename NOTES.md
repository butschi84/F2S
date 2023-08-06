# f2soperator

## testing locally

```
export Prometheus_URL=localhost:9090
export KUBECONFIG=~/.kube/config
go run main.go
```

# Release

## F2S

```
- tag as 'f2s-0.1.1'
- pipeline will run on github actions
- pipeline will create branch 'helm-release-0.1.1'
- create a pull request to merge branch 'helm-release-0.1.1'
- tag with '0.1.1' after merge
```

## Fizzlet

```
tag as 'fizzlet-0.1.1'
pipeline will run on github actions
```