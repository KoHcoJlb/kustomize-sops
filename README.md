# Simple Kustomize plugin to decrypt sops-encrypted yaml files

This plugin is designed to be simpler alternative to [KSOPS](https://github.com/viaduct-ai/kustomize-sops). It's main
target is to be used with [ArgoCD](https://github.com/argoproj/argo-cd).

## Installation with kustomize

This repo contains Kustomize component you can use to install this plugin in argocd. Add this to your
argocd's `kustomization.yaml`

```yaml
components:
- https://github.com/KoHcoJlb/kustomize-sops//argocd?ref=v0.1.0
```

Also, you need to mount corresponding private keys or environment variables into `argocd-repo-server` container.

### Example for age (for other key types see [SOPS](https://github.com/mozilla/sops) documentation)

repo-server.yaml (patchesStrategicMerge)

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: argocd-repo-server
spec:
  template:
    spec:
      containers:
      - name: argocd-repo-server
        volumeMounts:
        - mountPath: /home/argocd/.config/sops/age/
          name: sops-age-keys
      volumes:
      - name: sops-age-keys
        secret:
          secretName: sops-age-keys
```

sops-age-keys.yaml (resources)

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: sops-age-keys
  namespace: argocd
stringData:
  keys.txt: |
    AGE-SECRET-KEY-1EXA40TQ9U7Q544USTEZGDFY2WJ6CFNQU5V0YTECE0QW63AYNT6DS2JWV2P
```

## Usage

Just add this to your kustomization

```yaml
transformers:
- https://github.com/KoHcoJlb/kustomize-sops//transformer
```

## Note

[MAC verification](https://github.com/mozilla/sops#51message-authentication-code) is disabled in this plugin as it
conflicts with kustomize transformations.

## Other projects

[KSOPS](https://github.com/viaduct-ai/kustomize-sops)
