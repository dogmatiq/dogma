package dogma

import (
	"fmt"
	"iter"
	"maps"
	"reflect"
	"slices"
	"strings"
	"sync/atomic"
)

// RegisterCommand adds a [Command] message type to Dogma's message registry,
// making it eligible for use with [HandlesCommand] and [ExecutesCommand].
//
// id uniquely identifies the message type. The value must be a canonical RFC
// 9562 UUID string, such as "65f9620a-65c1-434e-8292-60cd7938c4de", and is
// case-insensitive. The engine uses the ID to associate message data with the
// correct Go type.
func RegisterCommand[
	T interface {
		Command
		*E
	},
	E any, // E is the "element" type of the pointer type T.
](id string, _ ...RegisterCommandOption) {
	registerMessageType[Command, T](id)
}

// RegisterCommandOption is an option that modifies the behavior of
// [RegisterCommand].
//
// This type exists for forward-compatibility.
type RegisterCommandOption interface {
	futureRegisterCommandOption()
}

// RegisterEvent adds a [Event] message type to Dogma's message registry, making
// it eligible for use with [HandlesEvent] and [RecordsEvent].
//
// id uniquely identifies the message type. The value must be a canonical RFC
// 9562 UUID string, such as "65f9620a-65c1-434e-8292-60cd7938c4de", and is
// case-insensitive. The engine uses the ID to associate message data with the
// correct Go type.
func RegisterEvent[
	T interface {
		*E
		Event
	},
	E any, // E is the "element" type of the pointer type T.
](id string, _ ...RegisterEventOption) {
	registerMessageType[Event, T](id)
}

// RegisterEventOption is an option that modifies the behavior of
// [RegisterEvent].
//
// This type exists for forward-compatibility.
type RegisterEventOption interface {
	futureRegisterEventOption()
}

// RegisterTimeout adds a [Timeout] message type to the Dogma message registry,
// making it eligible for use with [SchedulesTimeout].
//
// id uniquely identifies the message type. The value must be a canonical RFC
// 9562 UUID string, such as "65f9620a-65c1-434e-8292-60cd7938c4de", and is
// case-insensitive. The engine uses the ID to associate message data with the
// correct Go type.
func RegisterTimeout[
	T interface {
		*E
		Timeout
	},
	E any, // E is the "element" type of the pointer type T.
](id string, _ ...RegisterTimeoutOption) {
	registerMessageType[Timeout, T](id)
}

// RegisterTimeoutOption is an option that modifies the behavior of
// [RegisterTimeout].
//
// This type exists for forward-compatibility.
type RegisterTimeoutOption interface {
	futureRegisterTimeoutOption()
}

// RegisteredMessageType contains information about an implementation of [Command],
// [Event], or [Timeout] that's in Dogma's message registry.
//
// Use [RegisterCommand], [RegisterEvent], or [RegisterTimeout] to add messages
// to the registry.
type RegisteredMessageType struct {
	nocmp
	id  string
	typ reflect.Type
	new func() Message
}

// ID returns an RFC 9562 UUID that uniquely identifies the message type. The
// engine uses the ID to associate message data with the correct Go type.
func (t RegisteredMessageType) ID() string {
	return t.id
}

// GoType returns the [reflect.Type] of the message type.
func (t RegisteredMessageType) GoType() reflect.Type {
	return t.typ
}

// New returns a new instance of the message type.
//
// If the message type uses pointer receivers, it returns a non-nil pointer to a
// new zero-value of the underlying type.
func (t RegisteredMessageType) New() Message {
	return t.new()
}

// RegisteredMessageTypeFor returns the [RegisteredMessageType] for T.
//
// ok is false if T isn't in the message type registry.
func RegisteredMessageTypeFor[T Message]() (t RegisteredMessageType, ok bool) {
	queryMessageRegistry(
		func(reg *messageTypes) {
			key := reflect.TypeFor[T]()
			t, ok = reg.ByType[key]
		},
	)

	return t, ok
}

// RegisteredMessageTypeOf returns the [RegisteredMessageType] for the given
// message instance.
//
// ok is false if the message's type isn't in the message type registry.
func RegisteredMessageTypeOf(m Message) (t RegisteredMessageType, ok bool) {
	if m == nil {
		panic("message cannot be nil")
	}

	queryMessageRegistry(
		func(reg *messageTypes) {
			key := reflect.TypeOf(m)
			t, ok = reg.ByType[key]
		},
	)

	return t, ok
}

// registeredMessageTypeFor is a variant of [RegisteredMessageTypeFor] that
// panics if T isn't in the message type registry.
func registeredMessageTypeFor[T Message]() RegisteredMessageType {
	if t, ok := RegisteredMessageTypeFor[T](); ok {
		return t
	}

	panic(fmt.Sprintf(
		"%s is not in the message type registry",
		qualifiedNameOf(reflect.TypeFor[T]()),
	))
}

