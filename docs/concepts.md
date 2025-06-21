# Dogma concepts

[Dogma] is a comprehensive suite of tools for building robust message-driven
applications in Go.

This document is a primer on Dogma's core concepts — it introduces the
terminology, mental models, and architecture patterns you need to understand
before building an application.

We recommend reading this document from start to finish, as each section builds
on the last. For reference material, see the [API documentation].

> [!TIP]
> Text styled **in bold** introduces a word or phrase with a specific meaning in
> the Dogma ecosystem.

## What _is_ a "message-driven" application?

Dogma uses **messages** to describe both what should happen and what has already
occurred. Your **application** reacts to messages one at a time — a message
comes in, the application performs some action, and new messages come out. Each
step is self-contained and straightforward to understand. This is what we mean
by a message-driven application.

## Messages

A message is a data structure that describes something specific your application
should do, or something it has already done.

Dogma defines three **kinds** of message:

- A **command** message is a request — it tells your application to do something
  immediately.

  For example, _add 10 widgets to Alex's shopping cart_.

- An **event** message is a fact — it represents something your application has
  already done.

  For example, _added 10 widgets to Alex's shopping cart_.

- A **timeout** message is a delayed notification — it tells your application
  that some relevant period of time has elapsed.

  For example, _it's been 24 hours since the first item was added to Alex's cart_.

Every distinct action — such as "add item to cart" or "complete purchase" —
corresponds to a different **message type**, each represented by a Go type that
implements one of the [`Command`], [`Event`], or [`Timeout`] interfaces.

> [!IMPORTANT]
> Messages must not rely on information that may be unavailable when they're
> handled. Note how each earlier example explicitly includes the subject —
> _Alex_ — rather than an inferred "current user". This approach captures intent
> within the message itself, avoiding ambiguity in the handling logic.

## Handlers

A **message handler** is a component of your application that acts upon
[messages] it receives by producing new messages, updating state, or interacting
with external systems.

Dogma defines four types of handler:

- An **aggregate message handler** manages a group of related entities, such as
  a shopping cart and the items within it. It handles commands by recording
  events that reflect changes to the entities' state. Each aggregate in your
  application has many **instances** — for example, a shopping cart aggregate
  might use an instance per customer. Aggregates are the main building block of
  your application's logic.

- A **process message handler** coordinates a workflow that involves multiple
  aggregate instances or time-sensitive logic. It handles event messages by
  executing commands and scheduling timeout messages to drive the workflow
  forward. Like aggregates, each process can have many instances.

- A **projection message handler** builds a view of the application's state by
  observing event messages. This view, often called a **read-model**, is
  typically stored in a database and optimised for querying or presentation.
  Projections don't produce messages of any kind.

- An **integration message handler** perform actions outside the application,
  such as sending emails or processing payments using a third-party API. Like
  aggregates, integrations handle command messages and record event messages to
  describe what occurred. From Dogma's perspective, integrations are stateless.

Each handler in your application corresponds to a Go type that implements one
of the [`AggregateMessageHandler`], [`ProcessMessageHandler`],
[`IntegrationMessageHandler`], or [`ProjectionMessageHandler`] interfaces.

> [!NOTE]
> In [domain-driven design], an "aggregate" is a group of related entities
> treated as a single unit — they're dealt with "in aggregate". The term doesn't
> refer to data aggregation or summarization.

## Scopes

When a [message handler] handles a [message], it does so within a specific
**scope**.

The scope has two main roles:

- It provides information about the incoming message — for example, the time
  when an event occurred.
- It defines the messaging operations that the handler can perform in response
  to the message, such as executing commands or recording events.

Dogma defines multiple scopes types, each represented by a separate Go
interface. For example, the [`AggregateCommandScope`] interface represents the
scope in which an aggregate handles a command message. Your handler receives a
scope with each incoming message — you don't need to implement these interfaces
yourself.

## Event sourcing

Dogma treats event [messages] as the primary source of truth. When your
application performs an action, it records an event that becomes part of its
permanent history — the foundation of an approach known as [event sourcing].

We can derive different representations of the application's _state_ from these
events at any time — a powerful capability that lets you evolve the structure of
your data and introduce new views without touching your business logic.

