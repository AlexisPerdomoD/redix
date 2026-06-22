# Redix

[![Go Version](https://img.shields.io/badge/Go-1.26.3-00ADD8?logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

> *"What I cannot create, I do not understand."* — Richard Feynman

Redix is a from-scratch Redis-protocol-compatible server built entirely on Go's standard library. It implements the RESP2 wire protocol with a concurrent connection handler, a thread-safe in-memory key-value store with TTL semantics and automatic expiration, and a modular append-only file (AOF) persistence layer. Zero external dependencies.

## Design rationale

Redis is one of the most influential pieces of infrastructure software of the past decade. Building a compatible server from scratch demonstrates how production-grade data systems work — from wire protocol design and concurrent I/O to memory management and persistence semantics.

Redix reflects several deliberate engineering decisions:

- **Protocol correctness.** A faithful implementation of the RESP2 serialization format, handling bulk strings, arrays, errors, integers, and null values as specified by the Redis protocol.
- **Concurrency from day one.** Goroutine-based accept loop with per-connection state machines; all store access serialized through `sync.RWMutex`.
- **Decoupled persistence.** The AOF module is wholly separate from the store. The store has no knowledge of persistence, enabling alternative strategies (RDB, replication) without storage-layer changes.
- **Progressive maturity.** Each layer — protocol parsing, command dispatch, storage, persistence — is independently testable and replaceable.

## Supported commands

| Command          | Status | Notes                           |
| ---------------- | ------ | ------------------------------- |
| `PING`           | ✓      | With optional argument echo     |
| `SET`            | ✓      | With default TTL                |
| `GET`            | ✓      |                                 |
| `DEL`            | ✓      | Variadic                        |
| `HSET`           | ✓      |                                 |
| `HGET`           | ✓      |                                 |
| `HDEL`           | ✓      | Variadic                        |
| `EXISTS`         | ✓      | Variadic                        |
| `EXPIRE`         | ✓      |                                 |
| `TTL`            | ✓      | Returns remaining seconds       |

## Architecture

```
┌──────────┐   ┌──────────────┐   ┌────────────┐   ┌───────────┐
│  Client  │──▶│   Server     │──▶│  Protocol  │──▶│  Handler  │
│  (TCP)   │   │  (internal/  │   │  (internal/ │   │ (internal/│
│          │   │   server)    │   │   protocol) │   │   resp)   │
└──────────┘   └──────────────┘   └────────────┘   └─────┬─────┘
                                                          │
                                    ┌─────────────────────┼──────────────────┐
                                    ▼                     ▼                  │
                            ┌──────────────┐     ┌──────────────┐           │
                            │    Store     │     │     AOF      │           │
                            │  (internal/  │     │  (internal/  │           │
                            │   store)     │     │    aof)      │           │
                            └──────┬───────┘     └──────┬───────┘           │
                                   │                     │                   │
                                   └─────────────────────┼───────────────────┘
                                                         │
                                              (decoupled — store is
                                               persistence-agnostic)
```

### Layer responsibilities

| Layer | Package | Responsibility |
|-------|---------|---------------|
| Server | `internal/server/` | TCP listener, connection state machine, accept loop. Each connection runs in its own goroutine. |
| Protocol | `internal/protocol/` | RESP2 type definitions, parser, and serializer. Handles byte-level wire encoding. |
| Command | `internal/resp/` | Command dispatch and handler implementations. Translates parsed RESP frames into store operations. |
| Store | `internal/store/` | Thread-safe in-memory key-value store. Supports strings, hashes, TTL, and automatic expired-key eviction via a background goroutine. |
| Persistence | `internal/aof/` | Append-only file persistence. Fully decoupled: the store has no awareness of persistence semantics. |

## Technical highlights

- **Concurrent I/O.** Goroutine-per-connection model with non-blocking accept. Per-connection read deadlines prevent resource exhaustion from idle clients.
- **Thread-safe storage.** `sync.RWMutex` guards all store operations. Readers never block writers; writes serialize only the critical section (hash-map mutation), keeping lock contention minimal for read-heavy workloads.
- **TTL with automatic eviction.** Every key carries an absolute Unix expiration timestamp. A background goroutine scans the entire keyspace every 30 seconds, collecting expired candidates under read lock and revalidating them under write lock before removal.
- **RESP2 wire protocol.** Hand-written parser covering bulk strings, simple strings, errors, integers, arrays, and null bulk strings — no generated code, no reflection.
- **Decoupled persistence.** AOF is designed as an observer that records normalized commands after execution. The store remains persistence-agnostic, making it straightforward to add alternative persistence strategies (RDB snapshots, replication streams).
- **Zero dependencies.** Standard library only. No frameworks, no ORMs, no third-party packages.

## Getting started

### Prerequisites

- Go 1.26.3 or later

### Build & run

```bash
make build   # compile to build/redix
make run     # build and start the server
```

The server listens on `:6379` by default. Test it with any Redis client:

```bash
redis-cli PING
# +PONG

redis-cli SET foo bar
# +OK

redis-cli GET foo
# "bar"
```

### Configuration

Redix is configured through environment variables:

| Variable                        | Default     | Description                              |
| ------------------------------- | ----------- | ---------------------------------------- |
| `REDIX_PORT`                    | `6379`      | TCP port to listen on                    |
| `REDIX_CONNECTION_IDLE_TIMEOUT` | (none)      | Connection idle timeout (e.g. `5m`)      |
| `REDIX_LOG_LEVEL`               | `INFO`      | Log level: `DEBUG`, `INFO`, `WARN`, `ERROR` |

### Development commands

| Command      | Action                        |
| ------------ | ----------------------------- |
| `make test`  | Run all tests                 |
| `make lint`  | Static analysis with `go vet` |
| `make fmt`   | Format code with `go fmt`     |
| `make clean` | Remove compiled binary        |

## Project structure

```
.
├── cmd/
│   └── redix/main.go          Entrypoint
└── internal/
    ├── aof/                   Append-only file persistence
    ├── config/                Environment-based configuration
    ├── protocol/              RESP type definitions, parser, and serializer
    ├── resp/                  Command dispatch and handler implementations
    ├── server/                TCP listener, connection state machine, accept loop
    └── store/                 In-memory key-value store with TTL and hashes
```

## Roadmap

- AOF persistence: crash recovery and file compaction (rewrite)
- RDB snapshot support
- List, set, and sorted set data structures
- Pub/Sub

## Contributing

Contributions are welcome. Please open an issue or pull request on [GitHub](https://github.com/AlexisPerdomo/redix).

## License

MIT
