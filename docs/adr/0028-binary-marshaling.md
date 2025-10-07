# 28. Self-marshaling types

Date: 2025-10-07

## Status

Accepted

- References [27. Message type registry](0027-message-type-registry.md)

## Context

In a production system, all implementations of `AggregateRoot`, `ProcessRoot`
and `Message` must be serialisable to be useful. Even though aggregate roots can
be reconstructed from historical events, it's still useful to persist a snapshot
of the aggregate root as an optimization.

Currently, this serialisation is handled by the `enginekit/marshaler` package.
It supports multiple "codecs", such as Protocol Buffers, JSON and CBOR, and
provides a HTTP-like system for falling back based on message media type.

While this system is entirely functional, it means that engine implementations
(`enginekit`, in practice) must support any possible encoding an application
may wish to use.

## Decision

We will change the `Message`, `AggregateRoot` and `ProcessRoot` interfaces to
adhere to Go's standard [`encoding.BinaryMarshaler`] and
[`encoding.BinaryUnmarshaler`] interfaces.

Engines will use each message type's UUID, introduced in [ADR-27], as a
persistence-safe identifier for the message type. A new zero-value message type
can be constructed based on the UUID before calling `UnmarshalBinary()`.

## Consequences

Application developers are now responsible for implementing the
`MarshalBinary()` and `UnmarshalBinary()` methods on their message types,
aggregate roots and process roots. This puts control of the encoding format in
the hands of the application developer, but does introduce some boilerplate. We
may need to develop code generation utilities to help with this. The
[`dogmatiq/primo`] module can generate the marshaling methods for Protocol
Buffers messages.

Engines no longer need to provide a marshaling mechanism. Accordingly,
`enginekit/marshaler` package is no longer needed and may be removed.

<!-- references -->

[ADR-27]: 0027-message-type-registry.md
[`encoding.BinaryMarshaler`]: https://pkg.go.dev/encoding#BinaryMarshaler
[`encoding.BinaryUnmarshaler`]: https://pkg.go.dev/encoding#BinaryUnmarshaler
[`dogmatiq/primo`]: http://github.com/dogmatiq/primo
