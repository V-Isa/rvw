# Contributing

Thanks for helping improve `rvw`.

## Development

Run checks from the repository root:

```sh
GOCACHE=$PWD/.cache/go-build go test ./...
GOCACHE=$PWD/.cache/go-build go vet ./...
GOCACHE=$PWD/.cache/go-build go build ./...
GOCACHE=$PWD/.cache/go-build GOLANGCI_LINT_CACHE=$PWD/.cache/golangci-lint golangci-lint run ./...
GOCACHE=$PWD/.cache/go-build $(go env GOPATH)/bin/govulncheck ./...
```

Run `go mod tidy` after adding, removing, or updating dependencies.

## Code Conventions

Use Go's standard `got, want` style for tests:

```go
t.Fatalf("Result() = %v, want %v", got, want)
t.Fatalf("Run() error = %v, want nil", err)
```

Runtime errors should be lowercase, omit trailing punctuation, and wrap inspectable causes with `%w`:

```go
return fmt.Errorf("start command: %w", err)
```

Only ignore returned errors explicitly when the operation cannot fail in practice, such as writing to a `strings.Builder` or `bytes.Buffer`, or when a terminal status message is intentionally best-effort.

Manage direct dependencies deliberately. Let indirect dependencies be updated by `go mod tidy` unless a security issue, bug fix, or direct dependency upgrade requires a targeted update.

## Contribution License

By contributing to this project, you agree that your contributions are licensed under the same MIT License that covers the project.

Do not submit code, documentation, assets, or other material unless you have the right to contribute it under the MIT License.

## AI-Assisted Contributions

AI-assisted contributions are allowed when the contributor has reviewed, understood, and tested the change.

Contributors are responsible for the submitted code regardless of whether AI tools helped produce it. Do not submit AI-generated code that you cannot explain, do not trust, or suspect may contain copied proprietary or incompatible licensed material.

Maintainers may ask for clarification, request a simpler human-reviewed rewrite, or close low-confidence generated contributions.
