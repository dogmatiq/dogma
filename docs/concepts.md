# Dogma Concepts

[Dogma] is a comprehensive suite of tools for building robust message-driven
applications in Go.

This document is a primer on Dogma’s core concepts - it introduces the
terminology, mental models, and architecture patterns you need to understand
before building an application.

We recommend reading this document from start to finish, as each section builds
on the last. For reference material, see the [API documentation].

> [!TIP]
> Text styled **in bold** introduces a word or phrase with a specific meaning
> in the Dogma ecosystem.

## What is a "message-driven" application?

Dogma uses **messages** to describe both what should happen and what has already
occurred. Your **application** reacts to messages one at a time - a message
comes in, application logic is applied, and new messages come out. Each step is
self-contained and easy to reason about. This is what we mean by a
message-driven application.

## Messages

A message is a data structure that describes something specific your application
should do, or something it has already done.

There are three **kinds** of message:

- A **command** message is a request - an action your application should take.

  For example, _add 10 widgets to Alex's shopping cart_.

- An **event** message is a fact - an action your application has already taken.

  For example, _10 widgets were added to Alex's shopping cart_.

- A **timeout** message is a delayed request - an action that should be taken later.

  For example, _send Alex a reminder about their incomplete purchase after 24 hours_.

Every distinct action - such as "add item to cart" or "complete purchase" - is
represented by a different **message type**. Each message type is a Go type that
implements one of the [`Command`], [`Event`], or [`Timeout`] interfaces.

> [!IMPORTANT]
> Messages must contain all the information necessary to act on them. Notice how
> each example above explicitly references the subject - Alex - rather than an
> implied "current user". This approach "captures intent" at the time the
> message is created, leaving no room for ambiguity when it is handled.

## Handlers

A **message handler** is a component of your application that acts upon
[messages] it receives by producing new messages, updating state, or interacting
with external systems.

There are four types of handler:

- An **aggregate message handler** manages a group of related entities, such
  as a shopping cart and the items within it. It handles commands by recording
  events that reflect changes to the entities' state. Each aggregate in your
  application has many **instances** - for example, a shopping cart aggregate
  might use an instance per customer. Aggregates are the main building block of
  your application's logic.

  > [!NOTE]
  > The word "aggregate" is often a source of confusion.
  >
  > The terminology comes from [domain-driven design], where it refers to a
  > group of related entities that are treated as a single unit - they are dealt
  > with "in aggregate". It does not refer to data aggregation or summarization.

- A **process message handler** coordinates a workflow that involves multiple
  aggregate instances. It handles event messages by executing commands to drive
  the workflow forward. It may also schedule timeout messages to trigger actions
  at specific times. Like aggregates, each process can have many instances.

- A **projection message handler** builds a view of the application's state by
  observing event messages. This view, called a **read-model**, is typically
  stored in a database and optimised for querying or presentation. Projections
  cannot produce new messages.

- An **integration message handler** performs an action outside the
  application, such as sending an email or processing a payment using a
  third-party API. Like aggregates, integrations handle command messages and
  record event messages to describe what occurred. From Dogma’s perspective,
  integrations are stateless.

Each handler in your application is represented by a Go type that implements one
of the [`AggregateMessageHandler`], [`ProcessMessageHandler`],
[`IntegrationMessageHandler`], or [`ProjectionMessageHandler`] interfaces.

## Scopes

When a [handler] handles a [message], it does so within a specific **scope**.

The scope has two main roles:

- It provides information about the incoming message - for example, the time at
  which an event was recorded.
- It defines the messaging operations that the handler can perform in response
  to the message, such as executing commands or recording events.

There are several kinds of scopes, each represented by a separate Go interface.
For example, the [`AggregateCommandScope`] interface represents the scope in
which an aggregate handles a command message. Your handler receives a scope with
each incoming message - you do not need to implement these interfaces yourself.

## Event-sourcing

