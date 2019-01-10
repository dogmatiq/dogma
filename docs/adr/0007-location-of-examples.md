# 7. Location of Examples

Date: 2019-01-03

## Status

Accepted

## Context

We need to decide whether Dogma's examples should reside in the `dogma`
repository itself, or a separate `examples` repository.

## Decision

We've decided to move the examples to a separate repository, so that we can
provide fully-functional examples that depend on modules/packages that we would
not want to have as dependants of `dogma` itself, such as `mysql`, etc.

## Consequences

As always, introducing more repositories increases the overhead of dependency
management for Dogmatiq maintainers. However, keeping any code that is not
directly related to interoperability outside of `dogma` removes the possiblity
of having to release a new major semver version due to a backwards-incompatible
change that only affected the ancillary code and not the Dogma API itself.
