# Handler type comparison

This document compares the four Dogma handler types, highlighting their roles in
message flow, data access, and external interaction.

If you're new to Dogma, start by familiarising yourself with the key [concepts].

> [!TIP]
> Consult the [glossary] for definitions of terms used in this document.

## Comparison matrix

|                        | Aggregate | Process  | Integration |  Projection  |
| ---------------------- | :-------: | :------: | :---------: | :----------: |
| Command messages       |     ↓     |    ↑     |      ↓      |              |
| Event messages         |     ↑     |    ↓     |      ↑      |      ↓       |
| Timeout messages       |           |    ⇅     |             |              |
| [State]                |     ✓     |    ✓     |             |              |
| [External data access] |           |    ~     |      ✓      |      ✓       |
| [Isolation boundary]   | instance  | instance |   command   | event-stream |
| [Event stream binding] | instance  |    —     |   command   |      —       |

#### Legend

- `✓` supported
- `~` supported, but not recommended
- `↓` consumes
- `↑` produces
- `⇅` produces &amp; consumes
- `—` not applicable

## Definitions

### State

Aggregates and processes maintain state across many separate instances.
The engine routes each message to a specific instance using an identifier
provided by the application.

Each aggregate instance records events, and each process instance holds data.
Together, this represents the application's state.

Although both integration and projection message handlers may use or store data,
it's not considered state within Dogma. See the [definition of state] in the
glossary for further clarification.

### External data access

External data is any information not managed by the engine — including
projection data, third-party APIs, the system clock, and so on.

- Aggregate message handlers must not read or write external data. All logic
  must rely solely on information within command messages and the instance's
  state. This restriction ensures that aggregates are deterministic — they
  produce the same events for the same sequence of commands.

- Process message handlers may read external data during message routing and
  handling, primarily to access the application's projection data, though this
  is not recommended. Wherever possible, all logic should rely solely on
  information within event messages and the instance's state. This restriction
  ensures that processes are deterministic — they produce the same behavior for
  the same sequence of events.

- Integration message handlers may read and write external data — they exist
  specifically to interact with external systems.

- Projection message handlers may read and write external data. Where possible,
  all logic should rely solely on information within event messages and the
  existing projection data. This restriction ensures that projections are
  self-contained and deterministic — they produce the same data for the same
  sequence of events.

### Isolation boundary

The isolation boundary describes how the engine groups and orders a message
handler's side-effects.

- Aggregate and process message handlers use per-**instance** isolation. For
  each instance, the engine ensures that only one message can produce
  side-effects at a time. The handler always sees the current state, including
  the effects of any prior messages. Messages routed to different instances do
  not interfere with one another.

- Integration message handlers use per-**command** isolation. The engine does
  not guarantee any sequencing between command messages. It delivers each
  command message independently, and handlers must manage any necessary
  coordination or consistency themselves.

- Projection message handlers use per-**event-stream** isolation. The engine
  delivers event messages from a given stream in order and applies their
  side-effects sequentially. It may deliver events from different streams
  concurrently.

> [!IMPORTANT]
> These boundaries describe how the engine isolates side-effects, not how it
> schedules or executes code. Depending on the engine, it may handle messages in
> parallel — across threads, processes, or distributed nodes — but it still
> ensures that side-effects remain isolated and consistent.

### Event stream binding

Event stream binding describes how message handlers determine which event stream
to use when recording events. The engine maintains at least one event stream per
application, but may use multiple separate streams to improve throughput or
scalability.

Aggregate message handlers use per-**instance** binding. The engine persists all
events recorded by a specific instance to the same event stream. This guarantees
that consumers, such as projection message handlers, see events from each
instance in the order they occurred.

Integration message handlers use per-**command** binding. The engine persists
all events recorded by a specific command to the same event stream. This
guarantees that consumers see events from each command in the order they
occurred. The engine may choose a different event stream for each command.

Process and projection message handlers do not record event messages, so no
binding strategy applies.

<!-- anchors -->

[state]: #state
[isolation boundary]: #isolation-boundary
[external data access]: #external-data-access
[event stream binding]: #event-stream-binding

<!-- external links -->

[concepts]: concepts.md
[glossary]: glossary.md
[definition of state]: glossary.md#state
