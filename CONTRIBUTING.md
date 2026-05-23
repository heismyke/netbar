# Contributing

Thanks for helping improve `netbar`.

## Development

Requirements:

- Go 1.22 or newer

Run the test suite:

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

## Pull Requests

- Keep changes focused and easy to review.
- Add or update tests when behavior changes.
- Update `README.md` when user-facing flags, output, or install steps change.
- Do not commit local build output, editor files, or environment files.

## Commit Style

Use short imperative commit subjects, for example:

```text
Add tmux status formatter
Fix monitor shutdown race
```
