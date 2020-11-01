# 15. Routing Unrecognized Messages

Date: 2020-11-01

## Status

Accepted

- Amends [14. Applying Historical Events to Aggregate Instances](0014-apply-historical-events-to-aggregates.md)

## Context

ADR-14 relaxed the specification such that handler implementations were no
longer required to panic with an `UnrecognizedMessage` value when passed an
unexpected message type. However, it probably did so too broadly.

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
