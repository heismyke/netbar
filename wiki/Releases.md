# Releases

## v0.2.1

Patch release for interactive session sizing.

Highlights:

- Passes accurate `COLUMNS` and `LINES` to wrapped commands.
- Resizes the PTY immediately after entering session mode.
- Adds `-rows` and `-cols` overrides for terminals that report incorrect dimensions.

## v0.2.0

Adds interactive session mode.

Highlights:

- `netbar` opens a shell by default in interactive terminals.
- `netbar -- codex` runs Codex inside a netbar-managed session.
- Bottom terminal row shows network status.
- `-stream` keeps continuous line output available.

## v0.1.0

Initial open-source release.

Highlights:

- Connectivity checks.
- One-shot output.
- tmux-compatible output.
- Tests and CI.

See all releases:

https://github.com/heismyke/netbar/releases
