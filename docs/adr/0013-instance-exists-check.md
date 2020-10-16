# 13. Aggregate and Process Instance Existance Checks

Date: 2020-10-16

## Status

Accepted

## Context

When handling a message in an `AggregateMessageHandler` or
`ProcessMessageHandler`, there is a practical need to check if the instance
already exists before performing an operation on the scope, such as producing a
new message or destroying/ending the instance.

For example, it may be necessary to ignore a command that has been sent to an
aggregate that has been destroyed. There is currently no idiomatic way to do
this, short of calling `Create()` only to immediately call `Destroy()` if
creation succeeded.

## Decision

We have chosen to add an `Exists()` method to aggregate scopes, and an analogous
`HasBegun()` method to process scopes.

## Consequences

This is a fairly wide-reaching BC break for engine implementations, although
implementing it should be trivial.
