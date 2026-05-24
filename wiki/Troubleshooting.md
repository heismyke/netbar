# Troubleshooting

## `netbar --codex` Shows an Error

Use a space after `--`:

```sh
netbar -- codex
```

`--` separates netbar flags from the command to wrap.

## Codex Opens in Half the Screen

Force the terminal size:

```sh
stty size
```

If the output is:

```text
48 160
```

run:

```sh
netbar -rows 48 -cols 160 -- codex
```

## netbar Prints Lines Instead of Opening a Session

You may be running an old binary.

Check:

```sh
netbar -h
```

The help output should include:

- `-stream`
- `-rows`
- `-cols`

Upgrade:

```sh
go install github.com/heismyke/netbar/cmd/netbar@latest
hash -r
```

## `netbar` Command Not Found

Add Go's binary directory to `PATH`:

```sh
export PATH="$PATH:$(go env GOPATH)/bin"
```

## It Shows Offline While the Browser Works

`netbar` checks DNS and TCP connectivity against the configured probe host.

Try another host:

```sh
netbar -host 1.1.1.1:53 -once
```

## Can netbar Attach to an Existing Codex Session?

No. Start Codex through netbar:

```sh
netbar -- codex
```
