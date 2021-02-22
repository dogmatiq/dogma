# 16. Automatic Process Creation

Date: 2021-02-18

## Status

Accepted

## Context

After reviewing and reworking the aggregate API in
[ADR-16](0016-automatic-aggregate-creation.md) and
[ADR-17](0017-recreate-aggregate-after-destruction.md) we have conducted a
similar review of the process API in an effort to both simplify the API and
improve consistency between the aggregate and process APIs.

## Decision

- Remove `ProcessEventScope.Begin()`.
- Remove `Process[Event|Timeout]Scope.HasBegun()`.
- Remove `Process[Event|Timeout]Scope.Root()`.
- Pass the process root directly to `ProcessMessageHandler.Handle[Event|Timeout]()`.
- Allow `Process[EventTimeout]Scope.End()` to be called at any time.
- Allow `Process[EventTimeout]Scope.ExecuteCommand()` and `ScheduleTimeout()` to
  be called at any time. Doing so should "negate" any prior call to `End()` as
  though it were never called.
- Routing a command message to a process instance causes that instance to begin.

## Consequences

Largely, this should simplify implementations of both `ProcessMessageHandler` by
application developers and `Process[Event|Timeout]Scope` by engine developers.

The removal of `Begin()` does mean that engines no longer get an explicit
request to start a process instance. However, we can avoid starting unwanted
instances simply by never routing the command to them in the first place.

There could be some confusion if `End()` is called then subsequently "undone" by
executing a command or scheduling a timeout. In this case, any timeout messages
that had been scheduled before calling `End()` remain pending.
