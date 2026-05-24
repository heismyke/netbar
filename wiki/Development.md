# Development

## Clone

```sh
git clone https://github.com/heismyke/netbar.git
cd netbar
```

## Test

```sh
go test ./...
```

## Vet

```sh
go vet ./...
```

## Build

```sh
go build ./...
```

## Run from Source

```sh
go run ./cmd/netbar -once
go run ./cmd/netbar -- codex
```

## Release Checklist

1. Run `go test ./...`.
2. Run `go vet ./...`.
3. Run `go build ./...`.
4. Update docs if user-facing behavior changed.
5. Create a Git tag.
6. Push the tag.
7. Create a GitHub release.

## Contribution Guidelines

- Keep changes focused.
- Add tests for behavior changes.
- Update documentation for new flags or modes.
- Avoid committing local binaries or environment files.
