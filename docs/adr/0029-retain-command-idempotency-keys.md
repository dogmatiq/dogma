# 29. Retain command idempotency keys

Date: 2026-04-03

## Status

Proposed

> [!NOTE]
> This decision has not yet been accepted and is subject to change.

## Context

The `WithIdempotencyKey()` option allows users of `CommandExecutor` to
optionally attach a "key" to each command to prevent duplicate executions.

The idempotency key feature was an attempt to encode operation identity, a
potentially deep domain concept, into an opaque value so an engine could treat
repeated submissions of the same business operation as a single command.

We considered removing this feature because it creates pressure for engine
implementations to route all commands with an idempotency key through a single
shared pipeline or queue so they can provide global duplicate-detection
behavior. This constraint would limit horizontal scalability and reduce
opportunities to optimize command throughput.

Additionally, the engine cannot define business-level duplicate semantics
reliably. Those semantics belong in application domain logic, where each use
case can define what makes two operations equivalent.

## Decision

We will retain the `WithIdempotencyKey()` option and `IdempotencyKeyOption` in
the API.

Most engine implementations already have some notion of a "command ID". Such
engines can efficiently support idempotency keys by simply using the idempotency
key as the command ID. Engines that use UUIDs for command IDs can simply derive
a UUIDv5 from the idempotency key. This approach avoids the need for engines to
implement a shared command pipeline or global routing constraints to support
idempotency keys, while still allowing applications to express their intent to
deduplicate command submissions at the engine level.

This approach preserves the feature's utility while allowing each engine
implementation to choose its scalability strategy independently.

## Consequences

Applications can rely on `WithIdempotencyKey()` to express their intent to
deduplicate command submissions at the engine level, without forcing all engines
to implement a single shared command pipeline.

Engine implementations remain free to architect their command ingestion and
routing as they see fit. Engines built on [`envelopepb.Envelope`] already have
the infrastructure to support idempotency keys through command ID derivation.

[`envelopepb.Envelope`]: https://pkg.go.dev/github.com/dogmatiq/enginekit/protobuf/envelopepb#Envelope
