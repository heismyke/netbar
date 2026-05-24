# Installation

## Requirements

- Go 1.22 or newer.
- A terminal environment.

## Install

```sh
go install github.com/heismyke/netbar/cmd/netbar@latest
```

Make sure Go's binary directory is on your `PATH`:

```sh
export PATH="$PATH:$(go env GOPATH)/bin"
```

For zsh, add that line to `~/.zshrc`.

For bash, add it to `~/.bashrc` or `~/.bash_profile`.

## Upgrade

```sh
go install github.com/heismyke/netbar/cmd/netbar@latest
hash -r
```

`hash -r` clears your shell's command lookup cache so it sees the newest installed binary.

## Verify

```sh
netbar -h
netbar -once
```

The help output should include session-mode flags such as:

- `-stream`
- `-rows`
- `-cols`

## Install a Specific Version

```sh
go install github.com/heismyke/netbar/cmd/netbar@v0.2.1
```

## Common PATH Issue

If `netbar` is not found after installing:

```sh
echo "$(go env GOPATH)/bin"
```

Then add that directory to your `PATH`.
