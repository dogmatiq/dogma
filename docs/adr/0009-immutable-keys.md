# 9. Immutable Application and Handler Keys

Date: 2019-06-28

## Status

Accepted

- Amended by [20. Constraints on Identifier Values](0020-identifier-constraints.md)

## Context

Engine implementations require a mechanism for associating ancillary data with
Dogma applications and handlers.

For example, such data might include application state in the form of aggregate
and process roots, or historical events in an event sourcing system.

Currently, engine implementations rely on application and handler names as a key
for associated data. This is especially problematic for handlers as the name
initially chosen for a handler may become misleading over time as the handler's
implementation changes.

## Decision

We've decided to add an additional identifier to applications and handlers
called the "key".

The key's express purpose is for identifying associated data, and therefore has
more stringent requirements on its immutability than the name.

We further recommend the use of an RFC 4122 UUID as the format of all keys.
UUIDs can be generated at the time the application or handler is first
implemented. Many IDEs support generation of UUIDs.

Applications and handlers retain their names as a human-readable identifier.

## Consequences

The addition of the key provides a reliable mechanism for associating data with
applications and handlers. The consequences of changing a name or key are
clearer.

This does place some onus on the developer to generate a unique key and to
maintain the same key for the lifetime of the handler implementation.

<!-- references -->

[rfc 4122]: https://www.rfc-editor.org/rfc/rfc4122.html
