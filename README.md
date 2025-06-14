# Kubernetes playground

Playground for learning Kubernetes and Helm.


## Oh My Zsh Kubectl plugin

Add kubectl to ohmyzsh plugins=(... kubectl).

Ref.: https://github.com/ohmyzsh/ohmyzsh/tree/master/plugins/kubectl


## Basic commands

```bash
kubectl apply -f infra/
```

```bash
kubectl delete -f infra/frontend.yaml
```

**Note, if you delete infra/volume.yaml, volume data will be lost.**

```bash
kubectl delete -f infra/
```

## Helm commands

```bash
helm lint --debug ./helm-deploy
```

```bash
helm template --debug ./helm-deploy
```

```bash
helm install --dry-run --debug test ./helm-deploy
```

```bash
helm install test ./helm-deploy
```

```bash
helm status test
```

```bash
helm upgrade test
```

```bash
helm uninstall test
```
