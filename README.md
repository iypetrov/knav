# knav

A lightweight CLI tool, inspired by `kubectx`, for managing multiple kubeconfig files and safely switching between Kubernetes environments.

## Overview

`knav` simplifies working with multiple Kubernetes clusters by making context switching seamless and safer.

It also introduces **action restrictions**, helping prevent accidental execution of destructive commands—especially in sensitive environments like production.

## Requirements

Before using `knav`, make sure you have:

* [`kubectl`](https://kubernetes.io/docs/tasks/tools/)
* [`fzf`](https://github.com/junegunn/fzf)

## Installation

Clone the repository and build the tool:

```bash
make build && make local
```

## Configuration

On first run, `knav` automatically generates a configuration file at `~/.config/knav/config.yaml` with this content:

```yaml
current: local
targets:
  - name: local
    kubeconfigPath: ~/.kube/config
    restricted: false
```

You can override the config location by setting `KNAVCONFIG` to a directory that contains a `config.yaml` file.

### All available options

Top-level keys:

| Key       | Description                                                       |
| --------- | ----------------------------------------------------------------- |
| `current` | Name of the target that is currently active.                      |
| `targets` | List of target definitions (see below).                           |

Per-target keys:

| Key              | Required | Description                                                                                             |
| ---------------- | :------: | ------------------------------------------------------------------------------------------------------- |
| `name`           |    ✓     | Display name shown in the fuzzy picker and used as the value of `current`.                              |
| `kubeconfigPath` |    ✓     | Path to the kubeconfig file. `~` is expanded to `$HOME`.                                                |
| `envs`           |          | List of `{name, value}` pairs exported into the environment of the forwarded `kubectl` invocation.      |
| `restricted`     |          | If `true`, only verbs listed in `allowedActions` run without confirmation; everything else prompts.     |
| `allowedActions` |          | Whitelist of `kubectl` verbs allowed to run directly when the target is `restricted`.                   |
| `initScript`     |          | Shell commands executed **once** when you switch *to* this target via the picker. See below for usage.  |

### Full example

The example below demonstrates every option, with realistic patterns drawn from a daily workflow:

```yaml
current: local
targets:
  # A dev cluster, treated as read-only by default. All traffic goes
  # through a SOCKS proxy (e.g. opened by `cloudflared access tcp`).
  - name: local
    kubeconfigPath: ~/.kube/config
    envs:
      - name: HTTPS_PROXY
        value: socks5://127.0.0.1:1234
    restricted: true
    allowedActions:
      - get
      - top
      - describe
      - logs
      - api-resources
      - attach
      - exec
      - port-forward
      - cp
      - auth
      - debug
      - events
      - diff
      - wait
      - explain

  # A local cluster — no restrictions, no init.
  - name: dev-local 
    kubeconfigPath: ~/projects/dev/kubeconfigs/runtime/kubeconfig
    restricted: false

  # The init script makes sure you're logged in to you staging GCP account.
  - name: staging
    kubeconfigPath: ~/projects/remote/kubeconfig
    restricted: true
    allowedActions:
      - get
      - logs
      - describe
      - exec
      - port-forward
    initScript:
      - gcloud auth print-access-token > /dev/null || gcloud auth login
```
### Auto-export `KUBECONFIG` to match the active target

So that other tools (`flux`, `helm`, `k9s`, prompt themes like Powerlevel10k, ...) automatically point at the same cluster as `knav`:

```bash
autoload -Uz add-zsh-hook

function _auto_kubeconfig() {
  local new_config
  new_config="$(realpath "$(yq -r '.current as $c | .targets[] | select(.name == $c) | .kubeconfigPath' ~/.config/knav/config.yaml | sed "s|^~|$HOME|")")"

  if [[ "$KUBECONFIG" != "$new_config" ]]; then
    export KUBECONFIG="$new_config"
    p10k reload   # remove if you don't use Powerlevel10k
  fi
}

add-zsh-hook precmd _auto_kubeconfig

source <(knav completion zsh)
```

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.
