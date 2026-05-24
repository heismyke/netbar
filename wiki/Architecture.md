# Architecture

`netbar` has three core pieces:

## Connectivity Monitor

The monitor performs:

- DNS resolution check.
- TCP dial check.
- Latency measurement.
- Status calculation.

Statuses:

- `online`
- `degraded`
- `offline`

## Session Wrapper

Session mode starts a command inside a PTY and draws the network status on the bottom terminal row.

Responsibilities:

- Start the wrapped command.
- Forward input and output.
- Reserve terminal space.
- Redraw status updates.
- Handle terminal resize events.
- Restore the terminal on exit.

## Status Persistence

`netbar` stores the previous status in the user cache directory.

This enables transition labels such as:

```text
Back online
```

without needing a long-running daemon.

## Modes

- Session mode: default for interactive terminals.
- Stream mode: continuous line output.
- One-shot mode: one check and exit.
- tmux format: one-shot output with tmux styling.
