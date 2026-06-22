// Package aof provides append-only-file persistence for the Redix store.
//
// The AOF layer acts as a durable command log that records all mutating
// operations executed against the in-memory store. It enables recovery of
// state by replaying commands in order of execution.
//
// Design goals:
//
//   - Decouple persistence from storage implementation.
//   - Record normalized commands rather than raw protocol frames.
//   - Allow deterministic rebuild of state through replay.
//
// Architecture overview:
//
//  1. Write path:
//     - Commands are executed through the resp/handler layer.
//     - Each executed command is registered via AOF.Register.
//     - The AOF layer serializes and appends commands to disk.
//
//  2. Rewrite path (compaction):
//     - The current in-memory state (store.MemStore) is used as source of truth.
//     - A full snapshot is generated and replaces the existing AOF file.
//     - This prevents unbounded growth of the append-only log.
//
//  3. Recovery path:
//     - On startup, the AOF file is read sequentially.
//     - Each stored command is parsed and re-executed through the resp/handler.
//     - The store is rebuilt deterministically from persisted operations.
//
// Constraints:
//
//   - AOF does not parse raw RESP directly; it operates on normalized commands.
//   - The store remains agnostic of persistence concerns.
//   - The protocol layer is only used for encoding/decoding when required,
//     not for persistence semantics.
//
// This package is intended as a bridge between the execution layer (resp)
// and durable storage, enabling Redis-like persistence semantics.
package aof
