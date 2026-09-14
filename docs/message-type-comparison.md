# Message kind comparison

This document compares the three Dogma message kinds, highlighting their roles
in the application, how handlers produce and consume them, and the delivery
behavior you can rely on.

If you're new to Dogma, start by familiarizing yourself with the key [concepts].

> [!TIP]
> Consult the [glossary] for definitions of terms used in this document.

## Comparison matrix

|                       |         Command          |                 Event                 |       Deadline        |
| --------------------- | :----------------------: | :-----------------------------------: | :-------------------: |
| Represents            |        a request         |                a fact                 |   a moment in time    |
| [Producers]           | processes, external code |       aggregates, integrations        |       processes       |
| [Consumers]           | aggregates, integrations | processes, projections, external code |     its producer      |
| [Delivery]            |       immediately        |       immediately, historically       | at the scheduled time |
| [Ordering guarantees] |           none           |           per event-stream            |         none          |
| [Permanence]          |        transient         |               permanent               |       transient       |

## Definitions

### Producers

A message's producers are the components that create it and introduce it into
the application.

Command messages enter the application from two directions. **Process** message
handlers execute commands to drive their workflows forward, and **external
code** — such as an HTTP API or user interface — executes commands to act on
requests from users and other systems. Any number of producers may execute
commands of the same type.

Event messages, by contrast, come from a single source — exactly one handler,
either an **aggregate** or an **integration**, records each event type. That
handler is the sole authority for the fact that the event represents.

Deadline messages originate only from **process** message handlers, which
schedule each one for delivery at a specific future time. Any number of
processes may schedule deadlines of the same type.

### Consumers

A message's consumers are the handlers that receive it and act upon it.

The engine routes each command type to exactly one handler — an **aggregate**
or an **integration** — giving every action a single, unambiguous
implementation.

Any number of **process** and **projection** message handlers may consume each
event type. Recording an event once lets many parts of the application react to
it independently. **External code** can also observe events in a limited way —
when executing a command, it may attach an event observer to wait for a specific
event recorded as a result of that command, confirming that the requested action
occurred.

The engine delivers each deadline message back to **its producer** — the same
process instance that scheduled it. No other handler or instance ever sees it.

### Delivery

Delivery describes when a consumer receives each message.

The engine delivers each command message to its handler **immediately** — as
soon as possible after accepting it.

The engine delivers each event message to its consumers **immediately** —
asynchronously, shortly after it occurs. Because events are [permanent], the
engine can also deliver them **historically**. This capability serves
projections specifically: adding a new projection to the application, or
resetting an existing one, builds its data by consuming the complete event
history from the beginning of each stream.

The engine delivers each deadline message **at the scheduled time** — never
before, though possibly later, such as when the engine recovers from downtime
or retries after a failure.

### Ordering guarantees

Ordering guarantees describe the delivery sequence that consumers can rely on.

Command messages have no guaranteed order — the engine may deliver commands in
any order. Avoid logic that assumes one command arrives before another; when
one action must follow another, coordinate them with a process.

Event messages are ordered per **event-stream**. Consumers receive the events
on each [event stream] in the order they occurred, while events on different
streams may arrive in any order. [Event stream binding] determines which
events share a stream — in practice, events recorded against the same
aggregate instance, or by a single command, always arrive in order.

Deadline messages have no guaranteed order either — each arrives according to
its own scheduled time, but a delayed deadline may arrive after one scheduled
for a later time, and deadlines with the same scheduled time may arrive in any
order.

### Permanence

Permanence describes whether a message remains part of the application's
[state] after it's handled.

Event messages are **permanent**. Dogma treats them as the primary source of
truth — each recorded event becomes part of the application's history, the
foundation of [event sourcing].

Command and deadline messages are **transient**. They drive the application's
behavior, but form no part of its state once handled. The engine also cancels
a process instance's pending deadlines when the instance ends.

<!-- anchors -->

[consumers]: #consumers
[delivery]: #delivery
[ordering guarantees]: #ordering-guarantees
[permanence]: #permanence
[permanent]: #permanence
[producers]: #producers

<!-- external links -->

[concepts]: concepts.md
[event sourcing]: concepts.md#event-sourcing
[event stream]: glossary.md#event-stream
[event stream binding]: handler-type-comparison.md#event-stream-binding
[glossary]: glossary.md
[state]: glossary.md#state
