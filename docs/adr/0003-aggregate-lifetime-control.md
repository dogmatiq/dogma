# 3. Aggregate Lifetime Control

Date: 2018-12-10

## Status

Accepted

- Amended by [16. Automatic Aggregate Creation](0016-automatic-aggregate-creation.md)
- Amended by [17. Recreation of Aggregate Instances After Destruction](0017-recreate-aggregate-after-destruction.md)

## Context

We need a way to control the lifetime of an aggregate from the domain layer.

If you imagine an aggregate running atop a "CRUD-based" Dogma engine, it's easy
to see that creation and destruction of aggregate data are necessary operations.
This is less obvious in an event-sourced scenario where the notion of "deletion"
is not present.

We need to define these operations in such as a way that the domain implementor
can use the operations in a meaningful way within their domain, but engine
implementors are free to determine their own persistence semantics.

## Decision

We've opted to name the methods used to create and destroy aggregate instances
`AggregateScope.Create()` and `Destroy()`, respectively. Note that
`AggregateScope` has since been renamed to `AggregateCommandScope`.

`Create()` is a fairly self explanatory name. This is an idempotent operation.
The method returns `true` if the call actually resulted in the creation of the
instance; or `false` if the instance already exists.

`Destroy()` was chosen in preference to words such as "delete", as depending on
the engine implementation, no deletion necessarily occurs. It was chosen in
preference to "done", as it more clearly indicates that the aggregate instance
state will be "reset". This could be implemented in an ES-based engine by
recording internal events to represent the create/destroy operations, and only
loading those events that occurred since the most recent creation.

### Mandatory Events

It was also decided to require an event message be recorded in the same scope
as any successful call to `Create()` or `Destroy()`. This ensures that creation
and deletion is always represented by a domain event.

We decided against having `Create()` or `Destroy()` take an event as a
parameter for two reasons:

1. This would necessitate a further decision about `Create()` as to whether the
   event it is passed should be recorded in all cases or only if the instance
   does not already exist. Neither of which is appropriate in all cases.

2. If we decide to relax this requirement in the future, those methods would
   have to lose those event arguments, breaking backwards compatibilty.

## Consequences

The approach taken tends towards reducing the chance of BC breaks by giving
each method a single responsibility and placing the onus on the domain
implementor to combine their use correctly.

That said, the fact that `Create()` and `Destroy()` both mandate an event be
recorded, but do not enforce this in any way may provide to be a point of
confusion.

More experience is needed implementing real applications in Dogma before we are
willing to commit to adding event arguments to `Create()` and `Destroy()`.
