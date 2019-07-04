# 11. Message Timing Information

Date: 2019-07-05

## Status

Accepted

## Context

We need to decide whether message timing information should be exposed via the
API. Timing information refers to important points in time throughout the
lifecycle of a message.

The initial rationale for *not* exposing these timestamps was that any business
logic that depends on time in some way should include any time related
information within the message itself. We say that such logic "models time".

## Decision

The sections below focus on each of the message roles, their respective
timestamps of interest, and the decisions made in each case.

### Command Messages

We believe the initial requirement that the application "model time" is still
appropriate for command messages. The time at which the command message is
created is irrelevant; any time information relevant to the domain logic should
be included in the message itself.

**We have decided not to expose the command creation time.**

### Event Messages

The time at which an event is recorded is a fundamental property of the event
itself. Put another way, every event occurs at some time regardless of whether
the domain has any time-based logic.

Furthermore, the time at which event occurs may be relevant to some domain logic
that is *triggered* by the event, even if the aggregate that *produced* the event
has no time-based logic.

The inclusion of the "occurred time" as a fundamental property of the event is
supported by [Implementing Domain Driven
Design](https://www.amazon.com/Implementing-Domain-Driven-Design-Vaughn-Vernon/dp/0321834577),
Chapter 8, in the "Modeling Events" section.

**We have decided to include a `RecordedAt()` method on `ProcessEventScope` and
`ProjectionEventScope`.**

In actuality, a `Time()` method had already been added to `ProjectionEventScope`
without any supporting ADR, this method is to be renamed.

### Timeout Messages

The time at which a timeout message is scheduled to be handled is a fundamental
property of the timeout concept.

By using a timeout message the process is declaring that there is some
time-based component to its logic. It seems like an unnecessary hurdle to
require the application developer to include the scheduled time in the message
when the engine must necessarily keep track of this time anyway.

**We have decided to include a `ScheduledFor()` method on `ProcessTimeoutScope`.

## Consequences

As a result of this ADR it is easier for application developers to implement
logic that depends on time. This is especially true when new logic is added that
is triggered by domain logic that did not already have some time-based
component.

As a result, engine implements MUST record the time at which events occur. This
was not necessarily true before but it is likely that most engine
implementations would have done so anyway.
