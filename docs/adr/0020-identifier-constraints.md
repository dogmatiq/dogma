# 20. Constraints on Identifier Values

Date: 2023-06-14

## Status

Accepted

- Amends [9. Immutable Application and Handler Keys](0009-immutable-keys.md)
- Amends [12. Comparison of Identifiers](0012-identifier-comparison.md)

## Context

Identifiers — the names and keys used to identify applications and handlers —
are currently free-form. It would be beneficial if the engine could make some
assumptions about the format of identifiers so that they may be stored
efficiently.

## Decision

We will require all identifiers keys valid [RFC 9562] UUIDs. The `Identity()`
method on the "configurer" interfaces will continue to accept a string,
but that string must be an [RFC 9562] UUID in the canonical format:
`xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`.

We will also limit the maximum length of identifier names to 255 bytes. This
limit is specified in bytes rather than Unicode characters so that engines
implementations may easily reason about storage requirements.

## Consequences

Existing applications may need to be changed to use UUIDs or adhere to these
length limits. In practice, all known applications already meet these
requirements.

Engine implementations and `configkit` will need to validate identifiers and
reject those that do not meet the requirements.

<!-- references -->

[rfc 9562]: https://rfc-editor.org/rfc/rfc9562.html
