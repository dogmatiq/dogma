# Changing aggregate and process roots

This guide explains what can go wrong when you change the fields of an
[aggregate root] or a [process root], or how their values are calculated, and
how to make the change safely.

The engine calls the root's `MarshalBinary()` method to get its state as a byte
slice that can be persisted to the database. Later, when the instance handles
another message, the engine reads those bytes from the database and passes them
to the `UnmarshalBinary()` method of a new root. The bytes stay in the database
when you deploy a new version of your application, so the new version's
`UnmarshalBinary()` method can be called with data that was produced by the old
version's `MarshalBinary()` method.

Some changes to a root are safe: the new version can read data produced by the
old version and restore the root's state correctly. However, others are not
safe: the new version reads the old data without an error, but the root ends up
with the wrong values, and nothing reports the problem.

These problems are not unique to Dogma. Any application that stores data and
reads it back with newer code has them.

The format of a root's bytes is determined solely by its `MarshalBinary()`
method. You might write that method yourself, or generate it with [primo], which
generates the method for Protocol Buffers messages. The engine doesn't
understand or inspect the bytes in any way. This guide gives advice for roots
that use JSON or [Protocol Buffers].

## Contents

- [1. How the engine stores a root]
  - [1.1. Aggregates]
  - [1.2. Processes]
- [2. What can happen to old data]
- [3. Changes that cause a misread]
  - [3.1. Changes to what the state means]
  - [3.2. Changes to the format]
- [4. Changing an aggregate root]
- [5. Changing a process root]
  - [5.1. Make only changes that leave old state correct]
  - [5.2. Convert old state]
  - [5.3. Let existing instances finish with the old behavior]
  - [5.4. Add a new process handler]
- [6. Deploying the change]

## 1. How the engine stores a root

### 1.1. Aggregates

Strictly speaking, the engine doesn't need to store an aggregate root at all. It
can build the root from the historical events of the [aggregate instance]. It
calls the handler's `New()` method to get a new root, then passes each event to
the root's `ApplyEvent()` method, oldest first.

In practice, applying every event each time an engine loads an instance gets
slower as the instance records more events. So most engines also store a
**snapshot** — a copy of the root, from its `MarshalBinary()` method. To load an
instance from a snapshot, an engine calls `New()`, passes the snapshot to the
new root's `UnmarshalBinary()` method, and then applies only the events recorded
after the snapshot.

If `UnmarshalBinary()` returns an error, a common approach is for the engine to
ignore the snapshot and build the root from the instance's historical events
instead. This is what [runkit] does.

### 1.2. Processes

A process doesn't record events, so the engine can't build a process root from
historical events. Instead, the engine stores the root's state each time a
handler changes it by calling [`Mutate()`], and that stored state is the only
copy.

To load a [process instance], the engine calls `New()`, and passes the stored
state to the new root's `UnmarshalBinary()` method. When an instance ends, the
engine discards its state.

If `UnmarshalBinary()` returns an error, the engine can't load the instance, and
leaves the stored state as it is. The instance can't handle any messages until
you deploy a version that can read its state.

## 2. What can happen to old data

When the new version loads old data, `UnmarshalBinary()` has one of three
outcomes.

| Outcome                                                                                                                                                 | For an aggregate                                                                                                                                                  | For a process                                                                                                                                           |
| ------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **It reads the data correctly.** The root has the values that the new version would have produced from the same messages.                               | No problem.                                                                                                                                                       | No problem.                                                                                                                                             |
| **It returns an error.**                                                                                                                                | The engine ignores the snapshot and builds the root from the instance's historical events. Loading takes longer, but the root is correct.                         | The instance handles nothing until you deploy a version that can read its state. The stored state isn't harmed.                                         |
| **It misreads the data.** It returns no error, but the root's values differ from those that the new version would have produced from the same messages. | The handler makes its decisions from wrong values, and can record the wrong events. Recorded events are permanent. The next snapshot is a copy of the wrong root. | The handler makes its decisions from wrong values. The next call to `Mutate()` replaces the stored state, which was the only copy, with the wrong root. |

