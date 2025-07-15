# Dogma glossary

[A](#a) •
B •
[C](#c) •
[D](#d) •
[E](#e) •
F •
G •
[H](#h) •
[I](#i) •
J •
K •
L •
[M](#m) •
N •
[O](#o) •
[P](#p) •
Q •
R •
[S](#s) •
[T](#t) •
U •
[V](#v) •
[W](#w) •
X •
Y •
Z

## A

### Aggregate

A collection of related business entities that behave as a cohesive whole. For
example, a shopping cart and the items within it.

> [!NOTE]
> In [domain-driven design], an "aggregate" is a group of related entities
> treated as a single unit — they're dealt with "in aggregate". The term doesn't
> refer to data aggregation or summarization.

> [!TIP]
> "Aggregate" is often used informally to mean [aggregate message handler].

### Aggregate command scope

The [handler scope] in which an [aggregate message handler] handles a [command]
message by recording [event] messages that represent changes to an
[aggregate instance].

See [`dogma.AggregateCommandScope`].

### Aggregate instance

A distinct occurrence of an [aggregate] within an [application]. For example, a
shopping cart aggregate may use an instance for each customer.

### Aggregate message handler

A [message handler] that manages the [state] and behavior of an [aggregate] by
handling [command] messages and recording [event] messages.

See [`dogma.AggregateMessageHandler`].

### Aggregate root

The primary entity within an [aggregate] through which all [state] changes
occur. For example, a shopping cart aggregate that consists of a cart and its
items may use the cart itself as its root.

See [`dogma.AggregateRoot`].

### Application

A collection of related [message handlers] that together define the logic for a
specific business domain.

See [`dogma.Application`].

## C

### Checkpoint offset

The [offset] of the next [event] message that a [projection message handler]
expects to consume from a specific [event stream].

### Command

A [message] that represents a request for the [application] to perform an
action.

See [`dogma.Command`].

### Configurer

An interface used by an [application] or [message handler] to configure its
[identity] and declare its [routes] to the [engine].

See:

- [`dogma.ApplicationConfigurer`]
- [`dogma.AggregateConfigurer`]
- [`dogma.ProcessConfigurer`]
- [`dogma.ProjectionConfigurer`]
- [`dogma.IntegrationConfigurer`]

## D

### Domain-driven design

A software development approach that emphasizes modeling software around the
core concepts and behavior of the business domain it supports.

See [domain-driven design](https://en.wikipedia.org/wiki/Domain-driven_design)
for more information.

## E

### Engine

A runtime component that executes a Dogma [application] by delivering [messages]
to its [message handlers] and persisting application [state].

### Event

A [message] that describes an action that the [application] has already
performed.

See [`dogma.Event`].

### Event sourcing

An architectural pattern in which one or more [event streams] provide the
authoritative representation of an [application]'s [state].

See [event sourcing](https://martinfowler.com/eaaDev/EventSourcing.html) for
more information.

### Event stream

An immutable, ordered sequence of [event] messages.

## H

### Handler route

A declaration by an [application], made using a [configurer], that it includes a
specific [message handler], and incorporates that handler's [message routes].

See:

- [`dogma.ViaAggregate()`]
- [`dogma.ViaProcess()`]
- [`dogma.ViaProjection()`]
- [`dogma.ViaIntegration()`]

### Handler scope

The context within which a [message handler] executes application logic. For
example, when handling an incoming [message], the scope provides information
about that message and defines the messaging operations available to the
handler.

See [`dogma.HandlerScope`].

## I

### Identity

A human-readable name and machine-readable key (UUID), specified via a
[configurer], that uniquely identifies each [application] and [message handler].

### Instance

See [aggregate instance] or [process instance].

### Integration

See [integration message handler].

### Integration command scope

The [handler scope] in which an [integration message handler] handles a
[command] message by interacting with an external system.

See [`dogma.IntegrationCommandScope`].

### Integration message handler

A [stateless] [message handler] that interacts with external systems, such as a
third-party payment API, by handling [command] messages and recording [event]
messages to represent outcomes.

See [`dogma.IntegrationMessageHandler`].

## M

### Message

A data structure that represents a request for — or the prior occurrence of — an
action within an [application].

See:

- [command]
- [event]
- [timeout]
- [`dogma.Message`]

### Message handler

A component of an [application] that acts upon [messages] by producing
new messages or manipulating application [state].

See:

- [aggregate message handler]
- [process message handler]
- [projection message handler]
- [integration message handler]

### Message kind

The category of a [message] that defines its role in the application, one of
[command], [event], or [timeout].

### Message route

A declaration by a [message handler], made using a [configurer], that it
produces or consumes a specific [message type].

See:

- [`dogma.HandlesCommand()`] and [`dogma.ExecutesCommand()`]
- [`dogma.HandlesEvent()`] and [`dogma.RecordsEvent()`]
- [`dogma.SchedulesTimeout()`]

### Message type

A property of a [message] that identifies its specific meaning within the
business domain. Typically represented as a distinct Go type, such as
`AddItemToCart` or `OrderPlaced`.

## O

### Offset

The zero-based position of an [event] message within an [event stream].

## P

### Process

A [stateful] workflow within an [application] that coordinates business logic
involving multiple [aggregate instances], [integrations] or time-sensitive
logic.

> [!TIP]
> "Process" is often used informally to mean [process message handler].

### Process event scope

The [handler scope] in which a [process message handler] handles an [event]
message by updating [state], executing [command] messages, or scheduling [timeout]
messages to advance the [process instance]'s workflow.

See [`dogma.ProcessEventScope`].

### Process instance

A unique occurrence of a [process] within an [application], encapsulating the
[state] of a specific execution of a workflow.

### Process message handler

A [message handler] that orchestrates a [process] by handling [event] messages,
executing [command] messages and scheduling [timeout] messages.

See [`dogma.ProcessMessageHandler`].

### Process root

The primary entity within a [process]'s [state] through which all changes occur.
Named by analogy to [aggregate root].

See [`dogma.ProcessRoot`].

### Process timeout scope

The [handler scope] in which a [process message handler] handles a [timeout]
message at its scheduled time, by updating [state], executing [command]
messages, or scheduling more [timeout] messages to advance the [process instance]'s
workflow.

See [`dogma.ProcessTimeoutScope`].

### Projection

A read-optimized view of a subset of the application's [state] constructed by
observing [event] messages, typically persisted to a database.

> [!TIP]
> "Projection" is often used informally to mean [projection message handler],
> while the resulting view is commonly called a "read-model."

### Projection compact scope

The [handler scope] in which a [projection message handler] compacts its data
by removing or summarizing older or redundant data.

See [`dogma.ProjectionCompactScope`].

### Projection event scope

The [handler scope] in which a [projection message handler] applies an [event]
message to its [projection].

See [`dogma.ProjectionEventScope`].

### Projection message handler

A [stateless] [message handler] that builds a [projection] by handling [event]
messages.

See [`dogma.ProjectionMessageHandler`].

### `projectionkit`

A Go module that provides tools for building [projections] using a range of
popular self-hosted and cloud-based database systems.

See [`dogmatiq/projectionkit`].

## R

### Read-model

See [projection].

### Route

See [message route] or [handler route].

## S

### Saga

A software pattern for managing distributed transactions by applying a sequence
of operations and, if necessary, rolling back changes through compensating
actions.

See [process].

### State

The authoritative representation of an [application]'s current condition, as
managed by the [engine]. It consists of the [event] messages produced by
[aggregate message handlers] and the data stored within [process instances].

> [!IMPORTANT]
> Although "state" is often used informally to mean any application data, here
> it refers strictly to data managed by the [engine]. Data stored in a [projection]
> or an external system accessed by an [integration message handler] isn't
> considered state.

### Stateful

Describes a [message handler] or other component that has [state].

### Stateless

Describes a [message handler] or other component that doesn't have [state].

### Stream

See [event stream].

## T

### `testkit`

A Go module that provides high-level tools for testing Dogma [applications].

See [`dogmatiq/testkit`].

### Timeout

A [message] that describes an action that a [process] should take at a specific
time in the future.

See [`dogma.Timeout`].

## V

### Validation scope

The context within which a [message] executes its data validation logic.

### Veracity

An upcoming [engine] implementation built for horizontal scalability and
distributed workloads.

See [`dogmatiq/veracity`].

### Verity

An [engine] designed for typical application loads in smaller deployments.

See [`dogmatiq/verity`].

## W

### Workflow

See [process].

<!-- anchors -->

[aggregate instance]: #aggregate-instance
[aggregate instances]: #aggregate-instance
[aggregate message handler]: #aggregate-message-handler
[aggregate message handlers]: #aggregate-message-handler
[aggregate root]: #aggregate-root
[aggregate]: #aggregate
[application]: #application
[applications]: #application
[command]: #command
[configurer]: #configurer
[domain-driven design]: #domain-driven-design
[engine]: #engine
[event stream]: #event-stream
[event streams]: #event-stream
[event]: #event
[handler route]: #handler-route
[handler scope]: #handler-scope
[identity]: #identity
[integration message handler]: #integration-message-handler
[integration message handlers]: #integration-message-handler
[integrations]: #integration
[message handler]: #message-handler
[message handlers]: #message-handler
[message route]: #message-route
[message routes]: #message-route
[message type]: #message-type
[message]: #message
[messages]: #message
[offset]: #offset
[process instance]: #process-instance
[process instances]: #process-instance
[process message handler]: #process-message-handler
[process]: #process
[projection message handler]: #projection-message-handler
[projection]: #projection
[projections]: #projection
[routes]: #route
[state]: #state
[stateful]: #stateful
[stateless]: #stateless
[timeout]: #timeout

<!-- go modules -->

[`dogmatiq/projectionkit`]: https://pkg.go.dev/github.com/dogmatiq/projectionkit
[`dogmatiq/testkit`]: https://pkg.go.dev/github.com/dogmatiq/testkit
[`dogmatiq/veracity`]: https://pkg.go.dev/github.com/dogmatiq/veracity
[`dogmatiq/verity`]: https://pkg.go.dev/github.com/dogmatiq/verity

<!-- API references -->

[`dogma.AggregateCommandScope`]: https://pkg.go.dev/github.com/dogmatiq/dogma#AggregateCommandScope
[`dogma.AggregateConfigurer`]: https://pkg.go.dev/github.com/dogmatiq/dogma#AggregateConfigurer
[`dogma.AggregateMessageHandler`]: https://pkg.go.dev/github.com/dogmatiq/dogma#AggregateMessageHandler
[`dogma.AggregateRoot`]: https://pkg.go.dev/github.com/dogmatiq/dogma#AggregateRoot
[`dogma.Application`]: https://pkg.go.dev/github.com/dogmatiq/dogma#Application
[`dogma.ApplicationConfigurer`]: https://pkg.go.dev/github.com/dogmatiq/dogma#ApplicationConfigurer
[`dogma.Command`]: https://pkg.go.dev/github.com/dogmatiq/dogma#Command
[`dogma.Event`]: https://pkg.go.dev/github.com/dogmatiq/dogma#Event
[`dogma.ExecutesCommand()`]: https://pkg.go.dev/github.com/dogmatiq/dogma#ExecutesCommand
[`dogma.HandlerScope`]: https://pkg.go.dev/github.com/dogmatiq/dogma#HandlerScope
[`dogma.HandlesCommand()`]: https://pkg.go.dev/github.com/dogmatiq/dogma#HandlesCommand
[`dogma.HandlesEvent()`]: https://pkg.go.dev/github.com/dogmatiq/dogma#HandlesEvent
[`dogma.IntegrationCommandScope`]: https://pkg.go.dev/github.com/dogmatiq/dogma#IntegrationCommandScope
[`dogma.IntegrationConfigurer`]: https://pkg.go.dev/github.com/dogmatiq/dogma#IntegrationConfigurer
[`dogma.IntegrationMessageHandler`]: https://pkg.go.dev/github.com/dogmatiq/dogma#IntegrationMessageHandler
[`dogma.Message`]: https://pkg.go.dev/github.com/dogmatiq/dogma#Message
[`dogma.ProcessConfigurer`]: https://pkg.go.dev/github.com/dogmatiq/dogma#ProcessConfigurer
[`dogma.ProcessEventScope`]: https://pkg.go.dev/github.com/dogmatiq/dogma#ProcessEventScope
[`dogma.ProcessMessageHandler`]: https://pkg.go.dev/github.com/dogmatiq/dogma#ProcessMessageHandler
[`dogma.ProcessRoot`]: https://pkg.go.dev/github.com/dogmatiq/dogma#ProcessRoot
[`dogma.ProcessTimeoutScope`]: https://pkg.go.dev/github.com/dogmatiq/dogma#ProcessTimeoutScope
[`dogma.ProjectionCompactScope`]: https://pkg.go.dev/github.com/dogmatiq/dogma#ProjectionCompactScope
[`dogma.ProjectionConfigurer`]: https://pkg.go.dev/github.com/dogmatiq/dogma#ProjectionConfigurer
[`dogma.ProjectionEventScope`]: https://pkg.go.dev/github.com/dogmatiq/dogma#ProjectionEventScope
[`dogma.ProjectionMessageHandler`]: https://pkg.go.dev/github.com/dogmatiq/dogma#ProjectionMessageHandler
[`dogma.RecordsEvent()`]: https://pkg.go.dev/github.com/dogmatiq/dogma#RecordsEvent
[`dogma.SchedulesTimeout()`]: https://pkg.go.dev/github.com/dogmatiq/dogma#SchedulesTimeout
[`dogma.Timeout`]: https://pkg.go.dev/github.com/dogmatiq/dogma#Timeout
[`dogma.ViaAggregate()`]: https://pkg.go.dev/github.com/dogmatiq/dogma#ViaAggregate
[`dogma.ViaIntegration()`]: https://pkg.go.dev/github.com/dogmatiq/dogma#ViaIntegration
[`dogma.ViaProcess()`]: https://pkg.go.dev/github.com/dogmatiq/dogma#ViaProcess
[`dogma.ViaProjection()`]: https://pkg.go.dev/github.com/dogmatiq/dogma#ViaProjection
