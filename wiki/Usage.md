# Usage

## Open a netbar Session

```sh
netbar
```

This opens your default shell inside a netbar-managed terminal session and reserves the bottom row for network status.

## Run a Command Inside netbar

Use `--` to separate netbar flags from the command you want to run:

```sh
netbar -- zsh
netbar -- bash
netbar -- codex
```

The spacing matters. This is correct:

```sh
netbar -- codex
```

This is not correct:

```sh
netbar --codex
```

## One-Shot Status

```sh
netbar -once
```

Example:

```text
Online host=8.8.8.8:53 latency=42ms checked_at=2026-05-24T00:00:00+01:00
```

## Stream Status Lines

```sh
netbar -stream
```

This prints continuous line-by-line updates instead of opening an interactive session.

## Custom Probe Target

```sh
netbar -host 1.1.1.1:53
```

## Custom Check Interval

```sh
netbar -interval 5s
```

## Force Terminal Size

If a wrapped full-screen CLI starts too narrow:

```sh
stty size
```

If `stty size` prints:

```text
48 160
```

Run:

```sh
netbar -rows 48 -cols 160 -- codex
```
