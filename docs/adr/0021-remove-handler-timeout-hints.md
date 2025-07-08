# 21. Remove handler timeout hints

Date: 2024-08-17

## Status

Accepted

Supersedes [10. Handler Timeout Hints](0010-handler-timeout-hints.md)

## Context

Handler timeouts hints have been present for a number of years, but have not
been widely adopted. Nor have engine implementations made use of them in any way
that offers capabilities beyond what is possible by having the handler
implementation impose its own context deadline.

## Decision

We will remove the `TimeoutHint()` methods from the handler interfaces, and the
`NoTimeoutBehavior` type.

## Consequences

- Application and engine implementations are marginally simpler.
- Engine's cannot tell ahead of time how long a message handler expects to take
  to handler a specific message.
