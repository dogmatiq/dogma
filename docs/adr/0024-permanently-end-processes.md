# 24. Permanently End Processes

Date: 2024-08-07

## Status

Accepted

- Amends [19. Automatic Process Creation](0019-automatic-process-creation.md)

## Context

Routing an event to a process instance that has already ended causes the process
to be "restarted" with a new `ProcessRoot` instance.

This behavior is surprising because it is not possible to distinguish a
"restarted" process instance from one that has never existed. Even if it were
possible, it is unclear how to handle events within an ended process.

## Decision

We will change the specification such that an engine must not call
`ProcessMessageHandler.HandleEvent()` for a process instance that has ended.

This ADR does not propose any changes to the behavior of timeout messages, which
are already discarded when a process ends.

## Consequences

Engines must now track whether a process instance has ended and suppress event
delivery to that instance, whereas current implementations typically discard all
information about the instance when it ends.

Ambiguity around how to handle events for ended processes is eliminated by
removing the possibility altogether. Handling of events now aligns with the
existing treatment of timeouts, which are already discarded once a process ends.

This change resolves the immediate ambiguity but may not be sufficient in the
long term. In future, it may be necessary for a process to explicitly declare
which event types can start a new instance, or provide some other fine-grained
control over event routing.
