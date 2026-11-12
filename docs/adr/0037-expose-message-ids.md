# 37. Expose message IDs

Date: 2026-11-13

## Status

Accepted

- References [36. Native UUID type][ADR-36]

## Context

The Dogma API doesn't define a message ID — a value that uniquely identifies a
message. An engine may give each message an ID for its own use, but a handler
can't read it.

An integration or projection can send information to another system. If a
failure occurs after the handler sends the information but before the engine
records that the handler has handled the message, the engine delivers the
message again, and the handler sends the same information a second time.

A projection that stores its data and its checkpoint offsets in the same
database doesn't have this problem. A checkpoint offset is the position, within
an event stream, of the next event that the projection expects to handle. The
projection updates its data and its checkpoint offset in one transaction, so
when the engine delivers an event again, the checkpoint offset shows that the
projection has already handled that event. A projection that sends its data to
another system can't do this, because it can't update both the other system and
its own checkpoint offset in one transaction.

[dogmatiq/protobuf] defines an envelope — a [Protocol Buffers] message that
holds a Dogma message together with information about it, such as its ID. The
envelope stores each message's ID as a [UUID]. [enginekit], which provides tools
for building engines, and [runkit], an engine, use this envelope to store and
transport messages. [testkit], which runs applications in tests, doesn't use it,
and gives each message a sequential number as its ID, such as "1" or "2".

## Decision

We will add a `MessageID()` method to every scope in which a handler handles a
message:

- `AggregateCommandScope`
- `ProcessEventScope`
- `ProcessDeadlineScope`
- `IntegrationCommandScope`
- `ProjectionEventScope`

The engine must give each message an ID that no other message has. `MessageID()`
returns the same ID every time the engine delivers that message, to any handler.

Only integrations and projections need the ID today. We will add `MessageID()`
to the aggregate and process scopes as well, so that every scope in which a
handler handles a message has the same method.

Each message ID must be a UUID, and `MessageID()` returns a [`uuid.UUID`]. We
chose UUIDs because the envelope stores message IDs as UUIDs, so engines that
use it already meet this requirement. How the API represents UUIDs is a separate
decision, recorded in [ADR-36].

### Dismissed alternatives

We considered recommending that projections build a key from
`ProjectionEventScope.StreamID()` and `ProjectionEventScope.Offset()`, and send
it with the information so that the other system can discard the second copy. No
two events have the same stream ID and offset, so this would need no change to
the API or to engines. We dismissed it because the pair reflects how the engine
organizes events into streams, rather than identifying the event itself, and
because it only works for projections. Commands aren't in event streams, so an
integration has no such pair.

## Consequences

Integrations and projections can use a message's ID as a deduplication key, or
as part of one, when they send information about that message to another system.
A deduplication key is a value that is the same each time the handler sends
information about the same message, and different for different messages. The
other system uses it to discard the second copy, so sending the same information
twice is [idempotent], meaning that it has the same effect as sending it once.

Dogma now has a key for discarding duplicates in each direction. Code outside an
application passes an idempotency key to `WithIdempotencyKey()` so that the
engine can discard a duplicate command that enters the application. A handler
passes a message's ID to another system so that the other system can discard
duplicate information that leaves the application.

Every engine, including the one in [testkit], must give each message a UUID and
implement `MessageID()` on the five scopes. An engine that uses the key passed
to `WithIdempotencyKey()` as the command's ID can't keep doing so, because that
key can be any non-empty string. It can derive a UUID from the key instead.

[testkit] must switch from sequential numbers to UUIDs for its message IDs, so
any of its output that shows message IDs will show UUIDs instead of short
numbers.

This ADR introduces one term to the glossary: message ID.

<!-- references -->

[ADR-36]: 0036-native-uuid-type.md
[dogmatiq/protobuf]: https://github.com/dogmatiq/protobuf
[enginekit]: https://github.com/dogmatiq/enginekit
[idempotent]: https://en.wikipedia.org/wiki/Idempotence
[Protocol Buffers]: https://protobuf.dev/
[runkit]: https://github.com/dogmatiq/runkit
[testkit]: https://github.com/dogmatiq/testkit
[UUID]: https://www.rfc-editor.org/rfc/rfc9562.html
[`uuid.UUID`]: https://pkg.go.dev/uuid#UUID
