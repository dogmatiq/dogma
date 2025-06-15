<div align="center">

# Dogma

Build message-based, event-sourced applications in Go.

[![Documentation](https://img.shields.io/badge/go.dev-documentation-007d9c?&style=for-the-badge)](https://pkg.go.dev/github.com/dogmatiq/dogma)
[![Latest Version](https://img.shields.io/github/tag/dogmatiq/dogma.svg?&style=for-the-badge&label=semver)](https://github.com/dogmatiq/dogma/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/dogmatiq/dogma/ci.yml?style=for-the-badge&branch=main)](https://github.com/dogmatiq/dogma/actions/workflows/ci.yml)
[![Code Coverage](https://img.shields.io/codecov/c/github/dogmatiq/dogma/main.svg?style=for-the-badge)](https://codecov.io/github/dogmatiq/dogma)

</div>

## Overview

Dogma is a comprehensive suite of tools for building robust message-driven
applications in Go.

It provides an abstraction for describing your application’s business logic with
strict separation from components responsible for message delivery and
persistence.

## Features

- **Event-sourced by design** – Every state change is persisted as an immutable
  domain event. This enables full auditability and allows read-optimized views
  to be built or rebuilt from the event history at any time.

- **Grounded in [Domain-Driven Design]** – Dogma adopts core DDD concepts to
  guide how applications are decomposed and how messages flow between
  components.

- **High-level testing** – The [testkit] module encourages verification of
  application behavior by making assertions about domain events rather than
  inspecting state. It integrates seamlessly with Go’s standard [testing]
  package.

- **Native introspection** – Dogma's static analysis tools visualize message
  flow and application structure, enabling discovery of domain events across
  large codebases and multi-application projects.

- **Domain and infrastructure separation** – Domain logic is cleanly and
  strictly separated from infrastructure concerns such as message delivery and
  persistence.

- **Type-agnostic** – Messages and application state can be any Go type that
  marshals to a byte slice, with built-in support for JSON and Protocol Buffers.

- **Flexible persistence** – Support for a range of storage options such as
  PostgreSQL and cloud services like DynamoDB, enabling use across diverse
  environments.

## Repositories

This repository contains the Go interfaces that form the contract between the
application and the engine.

- [testkit]: Utilities for black-box testing of Dogma applications.
- [projectionkit]: Utilities for building [projections](#projection) in popular database systems.
- [example]: An example Dogma application that implements basic banking features.

## Why "Dogma"?

The name _Dogma_ is a tongue-in-cheek nod to the project's strong opinions about
how message-driven applications should be structured. It's not about rigid
rule-following, but about embracing consistent patterns that enable rich tooling
and clarity in complex systems — without sacrificing flexibility where it
matters.

## Concepts

Dogma leans heavily on the concepts of [Domain-Driven Design]. It's designed to
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

Messages must implement the appropriate interface; one of [`dogma.Command`],
[`dogma.Event`] or [`dogma.Timeout`].

### Message handler

A message **handler** is a component of an application that acts upon messages
it receives.

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

Each application is represented by an implementation of the
[`dogma.Application`] interface.

### Engine

An engine is a Go module that delivers messages to an
[application](#application) and persists the application's state.

A Dogma application can run on any Dogma engine. The choice of engine brings
with it a set of guarantees about how the application behaves, for example:

- **Consistency**: Different engines may provide different levels of
  consistency guarantees, such as [immediate consistency] or [eventual
  consistency].

- **Persistence**: The engine may offer a choice of persistence mechanisms for
  application state, such as in-memory, on-disk, or in a remote database.

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

Aggregates can be a difficult concept to grasp. The book [Domain-Driven Design
Distilled], by Vaughn Vernon offers a suitable introduction to aggregates and
the other elements of domain-driven design.

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
logic. The engine routes timeout messages back to the process instance that
produced them at the specified time.

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
view is often referred to as a "read model".

The [projectionkit] module provides tools for building read-models in popular
database systems, such as PostgreSQL, MySQL, DynamoDB and others.

<!-- references -->

[`dogma.aggregatemessagehandler`]: https://pkg.go.dev/github.com/dogmatiq/dogma?tab=doc#AggregateMessageHandler
[`dogma.application`]: https://pkg.go.dev/github.com/dogmatiq/dogma?tab=doc#Application
[`dogma.integrationmessagehandler`]: https://pkg.go.dev/github.com/dogmatiq/dogma?tab=doc#IntegrationMessageHandler
[`dogma.processmessagehandler`]: https://pkg.go.dev/github.com/dogmatiq/dogma?tab=doc#ProcessMessageHandler
[`dogma.projectionmessagehandler`]: https://pkg.go.dev/github.com/dogmatiq/dogma?tab=doc#ProjectionMessageHandler
[`dogma.command`]: https://pkg.go.dev/github.com/dogmatiq/dogma?tab=doc#Command
[`dogma.event`]: https://pkg.go.dev/github.com/dogmatiq/dogma?tab=doc#Event
[`dogma.timeout`]: https://pkg.go.dev/github.com/dogmatiq/dogma?tab=doc#Timeout
[api documentation]: https://pkg.go.dev/github.com/dogmatiq/dogma
[cqrs]: https://martinfowler.com/bliki/CQRS.html
[domain-driven design distilled]: https://www.amazon.com/Domain-Driven-Design-Distilled-Vaughn-Vernon/dp/0134434420
[domain-driven design]: https://en.wikipedia.org/wiki/Domain-driven_design
[event sourcing]: https://martinfowler.com/eaaDev/EventSourcing.html
[eventual consistency]: https://en.wikipedia.org/wiki/Eventual_consistency
[example]: https://github.com/dogmatiq/example
[immediate consistency]: http://www.informit.com/articles/article.aspx?p=2020371&seqNum=2
[projectionkit]: https://github.com/dogmatiq/projectionkit
[rfc 2119]: https://tools.ietf.org/html/rfc2119
[testing]: https://pkg.go.dev/testing
[testkit]: https://github.com/dogmatiq/testkit
[veracity]: https://github.com/dogmatiq/veracity
