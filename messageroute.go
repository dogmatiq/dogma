package dogma

import (
	"fmt"
	"reflect"
)

// HandlesCommand configures an [AggregateMessageHandler] or
// [IntegrationMessageHandler] as a consumer of [Command] messages of type T.
//
// Pass the returned [MessageRoute] to [AggregateConfigurer].Routes or
// [IntegrationConfigurer].Routes.
//
// The engine panics if the application has multiple handlers that handle T.
func HandlesCommand[T Command](...HandlesCommandOption) HandlesCommandRoute {
	return HandlesCommandRoute{typeOf[Command, T]()}
}

// RecordsEvent configures an [AggregateMessageHandler] or
// [IntegrationMessageHandler] as a producer of [Event] messages of type T.
//
// Pass the returned [MessageRoute] to [AggregateConfigurer].Routes or
// [IntegrationConfigurer].Routes.
//
// The engine panics if the application has multiple handlers that record T.
func RecordsEvent[T Event](...RecordsEventOption) RecordsEventRoute {
	return RecordsEventRoute{typeOf[Event, T]()}
}

// HandlesEvent configures a [ProcessMessageHandler] or [ProjectionMessageHandler] as a
// consumer of [Event] messages of type T.
//
// Pass the returned [MessageRoute] to [ProcessConfigurer].Routes or
// [ProjectionConfigurer].Routes.
//
// An application may have multiple handlers that handle T.
func HandlesEvent[T Event](...HandlesEventOption) HandlesEventRoute {
	return HandlesEventRoute{typeOf[Event, T]()}
}

// ExecutesCommand configures a [ProcessMessageHandler] as a producer of
// [Command] messages of type T.
//
// Pass the returned [MessageRoute] to [ProcessConfigurer].Routes.
//
// The application may have multiple handlers that execute T.
func ExecutesCommand[T Command](...ExecutesCommandOption) ExecutesCommandRoute {
	return ExecutesCommandRoute{typeOf[Command, T]()}
}

// SchedulesTimeout configures a [ProcessMessageHandler] as a scheduler of
// [Timeout] messages of type T.
//
// Pass the returned [MessageRoute] to [ProcessConfigurer].Routes.
//
// The application may have multiple handlers that schedule T.
func SchedulesTimeout[T Timeout](...SchedulesTimeoutOption) SchedulesTimeoutRoute {
	return SchedulesTimeoutRoute{typeOf[Timeout, T]()}
}

type (
	// MessageRoute is an interface for types that describe a relationship between a
	// message handler and a specific message type.
	MessageRoute interface{ isMessageRoute() }

	// HandlesCommandRoute is a [HandlerRoute] that represents a handler's
	// ability to handle [Command] messages of a specific type.
	//
	// Avoid constructing values of this type directly; use [HandlesCommand]
	// instead.
	HandlesCommandRoute struct{ Type reflect.Type }

	// ExecutesCommandRoute is a [HandlerRoute] that represents a handler's
	// ability to execute [Command] messages of a specific type.
	//
	// Avoid constructing values of this type directly; use [ExecutesCommand]
	// instead.
	ExecutesCommandRoute struct{ Type reflect.Type }

	// HandlesEventRoute is a [HandlerRoute] that represents a handler's
	// ability to handle [Event] messages of a specific type.
	//
	// Avoid constructing values of this type directly; use [HandlesEvent]
	// instead.
	HandlesEventRoute struct{ Type reflect.Type }

	// RecordsEventRoute is a [HandlerRoute] that represents a handler's
	// ability to record [Event] messages of a specific type.
	//
	// Avoid constructing values of this type directly; use [RecordsEvent]
	// instead.
	RecordsEventRoute struct{ Type reflect.Type }

	// SchedulesTimeoutRoute is a [HandlerRoute] that represents a handler's
	// ability to schedule [Timeout] messages of a specific type.
	//
	// Avoid constructing values of this type directly; use [SchedulesTimeout]
	// instead.
	SchedulesTimeoutRoute struct{ Type reflect.Type }
)

type (
	// HandlesCommandOption is an option that modifies the behavior of
	// [HandlesCommand].
	//
	// This type exists for forward-compatibility.
	HandlesCommandOption interface {
		futureHandlesCommandOption()
	}

	// ExecutesCommandOption is an option that modifies the behavior of
	// [ExecutesCommand].
	//
	// This type exists for forward-compatibility.
	ExecutesCommandOption interface {
		futureExecutesCommandOption()
	}

	// HandlesEventOption is an option that modifies the behavior of
	// [RecordsEvent].
	//
	// This type exists for forward-compatibility.
	HandlesEventOption interface {
		futureHandlesEventOption()
	}

	// RecordsEventOption is an option that modifies the behavior of
	// [ExecutesCommand].
	//
	// This type exists for forward-compatibility.
	RecordsEventOption interface {
		futureRecordsEventOption()
	}

	// SchedulesTimeoutOption is an option that modifies the behavior of
	// [SchedulesTimeout].
	//
	// This type exists for forward-compatibility.
	SchedulesTimeoutOption interface {
		futureSchedulesTimeoutOption()
	}
)

// typeOf returns the [reflect.Type] for C, which must be a concrete
// implementation of the interface I.
func typeOf[I, C Message]() reflect.Type {
	concrete := reflect.TypeFor[C]()

	if concrete.Kind() == reflect.Pointer {
		iface := reflect.TypeFor[I]()
		elem := concrete.Elem()

		if elem.Implements(iface) {
			panic(fmt.Sprintf(
				"%s implements %s using non-pointer receivers, use %s instead",
				concrete,
				iface,
				elem,
			))
		}
	}

	return concrete
}
