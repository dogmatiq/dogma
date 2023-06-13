# 8. Location of Testing Features

Date: 2019-01-03

## Status

Accepted

## Context

We need to decide whether Dogma's testing features should reside in the `dogma`
module itself, or a separate `dogmatest` module.

### Benefits to separate module

1. They can be versioned separately. A `dogma` release can be made without
   having to update the testing tools. This is a _pro_ for the releaser, but
   possibly a _con_ for the users.

1. Users that opt not to test, or not to test using our tools do not need to
   download `dogmatest` at all. This is not a particularly strong argument.

1. We can make BC breaking changes to `dogmatest`, without having to make
   a new major release of `dogma`. We would need to document clearly which
   major `dogma` versions are supported by which major `dogmatest` versions.

### Benefits to a single module

1. The user does not need to determine which `dogmatest` version to use with
   their application; compatible tools are always available right from `dogma`. If
   we want to encourage users to test their domain logic using these utilities;
   they should be easy to access.

1. As mentioned in [#16](github.com/dogmatiq/dogma), splitting the testing
   tools would mean that the example code would also need to be moved elsewhere.
   However, we have since already decided to move the example code to a separate
   module in [ADR-7](0007-location-of-examples.md).

## Decision

For the same reasons [ADR-7](0007-location-of-examples.md) we've decided to
implement testing tools in their own module.

## Consequences

As always, introducing more repositories increases the overhead of dependency
management for Dogmatiq maintainers. However, keeping any code that is not
directly related to interoperability outside of `dogma` removes the possiblity
of having to release a new major semver version due to a backwards-incompatible
change that only affected the ancillary code and not the Dogma API itself.

We will need to make it very clear which `dogmatest` versions work with which
`dogma` versions in the documentation for both projects. We may end up in a
scenario where we bump the major version of `dogmatest`, without such a bump
being made to `dogma`. However, we should probably at least _prefer_ to keep
their major versions in lock-step.
