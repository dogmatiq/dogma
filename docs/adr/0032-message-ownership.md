# 32. Message ownership

Date: 2026-04-25

## Status

Proposed

> [!NOTE]
> This decision has not yet been accepted and is subject to change.

## Context

Dogma messages are Go interface values backed by pointer-to-struct types.
These structs often contain slices, maps, and other reference types that share
underlying memory when copied by value.

When the engine passes a message to a handler — or when a handler passes a
message to a scope method — the two sides hold references to the same object
graph. If either side mutates the message or any value reachable from it, the
other side observes the change. This creates a class of subtle, hard-to-diagnose
bugs.

The problem is most pronounced in `AggregateRoot.ApplyEvent`. During replay the
engine decodes each historical event from storage, producing an independent
value. During live handling, however, the engine calls `ApplyEvent` with the
same in-memory value the handler passed to `RecordEvent`. Code that assigns a
field from the event into the root's state works correctly during replay — the
decoded event shares no memory with anything else — but silently aliases memory
during live handling. The root and the engine's copy of the event end up
pointing at the same slice or map, and a later mutation to one corrupts the
other.

## Decision

We will establish a message ownership rule: whenever a message value crosses the
boundary between the engine and the application, the application owns the
message. Ownership covers the entire object graph reachable from the message,
not just the top-level value.

The engine is responsible for ensuring that every message value it passes to
application code is independent of any value the engine retains, and that every
message value it receives from application code is not affected by subsequent
use of that value by the application. The application never needs to clone or
otherwise protect a message value — it can freely read, mutate, or retain any
message it holds.

This applies at every boundary where messages cross between the engine and the
application:

- Any handler method that receives a message from the engine, including
  `AggregateRoot.ApplyEvent()` — even during live handling, where the event must
  be independent of the value the handler passed to `RecordEvent()`.
- Any scope method that accepts a message from the handler.
- `CommandExecutor.ExecuteCommand()`, where external application code dispatches
  a command into the engine.
- Event observer callbacks registered via `WithEventObserver()`.

This ADR does not prescribe how the engine achieves independence. An engine
might marshal the message and decode it later, perform a deep copy, or use any
other mechanism — including doing nothing at all if it can prove the message is
not shared. The contract is defined in terms of the observable guarantee, not
the implementation.

## Consequences

Application developers can write handler and root logic without considering
whether a message value is shared with the engine. In particular, `ApplyEvent`
implementations can assign fields from the event directly into root state
regardless of whether the call occurs during replay or live handling.

The engine bears the full cost of ensuring independence. In practice, engines
are likely to achieve this by marshaling each message to its binary
representation immediately and releasing the in-memory value — but this is an
implementation choice, not a requirement.

<!-- references -->
