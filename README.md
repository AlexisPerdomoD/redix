# Redix

A lightweight, Redis-protocol-compatible server built from scratch in Go.

> _"What I cannot create, I do not understand."_ — Richard Feynman

Redix is a from-scratch reimplementation of a Redis-compatible server that prioritizes clean, didactic code without sacrificing performance. The project aims to deepen practical knowledge of Go's standard library — networking, concurrency, and systems programming — while producing a lightweight but production-viable alternative to a conventional Redis server. Zero external dependencies.

## Goals

- **Full RESP2/RESP3 protocol** — Correctly parse and serialize the Redis wire protocol.
- **Complete command set** — Core commands (`GET`, `SET`, `DEL`, expiration), plus Redis data structures: strings, lists, hashes, sets, sorted sets.
- **Production-grade persistence** — Append-only file (AOF) and snapshotting (RDB).
- **Concurrent by design** — Goroutine-based connection handling with proper synchronization from day one.
- **Performance parity** — Throughput and latency competitive with Redis for the implemented feature set.

## Design

- **Zero external dependencies.** Everything is built on Go's standard library.
- **Performance-oriented.** Concurrency model and data structures designed for throughput from the start.
- **Didactic code.** Readability and clarity are first-class concerns — the codebase should be teachable without being naive.
- **Progressive implementation.** Each component is built and understood in isolation before integration.

## Getting started

### Prerequisites

- Go 1.26.3 or later

### Build & run

```bash
make build   # compile to build/redix
make run     # build and start the server
```

The server listens on `:6379` (default Redis port). Test it with any Redis client:

```bash
redis-cli PING
# +PONG
```

### Other commands

| Command      | Action                        |
| ------------ | ----------------------------- |
| `make test`  | Run all tests                 |
| `make lint`  | Static analysis with `go vet` |
| `make fmt`   | Format code with `go fmt`     |
| `make clean` | Remove compiled binary        |

## Project structure

```
cmd/redix/main.go          Entrypoint
internal/
├── server/                TCP listener, connection state machine, accept loop
├── protocol/              RESP type definitions, parser, and wire writers
├── resp/                  Command dispatch and handler implementations
└── store/                 In-memory key-value store
```

Internal packages are progressively being built out, bottom-up, starting from the protocol layer.

## License

MIT
