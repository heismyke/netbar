# netbar

`netbar` is a tiny terminal connectivity monitor for AI CLI workflows.

It periodically checks DNS resolution and TCP connectivity, then prints a compact status such as `Online`, `Offline`, `Degraded`, or `Back online`. You can run it directly, wire it into terminal status bars, or use it with tools like tmux.

## Why

AI coding tools, package managers, deploy commands, and remote shells all depend on the network. When the network drops, terminal tools often look like they are hanging. `netbar` gives you a small network status signal so you can tell whether the issue is your command or your connection.

## Requirements

- Go 1.22 or newer

## Install

```sh
go install github.com/heismyke/netbar/cmd/netbar@latest
```

Make sure Go's binary directory is on your `PATH`:

```sh
export PATH="$PATH:$(go env GOPATH)/bin"
```

## Usage

Run continuously with the default probe target:

```sh
netbar
```

Probe a custom TCP endpoint or interval:

```sh
netbar -host 1.1.1.1:53 -interval 5s
```

Run a single check for terminal integrations:

```sh
netbar -once
```

Render a tmux-compatible status segment:

```sh
netbar -once -format tmux
```

Example output:

```text
Online host=8.8.8.8:53 latency=42ms checked_at=2026-05-24T00:00:00+01:00
```

Statuses:

- `online`: DNS and TCP checks pass within the latency threshold.
- `degraded`: DNS and TCP checks pass, but TCP latency is above the threshold.
- `offline`: DNS or TCP checks fail.

## tmux

tmux is optional. This is one way to display `netbar` like a terminal status bar.

Add this to `~/.tmux.conf`:

```tmux
set -g status on
set -g status-interval 5
set -g status-right '#(netbar -once -format tmux) %H:%M'
```

Reload tmux:

```sh
tmux source-file ~/.tmux.conf
```

`netbar` remembers the previous status in the user cache directory, so when the status changes from `offline` or `degraded` to `online`, it displays `Back online`.

## Flags

| Flag | Default | Description |
| --- | --- | --- |
| `-host` | `8.8.8.8:53` | TCP host and port to probe. |
| `-interval` | `3s` | Poll interval for continuous mode. |
| `-once` | `false` | Run one check and exit. Useful for tmux. |
| `-format` | `plain` | Output format: `plain` or `tmux`. |
| `-state-file` | user cache directory | File used to remember the previous status. |

## Development

Run all tests:

```sh
go test ./...
```

Run static checks:

```sh
go vet ./...
```

Build locally:

```sh
go build ./...
```

## Contributing

Issues and pull requests are welcome. Before opening a pull request, run:

```sh
go test ./...
go vet ./...
```

## License

MIT