Consult the [API documentation] for details on how each handler type makes use
of this event history.

## Applications

We've made frequent mention of your **application**, but what exactly is an
application in Dogma?

In practice, your application likely includes a user interface, APIs,
authentication, and more. In Dogma, "application" has a narrower meaning — it
refers to a collection of [message handlers] that work together to implement a
particular business domain. The other components of your software are outside
Dogma's responsibility — it imposes no constraints on how they're built or what
technologies you use.

Although the handlers within your application declare which [messages] they
consume and produce, the message types themselves aren't, strictly speaking,
part of any one application. Your broader application may consist of multiple
Dogma applications that communicate using messages.

In code, each application corresponds to a Go type that implements the
[`Application`] interface.

> [!TIP]
> If you're familiar with [domain-driven design], a Dogma application roughly
> aligns with the concept of a bounded context.

## Engines

We've discussed the _interfaces_ that you implement and use to build your
[application]. To actually _run_ your application, you need an **engine**.

Engines aren't part of the [`dogmatiq/dogma`] Go module, they're separate
modules that implement the runtime behaviour described by Dogma's interfaces.
You can choose an engine that suits your environment, or build your own.

The ecosystem currently offers three official engines:

- [`dogmatiq/verity`] — The original Dogma engine, designed for typical
  application loads in smaller deployments. While production-ready, it doesn't
  support scaling of a single application across multiple machines.

- [`dogmatiq/veracity`] — The next-generation Dogma engine built for
  horizontal scalability and distributed workloads. The Dogma maintainers intend
  for Veracity to fully replace Verity, becoming _the_ production Dogma engine.

- [`dogmatiq/testkit`] — A set of tools for testing Dogma applications. It
  includes an in-memory engine that executes and inspects application behavior
  without persisting state.

## What's next?

Now that you have a high-level understanding of Dogma's concepts, you can
explore the following resources:

- [API documentation] — detailed information about Dogma's API.
- [Handler type comparison] — a comparison of Dogma's four message handler
  types.
- [Glossary] — a central reference for Dogma's terminology.
- [`dogmatiq/example`] — a minimal example application with basic banking
  features.

<!-- anchors -->

[message]: #messages
[messages]: #messages
[message handler]: #handlers
[message handlers]: #handlers
[application]: #applications

<!-- other documentation  -->

[dogma]: https://github.com/dogmatiq/dogma?tab=readme-ov-file#readme
[glossary]: glossary.md
[handler type comparison]: handler-type-comparison.md

<!-- go modules -->

[`dogmatiq/example`]: https://github.com/dogmatiq/example
[`dogmatiq/testkit`]: https://github.com/dogmatiq/testkit
[`dogmatiq/veracity`]: https://github.com/dogmatiq/veracity
[`dogmatiq/verity`]: https://github.com/dogmatiq/verity

<!-- API references -->

[api documentation]: https://pkg.go.dev/github.com/dogmatiq/dogma
[`AggregateCommandScope`]: https://pkg.go.dev/github.com/dogmatiq/dogma#AggregateCommandScope
[`AggregateMessageHandler`]: https://pkg.go.dev/github.com/dogmatiq/dogma#AggregateMessageHandler
[`Application`]: https://pkg.go.dev/github.com/dogmatiq/dogma#Application
[`Command`]: https://pkg.go.dev/github.com/dogmatiq/dogma#Command
[`dogmatiq/dogma`]: https://pkg.go.dev/github.com/dogmatiq/dogma
[`Event`]: https://pkg.go.dev/github.com/dogmatiq/dogma#Event
[`IntegrationMessageHandler`]: https://pkg.go.dev/github.com/dogmatiq/dogma#IntegrationMessageHandler
[`ProcessMessageHandler`]: https://pkg.go.dev/github.com/dogmatiq/dogma#ProcessMessageHandler
[`ProjectionMessageHandler`]: https://pkg.go.dev/github.com/dogmatiq/dogma#ProjectionMessageHandler
[`Timeout`]: https://pkg.go.dev/github.com/dogmatiq/dogma#Timeout

<!-- external references -->

[domain-driven design]: https://en.wikipedia.org/wiki/Domain-driven_design
[event sourcing]: https://martinfowler.com/eaaDev/EventSourcing.html