// RegisteredMessageTypeByID returns the [RegisteredMessageType] with the given
// ID.
//
// The ID is a canonical RFC 9562 UUID string, such as
// "65f9620a-65c1-434e-8292-60cd7938c4de", and is case-insensitive.
//
// ok is false if there is no such message type in the registry.
func RegisteredMessageTypeByID(id string) (t RegisteredMessageType, ok bool) {
	id, err := normalizeUUID(id)
	if err != nil {
		panic(err.Error())
	}

	queryMessageRegistry(
		func(reg *messageTypes) {
			t, ok = reg.ByID[id]
		},
	)

	return t, ok
}

// RegisteredMessageTypes returns an iterator that yields information about each
// message in the Dogma message registry.
//
// Use [RegisterCommand], [RegisterEvent], or [RegisterTimeout] to add messages
// to the registry.
func RegisteredMessageTypes() iter.Seq[RegisteredMessageType] {
	return func(yield func(RegisteredMessageType) bool) {
		queryMessageRegistry(
			func(reg *messageTypes) {
				for _, t := range reg.Slice {
					if !yield(t) {
						return
					}
				}
			},
		)
	}
}

// messageTypes encapsulates the Dogma message registry.
type messageTypes struct {
	ByID   map[string]RegisteredMessageType
	ByType map[reflect.Type]RegisteredMessageType
	Slice  []RegisteredMessageType
}

// messageTypeRegistry is a global registry of types that implement [Command],
// [Event], and [Timeout].
//
// The messageTypes value is immutable. Every addition to the registry creates a
// new messageTypes value that atomically replaces the old value. The registry
// doesn't support removal of registered types.
//
// Assuming that additions to the registry typically occur within init()
// functions, and that additions after module initialization are rare, this
// approach should offer good read performance (just an atomic load) with a
// minor penalty during initialization due to repeated map and slice clones -
// while still providing thread safety.
var messageTypeRegistry atomic.Pointer[messageTypes]

func queryMessageRegistry(fn func(*messageTypes)) {
	reg := messageTypeRegistry.Load()
	if reg == nil {
		reg = &messageTypes{}
	}
	fn(reg)
}

func registerMessageType[
	K Message,
	T interface {
		Message
		*E
	},
	E any,
](id string) {
	typ := reflect.TypeFor[T]()

	id, err := normalizeUUID(id)
	if err != nil {
		panic(fmt.Sprintf(
			"cannot register %s: %s",
			qualifiedNameOf(typ),
			err,
		))
	}

	mergeMessageType(RegisteredMessageType{
		id:  id,
		typ: typ,
		new: func() Message {
			return T(new(E))
		},
	})
}

func mergeMessageType(t RegisteredMessageType) {
	for {
		// Read the existing registry and construct its replacement.
		existing := messageTypeRegistry.Load()
		replacement := &messageTypes{}

		if existing == nil {
			// The registry is empty, create new data structures.
			replacement.ByID = map[string]RegisteredMessageType{}
			replacement.ByType = map[reflect.Type]RegisteredMessageType{}
		} else {
			// The registry has messages. Check for existing registrations with the
			// same ID or Go type.
			if x, ok := existing.ByType[t.typ]; ok {
				if x.id == t.id {
					panic(fmt.Sprintf(
						"cannot register %s: it is already registered",
						qualifiedNameOf(t.typ),
					))
				}

				panic(fmt.Sprintf(
					"cannot register %s: it is already registered as %q",
					qualifiedNameOf(t.typ),
					x.id,
				))
			}

			if x, ok := existing.ByID[t.id]; ok {
				panic(fmt.Sprintf(
					"cannot register %s: %q is already associated with %s",
					qualifiedNameOf(t.typ),
					t.id,
					qualifiedNameOf(x.typ),
				))
			}

			// Clone existing data structures to avoid data races with other
			// goroutines that may be reading from the registry.
			replacement.ByID = maps.Clone(existing.ByID)
			replacement.ByType = maps.Clone(existing.ByType)
			replacement.Slice = slices.Clone(existing.Slice)
		}

		// Add the new type to the registry.
		replacement.ByID[t.id] = t
		replacement.ByType[t.typ] = t
		replacement.Slice = append(replacement.Slice, t)

		if messageTypeRegistry.CompareAndSwap(existing, replacement) {
			return
		}

		// The swap failed, which means that another goroutine has
		// modified the registry since this goroutine loaded it.
	}
}

func qualifiedNameOf(t reflect.Type) string {
	var name strings.Builder

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
		name.WriteString("*")
	}

	if p := t.PkgPath(); p != "" {
		name.WriteString(p)
		name.WriteString(".")
	}

	name.WriteString(t.Name())

	return name.String()
}
