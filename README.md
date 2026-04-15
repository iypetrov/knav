# knav
A simple CLI tool, inspired by kubectx, that makes it easy to manage multiple kubeconfig files when working with kubectl.

# Example

```yaml
# ~/.config/knav/config.yaml
current: local
targets:
    - name: local
      kubeconfigPath: ~/.kube/config
      envs:
        - name: HTTPS_PROXY
          value: socks5://127.0.0.1:1234
      restricted: true
      allowedActions:
        - get
        - logs
    - name: gardener-local
      kubeconfigPath: ~/Projects/oss/gardener/example/gardener-local/kind/local/kubeconfig
      restricted: false
```
