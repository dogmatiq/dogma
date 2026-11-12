# 36. Native UUID type

Date: 2026-11-13

## Status

Accepted

- References [20. Constraints on Identifier Values][ADR-20]
- References [27. Message type registry][ADR-27]
- Referenced by [37. Expose message IDs][ADR-37]

## Context

Dogma uses [RFC 9562] UUIDs as identifiers. [ADR-20] requires the key that
identifies each application and handler to be a UUID, and [ADR-27] gives each
message type a UUID. Each event stream that a projection consumes is also
identified by a UUID.

The API represents all of these UUIDs as a `string` in the canonical form, such
as "65f9620a-65c1-434e-8292-60cd7938c4de". A `string` can hold any text, so its
type doesn't guarantee that the value is a UUID, and code that accepts one has
to validate it. UUIDs are also case-insensitive, so two strings that represent
the same UUID aren't necessarily equal. The documentation for
`ProjectionMessageHandler.CheckpointOffset()` deals with this by requiring the
engine to pass the stream ID in lowercase.

Go 1.27 adds the [`uuid`] package to the standard library. Its [`uuid.UUID`]
type holds the 16 bytes of a UUID, so every value of that type is a UUID, and
two values that represent the same UUID are equal when compared with `==`.

Applications usually write the UUIDs that identify applications, handlers, and
message types as string literals. They pass them to `Identity()` in each
`Configure()` method, and to `RegisterCommand()`, `RegisterEvent()` and
`RegisterDeadline()`.

## Decision

We will require Go 1.27, and use `uuid.UUID` for each UUID that Dogma or the
engine supplies at run time:

- `ProjectionEventScope.StreamID()` returns a `uuid.UUID`.
- `ProjectionMessageHandler.CheckpointOffset()` takes the stream ID as a
  `uuid.UUID`.
- `RegisteredMessageType.ID()` returns a `uuid.UUID`.
- `RegisteredMessageTypeByID()` takes the message type ID as a `uuid.UUID`.

We will keep `string` for each UUID that the application writes in its source
code:

- the key passed to `Identity()` in a `Configure()` method
- the message type ID passed to `RegisterCommand()`, `RegisterEvent()` and
  `RegisterDeadline()`

If these UUIDs were `uuid.UUID` values, the application would have to wrap each
one in a call to [`uuid.MustParse()`]. That would add noise to every
`Configure()` method and every registration, without catching any mistake that
Dogma or the engine doesn't already catch at startup.

The registration functions won't accept the other forms that [`uuid.Parse()`]
accepts, such as a UUID with a `urn:uuid:` prefix.

## Consequences

The type of each of these UUIDs guarantees that it holds a well-formed UUID, and
code can compare two of them with `==`. Engines no longer need to convert stream
IDs to lowercase.

Applications and engines that use Dogma must use Go 1.27 or later.

Every projection handler and every engine must change to match the new
signatures of `StreamID()` and `CheckpointOffset()`, including the adaptors in
[projectionkit]. Engines and tools that call `RegisteredMessageType.ID()` or
`RegisteredMessageTypeByID()`, such as the marshaler in [enginekit], must also
change.

`Configure()` methods and registration calls in existing applications don't
change, and engines continue to validate the keys passed to `Identity()`.

<!-- references -->

[ADR-20]: 0020-identifier-constraints.md
[ADR-27]: 0027-message-type-registry.md
[ADR-37]: 0037-expose-message-ids.md
[enginekit]: https://github.com/dogmatiq/enginekit
[projectionkit]: https://github.com/dogmatiq/projectionkit
[RFC 9562]: https://www.rfc-editor.org/rfc/rfc9562.html
[`uuid`]: https://pkg.go.dev/uuid
[`uuid.MustParse()`]: https://pkg.go.dev/uuid#MustParse
[`uuid.Parse()`]: https://pkg.go.dev/uuid#Parse
[`uuid.UUID`]: https://pkg.go.dev/uuid#UUID
