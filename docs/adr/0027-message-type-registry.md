# 27. Message type registry

Date: 2025-07-13

## Status

Accepted

## Context

The concept of a global message type registry was proposed as a possible
solution to the following problems:

- **Message type identification** — The `enginekit/marshaler` package currently
  uses a "portable name" to identify message types for serialization and
  deserialization. The portable name is the unqualified Go type name, or the
  fully qualified Protocol Buffers message name when using Protocol Buffers.
  This approach couples the serialized representation to the Go representation,
  which makes it **difficult to rename or relocate types after they have been
  persisted**.

- **Message route coupling** — The only mechanism an engine has to find
  information about messages is through the handler's message route
  configuration. If a message type is removed from the route configuration, for
  example, if an aggregate stops producing new events of a certain type, the
  engine has no information about that type. This can make it **difficult or
  impossible to handle historical events that have already been persisted**.

## Decision

We will introduce a global message type registry that requires explicit
registration of all message types before they can be used in handler routes.

Each registered message type is associated with an [RFC 4122] UUID, which
uniquely identifies the type when serialized. Engines will use this UUID instead
of the "portable name".

## Consequences

Message types can be renamed freely as long as the UUID remains the same.

Engines can access message type metadata independently of routing configuration
by querying the global registry.

Application developers must deal with the additional boilerplate of registering
each message type before it can be used in a handler's route configuration.
Typically, this would occur in an `init()` function.

<!-- references -->

[issue #176]: https://github.com/dogmatiq/dogma/issues/176
[ADR-9]: 0009-immutable-keys.md

[ADR-20]: 0020-identifier-constraints.md -->

[RFC 4122]: https://datatracker.ietf.org/doc/html/rfc4122
