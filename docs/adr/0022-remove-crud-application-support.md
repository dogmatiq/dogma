# 22. Remove CRUD application support

Date: 2025-06-13

## Status

Accepted

- Referenced By [23. Message Order Guarantees](0023-message-order-guarantees.md)

## Context

Dogma has historically supported both CRUD-style and event-sourced applications.
However, in practice, all commerical use that we're aware of has been
for event-sourced applications.

## Decision

We will remove support for CRUD-style applications, positioning Dogma as an
exclusively event sourcing framework.

Some points in support of this decision:

- There's no official engine implementation that is event-sourced.
- Supporting multiple paradigms weakens the guarantees of the API.
- The maintainers are more interested in event sourcing applications.

Since the API already supports event-sourced applications well, this can be
implemented as a pure conceptual and documentation change initially.

## Consequences

- Supporting only event sourcing may mean only one engine implementation.
- Dogma's value proposition becomes more focused and easier to communicate.
- The API can be tailored specifically for event sourcing without compromise.
  For example, We can specify behavioral guarantees much more explicitly,
  particularly around message delivery order.

This decision also opens the door to API improvements:

- Removal of `AggregateCommandScope.Destroy()`. The concept of destroying
  something with immutable history has already been confusing in practice, and
  represents an irreversible action that is contrary to the principles of
  event sourcing.
- Simpler projection progress tracking. The `ProjectionMessageHandler` interface
  "resource versioning system" could be modeled directly as event stream
  offsets.
- Richer control over event consumption. Event-consuming handlers could grow
  features for specifying whether they are interested in historical events,
  since we know they are available, or only new events.

### Migration path

Although unlikely to exist in the wild, existing applications that were built
with CRUD patterns in mind will need to be refactored to use event sourcing
patterns. However, since the current API already supports event sourcing, most
applications should already be compatible or easily adaptable, as this was the
original intent of the abstraction.
