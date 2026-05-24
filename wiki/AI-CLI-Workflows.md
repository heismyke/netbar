# AI CLI Workflows

`netbar` is useful when working with AI CLI tools because network failures can look like slow model responses, stalled package installs, or hung commands.

## Codex

Start Codex inside netbar:

```sh
netbar -- codex
```

If Codex renders too narrowly:

```sh
stty size
netbar -rows 48 -cols 160 -- codex
```

## Other AI CLIs

The same pattern works for other interactive CLIs:

```sh
netbar -- your-ai-cli
```

If the command takes arguments:

```sh
netbar -- your-ai-cli --some-flag value
```

Everything after `--` belongs to the wrapped command.

## Recommended Workflow

1. Open a fresh terminal tab.
2. Start the AI CLI through netbar.
3. Keep working normally.
4. Watch the bottom row for network state changes.

## What netbar Cannot Do

`netbar` cannot attach underneath an already-running Codex session. Start the session through netbar from the beginning.
