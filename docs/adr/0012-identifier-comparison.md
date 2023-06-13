# 12. Comparison of Identifiers

Date: 2019-12-17

## Status

Accepted

- Amended by [20. Constraints on Identifier Values](0020-identifier-constraints.md)

## Context

Identifiers (the names and keys used to identify applications and handlers) must
be compared by engines to determine if two such entities are to be considered
equivalent.

The documentation specifies that such keys must be non-empty UTF-8 strings
consisting of printable characters without whitespace, but it did not previously
specify how such strings would be compared.

These identifiers are either mostly or entirely immutable and generated as part
of the source code. They do not need to be parsed and validated from user input.

## Decision

In keeping with current behavior, we've decided to specify byte-wise comparison
semantics for identifiers.

## Consequences

Existing tooling and engine implementations, such as `configkit` and the
in-memory engine implementation in `testkit` do not need to be changed.
Identifiers can be compared using Go's standard comparison operators.

It is possible that two identifiers may appear to be equal but consist of
different byte sequences. For example, when a character with a diacritic mark is
represented with a single codepoint, versus with combining characters. This is
considered unlikely to present problems in practice, but nonetheless it is worth
mentioning that the onus is on the application developer to normalize any UTF-8
identifiers.
