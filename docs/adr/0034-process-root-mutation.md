# 34. Process root mutation

Date: 2026-05-03

## Status

Proposed

> [!NOTE]
> This decision has not yet been accepted and is subject to change.

## Context

A process message handler maintains a `ProcessRoot` for each active process
instance, accumulating state across multiple event and timeout messages. The
engine persists each instance's root between invocations — by marshaling it with
`MarshalBinary()` after each successful `HandleEvent()` or `HandleTimeout()`
call — so that later invocations start from the accumulated state rather than
from scratch.

Currently, the handler mutates the root directly. `HandleEvent()` and
`HandleTimeout()` both receive the `ProcessRoot` and may update it freely by
calling methods on it or assigning to its fields.

This means the engine cannot determine whether a given invocation actually
changed the root without resorting to a potentially expensive deep equality
check against some snapshot of the previous state. The engine must marshal and
persist the root after every successful invocation — even those that only
executed commands, scheduled timeouts, or inspected state without making any
changes.

Aggregates do not have this problem; every call to `RecordEvent()` gives the
engine an unambiguous signal that the aggregate instance's state has changed.

## Decision

We will add a `Mutate()` method to `ProcessScope`. It accepts a single
callback that is passed a mutable reference to the root:

```go
Mutate(fn func(mut R))
```

The `ProcessRoot` parameter `r` of `HandleEvent()` and `HandleTimeout()` becomes
read-only from the handler's perspective. Handlers may no longer mutate `r`
directly; instead, they must call `Mutate()` with a callback that performs the
mutation against the `mut` parameter.

The engine may elide the call to `fn` in some circumstances — for example, when
the root is a `StatelessProcessRoot` and persistence is known to be a
no-op. Because `fn` is not guaranteed to execute, it must not perform any side
effects beyond mutating `mut`. In particular:

- It must not call any methods on `s`, the `ProcessScope`.
- It must not use `r`, even for read-only purposes.

Since the call to `fn` may be elided, all conditional logic — deciding whether
and how to mutate — should be evaluated before calling `Mutate()`. The callback
should express an unconditional change; however, it is acceptable for the net
effect of `fn` to be a no-op — for example, setting a field to the value it
already holds.

The engine must ensure that `r` reflects the post-mutation state after `fn`
returns, so that the handler can continue to read from it.

`Mutate()` may be called more than once within a single `HandleEvent()` or
`HandleTimeout()` invocation. Each call is an independent mutation window.

Calling `Mutate()` after `End()` panics, mirroring the behavior of
`ExecuteCommand()` and `ScheduleTimeout()`.

### Dismissed alternatives

We considered a zero-argument `s.Mutate()` — a dirty flag the handler sets after
mutating `r` directly, as it does now. This would allow the engine to skip
persisting the root when the flag is not set, but the mutation and the signal
remain decoupled, which invites mistakes.

## Consequences

The engine gains a reliable signal for avoiding unnecessary persistence: if
`Mutate()` was not called, the root does not need to be persisted for that
invocation, though commands and timeout messages produced in the same scope are
still dispatched. The engine also gains a natural observation point around each
mutation — it can capture the root's state before and after the callback — which
enables structured telemetry such as state-transition traces and audit logs.

This is a breaking change to `ProcessScope`. Engine implementations must add
`Mutate()` to their scope implementations, and application code must be updated
to call `Mutate()` instead of mutating `r` directly.

Existing handler code that mutates the root directly continues to compile, but
violates the updated contract. `testkit` may detect such violations by comparing
the root's serialized state before and after each invocation.

The relationship between `r` and `mut` is unspecified — they may be clones or
refer to the same memory. Because `r` is off-limits while `fn` is executing, the
engine may safely pass `r` to `fn`.

<!-- references -->
