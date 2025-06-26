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

Calling `End()` is now always final. The scope on which end is called will not
be able to perform any additional operations. Any future messages that are
routed to the ended process instance will be ignored.

Engines will need to keep track of ended processes indefinitely, presenting an
additional data management concern.
