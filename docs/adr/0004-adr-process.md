# 4. ADR Process

Date: 2018-12-07

## Status

Accepted

## Context

We need a documented process for proposing, discussing and ultimate accepting or
rejecting ADRs via pull requests.

## Decision

We have decided to adopt a process that favours rapid changes, at least while
the project is in its infancy. To this end, we will allow ADRs in the `proposed`
status to be merged to `main` via PRs. The PRs may contain related code changes.

Whenever a `proposed` ADR is merged to `main`, a GitHub issue is created to
capture discussion about the ADR. Such issues are given the `adr` label.

Any `proposed` ADRs remaining on `main` must be resolved either by approving
the ADR, or by rejecting it and reverting any associated code changes.

## Consequences

This approach allows project maintainers to proceed with changes quickly,
documenting their decisions, and giving us a standard way to discuss those
decisions.

This comes at the cost of being potentially "expensive" at release time. We
feel this is an compromise prior to the first major release.
