# 25. Prevent Reverting Ended Processes

Date: 2025-06-24

## Status

Proposed

- Amends [19. Automatic Process Creation](0019-automatic-process-creation.md)

## Context

Previously calling `End()` on a process and then later during the same event
handling, also executing a command or scheduling a timeout would revert/negate
the ending of the process.

With [ADR-24](0024-permanently-end-processes.md) permanently ending processes,
it would be best if `End()` was clear and final in all cases.

## Decision

We will change the specification such that an engine must not revert/negate an
ended process.

## Consequences

Engines must now panic if new actions are to be taken after being marked as
ending.

Ambiguity around how process lifecycle works when ending is eliminated by
removing the possibility altogether. Handling of events now aligns with the
existing treatment of timeouts, which are already discarded once a process
ends.
