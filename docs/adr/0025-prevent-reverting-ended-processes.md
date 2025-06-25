# 25. Prevent Reverting Ended Processes

Date: 2025-06-24

## Status

Accepted

- Amends [19. Automatic Process Creation](0019-automatic-process-creation.md)

## Context

Calling `ExecuteCommand()` or `ScheduleTimeout()` after calling `End()` on the
same `Process[Event|Timeout]Scope` instance causes the `End()` call to be
"reverted" or "negated" - having no effect.

In [ADR-24](0024-permanently-end-processes.md), the semantics of `End()` were
changed such that it takes effect permanently, meaning that no future events
will be routed to the ended instance (timeout messages already behaved this
way). In light of this change, it's misleading if `End()` can sometimes be
"undone".

## Decision

We will change the specification such that the engine must panic if the user
calls `ExecuteCommand()` or `ScheduleTimeout()` (or any future similar "write"
operation) after calling `End()` on the same `Process[Event|Timeout]Scope`
instance.

## Consequences

Engines must now track whether a process instance has ended and suppress event
delivery to that instance, whereas current implementations typically revert
ending the process and keep it alive.

Ambiguity around how to handle process effects after it has been ended is
eliminated by removing the possibility altogether. All effects must be produced
before calling `End()`.

This change resolves the immediate ambiguity but may not be sufficient in the
long term. In future, it may be necessary for a process to explicitly declare
ways to revert/negate ending or start a new instance.
