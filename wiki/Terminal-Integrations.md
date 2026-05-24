# Terminal Integrations

`netbar` supports integration modes for terminal prompts, scripts, and status bars.

## One-Shot Mode

```sh
netbar -once
```

This performs one connectivity check and exits.

## tmux-Compatible Output

```sh
netbar -once -format tmux
```

Example tmux configuration:

```tmux
set -g status on
set -g status-interval 5
set -g status-right '#(netbar -once -format tmux) %H:%M'
```

Reload:

```sh
tmux source-file ~/.tmux.conf
```

tmux is optional. It is one integration target, not the purpose of the project.

## Plain Output

```sh
netbar -once -format plain
```

## Stream Mode

```sh
netbar -stream
```

Use stream mode when you want continuous logs rather than an interactive session.
