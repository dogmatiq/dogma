# 23. Message order guarantees

Date: 2025-06-14

## Status

Accepted

- References [22. Remove CRUD Application Support](0022-remove-crud-application-support.md)

## Context

Given the removal of CRUD application support, we are able to make more
specific guarantees about the order in which messages are delivered to an
application's handlers.

## Decision

We will update the API documentation to describe engine behavior under the
following guarantees:

1. Command order is **undefined**.
2. Events from the same aggregate instance preserve their **recorded order**.
3. Events from the same scope preserve their **recorded order**.
4. Relative order of events from different scopes is otherwise **undefined**.
5. Timeouts from the same process instance follow a weak total order by
   "scheduled for" time.
6. Relative order of timeouts from different instances of the same process type,
   or from different process types, is **undefined**.

These rules define the **observable ordering semantics** that the engine must
exhibit. They do not imply any specific behavior around concurrency, retry
semantics, or other implementation details.

## Consequences

The guarantees around command and timeout ordering are unchanged and should pose
no issues for existing applications.

In practice, engines attempt to handle commands immediately, providing an
approximate chronological order. However, they are free to defer and retry
commands as necessary.

The use of a weak total order for timeouts makes the ordering of timeouts with
the same "scheduled for" time explicitly undefined. This matches the behavior of
current engine implementations and clarifies what was previously underspecified.

The guarantees around event ordering have been strengthened. Previously, there
were no guarantees about the order of events, regardless of their source. Where
events are now guaranteed to be observed in recorded order, this does not imply
that their "recorded at" time is monotonic. It simply reflects the system clock
at the time of recording.

Even where the behavior is unchanged in practice, these rules allow us to
document the guarantees more explicitly, making it clearer to application
developers that handler logic must be designed around these constraints.
