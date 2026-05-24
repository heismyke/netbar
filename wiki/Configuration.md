# Configuration

## Flags

| Flag | Default | Description |
| --- | --- | --- |
| `-host` | `8.8.8.8:53` | TCP host and port to probe. |
| `-interval` | `3s` | Poll interval for session and stream modes. |
| `-once` | `false` | Run one check and exit. |
| `-stream` | `false` | Print continuous line-by-line updates instead of opening a session. |
| `-format` | `plain` | Output format: `plain` or `tmux`. |
| `-state-file` | user cache directory | File used to remember previous status. |
| `-rows` | detected terminal rows | Override terminal rows in session mode. |
| `-cols` | detected terminal columns | Override terminal columns in session mode. |

## Probe Host

The default probe host is:

```text
8.8.8.8:53
```

Use another target:

```sh
netbar -host 1.1.1.1:53
```

## Interval

```sh
netbar -interval 5s
```

## State File

`netbar` stores the previous status so it can show `Back online`.

Override the state file:

```sh
netbar -state-file ~/.cache/netbar/status
```

## Terminal Dimensions

If needed:

```sh
netbar -rows 48 -cols 160 -- codex
```
