# 22. Remove CRUD Application Support

Date: 2025-06-13

## Status

Accepted

## Context

Dogma has historically supported both CRUD-style and event-sourced applications.
However, after careful consideration, the decision has been made to drop support
for CRUD applications, making Dogma an exclusively event-sourcing abstraction.

Several factors support this decision:

- We have never written a production CRUD engine (except for the `testkit`
  in-memory engine, which now has event-sourcing-like semantics anyway)
- Supporting only event-sourcing may mean we only need a single engine
  implementation, simplifying the Dogma ecosystem considerably
- Supporting multiple paradigms weakens the API by requiring compromises
  that serve both use cases rather than optimizing for one
- The development team has stronger interest and expertise in event-sourcing
  applications

## Decision

We will remove conceptual and documentation support for CRUD-style applications
from Dogma, positioning it as an exclusively event-sourcing framework.

Since the API already supports event-sourced applications well, this can be
implemented as a pure conceptual and documentation change initially.

## Consequences

### Positive Consequences

- **Simplified ecosystem**: With only event-sourcing to support, we can focus
  on a single engine implementation and optimization path
- **Stronger API**: The API can be optimized specifically for event-sourcing
  without compromises for CRUD patterns
- **Clearer messaging**: Dogma's value proposition becomes more focused and
  easier to communicate
- **Better guarantees**: We can specify behavioral guarantees much more
  explicitly, particularly around message delivery order

### Future API Improvements

This decision opens the door to several API improvements:

- **Remove `AggregateCommandScope.Destroy()`**: The concept of "destroying"
  something with immutable history is flawed and represents an irreversible
  action that dilutes some benefits of event-sourcing
- **Simplify projection versioning**: The `ProjectionMessageHandler` interface's
  resource versioning system could be simplified and modeled directly as event
  stream offsets
- **Richer event consumption**: Event-consuming handlers could grow features for
  specifying whether they are interested in historical events (since we know
  they are available) or only new events

### Message Delivery Order Guarantees

With exclusive focus on event-sourcing, Dogma can adopt explicit message
delivery order guarantees:

1. Commands are delivered in an undefined order
2. Events produced by integration handlers are observed in an undefined order
3. Events produced by a single aggregate **instance** are observed in the order
   they were recorded, relative to each other
4. Events produced by different instances of the same aggregate type, or by
   different aggregate types, have no relative order guarantees
5. Timeouts produced by a single process instance are observed by the relative
   order of their "scheduled for" time
6. Timeouts produced by different instances of the same process type, or by
   different process types, have no relative order guarantees (though in the
   common case they will be delivered as close to their "scheduled for" time
   as possible)

### Migration Path

Existing applications that were built with CRUD patterns in mind will need to
be refactored to use event-sourcing patterns. However, since the current API
already supports event-sourcing well, most applications should already be
compatible or easily adaptable.