Dogma treats event [messages] as the primary source of truth. When an event is
recorded, it becomes part of the application's permanent history. The
application's _state_ is derived from these events - it is a reflection of what
has occurred. This approach is known as [event-sourcing].

Consult the [api documentation] for details on how each handler type makes use
of this event history.

## Applications

We’ve made frequent mention of your **application**, but what exactly is an
application in Dogma?

In practice, your application likely includes a user interface, APIs,
authentication mechanisms, and more. Within Dogma, the term simply refers to a
collection of related [message handlers] that together implement the logic for a
particular business domain. The other components of your software are not
relevant to Dogma. Accordingly, it imposes no constraints on how they’re built
or what technologies you use.

Although the handlers within your application declare which [messages] they
consume and produce, the message types themselves are not, strictly speaking,
part of any one application. Your broader application may consist of multiple
Dogma applications, and messages can be used to communicate between them.

In code, each application is represented by a Go type that implements the
[`Application`] interface.

> [!TIP]
> If you're familiar with [domain-driven design], a Dogma application roughly
> aligns with the concept of a bounded context.

## Engines

So far, the code we've referenced has been limited to the interfaces that you
implement and use to build your [application]. To actually _run_ your
application, you need an **engine**.

Engines are not part of the [`dogma`] Go module - they’re external packages that
implement the runtime behaviour described by Dogma’s interfaces. You can choose
an engine that suits your environment, or build your own.

There are currently three official engines:

- [verity] - The original Dogma engine, designed for typical application loads
  in smaller deployments. While production-ready, it does not support horizontal
  scaling of individual applications; it uses a failover model instead.

- [veracity] _(under development)_ - The next-generation Dogma engine built for
  horizontal scalability and distributed workloads. In the long term, it will
  fully replace Verity, becoming _the_ production Dogma engine.

- [testkit] - A set of tools for testing Dogma applications. It includes an
  in-memory engine that executes and inspects application behavior without
  persisting state.

## What's next?

Now that you have a high-level understanding of Dogma's concepts, you can
explore the following resources:

- The [example] repository, which contains a simple banking application with
  features such as opening accounts and transferring funds.
- The [API documentation], for detailed information on Dogma's interfaces.

<!-- anchors -->

[message]: #messages
[messages]: #messages
[handler]: #handlers
[handlers]: #handlers
[message handler]: #handlers
[message handlers]: #handlers
[application]: #applications
[applications]: #applications

<!-- go modules -->

[dogma]: https://github.com/dogmatiq/dogma
[example]: https://github.com/dogmatiq/example
[testkit]: https://github.com/dogmatiq/testkit
[veracity]: https://github.com/dogmatiq/veracity
[verity]: https://github.com/dogmatiq/verity

<!-- references -->

[api documentation]: https://pkg.go.dev/github.com/dogmatiq/dogma
[`AggregateCommandScope`]: https://pkg.go.dev/github.com/dogmatiq/dogma#AggregateCommandScope
[`AggregateMessageHandler`]: https://pkg.go.dev/github.com/dogmatiq/dogma#AggregateMessageHandler
[`Application`]: https://pkg.go.dev/github.com/dogmatiq/dogma#Application
[`Command`]: https://pkg.go.dev/github.com/dogmatiq/dogma#Command
[`Event`]: https://pkg.go.dev/github.com/dogmatiq/dogma#Event
[`IntegrationMessageHandler`]: https://pkg.go.dev/github.com/dogmatiq/dogma#IntegrationMessageHandler
[`ProcessMessageHandler`]: https://pkg.go.dev/github.com/dogmatiq/dogma#ProcessMessageHandler
[`ProjectionMessageHandler`]: https://pkg.go.dev/github.com/dogmatiq/dogma#ProjectionMessageHandler
[`Timeout`]: https://pkg.go.dev/github.com/dogmatiq/dogma#Timeout

<!-- external references -->

[domain-driven design]: https://en.wikipedia.org/wiki/Domain-driven_design
[event-sourcing]: https://martinfowler.com/eaaDev/EventSourcing.html
