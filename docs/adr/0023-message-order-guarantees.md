# 23. Message Order Guarantees

Date: 2025-06-14

## Status

Accepted

- References [22. Remove CRUD Application Support](0022-remove-crud-application-support.md)

## Context

In light of the removal of CRUD application support, we are able to make more
specific guarantees about the order in which messages are delivered to an application's handlers.

## Decision

We will document the following claims about delivery order as part of the API.

1. Commands are delivered in an undefined order.
2. Events produced by integration handlers are observed in an undefined order.
3. Events produced by a single aggregate **instance** are observed in relative
   chronological order based on their "recorded at" time.
4. Events produced by different instances of the same aggregate type, or by
   different aggregate types, are observed in an undefined order.
5. Timeouts produced by a single process instance are observed in relative
   chronological order based on their "scheduled for" time.
6. Timeouts produced by different instances of the same process type, or by
   different process types, are observed in an undefined order.

## Consequences

The guarantees about the order of command and timeout messages have not changed,
and should present no issues for existing applications.

In practice, engines attempt to handle commands immediately, providing a weak
chronological order guarantee, but engines must be free to defer and retry
commands as necessary.

The guarantees about the order of events have been strengthened, in the case of
events produced by a single aggregate instance. Previously there were **no**
guarantees about the order of events, regardless of their source.

Even in the case where there is no practical change to the behavior, we can
now document these guarantees more explicitly, which may serve to make
application developers more aware of the requirement that they must
design their handlers around these guarantees.