A misread is the worst outcome. The rest of this guide shows which changes cause
a misread, and how to get one of the other two outcomes instead:

- For an aggregate, make `UnmarshalBinary()` reject any snapshot from the old
  version. See [changing an aggregate root].
- For a process, make `UnmarshalBinary()` read any old state correctly. See
  [changing a process root].

## 3. Changes that cause a misread

Every change that causes a misread is one of two kinds:

1. **Changes to what the state means** — from the same messages, the new version
   would produce different values than the old version did. A change of this
   kind causes a misread with any format.
2. **Changes to the format** — the new version's `UnmarshalBinary()` method
   reads different values from the bytes than the old version's
   `MarshalBinary()` method wrote. Whether a change is of this kind depends on
   the format.

### 3.1. Changes to what the state means

Here are some examples of this kind of change:

- Changing how a field's value is calculated. For example, the root of a
  shopping cart aggregate has a field that holds the number of items in the
  cart. The old version counts every item, and the new version leaves out the
  items that are gifts. Old data holds the value from the old calculation.
- Changing what a field's value measures. For example, a delay in seconds
  becomes a delay in milliseconds. Old data holds a number of seconds, and the
  new version reads it as a number of milliseconds.
- Adding a field that earlier messages would have set to something other than
  its default value.

The first two leave the root's fields exactly as they were, which makes them
easy to miss.

Adding a field is the most common of the three. Every field in a root gets its
value from the messages handled by the instance. Old data doesn't contain a new
field, so the new field gets its default value instead, which is usually zero.
Adding a field is only safe if that default value is the right value for every
instance that already exists.

For example, suppose that the new version introduces a `GiftWrapRequested`
event, and adds a field that counts how many of those events there have been.
Adding the field is safe. No `GiftWrapRequested` event existed before this
version, so zero is the right value for every existing instance.

Here is an example of adding a field that isn't safe, in the root of a shopping
cart aggregate. The change adds the `TotalQuantity` field:

```go
type cart struct {
	// Items maps the ID of each product in the cart to its quantity.
	Items map[string]int

	// TotalQuantity is the sum of the quantities in Items. [cart.ApplyEvent]
	// adds to it for each item that is added to the cart.
	TotalQuantity int
}
```

A snapshot from the old version has no `TotalQuantity`. Suppose that it's the
snapshot of a cart with eight items. When the new version loads the snapshot,
`TotalQuantity` is zero instead of eight, and the events recorded after the
snapshot add to that zero. If the handler uses `TotalQuantity` to limit each
cart to ten items, this cart can now grow to eighteen.

### 3.2. Changes to the format

Most of these changes depend on the format, but two cause a misread with any
format:

- Reusing a removed field's name or number. JSON identifies each field in the
  bytes by its name, and Protocol Buffers identifies each field by its number.
  If a new field has the same name or number as a field that you removed, old
  data can still hold the removed field's value, and the new field gets that
  value.
- Upgrading a dependency. If a field's type comes from another package, a new
  version of that package can marshal the type differently. The new version of
  your application then reads the old data differently, even though your own
  code is the same.

JSON and Protocol Buffers are both designed so that code can read data written
for a different definition of the same type. So most changes to a root's fields
don't cause an error, even when they cause a misread. The tables below show what
each format does with old data after each kind of change to a root's fields.

Each result is one of the outcomes in [what can happen to old data]: the new
version reads the data correctly, returns an error, or misreads the data. The
_Safe_ column shows whether the change avoids a misread:

- ✅ — the new version reads the data correctly, or `UnmarshalBinary()` returns
  an error. An error is the outcome you want for an aggregate, and is far better
  than a misread for a process.
- ❌ — the new version can misread the data without an error.
- ⚠️ — it depends on the types involved.

#### 3.2.1. JSON

This table assumes that the root's `UnmarshalBinary()` method calls
[`json.Unmarshal()`] to read the bytes into the root.

