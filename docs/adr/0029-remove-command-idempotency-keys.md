# 29. Remove command idempotency keys

Date: 2026-04-03

## Status

Accepted

## Context

The `WithIdempotencyKey()` option allows users of `CommandExecutor` to
optionally attach a "key" to each command to prevent duplicate executions.

The idempotency key feature was an attempt to encode operation identity, a
potentially deep domain concept, into an opaque value so an engine could treat
repeated submissions of the same business operation as a single command.

This requirement creates pressure for engine implementations to route all
commands with an idempotency key through a single shared pipeline or queue so
they can provide global duplicate-detection behavior. This constraint limits
horizontal scalability and reduces opportunities to optimize command throughput.

The engine can observe command submissions, but it cannot define business-level
duplicate semantics reliably. Those semantics belong in application domain
logic, where each use case can define what makes two operations equivalent.

## Decision

We will remove `WithIdempotencyKey()` and `IdempotencyKeyOption` from the API.

## Consequences

Engines are no longer required to implement command deduplication semantics at
submission time. This removes pressure toward a single command pipeline and
allows engine implementations to use more scalable command ingestion and routing
strategies.

Applications that currently rely on command idempotency keys must move
idempotency checks into business logic, such as by modeling operation identity
within aggregates or processes.

The command execution API surface is simpler, with fewer options that imply
cross-cutting guarantees engines cannot provide efficiently.

This decision does not affect the existing guarantee that a command's
side-effects occur exactly once, even if the engine invokes the handler more
than once internally.
