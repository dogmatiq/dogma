# 15. Routing Unrecognized Messages

Date: 2020-11-01

## Status

Accepted

## Context

[ADR-14](0014-apply-historical-events-to-aggregates.md) relaxed the
specification such that `AggregateRoot.ApplyEvent()` implementations were no
longer required to panic with an `UnrecognizedMessage` value when passed an
unexpected message type.

Prompted by this requirement, we relaxed the requirement for ALL handler
methods, which was likely too broad of a change.

Specifically, unlike when handling a message, the routing methods
`AggregateMessageHandler.RouteCommandToInstance()` and
`ProcessMessageHandler.RouteEventToInstance()` do not have the option of "doing
nothing" when passed an unexpected message type.

## Decision

Reinstate the hard requirement that the handlers MUST panic with
`UnexpectedMessage` when asked to route a message type that was not configured
as being consumed by that handler.

## Consequences

This reduces ambiguity about how the handler should be implemented, and allows
tools like "testkit" and "dogmavet" to make stricter assertions and more direct
suggestions.
