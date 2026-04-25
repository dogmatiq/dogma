# 33. Generic aggregate and process handlers

Date: 2026-04-25

## Status

Proposed

- References [6. Stateless aggregates and processes][ADR-6]

> [!NOTE]
> This decision has not yet been accepted and is subject to change.

## Context

`AggregateMessageHandler` and `ProcessMessageHandler` both operate on
application-defined root types that implement the `AggregateRoot` and
`ProcessRoot` interfaces. Today, handler methods accept and return the root as
the interface type rather than the application's concrete type.

This means every handler implementation must type-assert the root it receives.
For example, `AggregateMessageHandler.HandleCommand()` receives an
`AggregateRoot` that the handler must assert to its concrete type before use.
The compiler cannot verify these assertions — a mismatch between the concrete
type that `New()` returns and the type that `HandleCommand()` asserts is a
runtime panic, not a compile-time error.

## Decision

We will add a type parameter `R` to the following interfaces, constrained by
`AggregateRoot` or `ProcessRoot` as appropriate:

- `AggregateMessageHandler`
- `AggregateCommandScope`
- `ProcessMessageHandler`
- `ProcessScope`
- `ProcessEventScope`
- `ProcessTimeoutScope`

Engine implementations are unaware of the application's concrete root types, so
they need a way to represent handlers without that information. Rather than
introducing separate "untyped" interfaces, we will use the same generic types
with the constraint interface as the type parameter. For example, an
application defines `AggregateMessageHandler[ShoppingCart]`, but the engine
represents the same handler as `AggregateMessageHandler[AggregateRoot]`.

### Dismissed alternatives

**Separate untyped interfaces.** We considered introducing a parallel set of
interfaces — `UntypedAggregateMessageHandler`, `UntypedAggregateCommandScope`,
and their process equivalents — for engines to use instead of the generic
versions. Separate interfaces would make the untyped nature explicit in type
names, but this would double the number of interfaces with no practical benefit,
since the scope interfaces would be identical to their generic counterparts and
the handler interfaces would still need adaptors.

## Consequences

Each handler's `New` method returns `R`, and its message-handling methods
receive `R` directly. The compiler enforces that the root type is consistent
across all of a handler's methods, so type assertions are no longer necessary
unless the handler deliberately uses the interface type as its type parameter.

This is a breaking change to the handler and scope interfaces. Existing code
must add type parameters to references to these types. Handlers that do not
want to adopt a concrete root type can use `AggregateRoot` or `ProcessRoot` as
the type parameter to preserve their current behavior — the generic interfaces
work without a concrete root type.

`NoTimeoutMessagesBehavior` becomes parameterized by `R` because its
`HandleTimeout()` method receives a `ProcessTimeoutScope[R]`.

`StatelessProcessRoot` changes from a package-level variable to an exported
struct type. The current design — a variable holding an interface value — cannot
serve as a type parameter. [ADR-6] established the stateless root pattern, and
the struct form preserves the same semantics while being usable as a type
argument.

Integration and projection handlers are unaffected by this change. They do not
have root types and remain non-generic.

No current methods on the scope interfaces use `R`, but parameterizing them
means methods that use `R` can be added later without a breaking change.

<!-- references -->

[ADR-6]: 0006-stateless-aggregates-and-processes.md
