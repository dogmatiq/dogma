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

In Dogma, an **application** implements business logic by consuming and
producing messages . The application is strictly separated from the **engine**,
which handles message delivery and data persistence.

## Features

- **Built for [Domain Driven Design]**: The API uses DDD terminology to help
  developers align their understanding of the application's business logic with
  its implementation.

- **Flexible message format**: Supports any Go type that can be serialized as a
  byte slice, with built-in support for JSON and Protocol Buffers.

- **First-class testing**: Dogma's [testkit] module runs isolated behavioral
  tests of your application.

- **Engine-agnostic applications**: Choose the engine with the best messaging
  and persistence semantics for your application.

- **Built-in introspection**: Analyze application code to visualize how messages
  traverse your applications.

## Related repositories

- [testkit]: Utilities for black-box testing of Dogma applications.
- [projectionkit]: Utilities for building [projections](#projection) in popular database systems.
- [example]: An example Dogma application that implements basic banking features.

## Concepts

Dogma leans heavily on the concepts of [Domain Driven Design]. It's designed to
provide a suitable platform for applications that make use of design patterns
such as Command/Query Responsibility Segregation ([CQRS]), [Event Sourcing] and
[Eventual Consistency].

The following concepts are core to Dogma's design, and should be well understood
by any developer wishing to build an application:

- [Message](#message)
- [Message handler](#message-handler)
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

A timeout helps model business logic that depends on the passage of time.

Messages must implement the appropriate interface: `Command`, `Event` or
`Timeout`.

### Message handler

A message **handler** is part of an application that acts upon messages it
receives.

Each handler specifies the message types it expects to receive. These message
are **routed** to the handler by the [engine](#engine).

Command messages are always routed to a single handler. Event messages may be
routed to any number of handlers, including zero. Timeout messages are always
routed back to the handler that produced them.

Dogma defines four handler types, one each for [aggregates](#aggregate),
[processes](#process), [integrations](#integration) and
[projections](#projection). These concepts are described in more detail below.

### Application

An **application** is a collection of [message handlers](#message-handler) that
work together as a unit. Typically, each application encapsulates a specific
business (sub-)domain or "bounded-context".

The application is represented by an implementation of the [`dogma.Application`]
interface.

### Engine

An engine is a Go module that delivers messages to an
[application](#application) and persists the application's state.

A Dogma application can run on any Dogma engine. The choice of engine brings
with it a set of guarantees about how the application behaves, for example:

- **Consistency**: Different engines may provide different levels of
  consistency guarantees, such as [immediate consistency] or [eventual
  consistency].

- **Message delivery**: One engine may deliver messages in the same order that
  they were produced, while another may process messages out of order or in
  batches.

- **Persistence**: The engine may offer a choice of persistence mechanisms for
  application state, such as in-memory, on-disk, or in a remote database.

- **Data model**: The engine may provide a choice of data models for
  application state, such as relational or document-oriented.

- **Scalability**: The engine may provide a choice of scalability models, such
  as single-node or multi-node.

This repository is not itself an engine implementation. It defines the API that
engines and applications use to interact.

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

The engine manages each aggregate instance's state. State changes are
["immediately consistent"][immediate consistency] meaning that the changes made
by one command are always visible to future commands routed to the same
instance.

Aggregates can be a difficult concept to grasp. The book [Domain Driven Design
Distilled], by Vaughn Vernon offers a suitable introduction to aggregates and
the other elements of domain driven design.

### Process

A **process** automates a long running business process. In particular, they can
coordinate changes across multiple [aggregate](#aggregate) instances, or between
aggregates and [integrations](#integration).

Like aggregates, processes encapsulate related logic and state. Each
**instance** of a process represents a unique occurrence of that process within
the application.

Each process has an associated implementation of the
[`dogma.ProcessMessageHandler`] interface. The [engine](#engine) routes event
[messages](#message), which produces commands to execute.

A process may use timeout messages to model business processes with time-based
logic. The engine always routes timeout messages back to the process instance
that produced them.

Processes use command messages to make changes to the application's state.
Because each command represents a _separate_ atomic change, the results of a
process are ["eventually consistent"][eventual consistency].

### Integration

An **integration** is a message handler that interacts with some external
non-message-based system.

Each integration is an implementation of the [`dogma.IntegrationMessageHandler`]
interface. The [engine](#engine) routes command [messages](#message) to the
handler which interacts with some external system. Integrations may optionally
produce event messages that represent the results of their interactions.

Integrations are stateless from the perspective of the engine.

### Projection

A **projection** builds a partial view of the application's state from the
events that occur.

Each projection is an implementation of the [`dogma.ProjectionMessageHandler`]
interface. The [engine](#engine) routes event [messages](#message) to the
handler which typically updates a read-optimized database of some kind. This
view is often referred to as a "read model" or "query model".

The [projectionkit] module provides engine-agnostic tools for building
projections in popular database systems, such as PostgreSQL, MySQL, DynamoDB and
others.

<!-- references -->

[`dogma.AggregateMessageHandler`]: https://pkg.go.dev/github.com/dogmatiq/dogma?tab=doc#AggregateMessageHandler
[`dogma.IntegrationMessageHandler`]: https://pkg.go.dev/github.com/dogmatiq/dogma?tab=doc#IntegrationMessageHandler
[`dogma.ProcessMessageHandler`]: https://pkg.go.dev/github.com/dogmatiq/dogma?tab=doc#ProcessMessageHandler
[`dogma.ProjectionMessageHandler`]: https://pkg.go.dev/github.com/dogmatiq/dogma?tab=doc#ProjectionMessageHandler
[api documentation]: https://pkg.go.dev/github.com/dogmatiq/dogma
[cqrs]: https://martinfowler.com/bliki/CQRS.html
[domain driven design distilled]: https://www.amazon.com/Domain-Driven-Design-Distilled-Vaughn-Vernon/dp/0134434420
[domain driven design]: https://en.wikipedia.org/wiki/Domain-driven_design
[event sourcing]: https://martinfowler.com/eaaDev/EventSourcing.html
[eventual consistency]: https://en.wikipedia.org/wiki/Eventual_consistency
[example]: https://github.com/dogmatiq/example
[immediate consistency]: http://www.informit.com/articles/article.aspx?p=2020371&seqNum=2
[projectionkit]: https://github.com/dogmatiq/projectionkit
[rfc 2119]: https://tools.ietf.org/html/rfc2119
[testkit]: https://github.com/dogmatiq/testkit
[veracity]: https://github.com/dogmatiq/veracity
