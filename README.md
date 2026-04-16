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

You can define multiple kubeconfig targets like this:

```yaml
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

If you are using zsh, you can add this to your `~/.zshrc` file, so the right kubeconfig to be propagated automatically.

```bash
autoload -Uz add-zsh-hook

function _auto_kubeconfig() {
  local new_config
  new_config="$(realpath "$(yq -r '.current as $c | .targets[] | select(.name == $c) | .kubeconfigPath' ~/.config/knav/config.yaml | sed "s|^~|$HOME|")")"

  if [[ "$KUBECONFIG" != "$new_config" ]]; then
    export KUBECONFIG="$new_config"
    p10k reload
  fi
}

add-zsh-hook precmd _auto_kubeconfig
```

## Key Features

### 🔄 Easy Context Switching

Quickly switch between different Kubernetes environments without manually exporting or modifying kubeconfig files by just running `knav` - it will give you the option to fuzzy search you target config file. It works very similarly as `kubectx`.

### 🔒 Safe Mode (Restricted Targets)

Prevent accidental destructive actions in critical environments.

When `restricted: true` is enabled:
- Only commands listed in `allowedActions` (e.g. `get`, `logs`, `describe`) run without interruption.
- Any other command requires explicit confirmation before execution.

This is especially useful when you have admin access to production clusters but want an extra safety layer.

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.
