# 14. Applying Historical Events to Aggregate Instances

Date: 2020-10-19

## Status

Accepted

## Context

Event sourcing engines need to call `AggregateRoot.ApplyEvent()` with
"historical" event types. That is, event types that have already been recorded
against the instance but are no longer configured for production by that
aggregate.

The current specification language prohibits this, as per the `ApplyEvent()`
documentation:

> It MUST NOT be called with a message of any type that has not been
> configured for production by a prior call to Configure().

Additionally, without adding some new features to `AggregateConfigurer` it is
impossible to declare an event as historical, meaning that there is no way to
discover historical event types from the configuration.

## Decision

We have chosen to relax the language in the specification to allow calling
`ApplyEvent()` with any historical event types in addition to those configured
for production.

## Consequences

Existing event sourcing engine implementations are no longer violating the spec.

As this is a documentation change only it does not provide engines with any
information they need to determine if an event type is historical. This should
be a non-issue as the engine itself will have its own mechanism for loading
historical events. We may expand the functionality of `AggregateConfigurer` in
the future to allow declaration of historical event types.
