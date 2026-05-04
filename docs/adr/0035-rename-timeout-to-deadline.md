# 35. Rename timeout to deadline

Date: 2026-05-04

## Status

Proposed

> [!NOTE]
> This decision has not yet been accepted and is subject to change.

## Context

Dogma uses the term "timeout" for messages that notify a process handler that a
specific point in time has been reached. The term was adopted from [NServiceBus]
early in the project's history because no better alternative came to mind at the
time.

In practice, the term has three problems:

1. In Go, "deadline" is the established term for a specific wall-clock time at
   which something should happen. The `context` package and `net.Conn` both use
   "deadline" for this purpose, while "timeout" refers to a duration after
   which an operation is abandoned. Dogma's "timeout" messages behave like Go
   deadlines — they are scheduled for a specific `time.Time`, not after a
   `time.Duration`.

2. "Timeout" carries connotations of failure or abandonment — an operation that
   ran too long. While some Dogma uses fit this mold (for example, canceling an
   unpaid order), many do not. A process that schedules a reminder for 30 days
   before a contract renewal isn't "timing out" — it's reaching a deadline.

3. When discussing process behavior with domain experts, "deadline" maps more
   naturally to business language. "When does this invoice become overdue? The
   deadline for payment is January seventh" is a conversation that works. "The
   timeout for payment is January seventh" is not.

## Decision

We will rename all API symbols and documentation from "timeout" to "deadline".
This is a purely terminological change — it has no effect on runtime behavior.

The complete set of symbol renames is:

| Old                         | New                          |
| --------------------------- | ---------------------------- |
| `Timeout`                   | `Deadline`                   |
| `TimeoutValidationScope`    | `DeadlineValidationScope`    |
| `HandleTimeout()`           | `HandleDeadline()`           |
| `ScheduleTimeout()`         | `ScheduleDeadline()`         |
| `ProcessTimeoutScope`       | `ProcessDeadlineScope`       |
| `NoTimeoutMessagesBehavior` | `NoDeadlineMessagesBehavior` |
| `SchedulesTimeout()`        | `SchedulesDeadline()`        |
| `SchedulesTimeoutRoute`     | `SchedulesDeadlineRoute`     |
| `SchedulesTimeoutOption`    | `SchedulesDeadlineOption`    |
| `RegisterTimeout()`         | `RegisterDeadline()`         |
| `RegisterTimeoutOption`     | `RegisterDeadlineOption`     |

We considered `SetDeadline` instead of `ScheduleDeadline` for the
`ProcessScope` method. `SetDeadline` is more idiomatic English and matches
Go's `net.Conn.SetDeadline`. However, Go's `SetDeadline` overwrites the
previous value, whereas Dogma's method adds a new pending message — each call
produces an independent delivery. `ScheduleDeadline` preserves the existing
verb, reinforces the message-delivery semantics, and keeps the rename
mechanical.

### Documentation conventions

In addition to the symbol renames, we will adopt these conventions for how
"deadline" is used in prose:

- Deadlines are "reached", not "occurred" or "elapsed".
- Deadlines are scheduled "for" a time, not "to occur at" a time.
- Use "deadline" alone when the verb applies to the concept itself: schedule a
  deadline, reach a deadline, cancel a deadline.
- Use "deadline message" when the verb applies to the message artifact: deliver
  a deadline message, handle a deadline message, persist a deadline message.

### Historical records

Existing ADRs and CHANGELOG entries are left unchanged. They are point-in-time
records that used the vocabulary of their era.

## Consequences

Every downstream consumer of the Dogma API — engine implementations, `testkit`,
`enginekit`, and application code — will need a corresponding rename.

The glossary gains a new "deadline" entry replacing "timeout", and "process
deadline scope" replacing "process timeout scope".

<!-- references -->

[NServiceBus]: https://docs.particular.net/nservicebus/sagas/timeouts
