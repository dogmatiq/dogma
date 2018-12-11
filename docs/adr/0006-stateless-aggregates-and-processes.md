# 6. Stateless aggregates and processes

Date: 2018-12-11

## Status

Accepted

## Context

We need to decide how to represent aggregates and processes that do not have
any state.

This may sound nonsensical at first, but there are legitimate implementations of
both that do not require any state, or perhaps more correctly the only state is
whether or not they exist at all.

This is perhaps more likely to occur in a CQRS/ES environment where the state
associated with an aggregate is a "write-model". In this case, the only state
that needs to be maintained is that which is required to make decisions about
which events to produce.

## Decision

We've opted to have Dogma provide empty "root" implementations out of the box,
for both aggregates and processes. The implementations should be unexported
structs made available by exposed global variables in the `dogma` package.

Handler implementations can return these values from their `New()` methods to
indicate that they do not keep state.

As these values are valid implementations of the `AggregateRoot` / `ProcessRoot`
interfaces, engine implementations need not handle these impementations specially,
though they can opt to do so by treating them as "sentinel" values.

## Consequences

By including these implementations we ensure that there is a standard,
observable way to indicate statelessness without requiring vastly different
codepaths.

Stateless implementations still need to implement a `New()` method, which is
perhaps overly verbose, but this could be mitigated by also providing embeddable
structs with `New()` implementations that return the empty roots.
