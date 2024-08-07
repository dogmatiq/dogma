# 21. Permanently End Processes

Date: 2024-08-07

## Status

Proposed

- Amends [19. Automatic Process Creation](0019-automatic-process-creation.md)

## Context

Given a process instance that has ended, routing an event to that instance
causes the process instance to be "restarted". The state is a new instance of
the `ProcessRoot` obtained by calling `ProcessMessageHandler.New()`.

This behavior is surprising to users. It is also unclear how to properly handle
events that arrive after the process has ended.

## Decision

We will change the specification such that an engine must not call
`ProcessMessageHandler.HandleEvent()` for a process instance that has ended.

This ADR does not propose any changes to the behavior of timeout messages, which
are already discarded when a process ends.

## Consequences

TODO