| Change to the root's fields                                                                             | Safe | Result for old data                                                                                                                                                                                 |
| ------------------------------------------------------------------------------------------------------- | ---- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Add a field.                                                                                            | ❌   | **Misread, unless the value from `New()` is right.** The field keeps the value from `New()`, which is usually zero.                                                                                 |
| Change the value that `New()` gives a field.                                                            | ❌   | **Misread.** A field that isn't in the old data keeps the new value from `New()`. That includes a field that the old version left out because its value was zero and it has the `omitempty` option. |
| Remove a field.                                                                                         | ✅   | **Correct.** The stored value is ignored.                                                                                                                                                           |
| Rename a field that has no `json` tag.                                                                  | ❌   | **As for adding a field.** JSON identifies a value by the name of its field, so the result is the same as removing the field and adding another. A [`json` tag] with the old name prevents this.    |
| Use the name of a removed field for a new field.                                                        | ⚠️   | **Misread**, if the new field's type accepts the stored value. The new field gets the removed field's value. Otherwise, an **error**.                                                               |
| Change a field's type to one that still accepts the stored value. For example, from `int` to `float64`. | ✅   | **Correct.**                                                                                                                                                                                        |
| Change a field's type to one that doesn't accept the stored value. For example, from `int` to `string`. | ✅   | **Error.** The stored value can't be read as the new type.                                                                                                                                          |

#### 3.2.2. Protocol Buffers

This table assumes that the root's `UnmarshalBinary()` method calls
[`proto.Unmarshal()`] to read the bytes into the root's message. The
`UnmarshalBinary()` method that [primo] generates does exactly that.

> [!IMPORTANT]
> None of these changes can be relied on to cause an error. The rows marked ✅
> are safe only because the new version reads the data correctly. So with
> Protocol Buffers, old data causes an error only if your own code checks it. To
> add that check, see [changing an aggregate root].

| Change to the message definition                                                                | Safe | Result for old data                                                                                                                                                                                                                                         |
| ----------------------------------------------------------------------------------------------- | ---- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Add a field.                                                                                    | ❌   | **Misread, unless the default value is right.** The field has its [default value], which is usually zero. `proto.Unmarshal()` clears the message before it reads the data, so the field doesn't keep a value from `New()`.                                  |
| Remove a field.                                                                                 | ✅   | **Correct.** The stored value becomes an [unknown field]. Your code doesn't see it, but it's stored again each time the root is marshaled.                                                                                                                  |
| Rename a field.                                                                                 | ✅   | **Correct.** The binary format identifies a field by its number, not by its name.                                                                                                                                                                           |
| Change a field's number.                                                                        | ❌   | **As for adding a field.** The result is the same as removing the field and adding another.                                                                                                                                                                 |
| Use the number of a removed field for a new field.                                              | ❌   | **Misread**, if the types of the two fields have the same [wire type]. The new field gets the removed field's stored bytes, read as the new type. Otherwise, as for adding a field. To prevent this, mark the number as [reserved] when you remove a field. |
| Change a field's type to one with the same wire type. For example, from `int64` to `int32`.     | ❌   | **Misread for some values.** The stored bytes are read as the new type, so the value can change.                                                                                                                                                            |
| Change a field's type to one with a different wire type. For example, from `string` to `int64`. | ❌   | **As for adding a field.** The stored value becomes an unknown field.                                                                                                                                                                                       |
| Remove a value from an enum.                                                                    | ❌   | **Misread.** The field holds the stored number, which no longer has a name.                                                                                                                                                                                 |

[Updating a message type], in the Protocol Buffers documentation, describes each
of these changes in more detail.

## 4. Changing an aggregate root

The engine can always build an aggregate root again, so the safe thing to do
with any snapshot from the old version is to reject it. When `UnmarshalBinary()`
returns an error, the engine builds the root by applying the instance's
historical events with the new version.

The engine doesn't record which version of the application wrote a snapshot, so
the root needs its own way to recognize old snapshots, or the aggregate needs to
avoid them. The right strategy depends on your aggregate. For example, you
could:

- Store a version number in each snapshot, and reject any snapshot with a
  different number. You don't have to add a version number in advance. You can
  add it in the same deployment as the first change that would cause a misread.
  Snapshots stored before then have no version number, so the new version
  rejects them too.
