# Dogma Concepts

[Dogma] is a comprehensive suite of tools for building robust message-driven
applications in Go.

This document explains the core concepts you need to understand when writing
applications with [Dogma]. It complements the [API documentation], which
describes the interfaces and types you'll use to implement your application

## Message-driven applications

Dogma uses **messages** to describe what should, or already has, occurred. Your
application reacts to messages one at a time: a message comes in, some logic is
applied, and new messages come out. Each step is self-contained and easy to
understand.

For the most part, when referring to the **application**, we mean those
components that interact directly with Dogma — the message handlers and the
messages they exchange. Other components, such as your user interfaces, APIs and
authentication mechanisms are outside Dogma's scope.

## Messages

A **message** is a data structure that describes something specific your
application should do, or something it has already done.

There are three kinds of messages:

- A **command** describes an action that your application should take.<br />
  For example, _add 10 widgets to Alex's shopping cart_.

- An **event** describes an action that your application has already taken.<br />
  For example, _10 widgets were added to Alex's shopping cart_.

- A **timeout** describes an action that should be taken in the future.<br />
  For example, _send Alex a reminder about their incomplete purchase after 24 hours_.

Each distinct action — such as "add item to cart" or "complete purchase" — is
represented by a different **message type**. These types map directly to Go
types that implement one of the [`dogma.Command`], [`dogma.Event`], or
[`dogma.Timeout`] interfaces.

> [!Important]
> Note that each example explicitly includes the subject — Alex — rather than
> an implied "current user". Messages must contain all the information necessary
> to act on them. This approach "captures intent" at the time the message is
> created, leaving no room for ambiguity when the message is handled.

## Handlers

A **message handler** is a component of your application that acts upon messages
it receives by producing new messages, updating state, or interacting with
external systems.

There are four kinds of handlers, each with a different purpose.

- An **[aggregate message handler](#aggregate)** manages the state of a specific
  entity, such as a shopping cart or user account. It handles command messages
  and records event messages that reflect changes to the entity's state.
  Aggregates are the main building block of your application's logic.

- A **[process message handler](#process)** coordinates a workflow that involves
  multiple aggregates. It handles event messages by executing commands to drive
  the workflow forward. It may also schedule timeout messages to trigger actions
  at specific times in the future.

- A **[projection message handler](#projection)** builds a view of the
  application's state by observing event messages. This view, called a
  **read-model**, is typically stored in a database and optimised for querying
  or presentation. Projections do not produce new messages.

- An **[integration message handler](#integration)** performs an action outside
  the application, such as sending an email or processing a payment using a
  third-party API. Like aggregates, integrations handle command messages and
  record event messages to describe what occurred. From Dogma’s perspective,
  integrations are stateless — any required state exists outside the application
  itself.

## Event-sourcing

WIP/TODO

Dogma applications treat event messages as the primary source of truth.
Whenever an event is recorded, it becomes part of the application’s permanent
history.

Different handlers can use this history in different ways:

- [Aggregates](#aggregate) can reconstruct their state by replaying past events.
- [Processes](#process) ...

<!-- references -->

[dogma]: https://github.com/dogmatiq/dogma
[api documentation]: https://pkg.go.dev/github.com/dogmatiq/dogma
[`dogma.Command`]: https://pkg.go.dev/github.com/dogmatiq/dogma#Command
[`dogma.Event`]: https://pkg.go.dev/github.com/dogmatiq/dogma#Event
[`dogma.Timeout`]: https://pkg.go.dev/github.com/dogmatiq/dogma#Timeout
[messages]: #messages
