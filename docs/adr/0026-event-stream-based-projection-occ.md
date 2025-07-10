# 26. Event-stream based projection OCC

Date: 2025-07-11

## Status

Accepted

- References [22. Remove CRUD application support](0022-remove-crud-application-support.md)

## Context

The `ProjectionMessageHandler` interface uses a "resource versioning system" to
track which events have been applied to a projection. The `HandleEvent()` method
accepts a resource identifier, current version and a next version, and must only
apply the event if the current version matches the version stored in the
projection's "OCC store" for that resource.

The design of this system is deliberately abstract to allow for different engine
semantics. Specifically, it was built to allow projections to work with both
CRUD and event-sourced engines.

- CRUD engines, or engines without strictly ordered events, were expected to use
  event IDs as the resource, with an empty version for unhandled events and a
  non-empty version for handled events.
- Event-sourced engines were expected - and do - use the application or
  event-stream ID as the resource, the stream offset as the version.

The latter is the only implementation we have today, and with removal of CRUD
support in [ADR-22], it is the only behavior we expect going forward.

## Decision

We will change the `ProjectionMessageHandler` interface to use a more concrete
OCC API that refers specifically to event streams and offsets, rather than
generic resources and versions.

## Consequences

- The cognitive load of understanding the `ProjectionMessageHandler` interface
  is reduced, as it no longer requires understanding the abstract "resource
  versioning system".
- We can use more appropriate data types for the event stream identifiers and
  offsets, perhaps `string` and `uint64`. Currently resource identifiers and versions are
  both represented by opaque byte slices.
- We can build richer tooling that understands how to properly "reset" a
  projection to rebuild it from the start of the event streams.
- This decision is difficult to roll back. We _could_ leave the interface as
  is, and it will continue to work with event-sourced engines, but experience
  shows that it's been a source of confusion that we're probably better off
  eliminating.

[ADR-22]: 0022-remove-crud-application-support.md
