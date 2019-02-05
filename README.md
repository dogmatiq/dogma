# Dogma

[![Build Status](http://img.shields.io/travis/com/dogmatiq/dogma/master.svg)](https://travis-ci.com/dogmatiq/dogma)
[![Code Coverage](https://img.shields.io/codecov/c/github/dogmatiq/dogma/master.svg)](https://codecov.io/github/dogmatiq/dogma)
[![Latest Version](https://img.shields.io/github/tag/dogmatiq/dogma.svg?label=semver)](https://semver.org)
[![GoDoc](https://godoc.org/github.com/dogmatiq/dogma?status.svg)](https://godoc.org/github.com/dogmatiq/dogma)
[![Go Report Card](https://goreportcard.com/badge/github.com/dogmatiq/dogma)](https://goreportcard.com/report/github.com/dogmatiq/dogma)


Dogma is a specification and API for building message-based applications in Go.

Dogma attempts to define a practical standard for authoring message-based
applications in a manner agnostic to the mechanisms by which messages are
transported and application state is persisted.

## Related Repositories

- [dogmatiq/testkit] - utilities for blackbox testing of Dogma applications
- [dogmatiq/example] - an example Dogma application that implements the features of a simple bank

## Concepts

Dogma leans heavily on the concepts of [Domain Driven Design], and is intended
to provide a suitable platform for applications that may wish to make use of
various design patterns such as [CQRS], [Event Sourcing] and [Eventual Consistency].

### Message and Message Handler

A **message** is an application-defined unit of data that encapsulates a
**command** or **event** within a message-based application. A command messages
represents a request for the application to perform some action, whereas an
event message indicates that some action has already occurred. A single command
message can produce zero or more events.

Additionally, a **timeout** message can be used to perform actions within an
application at specific wall-clock times.

Messages are represented by the [`dogma.Message`](message.go) interface, which
is deliberately empty, allowing applications to use any Go type as a message.

A message **handler** is some portion of application-defined logic that acts
upon messages that it receives.

Handlers announce the messages they wish to receive based on the message's Go
type. Messages types that are received by a particular handler are said to be
**routed** to that handler.

Command messages are always routed to exactly one handler. Event messages may
be routed to zero or more handlers. Timeout messages behave differently, always
being routed back to the handler that produced them.

Each message represents a single atomic operation within the application.

Dogma defines four handler types, one each for [aggregates](#aggregate),
[processess](#process), [integrations](#integration) and
[projections](#projection). These concepts are described in more detail
below.

### Application

An **application** is a set of [message handlers](#message-and-message-handler)
that operate together as a unit. Applications are represented by the
[`dogma.Application`] interface.

### Engine

An **engine** is the platform upon which an [application](#application) is
executed. The engine is responsible for the delivery of messages and the
persistence of application data.

This module does not provide an engine implementation, but rather defines the
API that sits between the application and the engine. The API documentation for
each interface indicates whether the implementation is to be provided by the
engine or the application.

### Aggregate

An **aggregate** is a unit of application logic with associated state that
encodes the business invariants of a specific application. These invariants are
the "rules" of the business domain that must not be violated, even temporarily.
This concept is taken directly from [Domain Driven Design].

An aggregate receives command [messages](#message-and-message-handler) in order
to effect a change in a particular **instance** of that aggregate. Such state
changes are represented by event messages. By definition, changes to the state
of an aggregate instance are ["immediately consistent"][Immediate Consistency]
(aka "transactionally consistent"). This means that the results of a command
against a given instance are always visible to subsequent commands for that
instance.

Aggregate state is managed by the [engine](#engine), ensuring that state
transitions occur atomically with the recording of the events that describe
those state
transitions.

Aggregates can be quite a difficult concept to grasp. The book [Domain Driven
Design Distilled], by Vaugn Vernon offers a suitable introduction to aggregates
and the other elements of domain driven design.

Aggregates are represented by the [`dogma.AggregateMessageHandler`] interface.

### Process

A **process** is a unit of application logic with associated state that serves
to automate some long running business process. In particular, they can be used
to coordinate changes across multiple [aggregate](#aggregate) instances, or
between aggregates and [integrations](#integration).

Processes receive event [messages](#message-and-message-handler) and produce
command messages. Like aggregates, the received events are routed to a specific
instance.

Additionally, processes can produce timeout messages, which are routed back to
the same process instance at a specific time. Such messages are used to
implement processes that incorporate some time-based component.

Because a process coordinates changes within the application using multiple
messages, and each message represents a single atomic change to the application's
state, the changes made by a process are ["eventually consistent"][Eventual Consistency].
The precise guarantees regarding process consistency are specific to the [engine](#Engine) implementation.

Process state is managed by the engine, ensuring that state transitions occur
atomically with the enqueing of the commands the processes produce.

Processes are represented by the [`dogma.ProcessMessageHandler`] interface.

### Integration

An **integration** is a unit of application logic that integrates an
application with some non-message-based system.

Integrations receive command [messages](#message-and-message-handler) and
produce event messages. They do not have any state that is managed by the
[engine](#Engine).

Integrations are represented by the [`dogma.IntegrationMessageHandler`] interface.

### Projection

A **projection** is a unit of application logic that derives some specific
portion of application state from the events that occur. This state is often
referred to as a "read model".

Projections receive event [messages](#message-and-message-handler) and do not
produce messages of any kind.

They do not have any state that is modelled by the Dogma API, but [engine](#Engine)
implementations may provide mechanisms for persisting projection state in
various data stores, such as SQL databases, document stores, flat files, etc.

Projections are represented by the [`dogma.ProjectionMessageHandler`] interface.

<!-- references -->
[Domain Driven Design]: https://en.wikipedia.org/wiki/Domain-driven_design
[Domain Driven Design Distilled]: https://www.amazon.com/Domain-Driven-Design-Distilled-Vaughn-Vernon/dp/0134434420
[CQRS]: https://martinfowler.com/bliki/CQRS.html
[Event Sourcing]: https://martinfowler.com/eaaDev/EventSourcing.html
[Immediate Consistency]: http://www.informit.com/articles/article.aspx?p=2020371&seqNum=2
[Eventual Consistency]: https://en.wikipedia.org/wiki/Eventual_consistency
[API documentation]: https://godoc.org/github.com/dogmatiq/dogma
[RFC 2119]: https://tools.ietf.org/html/rfc2119

[dogmatiq/testkit]: https://github.com/dogmatiq/testkit
[dogmatiq/example]: https://github.com/dogmatiq/example

[`dogma.Application`]: https://godoc.org/github.com/dogmatiq/dogma#Application
[`dogma.AggregateMessageHandler`]: https://godoc.org/github.com/dogmatiq/dogma#AggregateMessageHandler
[`dogma.ProcessMessageHandler`]: https://godoc.org/github.com/dogmatiq/dogma#ProcessMessageHandler
[`dogma.IntegrationMessageHandler`]: https://godoc.org/github.com/dogmatiq/dogma#IntegrationMessageHandler
[`dogma.ProjectionMessageHandler`]: https://godoc.org/github.com/dogmatiq/dogma#ProjectionMessageHandler
