# Session Mode

Session mode is the default behavior when `netbar` is run in an interactive terminal.

```sh
netbar
```

In this mode, netbar:

- Starts a shell or command inside a PTY.
- Reserves the bottom terminal row.
- Draws network status on that row.
- Updates the status as connectivity changes.
- Restores the terminal when the wrapped command exits.

## Run a Shell

```sh
netbar
```

or:

```sh
netbar -- zsh
```

## Run an App

```sh
netbar -- codex
```

## Why netbar Must Wrap the Command

Terminal applications draw to the same screen. A separate process cannot reliably draw under an already-running full-screen terminal app.

That is why this works:

```sh
netbar -- codex
```

and this does not affect an existing Codex session:

```sh
netbar
```

## Terminal Size Overrides

Some terminal apps read `COLUMNS` and `LINES` during startup. If the terminal reports a narrow size, the app may render in only part of the screen.

Use:

```sh
stty size
netbar -rows 48 -cols 160 -- codex
```

Replace `48` and `160` with your real terminal size.
