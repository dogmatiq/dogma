<div align="center">

# Dogma

Build message-based applications in Go.

[![Documentation](https://img.shields.io/badge/go.dev-documentation-007d9c?&style=for-the-badge)](https://pkg.go.dev/github.com/dogmatiq/dogma)
[![Latest Version](https://img.shields.io/github/tag/dogmatiq/dogma.svg?&style=for-the-badge&label=semver)](https://github.com/dogmatiq/dogma/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/dogmatiq/dogma/ci.yml?style=for-the-badge&branch=main)](https://github.com/dogmatiq/dogma/actions/workflows/ci.yml)
[![Code Coverage](https://img.shields.io/codecov/c/github/dogmatiq/dogma/main.svg?style=for-the-badge)](https://codecov.io/github/dogmatiq/dogma)

</div>

## Overview

Dogma is a toolkit for building message-based applications in Go.

In Dogma, business logic is encapsulated in an **application** which consumes
and produces messages. The application is strictly separated from the
**engine**, which is responsible for message delivery and data persistence.

## Features

- **Built for [Domain Driven Design]** — The API uses DDD terminology to help
  developers align their understanding of the application's business logic with
  its implementation.
- **Flexible message format** — Supports any Go type that can be serialized to a
  byte slice, with built-in support for JSON and Protocol Buffers.
- **First-class testing** — Dogma's [testkit] module runs isolated behavioral tests of your application.
- **Multiple engine implementations** — Choose the engine with the best messaging and persistence semantics for your application.
- **Built-in introspection** — Analyze application code to visualize how messages traverse your applications.

## Related Repositories

- [dogmatiq/testkit] — utilities for blackbox testing of Dogma applications
- [dogmatiq/projectionkit] — utilities for building [projections](#projection) in various popular database systems
- [dogmatiq/example] — an example Dogma application that implements the features of a simple bank

## Concepts

Dogma leans heavily on the concepts of [Domain Driven Design], and is intended
to provide a suitable platform for applications that may wish to make use of
various design patterns such as [CQRS], [Event Sourcing] and [Eventual Consistency].

The following concepts are core to Dogma's design, and should be well understood
by any developer wishing to build an application:

- [Message](#message)
- [Message Handler](#message-handler)
- [Application](#application)
- [Engine](#engine)
- [Aggregate](#aggregate)
- [Process](#process)
- [Integration](#integration)
- [Projection](#projection)

### Message

A **message** is a data structure that represents a **command**, **event** or
**timeout** within an application.

A command is a request to make a single atomic change to the application's
state. An event indicates that the state has changed in some way. A single
command can produce any number of events, including zero.

A timeout is used to model business logic that depends on the passage of time.

Messages must implement the appropriate interface: `Command`, `Event` or
`Timeout`. These currently serve as aliases [`dogma.Message`](message.go), but
may diverge in the future.

### Message Handler

A message **handler** is part of an application that acts upon messages it
receives.

Handlers specify the Go type of the messages they can handle. These message
types are said to be **routed** to that handler.

Command messages are always routed to exactly one handler. Event messages may be
routed to any number of handlers, including zero. Timeout messages are always
routed back to the handler that produced them.

Dogma defines four handler types, one each for [aggregates](#aggregate),
[processess](#process), [integrations](#integration) and
[projections](#projection). These concepts are described in more detail
below.

### Application

An **application** is a collection of [message handlers](#message-handler) that
work together as a unit. Typically, each application is responsible for a
specific business (sub-)domain or "bounded-context".

### Engine

An engine is a Go module that delivers messages to an
[application](#application) and persists the application's state.

A Dogma application can be run on any engine. The choice of engine brings with
it a set of guarantees about how the application will behave, for example:

- **Consistency** — Different engines may provide different levels of
  consistency guarantees, such as [immediate consistency] or [eventual
  consistency].
- **Message delivery** — The engine may guarantee that messages are delivered to
  handlers in the order they are produced. Alternatively, messages may be
  processed out of order or in batches.
- **Persistence** — The engine may offer a choice of persistence mechanisms for
  application state, such as in-memory, on-disk, or in a remote database.
- **Data model** — The engine may provide a choice of data models for
  application state, such as relational or document-oriented.
- **Scalability** — The engine may provide a choice of scalability models, such
  as single-node or multi-node.

This repository is not an engine implementation. It defines the API that engines
and applications use to interact. The documentation for each interface indicates
whether the implementation is to be provided by the engine or the application.

One example of a Dogma engine is [Veracity].

### Aggregate

An **aggregate** is an entity that encapsulates a specific part of an
application's business logic and its associated state. Each **instance** of an
aggregate represents a unique occurrence of that entity within the application.

Each aggregate has an associated implementation of the
[`dogma.AggregateMessageHandler`] interface. The [engine](#engine) routes
command [messages](#message) to the handler to change the state of specific
instances. Such changes are represented by event messages.

An important responsibility of an aggregate is to enforce the invariants of the
business domain. These are the rules that must hold true at all times. For
example, in a hypothetical banking system, an aggregate representing a
customer's account balance must ensure that the balance never goes below zero.

Each aggregate instance's state is managed by the engine. State changes are
["immediately consistent"][immediate consistency] meaning that the changes made
by one command are always visible to subsequent commands routed to the same
instance.

Aggregates can be a difficult concept to grasp. The book [Domain Driven Design
Distilled], by Vaugn Vernon offers a suitable introduction to aggregates and the
other elements of domain driven design.

### Process

A **process** automates a long running business process. In particular, they can
be used to coordinate changes across multiple [aggregate](#aggregate) instances,
or between aggregates and [integrations](#integration).

Like aggregates, processes encapsulate related logic and state. Each
**instance** of a process represents a unique occurrence of that process within
the application.

Each process has an associated implementation of the
[`dogma.ProcessMessageHandler`] interface. The [engine](#engine) routes event
[messages](#message) to the handler to produce commands that are to be executed.

A process may use timeout messages to model business processes with time-based
logic. The engine always routes timeout messages back to the process instance
that produced them.

Processes use multiple atomic command messages to make changes to an
application's state. Therefore, the results of a process are considered
["eventually consistent"][eventual consistency].

### Integration

An **integration** is a message handler that interacts with some external
non-message-based system.

Each integration is an implementation of the [`dogma.ProcessMessageHandler`]
interface. The [engine](#engine) routes command [messages](#message) to the
handler which interacts with some external systems. Integrations may optionally
produce event messages that represent the results of their interactions.

Integrations are stateless from the perspective of the engine.

### Projection

A **projection** builds a partial view of the application's state from the
events that occur.

Each projection is an implementation of the [`dogma.ProjectionMessageHandler`]
interface. The [engine](#engine) routes event [messages](#message) to the
handler which typically updates a read-optimized database of some kind. This
view is often referred to as a "read model" or "query model", especially when
employing the [CQRS] pattern.

The [dogmatiq/projectionkit] module provides engine-agnostic tools for building
projections in various popular database systems, such as PostgreSQL, MySQL,
DynamoDB and others

<!-- references -->

[`dogma.aggregatemessagehandler`]: https://pkg.go.dev/github.com/dogmatiq/dogma?tab=doc#AggregateMessageHandler
[`dogma.application`]: https://pkg.go.dev/github.com/dogmatiq/dogma?tab=doc#Application
[`dogma.integrationmessagehandler`]: https://pkg.go.dev/github.com/dogmatiq/dogma?tab=doc#IntegrationMessageHandler
[`dogma.processmessagehandler`]: https://pkg.go.dev/github.com/dogmatiq/dogma?tab=doc#ProcessMessageHandler
[`dogma.projectionmessagehandler`]: https://pkg.go.dev/github.com/dogmatiq/dogma?tab=doc#ProjectionMessageHandler
[api documentation]: https://pkg.go.dev/github.com/dogmatiq/dogma
[cqrs]: https://martinfowler.com/bliki/CQRS.html
[dogmatiq/example]: https://github.com/dogmatiq/example
[dogmatiq/projectionkit]: https://github.com/dogmatiq/projectionkit
[dogmatiq/testkit]: https://github.com/dogmatiq/testkit
[domain driven design distilled]: https://www.amazon.com/Domain-Driven-Design-Distilled-Vaughn-Vernon/dp/0134434420
[domain driven design]: https://en.wikipedia.org/wiki/Domain-driven_design
[event sourcing]: https://martinfowler.com/eaaDev/EventSourcing.html
[eventual consistency]: https://en.wikipedia.org/wiki/Eventual_consistency
[immediate consistency]: http://www.informit.com/articles/article.aspx?p=2020371&seqNum=2
[rfc 2119]: https://tools.ietf.org/html/rfc2119
[testkit]: https://github.com/dogmatiq/testkit
[veracity]: https://github.com/dogmatiq/veracity
