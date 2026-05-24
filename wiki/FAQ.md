# FAQ

## Is netbar a tmux plugin?

No. tmux is only one optional integration. `netbar` is a general terminal connectivity monitor.

## Can I use netbar without tmux?

Yes:

```sh
netbar
netbar -- codex
netbar -once
```

## Why does netbar need to start Codex?

To draw a status row beneath Codex, netbar must own the terminal session and reserve that row before Codex starts.

## Can I run netbar inside an existing Codex session?

You can run commands, but netbar cannot attach a bottom status row to the already-running Codex UI.

Use:

```sh
netbar -- codex
```

from a fresh terminal.

## Does netbar require a daemon?

No. It runs as a normal CLI.

## What does Back online mean?

It means the current check is online and the previous stored status was degraded or offline.

## How do I update?

```sh
go install github.com/heismyke/netbar/cmd/netbar@latest
hash -r
```
