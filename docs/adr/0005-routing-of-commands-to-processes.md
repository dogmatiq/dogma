# 5. Routing of commands to processes

Date: 2018-12-11

## Status

Accepted

## Context

We need to decide whether command messages can be routed directly to processes.

## Decision

We have decided to disallow this behavior - processes may only handle events
and timeout messages.

If we were to allow processes to accept commands directly it may be tempting to
implement domain idempotency in the process instead of in an aggregate where
such logic belongs.

Furthermore, it's easier to allow commands to be routed to processes in a future
version than it is to remove it once it's in use.

## Consequences

It's simpler to describe how the various message and handler types interact.

A minor downside is that commands that unconditionally trigger a process still
have to pass through an aggregate that "does nothing".
