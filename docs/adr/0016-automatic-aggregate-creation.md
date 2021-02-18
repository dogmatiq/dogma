# 16. Automatic Aggregate Creation

Date: 2020-11-02

## Status

Accepted

- Amends [3. Aggregate Lifetime Control](0003-aggregate-lifetime-control.md)

## Context

After implementing several real engines and business domains we have found the
rules governing how `AggregateCommandScope.Create()`, `Root()`, `RecordEvent()`
and `Destroy()` can be called in relationship to one another tends to lead to
overly cumbersome implementations.

### Effect on handler implementations

From the handler implementors perspective there are numerous subtle interactions
between the aggregate instance's "state of existence" and the domain logic.

Perhaps most egregiously is the requirement that an instance already exist
before calling `Root()`. This means that the state of the aggregate root itself
can not be used to determine which events should be recorded without first
creating the instance. A less obvious problem is that when the domain logic is
such that `RecordEvent()` is not called within the same scope as a call to
`Create()` the engine must panic.

The net result of this is that the handler becomes equally concerned with the
notion of whether the instance exists or not than it is with the actual business
logic. The subtle interactions between the domain implementation (especially
when located within methods of the aggregate root) and the handler itself are
difficult to glean by reading the code, and hard to reason about.

In summary, the semantics of `Create()` do not help the handler implementor to
implement their domain logic in a clear and concise way.

### Effect on engine implementations

From the engine implementors perspective it seems that a non-trivial amount of
validation logic needs to be implemented within the scope to verify that the
application code is using the scope correctly, even though the engine does not
really benefit from these requirements.

For example, event sourced engine implementations tend to call
`AggregateMessageHandler.New()` to construct a new instance before applying
historical events via `ApplyEvent()`. This means that the root instance is
constructed in memory before `HandleCommand()` is even invoked. The requirement
that `Root()` panic if the instance does not exist does not save the engine from
constructing the in-memory root value.

## Decision

- Remove `AggregateCommandScope.Create()`.
- Remove `AggregateCommandScope.Exists()`.
- Allow `AggregateCommandScope.Root()` to be called at any time.
- Allow `AggregateCommandScope.Destroy()` to be called at any time.
- Remove `StatelessAggregateRoot` and `StatelessAggregateBehavior`. With the
  notion of "existence" being removed from the public API a stateless aggregate
  becomes nonsensical.
- Reinstate the hard requirement that the handlers MUST panic with
  `UnexpectedMessage` when asked to handle a message type that was not
  configured as being consumed by that handler. Removing the requirement to call
  `Create()` should simplify the dispatching logic sufficiently such that no
  extra logic is required to produce the panic.

In essence, the aggregate instance is "automatically created" the first time an
event is recorded.

## Consequences

Largely, this should simplify implementations of both `AggregateMessageHandler`
by application developers and `AggregateCommandScope` by engine developers.

Handler implementations should become more clearly based purely on domain logic.
Note that along with the removal of `StatelessAggregateRoot`, some aggregates
implementations will require an additional type declaration for their root
value.

The removal of `Create()` does mean that engines no longer get an explicit
request to create an instance. However, since all state changes must be done by
recording event we expect this to be a fairly trivial change to existing engine
implementations.
