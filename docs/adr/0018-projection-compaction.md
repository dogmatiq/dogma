# 18. Compacting Projection Data

Date: 2020-11-11

## Status

Accepted

## Context

Many projections produce data that is only required for a limited period of
time.

An application developer needs to consider how long projection data will be
retained and ideally implement measures to remove any unnecessary data.

Often such data can be removed when some future event occurs. However, in some
cases there is no future event that appropriately indicates the data is no
longer required.

We have encountered a genuine example of this when implementing an SQL
projection that inserts a row whenever a user performs a specific action. The
data is queried in order to enforce a "cool down" that prevents the user from
repeating that action again within a specific time frame.

The insert is triggered by the occurrence of an event, but the data becomes
unnecessary whenever the "cool down" time has elapsed.

In this particular use case the "cool down" was not part of the business logic,
but rather an API level restriction. Hence, processes/timeout messages were not
the appropriate solution.

## Decision

We have decided to add a `Compact()` method to `ProjectionMessageHandler`.

The implementation of `Compact()` can modify the projection's data by whatever
means is appropriate such that unnecessary data is removed but the projection
still serves its purpose.

## Consequences

### Engine Complexity

This change does introduce further complexity to engine implementations. It is,
of course, entirely possible to perform this kind of compaction outside the
engine. To keep this compexity to a minimum we have avoided building in any
facility for the projection handler to schedule compaction at specific times.

### Backwards Compatibility

Every existing projection implementation now requires an addtional method, even
if no compaction is required. The `NoCompactBehavior` struct can be embedded
within a `ProjectionMessageHandler` implementation to provide a no-op
implementation of `Compact()`.

### First Class Compaction

Making compaction a first-class feature of the projection interface encourages
application developers to think about the lifetime of the projection's data and
how it might be cleaned up; something that might easily be ignored.

### Testing

`dogmatiq/testkit` can be updated to perform compaction when events are routed
to projections, ensuring that the compaction code is actually invoked during
testing and hopefully that it does not interfere with the regular operation of
the projection.

### Optimal Scheduling of Compaction

Many aspects of an engine's behavior are engine-specific. By giving the engine
control over when compaction occurs engine implemenators may be able to make
informed decisions about when it's appropriate to perform compaction.

For example, an engine may perform compaction only when there is a period of low
event activity.

Engine implementations that use clustering and/or sharding techniques might
schedule compaction to occur on a single node.
