# 30. Observe command outcomes via events

Date: 2026-04-03

## Status

Accepted

- References [23. Message order guarantees](0023-message-order-guarantees.md)

## Context

Dogma is natively asynchronous. `CommandExecutor.ExecuteCommand()` returns when
the engine has taken ownership of a command. It does not block until the command
is handled.

That model is correct for Dogma's event-driven architecture, but applications
often need to expose synchronous APIs at their boundaries, which means observing
the effects of the commands they execute.

The current API provides no mechanism for this. Applications must resort to
polling a projection, which is a significant amount of boilerplate and
infrastructural overhead for a common and transient need.

Dogma should provide a first-class mechanism for this pattern without
introducing request-response semantics that hide its asynchronous nature: not
every command produces events, and a command may never lead to a particular
event.

## Decision

We will add a new generic option to `ExecuteCommand()`:

```go
func WithEventObserver[T Event](EventObserver[T]) ExecuteCommandOption { /* ... */ }
type EventObserver[T Event] func(ctx context.Context, event T) (satisfied bool, err error)
```

The type parameter `T` constrains which events are observed. The engine calls
the observer once for each event of type `T` produced while handling the
submitted command. The event does not have to be recorded directly by the
handler that executes the command; it may be recorded after a chain of further
commands and processes have run.

Multiple `WithEventObserver` options may be passed, each with its own type
parameter and observer function.

Without `WithEventObserver`, `ExecuteCommand()` retains its existing
fire-and-forget behavior. With it, `ExecuteCommand()` blocks until one of the
following occurs:

- Any observer reports completion by returning `satisfied == true`.
- Any observer returns a non-nil error.
- The caller's context is canceled or reaches its deadline.
- The engine determines that no further relevant events can occur.

We considered modeling this feature explicitly as a projection, because the
behavior is conceptually similar to building a temporary read model from events.
However, existing projection APIs are about long-lived handlers and persistent
state.

We also considered adding a new `CommandExecutor.ExecuteCommandSync()` method,
but ruled this approach out as it moves closer to the request-response semantics
we sought to avoid.

## Consequences

Applications gain a standard way to build synchronous boundaries on top of
Dogma without introducing request-response messaging into the API.

The command execution API remains unified. Callers choose whether to treat
`ExecuteCommand()` as fire-and-forget or as a bounded wait by passing
`WithEventObserver`. Adding or removing the option does not change the
fundamental nature of what `ExecuteCommand()` does to the application's state.

The ordering of observed events remains governed by
[23. Message order guarantees](0023-message-order-guarantees.md).

If the engine determines that no further relevant events can occur and no
observer returned `satisfied == true`, `ExecuteCommand()` returns
`ErrObserverNotSatisfied`. This applies whether the engine reached that
conclusion dynamically, after draining the causal work, or statically, by
inspecting the handler graph before any work runs.
