# Traefik Readiness Plugin

This plugin adds an readiness endpoint that is useful when Traefik is deployed in Kubernetes.

## Configuration

Traefik Helm chart values:

```yaml
deployment:
  healthchecksPort: 8082
  readinessProbePath: /ready # TODO: implement this in the Traefik Helm chart

experimental:
  plugins:
    readiness:
      moduleName: github.com/livekit/traefik-readiness-plugin
      version: v0.0.1-alpha.1

ingressRoute:
  healthcheck:
    enabled: true
    matchRule: PathPrefix(`/ping`) || PathPrefix(`/ready`)
    entryPoints:
      - ping
    middlewares:
      - name: readiness

ports:
  ping:
    port: 8082
    expose:
      default: true
    exposedPort: 8082
    protocol: TCP
```

Middleware:

```yaml
apiVersion: traefik.io/v1alpha1
kind: Middleware
metadata:
  name: readiness
  namespace: traefik
spec:
  plugin:
    readiness:
      ReadyPath: /ready
      ReadyCPULimit: 0.8
```
