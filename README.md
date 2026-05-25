# Redix

A lightweight, Redis-protocol-compatible server built from scratch in Go.

> _"What I cannot create, I do not understand."_ — Richard Feynman

Redix is a from-scratch reimplementation of a Redis-compatible server that prioritizes clean, didactic code without sacrificing performance. The project aims to deepen practical knowledge of Go's standard library — networking, concurrency, and systems programming — while producing a lightweight but production-viable alternative to a conventional Redis server. Zero external dependencies.

## Status

Early development. The current codebase accepts TCP connections and echoes bytes back. The RESP protocol parser, command router, and in-memory store are yet to be built.

## Features (planned)

- **Wire protocol** — Full Redis Serialization Protocol (RESP2/RESP3)
- **Core commands** — `GET`, `SET`, `DEL`, key expiration
- **Data structures** — Strings, lists, hashes, sets, sorted sets
- **Persistence** — Append-only file (AOF) and snapshotting (RDB)
- **Concurrency** — Event-loop or goroutine-based connection handling

## Getting started

### Prerequisites

- Go 1.26.3 or later

### Build & run

```bash
make build   # compile to build/redix
make run     # build and start the server
```

The server listens on `:6379` (default Redis port).

### Other commands

| Command      | Action                        |
| ------------ | ----------------------------- |
| `make test`  | Run all tests                 |
| `make lint`  | Static analysis with `go vet` |
| `make fmt`   | Format code with `go fmt`     |
| `make clean` | Remove compiled binary        |

## Design

- **Zero external dependencies.** Everything is built on Go's standard library.
- **Performance-oriented.** Concurrency model and data structures designed for throughput from the start.
- **Didactic code.** Readability and clarity are first-class concerns — the codebase should be teachable without being naive.
- **Progressive implementation.** Each component is built and understood in isolation before integration.

## Project structure

```
cmd/redix/main.go          Entrypoint
internal/
├── server/                TCP listener and connection state machine
├── protocol/              Command parsing and routing     (stub)
├── resp/                  RESP wire format                (stub)
└── store/                 In-memory key-value store       (stub)
```

## License

MIT
