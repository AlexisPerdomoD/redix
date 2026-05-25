# Redix — Agent Guide

A from-scratch Redis-protocol-compatible server in Go. Zero external
dependencies.

## Commands

| Command      | Action                   |
| ------------ | ------------------------ |
| `make build` | Compile to `build/redix` |
| `make run`   | Build + run              |
| `make test`  | `go test ./...`          |
| `make lint`  | `go vet ./...`           |
| `make fmt`   | `go fmt ./...`           |

No tests exist yet. `make test` trivially passes. No CI.

## Commit convention

Use conventional commits:

```
type(scope): short imperative description
```

| Type       | When to use               |
| ---------- | -------------------------- |
| `feat`     | New feature or command     |
| `fix`      | Bug fix                    |
| `docs`     | Documentation              |
| `refactor` | Code restructuring         |
| `test`     | Adding or updating tests   |
| `chore`    | Build, config, tooling     |

Examples:
```
feat(server): implement RESP2 parser
chore: add Makefile and .gitignore
```

## Architecture

- `cmd/redix/main.go` — single entrypoint
- `internal/server/` — TCP listener + connection state machine (only real implementation)
- `internal/protocol/` — stub
- `internal/resp/` — stub
- `internal/store/` — stub

The three stubs above contain only package declarations. Any RESP parsing,
protocol handling, or storage work must build them out from scratch.

## Quirks

- Go 1.26.3. No `toolchain` directive. Stdlib only, `log/slog` for logging.
- Module: `github.com/AlexisPerdomo/redix` (note the `d` in Perdomo).
- Single commit on `main`, no tags, no branches.
