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
// 4122 UUID string, such as "65f9620a-65c1-434e-8292-60cd7938c4de", and is
// case-insensitive. The engine uses the ID to associate message data with the
// correct Go type.
func RegisterCommand[T Command](id string, _ ...RegisterCommandOption) {
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
// 4122 UUID string, such as "65f9620a-65c1-434e-8292-60cd7938c4de", and is
// case-insensitive. The engine uses the ID to associate message data with the
// correct Go type.
func RegisterEvent[T Event](id string, _ ...RegisterEventOption) {
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
// 4122 UUID string, such as "65f9620a-65c1-434e-8292-60cd7938c4de", and is
// case-insensitive. The engine uses the ID to associate message data with the
// correct Go type.
func RegisterTimeout[T Timeout](id string, _ ...RegisterTimeoutOption) {
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

// ID returns an RFC 4122 UUID that uniquely identifies the message type. The
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
// It panics if T isn't in the message type registry. See [RegisterCommand],
// [RegisterEvent], and [RegisterTimeout].
func RegisteredMessageTypeFor[T Message]() (RegisteredMessageType, bool) {
	reg := messageTypeRegistry.Load()
	key := reflect.TypeFor[T]()

	t, ok := reg.ByType[key]
	return t, ok
}

// registeredMessageTypeFor is a variant of [RegisteredMessageTypeFor] that
// panics if T isn't in the message type registry.
func registeredMessageTypeFor[T Message]() RegisteredMessageType {
	if t, ok := RegisteredMessageTypeFor[T](); ok {
		return t
	}

	panic(fmt.Sprintf(
		"%s is not in the message type registry, use dogma.Register%s() to add it",
		qualifiedNameOf(reflect.TypeFor[T]()),
		messageKindFor[T]().Name(),
	))
}

// RegisteredMessageTypeByID returns the [RegisteredMessageType] with the given
// ID.
//
// The ID is a canonical RFC 4122 UUID string, such as
// "65f9620a-65c1-434e-8292-60cd7938c4de", and is case-insensitive.
func RegisteredMessageTypeByID(id string) (RegisteredMessageType, bool) {
	id, err := normalizeUUID(id)
	if err != nil {
		panic(err.Error())
	}

	t, ok := messageTypeRegistry.Load().ByID[id]
	return t, ok
}

// RegisteredMessageTypes returns an iterator that yields information about each
// message in the Dogma message registry.
//
// Use [RegisterCommand], [RegisterEvent], or [RegisterTimeout] to add messages
// to the registry.
func RegisteredMessageTypes() iter.Seq[RegisteredMessageType] {
	return slices.Values(messageTypeRegistry.Load().Slice)
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

func registerMessageType[K, T Message](id string) {
	typ := reflect.TypeFor[T]()

	id, err := normalizeUUID(id)
	if err != nil {
		panic(fmt.Sprintf(
			"cannot register %s: %s",
			qualifiedNameOf(typ),
			err,
		))
	}

	t := RegisteredMessageType{
		id:  id,
		typ: typ,
	}

	switch typ.Kind() {
	case reflect.Interface:
		panic(fmt.Sprintf(
			"cannot register %s: message type is an interface, expected a concrete type",
			qualifiedNameOf(typ),
		))

	case reflect.Pointer:
		elem := typ.Elem()
		kind := messageKindFor[T]()

		if elem.Implements(kind) {
			panic(fmt.Sprintf(
				"cannot register %s: message type uses non-pointer receivers, use %s (non-pointer) instead",
				qualifiedNameOf(typ),
				qualifiedNameOf(elem),
			))
		}

		t.new = func() Message {
			// There's no way to get the elem's type statically while still
			// supporting both pointer and non-pointer receivers, so the
			// implementation must use reflection to construct new instances.
			return reflect.New(elem).Interface().(Message)
		}

	default:
		t.new = func() Message {
			var zero T
			return zero
		}
	}

	mergeMessageType(t)
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

var (
	commandKind = reflect.TypeFor[Command]()
	eventKind   = reflect.TypeFor[Event]()
	timeoutKind = reflect.TypeFor[Timeout]()
)

func messageKindFor[T Message]() reflect.Type {
	t := reflect.TypeFor[T]()

	switch {
	case t.Implements(commandKind):
		return commandKind
	case t.Implements(eventKind):
		return eventKind
	case t.Implements(timeoutKind):
		return timeoutKind
	default:
		panic(fmt.Sprintf(
			"%s does not implement dogma.Command, dogma.Event, or dogma.Timeout",
			qualifiedNameOf(t),
		))
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
