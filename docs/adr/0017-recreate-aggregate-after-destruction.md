# 17. Recreation of Aggregate Instances After Destruction

Date: 2020-11-05

## Status

Accepted

- Amends [3. Aggregate Lifetime Control](0003-aggregate-lifetime-control.md)

## Context

After implementing several real business domains we have found that it can be
difficult to use `AggregateCommandScope.Destroy()` effectively.

Ideally, `Destroy()` could be called after recording any event that effectively
"resets" the aggregate's state, from the perspective of the business logic.

In practice, logic that results in a call to `Destroy()` may be followed by some
conditional logic that records a new event. In existing implementations this
results in a panic from the engine, though as of ADR-16, which removed
`Create()` this is no longer documented as required engine behavior.

## Decision

Remove the requirement that any call to `Destroy()` be preceeded by a call to
`RecordEvent()` within the same scope.

Calling `RecordEvent()` *after* `Destroy()` event should "negate" the call to
`Destroy()`, as though it were never called.

## Consequences

The complex interplay between `RecordEvent()` and `Destroy()` is removed,
allowing each to be understood and used in isolation.

Changes to existing engine implementations should be minimal, as they already
handle destruction by setting an in-memory flag that is checked after
`HandleCommand()` is invoked. Likely they can be changed to simply unset that
flag in `RecordEvent()`.

There is one subtle behavior that results from these changes, which is that
after calling `Destroy()` the `AggregateRoot` returned by
`AggregateCommandScope.Root()` is not "reset", as there is no mechanism to do
so. This potentially differs to how the root will appear when the next command
is received, when it will be equal to `AggregateMessageHandler.New()`.
