# 11. Message Timing Information

Date: 2019-07-05

## Status

Accepted

## Context

We need to decide whether message timing information should be exposed via the
API. In this context "timing information" refers to important points in time
throughout the lifecycle of a message.

The initial rationale for *not* exposing these timestamps was that any business
logic that depends on time in some way should explicitly include any timing
information within the message itself. We call such logic "time-based" and the
approach of including explicit timing information "modeling time".

## Decision

The sections below focus on each of the message roles, their respective
timestamps of interest, and the decisions made in each case.

### Command Messages

We believe the existing requirement that the application "model time" is still
appropriate for command messages. The time at which the command message is
created or enqueued is irrelevant; any time information relevant to the domain
logic should be included in the message itself.

**We have decided not to expose the command creation time.**

### Event Messages

The time at which an event is recorded is a fundamental property of the event
itself. Put another way, every event occurs at some time regardless of whether
the domain is time-based.

Furthermore, the time at which the event occurs may be relevant to some
ancillary domain logic that is *triggered* by the event, even if the aggregate
that *produced* the event has no time-based logic.

The inclusion of the "occurred time" as a fundamental property of the event is
supported by [Implementing Domain Driven
Design](https://www.amazon.com/Implementing-Domain-Driven-Design-Vaughn-Vernon/dp/0321834577),
Chapter 8, in the "Modeling Events" section.

**We have decided to include a `RecordedAt()` method on `ProcessEventScope` and `ProjectionEventScope`.**

In actuality, a `Time()` method had already been added to `ProjectionEventScope`
without any supporting ADR, this method is to be renamed.

### Timeout Messages

The time at which a timeout message is scheduled to be handled is a fundamental
property of the timeout concept.

By definition, the use of a timeout message indicates that there is time-based
logic. It seems like an unnecessary imposition to require the application
developer to include the scheduled time in the message.

**We have decided to include a `ScheduledFor()` method on `ProcessTimeoutScope`.**

## Consequences

As a result of this ADR it is easier for application developers to implement
time-based logic.

Engine implementations must now record the time at which events occur. This was
not necessarily true before but it is likely that most engines would have done
so anyway.
