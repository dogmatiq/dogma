# 10. Handler Timeout Hints

Date: 2019-06-28

## Status

Accepted

## Context

We need to decide on a mechanism for engine implementations to determine
suitable timeout durations to apply when handling a message.

For aggregate message handlers, which are not permitted to access external
resources, a fairly constant timeout duration should be discernable by the
engine developers.

For all other handler types, which may make network requests or perform CPU
intensive work, there is no one timeout duration that makes sense in all
circumstances.

## Decision

We have decided to allow process, integration and projection message handlers
to provide a timeout "hint" on a per-message basis by way of a
`TimeoutHint(dogma.Message) time.Duration` method.

By returning a zero-value duration, the handler indicates that it can provide no
useful "hint" and that the engine should choose a timeout by other means.

## Consequences

Engine portability is aided by including the timeout "hints" alongside the code
that is subject to those timeouts.

As always, such a change increases the complexity of the API.

We did discuss leaving the configuration of timeouts as an engine-specific
problem. However, we ultimately decided that the minor increase in API surface
area is outweighed by the portability benefits.