- Check for a field that old snapshots don't have. With JSON, a pointer field is
  nil when the data doesn't contain it, and with Protocol Buffers, so is a field
  marked `optional`. This only detects an added field, not a change to what a
  field means.
- Not store snapshots at all, by embedding [`dogma.NoSnapshotBehavior`] in the
  root. Then there are no old snapshots to misread. This suits an aggregate
  whose instances record only a few events.

Whichever check you choose, it goes in the root's `UnmarshalBinary()` method.
That's a problem if you're using [primo] to generate the root's
`MarshalBinary()` and `UnmarshalBinary()` methods, because there's no place in a
generated method for code of your own. Until primo supports this, you need a
workaround, such as:

- Compiling the root's `.proto` file without primo, and declaring the methods on
  the generated type yourself.
- Making the root a separate type that embeds the generated type, and giving the
  root its own methods. The embedded field is a pointer, so the handler's
  `New()` method must set it.

These proposals would remove the need for a workaround:

- [dogmatiq/primo#78] — primo would generate a version check.
- [dogmatiq/primo#80] — the methods that primo generates would call code of your
  own.
- [dogmatiq/dogma#226] — roots would have a `Validate()` method, which the
  engine would call after `UnmarshalBinary()`.
- [dogmatiq/dogma#234] — the engine would store a version number with each
  root's data, and act on it.

## 5. Changing a process root

You can recognize the state stored by the old version in the same ways as [an
aggregate's snapshots], such as a version number, or a check for a field that
old state doesn't have. But you can't reject old process state the way you can
reject an old snapshot. A process has no historical events, so the engine can't
rebuild a process root, and the stored state is the only copy. An error from
`UnmarshalBinary()` stops the instance until you deploy a version that can read
the state.

So you need another approach, and the right one depends on your process. For
example, you could:

- [Make only changes that leave old state correct], so that `UnmarshalBinary()`
  needs no extra code.
- [Convert old state] into the current form when `UnmarshalBinary()` reads it.
- [Let existing instances finish with the old behavior].
- [Add a new process handler].

The rest of this section describes each approach, and what to watch out for.

### 5.1. Make only changes that leave old state correct

With this approach, you make only changes that the new version reads correctly
without any extra code in `UnmarshalBinary()`. To do that, follow these rules:

- Only add a new field if its default value is correct for every existing
  instance. Old state doesn't contain the new field, so the field gets its
  default value. [Changes to the format] shows what that is for each format.
- Don't change what a field means, or how its value is calculated.
- Make only the changes to the format that the tables in [Changes to the format]
  mark ✅.

### 5.2. Convert old state

With this approach, the new version's `UnmarshalBinary()` method recognizes old
state and converts it into the current form.

For example, suppose that the root of a shopping cart checkout process has a
`Paid` flag, which is `true` once the customer has paid, and that a new version
replaces the flag with a `Status` enumeration, whose values are
`AWAITING_PAYMENT`, `PAID` and `SHIPPED`.

The old state has a `Paid` value, but no `Status` value. Without a conversion,
every instance in progress would get the default status, even an instance whose
customer has already paid.

To avoid this, `UnmarshalBinary()` can convert the old state, by setting
`Status` from the stored `Paid` value: `PAID` if it's `true`, and
`AWAITING_PAYMENT` if it's `false`.

This approach suits changes that don't follow the [rules of the previous
approach][make only changes that leave old state correct]. But a conversion can
only use what the old state contains. Suppose that the new version adds a field
whose value depends on the events and deadlines handled by the old version, and
the old state holds nothing about them. No conversion can give that field the
correct value, so you need one of the other approaches in this section.

If you use this approach, keep these things in mind:

- **Keep each conversion while old state might still need it**

  Old state stays stored until its instance ends, or until the handler next
  calls `Mutate()` for that instance. An instance that is waiting for an event
  can keep the same state for a long time.

  So conversions build up: after several versions of your application that each
  need a conversion, stored state can be in several old forms.
  `UnmarshalBinary()` then needs a conversion for each form, and a way to tell
  the forms apart, such as a version number.

  runkit can't tell you which forms of old state are still stored. So only
  remove a conversion when you're sure that no instance still has state in that
  form. [dogmatiq/dogma#234] proposes that engines report how many instances
  have state in each form.

- **Make each conversion deterministic**

  A conversion must give the same result from the same old state every time. For
  example, don't base it on the current time, or on data from outside the state.
  The engine stores the converted root only when the handler next calls
  `Mutate()`. Until then, `UnmarshalBinary()` converts the same old state each
  time the engine loads the instance. A conversion that gave a different result
  on each load could make the handler act inconsistently.

- **Return an error for state that the code doesn't recognize**

  This protects the old version, not the new one. During a rolling deployment,
  and after a rollback, the old version loads state written by the new version.
  If the old version's `UnmarshalBinary()` returns an error for state that it
  doesn't recognize, the instance stops, but the stored state isn't harmed, and
  the instance carries on when the new version loads it. Without the error, the
  old version could misread the state and store a wrong root in its place. The
  old version only has this check if you added it in an earlier deployment, as
  [deploying the change] explains.

- **Test each conversion with an example of the old state**

  Tests only ever see state written by the current version, including tests that
  use [testkit]. So unless a test supplies old state itself, a conversion first
  runs in production. To test a conversion, keep a sample of the bytes that the
  old version wrote, for example in a file, and check the root that
  `UnmarshalBinary()` produces from it.

If you're using primo to generate the root's marshaling methods, the conversion
needs one of the workarounds described for [an aggregate's snapshots].

With Protocol Buffers, keep a field in the message definition while any
conversion reads it, and mark the field as [deprecated]. When no conversion
reads a field any more, remove the field and mark its number as [reserved].

### 5.3. Let existing instances finish with the old behavior

With this approach, the handler keeps both behaviors. Instances that the old
version started carry on with the old behavior until they end. Instances that
the new version starts use the new behavior. The old behavior doesn't need the
information that's missing from old state, so nothing has to fill it in.

The handler needs to know which behavior each instance uses. One way is to keep
the state that only the new behavior needs in a separate field. The handler sets
that field when it starts a new instance, and `New()` leaves it unset. Old state
doesn't have the field, so an instance without it uses the old behavior.

With Protocol Buffers, a `oneof` that holds either the old state or the new
state makes the choice explicit. Some ways of moving existing fields into a
`oneof` aren't safe, so check [updating a message type] first.

When no existing instance still has old state, you can remove the old behavior.

### 5.4. Add a new process handler

With this approach, you give the new behavior to a new process handler with a
different identity. [Dogma delivers the complete history of events] to a handler
with a new identity, so the new handler receives every event that the old
handler received, as well as new ones. What it does with the old events depends
on whether you keep the old handler:

- **Keep the old handler for existing instances**

  The two handlers must agree on which instances each one handles. For example,
  both handlers can compare the time at which each event was recorded with a
  fixed cutoff time. The new handler ignores the events recorded before the
  cutoff, and the old handler ignores any event recorded after it that would
  start a new instance. Remove the old handler only when you're sure that none
  of its instances are still in progress.

- **Remove the old handler**

  In this case, the new handler takes over every workflow, so it handles the old
  events as well as new ones. Its instances are separate from the old handler's,
  so it starts its own instance for every workflow that the old handler ever
  handled, including workflows that ended long ago. Those instances execute
  their commands and schedule their deadlines again. A deadline scheduled for a
  time that has already passed is delivered straight away.

> [!CAUTION]
> Only remove the old handler if everything that the new handler repeats is
> idempotent, which means that doing it a second time has no further effect. In
> practice, every handler of the commands that the new handler executes must
> ignore any command whose effect has already happened.

## 6. Deploying the change

During a rolling deployment, the old version and the new version of your
application run at the same time, and load the same instances. For example,
[runkit] can run several engines for one application against the same database.
After a rollback, the old version loads data from the new version. In both
cases, the old version can misread that data in the ways described in this
guide.

A check, such as a version number, only protects a version of your application
that contains it. If the new version is the first to have the check, the old
version has none, and it can still misread data from the new version. There are
two ways to deal with this:

- Add the check in a deployment of its own, before the deployment that makes the
  change. Then the old version already has the check when the change is
  deployed.
- Stop every engine that runs the old version before you start any engine that
  runs the new version, and don't roll back to the old version afterward.

<!-- anchors -->

[1. how the engine stores a root]: #1-how-the-engine-stores-a-root
[1.1. aggregates]: #11-aggregates
[1.2. processes]: #12-processes
[2. what can happen to old data]: #2-what-can-happen-to-old-data
[3. changes that cause a misread]: #3-changes-that-cause-a-misread
[3.1. changes to what the state means]: #31-changes-to-what-the-state-means
[3.2. changes to the format]: #32-changes-to-the-format
[4. changing an aggregate root]: #4-changing-an-aggregate-root
[5. changing a process root]: #5-changing-a-process-root
[5.1. make only changes that leave old state correct]: #51-make-only-changes-that-leave-old-state-correct
[5.2. convert old state]: #52-convert-old-state
[5.3. let existing instances finish with the old behavior]: #53-let-existing-instances-finish-with-the-old-behavior
[5.4. add a new process handler]: #54-add-a-new-process-handler
[6. deploying the change]: #6-deploying-the-change
[add a new process handler]: #54-add-a-new-process-handler
[an aggregate's snapshots]: #4-changing-an-aggregate-root
[changes to the format]: #32-changes-to-the-format
[changing a process root]: #5-changing-a-process-root
[changing an aggregate root]: #4-changing-an-aggregate-root
[convert old state]: #52-convert-old-state
[deploying the change]: #6-deploying-the-change
[let existing instances finish with the old behavior]: #53-let-existing-instances-finish-with-the-old-behavior
[make only changes that leave old state correct]: #51-make-only-changes-that-leave-old-state-correct
[what can happen to old data]: #2-what-can-happen-to-old-data

<!-- references -->

[aggregate instance]: ../glossary.md#aggregate-instance
[aggregate root]: ../glossary.md#aggregate-root
[default value]: https://protobuf.dev/programming-guides/proto3/#default
[deprecated]: https://protobuf.dev/programming-guides/proto3/#options
[dogma delivers the complete history of events]: https://pkg.go.dev/github.com/dogmatiq/dogma#ProcessMessageHandler
[`dogma.NoSnapshotBehavior`]: https://pkg.go.dev/github.com/dogmatiq/dogma#NoSnapshotBehavior
[dogmatiq/dogma#226]: https://github.com/dogmatiq/dogma/issues/226
[dogmatiq/dogma#234]: https://github.com/dogmatiq/dogma/issues/234
[dogmatiq/primo#78]: https://github.com/dogmatiq/primo/issues/78
[dogmatiq/primo#80]: https://github.com/dogmatiq/primo/issues/80
[`json` tag]: https://pkg.go.dev/encoding/json#Marshal
[`json.Unmarshal()`]: https://pkg.go.dev/encoding/json#Unmarshal
[`Mutate()`]: https://pkg.go.dev/github.com/dogmatiq/dogma#ProcessScope.Mutate
[primo]: https://github.com/dogmatiq/primo
[process instance]: ../glossary.md#process-instance
[process root]: ../glossary.md#process-root
[`proto.Unmarshal()`]: https://pkg.go.dev/google.golang.org/protobuf/proto#Unmarshal
[protocol buffers]: https://protobuf.dev/
[reserved]: https://protobuf.dev/programming-guides/proto3/#fieldreserved
[runkit]: https://github.com/dogmatiq/runkit
[testkit]: https://github.com/dogmatiq/testkit
[unknown field]: https://protobuf.dev/programming-guides/proto3/#unknowns
[updating a message type]: https://protobuf.dev/programming-guides/proto3/#updating
[wire type]: https://protobuf.dev/programming-guides/encoding/#structure
